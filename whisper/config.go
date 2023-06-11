package whisper

import "fmt"

// Config is the whisper config.
type Config struct {
	Model      string
	AudioPath  string
	OutputPath string
	Threads    uint
	Language   string
	Debug      bool
}

// Validate validates the config.
func (c *Config) Validate() error {
	if c.AudioPath == "" {
		return fmt.Errorf("audio path is required")
	}

	if c.Model == "" {
		return fmt.Errorf("model is required")
	}

	return nil
}
