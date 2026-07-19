package pawnapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Literal holds a string or numeric value. Its zero value marshals as null.
type Literal struct {
	str    string
	num    string // decimal text, preserves formatting such as "0xFFFF" only if source used a number token; hex constants are stored as their decimal value.
	b      bool
	isStr  bool
	isNum  bool
	isBool bool
}

func StringLiteral(s string) Literal {
	return Literal{str: s, isStr: true}
}

func NumberLiteral(n int64) Literal {
	return Literal{num: strconv.FormatInt(n, 10), isNum: true}
}

func BoolLiteral(b bool) Literal {
	return Literal{b: b, isBool: true}
}

func (l Literal) IsZero() bool {
	return !l.isStr && !l.isNum && !l.isBool
}

func (l Literal) String() string {
	switch {
	case l.isStr:
		return l.str
	case l.isNum:
		return l.num
	case l.isBool:
		return strconv.FormatBool(l.b)
	default:
		return ""
	}
}

func (l Literal) MarshalJSON() ([]byte, error) {
	switch {
	case l.isStr:
		return json.Marshal(l.str)
	case l.isNum:
		return []byte(l.num), nil
	case l.isBool:
		return json.Marshal(l.b)
	default:
		return []byte("null"), nil
	}
}

func (l *Literal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*l = Literal{}
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("pawnapi: literal value: %w", err)
		}
		*l = Literal{str: s, isStr: true}
		return nil
	}
	if string(data) == "true" || string(data) == "false" {
		*l = Literal{b: string(data) == "true", isBool: true}
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("pawnapi: literal value: %w", err)
	}
	*l = Literal{num: n.String(), isNum: true}
	return nil
}
