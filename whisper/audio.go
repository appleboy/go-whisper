package whisper

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// sh executes shell command.
func sh(c string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", c)
	cmd.Env = os.Environ()
	o, err := cmd.CombinedOutput()
	return string(o), err
}

// AudioToWav converts audio to wav for transcribe.
func audioToWav(src, dst string) error {
	out, err := sh(fmt.Sprintf("ffmpeg -i %s -format s16le -ar 16000 -ac 1 -acodec pcm_s16le %s", src, dst))
	if err != nil {
		return fmt.Errorf("error: %w out: %s", err, out)
	}

	return nil
}

// SrtTimestamp converts time.Duration to srt timestamp.
func srtTimestamp(t time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d",
		t/time.Hour,
		(t%time.Hour)/time.Minute,
		(t%time.Minute)/time.Second,
		(t%time.Second)/time.Millisecond,
	)
}
