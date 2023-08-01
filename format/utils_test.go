package format

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
)

func TestEqualsAfterFormat(t *testing.T) {

	ast, err := parser.Parse("test.thrift", []byte(ThriftTestContent))
	assert.NoError(t, err)
	assert.NotNil(t, ast)

	formated, err := FormatDocument(ast.(*parser.Document))
	assert.NoError(t, err)
	type args struct {
		doc1 string
		doc2 string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "equals",
			args: args{

				doc1: ThriftTestContent,
				doc2: formated,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EqualsAfterFormat(tt.args.doc1, tt.args.doc2)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
