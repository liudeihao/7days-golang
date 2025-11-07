package main

import (
    "bytes"
    "testing"
)

func TestGetterFunc_Get(t *testing.T) {
    var f GetterFunc = func(key string) ([]byte, error) {
        return []byte(key), nil
    }
    expected := []byte("key")
    if v, _ := f.Get("key"); !bytes.Equal(v, expected) {
        t.Errorf("expected: %v, got: %v", expected, v)
    }
}
