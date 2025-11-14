package persistence

import (
	"fmt"
	"strings"
)

type Columns struct {
	readable []string
	writable []string
	alias    string
	idField  string
}

func NewColumns(readable []string, writable []string, alias, idField string) *Columns {
	return &Columns{
		readable: readable,
		writable: writable,
		alias:    alias,
		idField:  idField,
	}
}

func (c *Columns) ForSelect() []string {
	return c.readable
}

func (c *Columns) ForInsert() []string {
	return c.writable
}

func (c *Columns) GetIDField() string {
	return c.idField
}

func (c *Columns) OnConflict() string {
	if len(c.writable) == 0 || c.idField == "" {
		return ""
	}

	var statements []string
	for _, col := range c.writable {
		if col == c.idField {
			continue
		}
		statements = append(statements, fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}

	if len(statements) == 0 {
		return fmt.Sprintf("ON CONFLICT (%s) DO NOTHING", c.idField)
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", c.idField, strings.Join(statements, ","))
}
