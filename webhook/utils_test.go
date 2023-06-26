package webhook

import (
	"reflect"
	"testing"
)

func TestToHeaders(t *testing.T) {
	type args struct {
		headers []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "empty header",
			args: args{
				headers: []string{},
			},
			want: map[string]string{},
		},
		{
			name: "single header",
			args: args{
				headers: []string{"X-Drone-Token=1234"},
			},
			want: map[string]string{
				"X-Drone-Token": "1234",
			},
		},
		{
			name: "multiple headers",
			args: args{
				headers: []string{
					"X-Drone-Token=1234",
					"X-UUID=foobar",
				},
			},
			want: map[string]string{
				"X-Drone-Token": "1234",
				"X-UUID":        "foobar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToHeaders(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
