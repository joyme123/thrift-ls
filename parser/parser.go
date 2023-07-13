package parser

type IncludeCall func(include string) (filename string, content []byte, err error)

type Parser interface {
	Parse(filename string, content []byte) *ParseResult
	ParseRecursively(filename string, content []byte, maxDepth int, call IncludeCall) []*ParseResult
}

// PEGParser use PEG as a parser implementation
type PEGParser struct {
	parsed map[string]struct{}
}

func (p *PEGParser) Parse(filename string, content []byte) (*Document, []error) {
	if p.parsed == nil {
		p.parsed = make(map[string]struct{})
	}
	p.parsed[filename] = struct{}{}

	doc, err := Parse(filename, content)
	if err != nil {
		var errors []error
		errList, ok := err.(ErrorLister)
		if ok {
			errors = errList.Errors()
		} else {
			errors = append(errors, err)
		}

		var res *Document
		if doc != nil {
			res = doc.(*Document)
		}
		return res, errors
	}

	return doc.(*Document), nil
}

func (p *PEGParser) ParseRecursively(filename string, content []byte, maxDepth int, call IncludeCall) []*ParseResult {
	return p.parseRecursively(filename, content, 0, maxDepth, call)
}

func (p *PEGParser) parseRecursively(filename string, content []byte, curDepth int, maxDepth int, call IncludeCall) []*ParseResult {
	if curDepth > maxDepth && maxDepth > 0 {
		return nil
	}

	results := make([]*ParseResult, 0)

	doc, errs := p.Parse(filename, content)
	results = append(results, &ParseResult{
		Doc:    doc,
		Errors: errs,
	})

	if doc != nil {
		for _, include := range doc.Includes {
			f, c, err := call(include.Path)
			if err == nil {
				errs = append(errs, err)
			}
			if _, ok := p.parsed[f]; ok {
				continue
			}
			subRes := p.parseRecursively(f, c, curDepth+1, maxDepth, call)
			results = append(results, subRes...)
		}
	}

	return results
}

type ParseResult struct {
	Doc    *Document
	Errors []error
}
