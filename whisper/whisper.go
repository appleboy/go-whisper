package whisper

import (
	"fmt"
	"os"
	"path"
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
	cfg      *Config
	ctx      whisper.Context
	model    whisper.Model
	segments []whisper.Segment
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
	e.ctx.SetSpeedup(e.cfg.SpeedUp)
	e.ctx.SetTranslate(e.cfg.Translate)

	log.Info().Msgf("%s", e.ctx.SystemInfo())

	if e.cfg.Language != "" {
		_ = e.ctx.SetLanguage(e.cfg.Language)
	}

	log.Debug().Msg("start transcribe process")
	e.ctx.ResetTimings()
	if err := e.ctx.Process(data, nil); err != nil {
		return err
	}
	e.ctx.PrintTimings()

	for {
		segment, err := e.ctx.NextSegment()
		if err != nil {
			break
		}
		e.segments = append(e.segments, segment)
	}

	return nil
}

// getOutputPath is a method of the Engine struct that takes a format string as input.
// It returns the output path for the converted audio file based on the given format.
func (e *Engine) getOutputPath(format string) string {
	// Get the file extension of the audio file from the configuration.
	ext := filepath.Ext(e.cfg.AudioPath)
	// Get the base name of the audio file from the configuration.
	filename := filepath.Base(e.cfg.AudioPath)
	// Get the directory path of the audio file from the configuration.
	folder := filepath.Dir(e.cfg.AudioPath)
	// If the OutputFolder field in the configuration is not empty,
	// use it as the folder for the output file.
	if e.cfg.OutputFolder != "" {
		folder = e.cfg.OutputFolder
	}

	// Join the folder path, the base name of the audio file without its extension,
	// and the new format to create the output path for the converted audio file.
	return path.Join(folder, strings.TrimSuffix(filename, ext)+"."+format)
}

// Save saves the speech result to file.
func (e *Engine) Save(format string) error {
	outputPath := e.getOutputPath(format)
	log.Info().
		Str("output-path", outputPath).
		Str("output-format", format).
		Msg("save text to file")
	text := ""
	switch OutputFormat(format) {
	case FormatSrt:
		for i, segment := range e.segments {
			text += fmt.Sprintf("%d\n", i+1)
			text += fmt.Sprintf("%s --> %s\n", srtTimestamp(segment.Start), srtTimestamp(segment.End))
			text += segment.Text + "\n\n"

		}
	case FormatTxt:
		for _, segment := range e.segments {
			text += segment.Text
		}
	case FormatCSV:
		text = "start,end,text\n"
		for _, segment := range e.segments {
			text += fmt.Sprintf("%s,%s,\"%s\"\n", segment.Start, segment.End, segment.Text)
		}
	}

	if err := os.WriteFile(outputPath, []byte(text), 0o644); err != nil {
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
