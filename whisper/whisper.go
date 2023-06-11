package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
	"github.com/rs/zerolog/log"
)

// New for creating a new whisper engine.
func New(cfg *Config) (*Engine, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Engine{
		cfg: cfg,
	}, nil
}

// Engine is the whisper engine.
type Engine struct {
	cfg *Config
	ctx whisper.Context
}

// Transcribe converts audio to text.
func (e *Engine) Transcript() error {
	var data []float32
	var err error

	dir, err := os.MkdirTemp("", "whisper")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	convertedPath := filepath.Join(dir, "converted.wav")

	log.Debug().Msg("start convert audio to wav")
	if err := audioToWav(e.cfg.AudioPath, convertedPath); err != nil {
		return err
	}

	// open converted file
	fh, err := os.Open(convertedPath)
	if err != nil {
		return err
	}
	defer fh.Close()

	// Load the model
	model, err := whisper.New(e.cfg.Model)
	if err != nil {
		return err
	}
	defer model.Close()

	// Decode the WAV file - load the full buffer
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return err
	} else if dec.SampleRate != whisper.SampleRate {
		return fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}

	e.ctx, err = model.NewContext()
	if err != nil {
		return err
	}

	e.ctx.SetThreads(e.cfg.Threads)

	log.Info().Msgf("%s", e.ctx.SystemInfo())

	if e.cfg.Language != "" {
		_ = e.ctx.SetLanguage(e.cfg.Language)
	}

	log.Debug().Msg("start transcribe process")
	return e.ctx.Process(data, nil)
}

func (e *Engine) Save() error {
	log.Debug().Msg("start save process")
	text := ""
	for {
		segment, err := e.ctx.NextSegment()
		if err != nil {
			break
		}
		text += segment.Text
		log.Info().Msgf(
			"[%6s -> %6s] %s",
			segment.Start.Truncate(time.Millisecond),
			segment.End.Truncate(time.Millisecond),
			segment.Text,
		)
	}

	if e.cfg.OutputPath != "" {
		if err := os.WriteFile(e.cfg.OutputPath, []byte(text), 0o644); err != nil {
			return err
		}
	}

	return nil
}
