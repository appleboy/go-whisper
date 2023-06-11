package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/appleboy/go-whisper/config"
	"github.com/go-audio/wav"
	"github.com/rs/zerolog/log"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

type Config struct {
	Whisper *config.Whisper
}

// Validate validates the config.
func (c *Config) Validate() error {
	if c.Whisper == nil {
		return fmt.Errorf("whisper config is required")
	}

	if c.Whisper.AudioPath == "" {
		return fmt.Errorf("audio path is required")
	}

	if c.Whisper.Model == "" {
		return fmt.Errorf("model is required")
	}

	return nil
}

func sh(c string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Env = os.Environ()
	o, err := cmd.CombinedOutput()
	return string(o), err
}

// AudioToWav converts audio to wav for transcribe.
func audioToWav(src, dst string) error {
	out, err := sh(fmt.Sprintf("ffmpeg -i %s -format s16le -ar 16000 -ac 1 -acodec pcm_s16le %s", src, dst))
	if err != nil {
		return fmt.Errorf("error: %w out: %s", err, out)
	}

	return nil
}

// Transcribe converts audio to text.
func Transcript(cfg *Config) (string, error) {
	var data []float32

	l := log.With().
		Str("module", "transcript").
		Logger()

	dir, err := os.MkdirTemp("", "whisper")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	convertedPath := filepath.Join(dir, "converted.wav")

	l.Debug().Msg("start convert audio to wav")
	if err := audioToWav(cfg.Whisper.AudioPath, convertedPath); err != nil {
		return "", err
	}

	// open converted file
	fh, err := os.Open(convertedPath)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	// Load the model
	model, err := whisper.New(cfg.Whisper.Model)
	if err != nil {
		return "", err
	}
	defer model.Close()

	// Decode the WAV file - load the full buffer
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return "", err
	} else if dec.SampleRate != whisper.SampleRate {
		return "", fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return "", fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}

	context, err := model.NewContext()
	if err != nil {
		return "", err
	}

	context.SetThreads(cfg.Whisper.Threads)

	l.Info().Msgf("%s", context.SystemInfo())

	if cfg.Whisper.Language != "" {
		_ = context.SetLanguage(cfg.Whisper.Language)
	}

	l.Debug().Msg("start transcribe process")
	context.ResetTimings()
	if err := context.Process(data, nil); err != nil {
		return "", err
	}

	text := ""
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		text += segment.Text
		l.Info().Msgf(
			"[%6s -> %6s] %s",
			segment.Start.Truncate(time.Millisecond),
			segment.End.Truncate(time.Millisecond),
			segment.Text,
		)
	}

	if cfg.Whisper.OutputPath != "" {
		if err := os.WriteFile(cfg.Whisper.OutputPath, []byte(text), 0o644); err != nil {
			return text, err
		}
	}

	return text, nil
}
