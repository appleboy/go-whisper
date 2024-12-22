package youtube

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/appleboy/go-whisper/config"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
	"golang.org/x/net/http/httpproxy"
)

// Engine is the youtube engine.
type Engine struct {
	cfg   *config.Youtube
	video *youtube.Video
}

// Filename returns a sanitized filename.
func (e *Engine) Filename() string {
	if e.video == nil {
		return ""
	}

	return ytdl.SanitizeFilename(e.video.Title)
}

// Download downloads youtube video.
func (e *Engine) Download(ctx context.Context) (string, error) {
	proxyFunc := httpproxy.FromEnvironment().ProxyFunc()
	httpTransport := &http.Transport{
		Proxy: func(r *http.Request) (uri *url.URL, err error) {
			return proxyFunc(r.URL)
		},
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	if e.cfg.Insecure {
		httpTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	for i := 0; i < e.cfg.Retry; i++ {
		output, err := e.download(ctx, httpTransport)
		if err != nil {
			return "", err
		}
		if output != "" {
			return output, nil
		}
		time.Sleep(1 * time.Second)
	}

	return "", errors.New("youtube video can't download")
}

func (e *Engine) download(ctx context.Context, trans http.RoundTripper) (string, error) {
	folder, err := os.MkdirTemp("", "youtube")
	if err != nil {
		panic(err)
	}

	downloader := &ytdl.Downloader{}
	downloader.HTTPClient = &http.Client{Transport: trans}

	e.video, err = downloader.GetVideo(e.cfg.URL)
	if err != nil {
		panic(err)
	}

	mimetype := "audio/mp4"
	outputQuality := "tiny"

	formats := e.video.Formats
	if mimetype != "" {
		formats = formats.Type(mimetype)
	}
	if len(formats) == 0 {
		return "", errors.New("no formats found")
	}

	var format *youtube.Format
	itag, _ := strconv.Atoi(outputQuality)
	switch {
	case itag > 0:
		// When an itag is specified, do not filter format with mime-type
		formats = e.video.Formats.Itag(itag)
		if len(formats) == 0 {
			return "", fmt.Errorf("unable to find format with itag %d", itag)
		}

	case outputQuality != "":
		formats = formats.Quality(outputQuality)
		if len(formats) == 0 {
			return "", fmt.Errorf("unable to find format with quality %s", outputQuality)
		}

	default:
		// select the first format
		formats.Sort()
		format = &formats[0]
	}

	outputFile := path.Join(folder, "video.mp4")

	if err := downloader.Download(ctx, e.video, format, outputFile); err != nil {
		return "", err
	}
	if isFileExistsAndNotEmpty(outputFile) {
		return outputFile, nil
	}

	return "", errors.New("download file is empty")
}

// New for creating a new youtube engine.
func New(cfg *config.Youtube) (*Engine, error) {
	return &Engine{
		cfg: cfg,
	}, nil
}

// isFileExistsAndNotEmpty check file not zero byte file
func isFileExistsAndNotEmpty(name string) bool {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fileInfo.Size() > 0
}
