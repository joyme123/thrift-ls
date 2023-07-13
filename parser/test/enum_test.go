package test

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ParseEnum(t *testing.T) {
	demoContent := `enum TweetType {
    TWEET,         // 0
    RETWEET = 2,   // 2
    DM = 0xa,      // 10
    REPLY          // 11
    POST = 0o21    // 17
}`
	ast, err := parser.Parse("test.thrift", []byte(demoContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	enums := ast.(*parser.Document).Enums
	assert.Len(t, enums, 1)
	assert.Equal(t, "TweetType", enums[0].Name.Name)
	assert.Len(t, enums[0].Values, 5)
	assert.Equal(t, "TWEET", enums[0].Values[0].Name.Name)
	assert.Equal(t, int64(0), enums[0].Values[0].Value)
	assert.Equal(t, "RETWEET", enums[0].Values[1].Name.Name)
	assert.Equal(t, int64(2), enums[0].Values[1].Value)
	assert.Equal(t, "DM", enums[0].Values[2].Name.Name)
	assert.Equal(t, int64(10), enums[0].Values[2].Value)
	assert.Equal(t, "REPLY", enums[0].Values[3].Name.Name)
	assert.Equal(t, int64(11), enums[0].Values[3].Value)
	assert.Equal(t, "POST", enums[0].Values[4].Name.Name)
	assert.Equal(t, int64(17), enums[0].Values[4].Value)
}
