package whisper

import (
	"testing"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

func TestEngine_getOutputPath(t *testing.T) {
	type fields struct {
		cfg      *Config
		ctx      whisper.Context
		model    whisper.Model
		segments []whisper.Segment
	}
	type args struct {
		format string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "change wav to txt",
			fields: fields{
				cfg: &Config{
					AudioPath: "/test/1234/foo.wav",
				},
			},
			args: args{
				format: "txt",
			},
			want: "/test/1234/foo.txt",
		},
		{
			name: "change output folder",
			fields: fields{
				cfg: &Config{
					AudioPath:    "/test/1234/foo.wav",
					OutputFolder: "/foo/bar",
				},
			},
			args: args{
				format: "txt",
			},
			want: "/foo/bar/foo.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				cfg:      tt.fields.cfg,
				ctx:      tt.fields.ctx,
				model:    tt.fields.model,
				segments: tt.fields.segments,
			}
			if got := e.getOutputPath(tt.args.format); got != tt.want {
				t.Errorf("Engine.getOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
