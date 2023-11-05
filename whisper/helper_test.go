package whisper

import (
	"testing"
	"time"
)

func TestSrtTimestamp(t *testing.T) {
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				t: time.Duration(1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond),
			},
			want: "01:02:03,004",
		},
		{
			name: "test 2",
			args: args{
				t: time.Duration(10*time.Hour + 20*time.Minute + 30*time.Second + 40*time.Millisecond),
			},
			want: "10:20:30,040",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := srtTimestamp(tt.args.t); got != tt.want {
				t.Errorf("srtTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}
