package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocationFromCurrent(t *testing.T) {
	type args struct {
		c *current
	}
	tests := []struct {
		name string
		args args
		want Location
	}{
		{
			name: "oneline",
			args: args{
				c: &current{
					pos: position{
						line:   1,
						col:    1,
						offset: 0,
					},
					text: []byte("aaaaa"),
				},
			},
			want: Location{
				StartPos: Position{
					Line:   1,
					Col:    1,
					Offset: 0,
				},
				EndPos: Position{
					Line:   1,
					Col:    6,
					Offset: 5,
				},
			},
		},
		{
			name: "multiline with CRLF",
			args: args{
				c: &current{
					pos: position{
						line:   2,
						col:    0,
						offset: 1,
					},
					text: []byte("\r\naa\r\n\r\naaa"),
				},
			},
			want: Location{
				StartPos: Position{
					Line:   2,
					Col:    0,
					Offset: 1,
				},
				EndPos: Position{
					Line:   4,
					Col:    5,
					Offset: 12,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewLocationFromCurrent(tt.args.c))
		})
	}
}
