package whisper

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/appleboy/go-whisper/config"
	"github.com/appleboy/go-whisper/webhook"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
	"github.com/rs/zerolog/log"
)

type OutputFormat string

func (f OutputFormat) String() string {
	return string(f)
}

var (
	FormatTxt OutputFormat = "txt"
	FormatSrt OutputFormat = "srt"
	FormatCSV OutputFormat = "csv"
)

type request struct {
	Progress int `json:"progress"`
}

// New for creating a new whisper engine.
func New(cfg *config.Whisper, webhook *webhook.Client) (*Engine, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Engine{
		cfg:     cfg,
		webhook: webhook,
	}, nil
}

// Engine is the whisper engine.
type Engine struct {
	cfg      *config.Whisper
	webhook  *webhook.Client
	ctx      whisper.Context
	model    whisper.Model
	segments []whisper.Segment
	progress int
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

	if e.cfg.Language != "" && e.cfg.Language != "auto"{
		_ = e.ctx.SetLanguage(e.cfg.Language)
	}

	log.Debug().Msg("start transcribe process")
	e.ctx.ResetTimings()
	if err := e.ctx.Process(data, e.cbSegment(), e.cbProgress()); err != nil {
		return err
	}
	e.ctx.PrintTimings()

	return nil
}

// cbSegment is a method of the Engine struct that returns a function.
// The function takes a segment whisper.Segment as input and returns nothing.
// It appends the given segment to the segments field of the Engine struct.
// If the PrintSegment field in the configuration is true, it prints the segment.
// The segment is printed with the start and end time truncated to milliseconds.
func (e *Engine) cbSegment() func(segment whisper.Segment) {
	return func(segment whisper.Segment) {
		e.segments = append(e.segments, segment)
		if !e.cfg.PrintSegment {
			return
		}
		log.Info().Msgf(
			"[%6s -> %6s] %s",
			segment.Start.Truncate(time.Millisecond),
			segment.End.Truncate(time.Millisecond),
			segment.Text,
		)
	}
}

// cbProgress is a method of the Engine struct that returns a function.
// The function takes a progress int as input and returns nothing.
// It sets the progress field of the Engine struct to the given progress int.
// If the PrintProgress field in the configuration is true, it prints the progress.
func (e *Engine) cbProgress() func(progress int) {
	return func(progress int) {
		// If the progress is greater than 100, set it to 100.
		if progress > 100 {
			progress = 100
		}

		if e.progress == progress {
			return
		}
		e.progress = progress
		if e.cfg.PrintProgress {
			log.Info().Msgf("current progress: %d%%", progress)
		}

		// send webhook
		if e.webhook != nil {
			e.webhook.Send(context.Background(), &request{
				Progress: progress,
			})
		}
	}
}

// getOutputPath is a method of the Engine struct that takes a format string as input.
// It returns the output path for the converted audio file based on the given format.
func (e *Engine) getOutputPath(format string) string {
	// Get the file extension of the audio file from the configuration.
	ext := filepath.Ext(e.cfg.AudioPath)
	// Get the base name of the audio file from the configuration.
	filename := filepath.Base(e.cfg.AudioPath)
	// If the OutputFilename field in the configuration is not empty,
	if e.cfg.OutputFilename != "" {
		filename = e.cfg.OutputFilename
	}
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

// Save saves the text to a file.
// It takes a format string as input and returns an error.
// It gets the output path for the converted audio file based on the given format.
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
