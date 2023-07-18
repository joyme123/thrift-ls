package diagnostic

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/joyme123/thrift-ls/lsp/cache"
	"github.com/joyme123/thrift-ls/lsp/lsputils"
	"github.com/joyme123/thrift-ls/parser"
	log "github.com/sirupsen/logrus"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

type CycleCheck struct {
}

func (c *CycleCheck) Diagnostic(ctx context.Context, ss *cache.Snapshot, changeFiles []uri.URI) (DiagnosticResult, error) {
	includesMap := make(map[uri.URI][]Include)
	for _, file := range changeFiles {
		getIncludes(ctx, ss, file, &includesMap)
	}
	cyclePairs := cycleDetect(&includesMap)

	return cycleToDiagnosticItems(cyclePairs), nil
}

func cycleToDiagnosticItems(pairs []CyclePair) DiagnosticResult {
	diagnostics := make(DiagnosticResult)
	for i := range pairs {
		diagnostics[pairs[i].file] = append(diagnostics[pairs[i].file], cyclePairToDiagnostic(pairs[i]))
	}

	return diagnostics
}

func cyclePairToDiagnostic(pair CyclePair) protocol.Diagnostic {
	res := protocol.Diagnostic{
		Range:    lsputils.ASTNodeToRange(pair.include.include),
		Severity: protocol.DiagnosticSeverityWarning,
		Source:   "thrift-ls",
		Message:  fmt.Sprintf("cycle dependency in %s", pair.include.file),
	}
	return res
}

type Include struct {
	file    uri.URI
	include *parser.Include
}

type CyclePair struct {
	file    uri.URI
	include Include
}

func cycleDetect(includesMap *map[uri.URI][]Include) []CyclePair {
	cyclePairs := make([]CyclePair, 0)

	for uri, includes := range *includesMap {
		for _, incI := range includes {
			for _, incJ := range (*includesMap)[incI.file] {
				if uri == incJ.file {
					cyclePairs = append(cyclePairs, CyclePair{
						file:    uri,
						include: incI,
					})
				}
			}
		}
	}

	return cyclePairs
}

func getIncludes(ctx context.Context, ss *cache.Snapshot, file uri.URI, includesMap *map[uri.URI][]Include) error {
	pf, err := ss.Parse(ctx, file)
	if err != nil {
		log.Errorf("parse %s failed: %v", file, err)
		return err
	}
	if pf.AST() == nil {
		log.Errorf("parse ast failed: %v", pf.AggregatedError())
		return pf.AggregatedError()
	}

	includes := pf.AST().Includes
	for i := range includes {
		if includes[i].Path != nil && includes[i].Path.BadNode {
			continue
		}
		(*includesMap)[file] = append((*includesMap)[file], Include{
			file:    toURI(file, includes[i].Path.Value),
			include: includes[i],
		})

		includeURI := toURI(file, includes[i].Path.Value)
		if _, ok := (*includesMap)[includeURI]; ok {
			continue
		}
		getIncludes(ctx, ss, includeURI, includesMap)
	}

	return nil
}

func toURI(cur uri.URI, includePath string) uri.URI {
	filePath := cur.Filename()
	items := strings.Split(filePath, "/")
	basePath := strings.TrimSuffix(filePath, items[len(items)-1])

	path := filepath.Join(basePath, includePath)

	return uri.File(path)
}
