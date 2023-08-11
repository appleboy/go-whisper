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
	retry int
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

	i := 0
	for i < e.retry {
		output, err := e.download(httpTransport)
		if err != nil {
			return "", err
		}
		if output != "" {
			return output, nil
		}
		time.Sleep(1 * time.Second)
		i++
	}

	return "", errors.New("youtube video can't download")
}

func (e *Engine) download(trans http.RoundTripper) (string, error) {
	folder, err := os.MkdirTemp("", "youtube")
	if err != nil {
		panic(err)
	}

	downloader := &ytdl.Downloader{}
	downloader.HTTPClient = &http.Client{Transport: trans}
	if e.cfg.Debug {
		downloader.Debug = true
	}

	e.video, err = downloader.GetVideo(e.cfg.URL)
	if err != nil {
		panic(err)
	}

	mimetype := "video/3gpp"
	outputQuality := ""

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
		format = e.video.Formats.FindByItag(itag)
		if format == nil {
			return "", fmt.Errorf("unable to find format with itag %d", itag)
		}

	case outputQuality != "":
		format = formats.FindByQuality(outputQuality)
		if format == nil {
			return "", fmt.Errorf("unable to find format with quality %s", outputQuality)
		}

	default:
		// select the first format
		formats.Sort()
		format = &formats[0]
	}

	outputFile := path.Join(folder, "video.3gp")

	if err := downloader.Download(context.Background(), e.video, format, outputFile); err != nil {
		return "", err
	}
	if isFileExistsAndNotEmpty(outputFile) {
		return outputFile, nil
	}

	return "", nil
}

// New for creating a new youtube engine.
func New(cfg *config.Youtube) (*Engine, error) {
	return &Engine{
		cfg:   cfg,
		retry: 100,
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
