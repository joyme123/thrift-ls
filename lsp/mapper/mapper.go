package mapper

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"unicode/utf8"

	"github.com/joyme123/thrift-ls/lsp/types"
	"github.com/joyme123/thrift-ls/parser"
	"go.lsp.dev/uri"
)

type Mapper struct {
	fileURI uri.URI
	content []byte

	lineInit  sync.Once
	lineStart []int // line start 0-based byte offset. lsp: 0-based, parser: 1-based
	nonASCII  bool
}

// NewMapper ...
func NewMapper(fileURI uri.URI, content []byte) *Mapper {
	return &Mapper{
		fileURI: fileURI,
		content: content,
	}
}

func (m *Mapper) initLineStart() {
	m.lineInit.Do(func() {
		nlines := bytes.Count(m.content, []byte("\n"))
		m.lineStart = make([]int, 1, nlines+1) // initially []int{0}
		for offset, b := range m.content {
			if b == '\n' {
				m.lineStart = append(m.lineStart, offset+1)
			}
			if b >= utf8.RuneSelf {
				m.nonASCII = true
			}
		}
	})
}

// convert from utf16-based to rune-based position
func (m *Mapper) LSPPosToParserPosition(pos types.Position) (parser.Position, error) {
	m.initLineStart()
	line := int(pos.Line) + 1
	if line >= len(m.lineStart) {
		return parser.InvalidPosition, fmt.Errorf("invalid position line, completion line: %d, total line: %d", line, len(m.lineStart))
	}

	if !m.nonASCII {
		col := int(pos.Character) + 1
		offset := m.lineStart[pos.Line] + int(pos.Character)
		if offset >= len(m.content) {
			return parser.InvalidPosition, fmt.Errorf("invalid position offset: %d, total content: %d, %s", offset, len(m.content), string(m.content))
		}
		lineLength := m.lineStart[pos.Line+1] - m.lineStart[pos.Line]
		if col > lineLength {
			return parser.InvalidPosition, fmt.Errorf("invalid position column: %d, line length: %d, %s", col, lineLength, string(m.content))
		}

		return parser.Position{
			Line:   line,
			Col:    col,
			Offset: offset,
		}, nil
	}

	lineStart := m.lineStart[pos.Line]
	lineEnd := 0
	if int(pos.Line) == len(m.lineStart)-1 {
		lineEnd = len(m.content)
	} else {
		lineEnd = m.lineStart[pos.Line+1]
	}
	lineBytes := m.content[lineStart:lineEnd]

	utf16Col := -1
	bytesCol := -1
	for len(lineBytes) > 0 {
		if utf16Col >= int(pos.Character) {
			break
		}
		utf16Col++
		if lineBytes[0] < utf8.RuneSelf {
			lineBytes = lineBytes[1:]
			bytesCol++
			continue
		}

		r, size := utf8.DecodeRune(lineBytes)
		if r >= 0x10000 {
			utf16Col++
		}
		lineBytes = lineBytes[size:]
		bytesCol += size
	}

	runeLen := utf8.RuneCount(m.content[lineStart : lineStart+bytesCol+1])
	offset := lineStart + bytesCol
	if offset >= len(m.content) {
		return parser.InvalidPosition, errors.New("invalid position character")
	}
	if offset >= m.lineStart[pos.Line+1] {
		return parser.InvalidPosition, errors.New("invalid position character")
	}

	return parser.Position{
		Line:   line,
		Col:    runeLen,
		Offset: lineStart + bytesCol,
	}, nil
}
