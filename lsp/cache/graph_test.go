package cache

import (
	"testing"

	"github.com/joyme123/thrift-ls/parser"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/uri"
)

func Test_Graph(t *testing.T) {
	graph := NewIncludeGraph()

	// node1:
	//   file:///tmp/model/user.thrift
	//   include "../base.thrift"
	//   include "../addr.thrift"

	// node2:
	//   file:///tmp/base.thrift

	// node3:
	//   file:///tmp/addr.thrift
	//   include "./base.thrift"
	file1 := uri.New("file:///tmp/model/user.thrift")
	file2 := uri.New("file:///tmp/base.thrift")
	file3 := uri.New("file:///tmp/addr.thrift")
	graph.Set(file1, []*parser.Include{
		{Path: &parser.Literal{Value: "../base.thrift"}},
		{Path: &parser.Literal{Value: "../addr.thrift"}},
	})
	graph.Set(file2, nil)
	graph.Set(file3, []*parser.Include{
		{Path: &parser.Literal{Value: "./base.thrift"}},
	})

	expectNode1 := &IncludeNode{
		outdegree: []uri.URI{file3, file2},
	}
	expectNode2 := &IncludeNode{
		indegree: []uri.URI{file1, file3},
	}
	expectNode3 := &IncludeNode{
		indegree:  []uri.URI{file1},
		outdegree: []uri.URI{file2},
	}

	assert.Equal(t, expectNode1, graph.Get("file:///tmp/model/user.thrift"), "user.thrift")
	assert.Equal(t, expectNode2, graph.Get("file:///tmp/base.thrift"), "base.thrift")
	assert.Equal(t, expectNode3, graph.Get("file:///tmp/addr.thrift"), "addr.thrift")

	assert.Equal(t, expectNode1, expectNode1.Clone())
	assert.Equal(t, expectNode2, expectNode2.Clone())
	assert.Equal(t, expectNode3, expectNode3.Clone())

	graph.Remove(file2)
	assert.Equal(t, expectNode1, graph.Get("file:///tmp/model/user.thrift"), "user.thrift")
	assert.Equal(t, expectNode2, graph.Get("file:///tmp/base.thrift"), "base.thrift")
	assert.Equal(t, expectNode3, graph.Get("file:///tmp/addr.thrift"), "addr.thrift")

	graph.Remove(file1)
	expectNode1 = nil
	expectNode2 = &IncludeNode{
		indegree: []uri.URI{file3},
	}
	expectNode3 = &IncludeNode{
		outdegree: []uri.URI{file2},
	}
	assert.Equal(t, expectNode1, graph.Get("file:///tmp/model/user.thrift"), "user.thrift")
	assert.Equal(t, expectNode2, graph.Get("file:///tmp/base.thrift"), "base.thrift")
	assert.Equal(t, expectNode3, graph.Get("file:///tmp/addr.thrift"), "addr.thrift")

	graph.Remove(file3)
	assert.Nil(t, graph.Get("file:///tmp/model/user.thrift"), "user.thrift")
	assert.Nil(t, graph.Get("file:///tmp/base.thrift"), "base.thrift")
	assert.Nil(t, graph.Get("file:///tmp/addr.thrift"), "addr.thrift")
}
