package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	cfg   *Config
	ctx   whisper.Context
	model whisper.Model
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
	e.model, err = whisper.New(e.cfg.Model)
	if err != nil {
		return err
	}

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

	e.ctx, err = e.model.NewContext()
	if err != nil {
		return err
	}

	e.ctx.SetThreads(e.cfg.Threads)

	log.Info().Msgf("%s", e.ctx.SystemInfo())

	if e.cfg.Language != "" {
		_ = e.ctx.SetLanguage(e.cfg.Language)
	}

	log.Debug().Msg("start transcribe process")
	e.ctx.ResetTimings()
	return e.ctx.Process(data, nil)
}

// getOutputPath function determines the output path for the engine's output.
// If a specific output path is provided, it uses that path. Otherwise,
// it derives the output path by removing the file extension from the AudioPath
// and appending the specified OutputFormat.
func (e *Engine) getOutputPath() string {
	if e.cfg.OutputPath != "" {
		return e.cfg.OutputPath
	}

	ext := filepath.Ext(e.cfg.AudioPath)
	base := strings.TrimSuffix(e.cfg.AudioPath, ext)

	return base + "." + e.cfg.OutputFormat
}

// Save saves the speech result to file.
func (e *Engine) Save() error {
	outputPath := e.getOutputPath()
	log.Debug().
		Str("output-path", outputPath).
		Str("output-format", e.cfg.OutputFormat).
		Msg("start save to file process")
	e.ctx.PrintTimings()
	text := ""
	switch OutputFormat(e.cfg.OutputFormat) {
	case FormatSrt:
		n := 1
		for {
			segment, err := e.ctx.NextSegment()
			if err != nil {
				break
			}
			text += fmt.Sprintf("%d\n", n)
			text += fmt.Sprintf("%s --> %s\n", srtTimestamp(segment.Start), srtTimestamp(segment.End))
			text += segment.Text + "\n\n"
			n++
			log.Info().Msgf(
				"[%6s -> %6s] %s",
				segment.Start.Truncate(time.Millisecond),
				segment.End.Truncate(time.Millisecond),
				segment.Text,
			)
		}
	case FormatTxt:
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
	}

	if err := os.WriteFile(e.getOutputPath(), []byte(text), 0o644); err != nil {
		return err
	}

	return nil
}

// Close closes the engine.
func (e *Engine) Close() error {
	if e.ctx == nil {
		return nil
	}

	return e.model.Close()
}

func cb(segment whisper.Segment) {
	log.Info().Msgf(
		"[%6s -> %6s] %s",
		segment.Start.Truncate(time.Millisecond),
		segment.End.Truncate(time.Millisecond),
		segment.Text,
	)
	return
}
