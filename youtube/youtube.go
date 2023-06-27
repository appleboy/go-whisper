package youtube

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/kkdai/youtube/v2"
	ytdl "github.com/kkdai/youtube/v2/downloader"
	"golang.org/x/net/http/httpproxy"
)

func DownloadVideo(u string) (string, error) {
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

	downloader := &ytdl.Downloader{
		// OutputDir: folder,
	}
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

	video, err := downloader.GetVideo(u)
	if err != nil {
		panic(err)
	}

	mimetype := "mp4"
	outputQuality := "medium"

	formats := video.Formats
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
		format = video.Formats.FindByItag(itag)
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

	outputFile := path.Join(folder, "video.mp4")
	if err := downloader.Download(context.Background(), video, format, outputFile); err != nil {
		return "", err
	}

	return outputFile, nil
}
