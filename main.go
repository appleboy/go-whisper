package main

import (
	"os"
	"strconv"
	"time"

	"github.com/appleboy/go-whisper/config"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// Version set at compile-time
var (
	Version string
)

func main() {
	isTerm := isatty.IsTerminal(os.Stdout.Fd())
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !isTerm,
		},
	)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	app := cli.NewApp()
	app.Name = "Transcript Using Whisper API"
	app.Usage = "Speech-to-Text."
	app.Copyright = "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Bo-Yi Wu"
	app.Authors = []*cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "model",
			Usage:   "model name",
			EnvVars: []string{"PLUGIN_MODEL", "INPUT_MODEL"},
		},
		&cli.StringFlag{
			Name:    "audio-path",
			Usage:   "audio path",
			EnvVars: []string{"PLUGIN_AUDIO_PATH", "INPUT_AUDIO_PATH"},
		},
		&cli.StringFlag{
			Name:    "output-path",
			Usage:   "output path",
			EnvVars: []string{"PLUGIN_OUTPUT_PATH", "INPUT_OUTPUT_PATH"},
		},
		&cli.StringFlag{
			Name:    "language",
			Usage:   "language",
			EnvVars: []string{"PLUGIN_LANGUAGE", "INPUT_LANGUAGE"},
			Value:   "auto",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "enable debug mode",
			EnvVars: []string{"PLUGIN_DEBUG", "INPUT_DEBUG"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("can't run app")
	}
}

func run(c *cli.Context) error {
	if c.Bool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
	}

	cfg := &Config{
		Whisper: &config.Whisper{
			Model:      c.String("model"),
			AudioPath:  c.String("audio-path"),
			OutputPath: c.String("output-path"),
			Debug:      c.Bool("debug"),
			Language:   c.String("language"),
		},
	}

	if cfg.Whisper.Debug {
		spew.Dump(cfg)
	}

	_, err := Transcript(cfg)

	return err
}
