package sitter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type (
	Query struct {
		t *TreeSitter
		q uint64
	}

	QueryCursor struct {
		t  *TreeSitter
		qc uint64
	}

	QueryCapture struct {
		ID   uint32
		Node *Node
	}

	QueryMatch struct {
		ID           uint32
		PatternIndex uint16
		Captures     []QueryCapture
	}
)

const (
	QueryErrorNone uint32 = iota
	QueryErrorSyntax
	QueryErrorNodeType
	QueryErrorField
	QueryErrorCapture
	QueryErrorStructure
	QueryErrorLanguage
)

func (ts *TreeSitter) NewQuery(pattern string, l *Language) (*Query, error) {
	errOffPtr, err := ts.call(_malloc, 4)
	if err != nil {
		return nil, fmt.Errorf("allocating query error offset: %w", err)
	}
	errTypePtr, err := ts.call(_malloc, 4)
	if err != nil {
		return nil, fmt.Errorf("allocating query error type: %w", err)
	}
	patternPtr, patternSize, freePattern, err := ts.allocateString(pattern)
	if err != nil {
		return nil, fmt.Errorf("allocating pattern string: %w", err)
	}
	defer freePattern()
	queryPtr, err := ts.call(_queryNew, l.l, patternPtr, patternSize, errOffPtr[0], errTypePtr[0])
	if err != nil {
		return nil, fmt.Errorf("creating query: %w", err)
	}
	errorOffset, ok := ts.mod.Memory().ReadUint32Le(uint32(errOffPtr[0]))
	if !ok {
		return nil, errors.New("invalid query error offset")
	}
	errorType, ok := ts.mod.Memory().ReadUint32Le(uint32(errTypePtr[0]))
	if !ok {
		return nil, errors.New("invalid query error type")
	}

	if errorType != QueryErrorNone {
		// search for the line containing the offset
		line := 1
		line_start := 0
		for i, c := range pattern {
			line_start = i
			if uint32(i) >= errorOffset {
				break
			}
			if c == '\n' {
				line++
			}
		}
		column := int(errorOffset) - line_start
		errorTypeToString := QueryErrorTypeToString(errorType)

		var message string
		switch errorType {
		// errors that apply to a single identifier
		case QueryErrorNodeType:
			fallthrough
		case QueryErrorField:
			fallthrough
		case QueryErrorCapture:
			// find identifier at input[errorOffset]
			// and report it in the error message
			s := string(pattern[errorOffset:])
			identifierRegexp := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*`)
			m := identifierRegexp.FindStringSubmatch(s)
			if len(m) > 0 {
				message = fmt.Sprintf("invalid %s '%s' at line %d column %d",
					errorTypeToString, m[0], line, column)
			} else {
				message = fmt.Sprintf("invalid %s at line %d column %d",
					errorTypeToString, line, column)
			}

		// errors the report position
		case QueryErrorSyntax:
			fallthrough
		case QueryErrorStructure:
			fallthrough
		case QueryErrorLanguage:
			fallthrough
		default:
			s := string(pattern[errorOffset:])
			lines := strings.Split(s, "\n")
			whitespace := strings.Repeat(" ", column)
			message = fmt.Sprintf("invalid %s at line %d column %d\n%s\n%s^",
				errorTypeToString, line, column,
				lines[0], whitespace)
		}

		return nil, errors.New(message)
	}

	return &Query{ts, queryPtr[0]}, nil
}

func (q *Query) CaptureNameForID(id uint32) (string, error) {
	strlenPtr, err := q.t.call(_malloc, 4)
	if err != nil {
		return "", fmt.Errorf("allocating string length: %w", err)
	}
	namePtr, err := q.t.call(_queryCaptureNameForID, q.q, uint64(id), strlenPtr[0])
	if err != nil {
		return "", fmt.Errorf("getting capture name for id: %w", err)
	}
	strlen, ok := q.t.mod.Memory().ReadUint32Le(uint32(strlenPtr[0]))
	if !ok {
		return "", errors.New("invalid str len")
	}
	captureName, ok := q.t.mod.Memory().Read(uint32(namePtr[0]), strlen)
	if !ok {
		return "", errors.New("invalid capture name")
	}
	return string(captureName), nil
}

func (ts *TreeSitter) NewQueryCursor() (*QueryCursor, error) {
	qc, err := ts.call(_queryCursorNew)
	if err != nil {
		return nil, fmt.Errorf("creating query cursor: %w", err)
	}
	return &QueryCursor{ts, qc[0]}, nil
}

func (qc *QueryCursor) Exec(q *Query, n *Node) error {
	_, err := qc.t.call(_queryCursorExec, qc.qc, q.q, n.ptr)
	return err
}

func (ts *TreeSitter) allocateQueryMatch() (uint64, error) {
	// allocate tsquerymatch 12 bytes
	nodePtr, err := ts.call(_malloc, uint64(12))
	if err != nil {
		return 0, fmt.Errorf("allocating query match: %w", err)
	}
	return nodePtr[0], nil
}

func (qc QueryCursor) NextMatch() (*QueryMatch, bool, error) {
	queryMatchPtr, err := qc.t.allocateQueryMatch()
	if err != nil {
		return nil, false, err
	}
	hasNextMatch, err := qc.t.call(_queryCursorNextMatch, qc.qc, queryMatchPtr)
	if err != nil {
		return nil, false, fmt.Errorf("getting query cursor next match: %w", err)
	}
	if hasNextMatch[0] == 0 {
		return nil, false, nil
	}

	queryMatchID, ok := qc.t.mod.Memory().ReadUint32Le(uint32(queryMatchPtr))
	if !ok {
		return nil, false, errors.New("invalid query match id")
	}
	queryMatchPatternIndex, ok := qc.t.mod.Memory().ReadUint16Le(uint32(queryMatchPtr) + 4)
	if !ok {
		return nil, false, errors.New("invalid query match pattern index")
	}
	queryMatchCaptureCount, ok := qc.t.mod.Memory().ReadUint16Le(uint32(queryMatchPtr) + 6)
	if !ok {
		return nil, false, errors.New("invalid query match pattern index")
	}
	queryMatchCapturesPtr, ok := qc.t.mod.Memory().ReadUint32Le(uint32(queryMatchPtr) + 8)
	if !ok {
		return nil, false, errors.New("invalid query match captures pointer")
	}
	qcs := make([]QueryCapture, queryMatchCaptureCount)
	addr := queryMatchCapturesPtr
	for i := range queryMatchCaptureCount {
		captureIndex, ok := qc.t.mod.Memory().ReadUint32Le(addr + 24)
		if !ok {
			return nil, false, errors.New("invalid capture index")
		}
		qcs[i] = QueryCapture{
			ID:   captureIndex,
			Node: newNode(qc.t, uint64(addr)),
		}
		addr += 28
	}
	return &QueryMatch{
		ID:           queryMatchID,
		PatternIndex: queryMatchPatternIndex,
		Captures:     qcs,
	}, true, nil
}

func QueryErrorTypeToString(errorType uint32) string {
	switch errorType {
	case QueryErrorNone:
		return "none"
	case QueryErrorNodeType:
		return "node type"
	case QueryErrorField:
		return "field"
	case QueryErrorCapture:
		return "capture"
	case QueryErrorSyntax:
		return "syntax"
	default:
		return "unknown"
	}
}
