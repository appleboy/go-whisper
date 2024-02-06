package config

import "fmt"

// Whisper is the configuration for whisper.
type Whisper struct {
	Model        string
	AudioPath    string
	Threads      uint
	Language     string
	Debug        bool
	SpeedUp      bool
	Translate    bool
	Prompt       string
	MaxContext   uint
	BeamSize     uint
	EntropyThold float64

	PrintProgress bool
	PrintSegment  bool

	OutputFolder   string
	OutputFilename string
	OutputFormat   []string
}

// Validate checks if the Whisper configuration is valid.
// It returns an error if the audio path or model is missing.
func (c *Whisper) Validate() error {
	if c.AudioPath == "" {
		return fmt.Errorf("audio path is required")
	}

	if c.Model == "" {
		return fmt.Errorf("model is required")
	}

	return nil
}

// Webhook represents a webhook configuration with URL, Insecure and Headers.
type Webhook struct {
	URL      string
	Insecure bool
	Headers  []string
}

// Setting is the configuration for whisper.
type Setting struct {
	Whisper Whisper
	Webhook Webhook
	Youtube Youtube
}

// Youtube represents the configuration for a YouTube video.
type Youtube struct {
	URL      string // URL is the YouTube video URL.
	Insecure bool   // Insecure specifies whether to skip SSL verification.
	Debug    bool   // Debug specifies whether to enable debug mode.
	Retry    int    // Retry specifies the number of times to retry on failure.
}
