package parser_sqlite

import (
	"fmt"

	"github.com/CovenantSQL/sqlparser"
	"github.com/gokch/ornn/config"
	"github.com/gokch/ornn/parser"
)

func New(sch *config.Schema) parser.Parser {
	return &Parser{
		sch: sch,
	}
}

type Parser struct {
	sch *config.Schema
}

func (t *Parser) Parse(sql string) (*parser.ParsedQuery, error) {
	stmtNodes, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}

	switch stmt := stmtNodes.(type) {
	case *sqlparser.Select:
		_ = stmt
	case *sqlparser.Insert:
		_ = stmt
	default:
		return nil, fmt.Errorf("parser error | not support query statement %T", stmt)
	}

	parsedQuery := &parser.ParsedQuery{}
	return parsedQuery, nil
}