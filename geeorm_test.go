package geeorm

import "testing"

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	e, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestEngine(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()
}
