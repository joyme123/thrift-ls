package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type args struct {
		fh FileHandle
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				fh: &Overlay{
					uri: "file:///tmp/types.thrift",
					content: []byte(`
#include "base.thrift"
struct Xtruct3
{
  1:  string string_thing,
  4:  i32    changed,
  9:  i32    i32_thing,
  11: i64    i64_thing
}
					`),
					version: 0,
					saved:   false,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "invalid ast",
			args: args{
				fh: &Overlay{
					uri: "file:///tmp/types.thrift",
					content: []byte(`
#include "base.thrift"
struct Xtruct3
{
  1:  string string_thing,
  4:  i32    changed,
  9:  i32    i32_thing,
  11: i64    i64_thing,
  12: 
}
					`),
					version: 0,
					saved:   false,
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.fh)
			tt.assertion(t, err)
			t.Logf("got: %v\n", got)
		})
	}
}
