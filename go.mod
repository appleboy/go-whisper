module github.com/appleboy/go-whisper

go 1.20

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/ggerganov/whisper.cpp/bindings/go v0.0.0-20230606002726-57543c169e27
	github.com/go-audio/wav v1.1.0
	github.com/joho/godotenv v1.5.1
	github.com/kkdai/youtube/v2 v2.8.3
	github.com/mattn/go-isatty v0.0.19
	github.com/rs/zerolog v1.30.0
	github.com/urfave/cli/v2 v2.25.7
	golang.org/x/net v0.14.0
)

require (
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/bitly/go-simplejson v0.5.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/dop251/goja v0.0.0-20230812105242-81d76064690d // indirect
	github.com/go-audio/audio v1.0.0 // indirect
	github.com/go-audio/riff v1.0.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/google/pprof v0.0.0-20230821062121-407c9e7a662f // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/vbauerster/mpb/v5 v5.4.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
)

replace github.com/ggerganov/whisper.cpp/bindings/go => github.com/appleboy/whisper.cpp/bindings/go v0.0.0-20230808024901-03650882d33b
