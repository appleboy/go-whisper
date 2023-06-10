package config

type Whisper struct {
	Model      string
	AudioPath  string
	OutputPath string
	Threads    uint
	Language   string
	Debug      bool
}

type Config struct{}
