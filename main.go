package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/appleboy/go-whisper/config"
	"github.com/appleboy/go-whisper/webhook"
	"github.com/appleboy/go-whisper/whisper"
	"github.com/appleboy/go-whisper/youtube"

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
	app.Name = "Speech-to-Text Using Whisper API"
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
			Usage:   "Model is the interface to a whisper model",
			EnvVars: []string{"PLUGIN_MODEL", "INPUT_MODEL"},
		},
		&cli.StringFlag{
			Name:    "audio-path",
			Usage:   "audio path",
			EnvVars: []string{"PLUGIN_AUDIO_PATH", "INPUT_AUDIO_PATH"},
		},
		&cli.StringFlag{
			Name:    "output-folder",
			Usage:   "output folder",
			EnvVars: []string{"PLUGIN_OUTPUT_FOLDER", "INPUT_OUTPUT_FOLDER"},
		},
		&cli.StringSliceFlag{
			Name:    "output-format",
			Usage:   "output format, support txt, srt, csv",
			EnvVars: []string{"PLUGIN_OUTPUT_FORMAT", "INPUT_OUTPUT_FORMAT"},
			Value:   cli.NewStringSlice("txt"),
		},
		&cli.StringFlag{
			Name:    "output-filename",
			Usage:   "output filename",
			EnvVars: []string{"PLUGIN_OUTPUT_FILENAME", "INPUT_OUTPUT_FILENAME"},
		},
		&cli.StringFlag{
			Name:    "language",
			Usage:   "Set the language to use for speech recognition",
			EnvVars: []string{"PLUGIN_LANGUAGE", "INPUT_LANGUAGE"},
			Value:   "auto",
		},
		&cli.UintFlag{
			Name:    "threads",
			Usage:   "Set number of threads to use",
			EnvVars: []string{"PLUGIN_THREADS", "INPUT_THREADS"},
			Value:   uint(runtime.NumCPU()),
		},
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "enable debug mode",
			EnvVars: []string{"PLUGIN_DEBUG", "INPUT_DEBUG"},
		},
		&cli.BoolFlag{
			Name:    "speedup",
			Usage:   "speed up audio by x2 (reduced accuracy)",
			EnvVars: []string{"PLUGIN_SPEEDUP", "INPUT_SPEEDUP"},
		},
		&cli.BoolFlag{
			Name:    "translate",
			Usage:   "translate from source language to english",
			EnvVars: []string{"PLUGIN_TRANSLATE", "INPUT_TRANSLATE"},
		},
		&cli.BoolFlag{
			Name:    "print-progress",
			Usage:   "print progress",
			EnvVars: []string{"PLUGIN_PRINT_PROGRESS", "INPUT_PRINT_PROGRESS"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "print-segment",
			Usage:   "print segment",
			EnvVars: []string{"PLUGIN_PRINT_SEGMENT", "INPUT_PRINT_SEGMENT"},
		},
		&cli.StringFlag{
			Name:    "webhook-url",
			Usage:   "webhook url",
			EnvVars: []string{"PLUGIN_WEBHOOK_URL", "INPUT_WEBHOOK_URL"},
		},
		&cli.BoolFlag{
			Name:    "webhook-insecure",
			Usage:   "webhook insecure",
			EnvVars: []string{"PLUGIN_WEBHOOK_INSECURE", "INPUT_WEBHOOK_INSECURE"},
		},
		&cli.StringSliceFlag{
			Name:    "webhook-headers",
			Usage:   "webhook headers",
			EnvVars: []string{"PLUGIN_WEBHOOK_HEADERS", "INPUT_WEBHOOK_HEADERS"},
		},
		&cli.StringFlag{
			Name:    "youtube",
			Usage:   "youtube url",
			EnvVars: []string{"PLUGIN_YOUTUBE", "INPUT_YOUTUBE"},
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

	cfg := config.Setting{
		Whisper: config.Whisper{
			Model:         c.String("model"),
			AudioPath:     c.String("audio-path"),
			Language:      c.String("language"),
			Threads:       c.Uint("threads"),
			SpeedUp:       c.Bool("speedup"),
			Translate:     c.Bool("translate"),
			PrintProgress: c.Bool("print-progress"),
			PrintSegment:  c.Bool("print-segment"),

			OutputFolder:   c.String("output-folder"),
			OutputFormat:   c.StringSlice("output-format"),
			OutputFilename: c.String("output-filename"),

			Debug: c.Bool("debug"),
		},

		Webhook: config.Webhook{
			URL:      c.String("webhook-url"),
			Insecure: c.Bool("webhook-insecure"),
			Headers:  c.StringSlice("webhook-headers"),
		},
	}

	if cfg.Whisper.Debug {
		spew.Dump(cfg)
	}

	// check youtube url
	if c.String("youtube") != "" {
		videoPath, err := youtube.DownloadVideo(c.String("youtube"))
		if err != nil {
			return err
		}
		cfg.Whisper.AudioPath = videoPath
	}

	e, err := whisper.New(
		&cfg.Whisper,
		webhook.NewClient(
			cfg.Webhook.URL,
			cfg.Webhook.Insecure,
			webhook.ToHeaders(cfg.Webhook.Headers),
		),
	)
	if err != nil {
		return err
	}

	if err := e.Transcript(); err != nil {
		return err
	}
	defer e.Close()

	for _, ext := range cfg.Whisper.OutputFormat {
		if err := e.Save(ext); err != nil {
			return err
		}
	}

	return nil
}
