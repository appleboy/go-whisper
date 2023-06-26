package config

import "fmt"

// Whisper is the configuration for whisper.
type Whisper struct {
	Model     string
	AudioPath string
	Threads   uint
	Language  string
	Debug     bool
	SpeedUp   bool
	Translate bool

	PrintProgress bool
	PrintSegment  bool

	OutputFolder string
	OutputFormat []string
}

// Validate validates the config.
func (c *Whisper) Validate() error {
	if c.AudioPath == "" {
		return fmt.Errorf("audio path is required")
	}

	if c.Model == "" {
		return fmt.Errorf("model is required")
	}

	return nil
}

// Webhook is the configuration for webhook.
type Webhook struct {
	URL      string
	Insecure bool
	Headers  []string
}

// Setting is the configuration for whisper.
type Setting struct {
	Whisper Whisper
	Webhook Webhook
}
