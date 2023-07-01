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
	folder, err := os.MkdirTemp("", "youtube")
	if err != nil {
		panic(err)
	}

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

	downloader := &ytdl.Downloader{}
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

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

	return outputFile, nil
}

// New for creating a new youtube engine.
func New(cfg *config.Youtube) (*Engine, error) {
	return &Engine{
		cfg: cfg,
	}, nil
}
