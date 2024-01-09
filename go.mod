module github.com/appleboy/go-whisper

go 1.21

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/ggerganov/whisper.cpp/bindings/go v0.0.0-20230606002726-57543c169e27
	github.com/go-audio/wav v1.1.0
	github.com/joho/godotenv v1.5.1
	github.com/kkdai/youtube/v2 v2.9.0
	github.com/mattn/go-isatty v0.0.20
	github.com/rs/zerolog v1.31.0
	github.com/urfave/cli/v2 v2.26.0
	golang.org/x/net v0.19.0
)

require (
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/bitly/go-simplejson v0.5.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/dop251/goja v0.0.0-20231027120936-b396bb4c349d // indirect
	github.com/go-audio/audio v1.0.0 // indirect
	github.com/go-audio/riff v1.0.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/google/pprof v0.0.0-20231205033806-a5a03c77bf08 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/vbauerster/mpb/v5 v5.4.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/ggerganov/whisper.cpp/bindings/go => github.com/appleboy/whisper.cpp/bindings/go v0.0.0-20240109015431-0bbb5371d174
