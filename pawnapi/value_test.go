package pawnapi

import (
	"encoding/json"
	"testing"
)

func TestLiteral_NumberRoundTrip(t *testing.T) {
	l := NumberLiteral(65535)
	buf, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "65535" {
		t.Fatalf("got %s, want 65535", buf)
	}
	var got Literal
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Fatal(err)
	}
	if got.String() != "65535" {
		t.Fatalf("got %s, want 65535", got.String())
	}
}

func TestLiteral_StringRoundTrip(t *testing.T) {
	l := StringLiteral("hello")
	buf, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	var got Literal
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Fatal(err)
	}
	if got.String() != "hello" {
		t.Fatalf("got %q, want %q", got.String(), "hello")
	}
}

func TestLiteral_BoolRoundTrip(t *testing.T) {
	l := BoolLiteral(true)
	buf, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "true" {
		t.Fatalf("got %s, want true", buf)
	}
	var got Literal
	if err := json.Unmarshal(buf, &got); err != nil {
		t.Fatal(err)
	}
	if got.String() != "true" {
		t.Fatalf("got %s, want true", got.String())
	}
}

func TestLiteral_Zero(t *testing.T) {
	var l Literal
	if !l.IsZero() {
		t.Fatal("zero value Literal should report IsZero")
	}
	buf, err := json.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf) != "null" {
		t.Fatalf("got %s, want null", buf)
	}
}
