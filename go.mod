module github.com/appleboy/go-whisper

go 1.20

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/ggerganov/whisper.cpp/bindings/go v0.0.0-20230606002726-57543c169e27
	github.com/go-audio/wav v1.1.0
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-isatty v0.0.19
	github.com/rs/zerolog v1.29.1
	github.com/urfave/cli/v2 v2.25.6
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/go-audio/audio v1.0.0 // indirect
	github.com/go-audio/riff v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.8.0 // indirect
)

replace github.com/ggerganov/whisper.cpp/bindings/go => github.com/appleboy/whisper.cpp/bindings/go v0.0.0-20230617020330-4d2f9dd8c28e
