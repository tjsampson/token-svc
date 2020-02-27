package postgres

import (
	"bytes"
	"fmt"
)

// Clause is used to for SQL CLAUSE
type Clause struct {
	prefix string
	value  string
	params []interface{}
}

// WhereClause is a type of SQL Clause
type WhereClause struct {
	clauses []Clause
}

// Where returns a WhereClause
func Where() *WhereClause {
	return &WhereClause{}
}

// AddAndClause appends an AND CLAUSE
func (w *WhereClause) AddAndClause(value string, params ...interface{}) *WhereClause {
	w.clauses = append(w.clauses, Clause{prefix: "AND", value: value, params: params})
	return w
}

// AddOrClause appends an OR CLAUSE
func (w *WhereClause) AddOrClause(value string, params ...interface{}) *WhereClause {
	w.clauses = append(w.clauses, Clause{prefix: "OR", value: value, params: params})
	return w
}

// Value returns the string value of the CLAUSES
func (w *WhereClause) Value() string {
	if len(w.clauses) == 0 {
		return ""
	}

	buffer := bytes.Buffer{}
	positionCounter := 1

	getPositionalCounter := func(start int, count int) []interface{} {
		returnValues := []interface{}{}
		for index := start; index < count+start; index++ {
			returnValues = append(returnValues, fmt.Sprintf("$%d", index))
		}
		return returnValues
	}

	for index, element := range w.clauses {
		positionCounters := getPositionalCounter(positionCounter, len(element.params))
		tempValue := fmt.Sprintf(element.value, positionCounters...)
		positionCounter += len(element.params)
		if index == 0 {
			buffer.WriteString(fmt.Sprintf("WHERE %s", tempValue))
		} else {
			buffer.WriteString(fmt.Sprintf(" %s %s", element.prefix, tempValue))
		}
	}

	return buffer.String()
}

// Params contains the WhereClause Params
func (w WhereClause) Params() []interface{} {
	returnParams := []interface{}{}
	for _, element := range w.clauses {
		returnParams = append(returnParams, element.params...)
	}
	return returnParams
}
