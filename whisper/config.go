package whisper

import "fmt"

type OutputFormat string

func (f OutputFormat) String() string {
	return string(f)
}

var (
	FormatTxt OutputFormat = "txt"
	FormatSrt OutputFormat = "srt"
	FormatCSV OutputFormat = "csv"
)

// Config is the whisper config.
type Config struct {
	Model     string
	AudioPath string
	Threads   uint
	Language  string
	Debug     bool
	SpeedUp   bool
	Translate bool

	PrintProgress bool

	OutputFolder string
	OutputFormat []string
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
