package lru

import "testing"

type String string

func (s String) Len() int {
    return len(s)
}

func TestGet(t *testing.T) {
    lru := New(0, nil)
    lru.Add("key1", String("1234"))
    if v, ok := lru.Get("key1"); !ok || v != String("1234") {
        t.Fatalf("expected %s, got %s", "1234", v)
    }
    lru = New(100, nil)
    if _, ok := lru.Get("key2"); ok {
        t.Fatal("expected false, got", ok)
    }
}

func TestAdd(t *testing.T) {
    lru := New(0, nil)
    lru.Add("key1", String("1234"))
    lru.Add("key1", String("4567"))
    if v, ok := lru.Get("key1"); !ok || v != String("4567") {
        t.Fatalf("expected %s, got %s", "4567", v)
    }
}

func TestRemoveOldest(t *testing.T) {
    lru := New(0, nil)
    lru.Add("key1", String("1234"))
    lru.Add("key2", String("4567"))
    lru.RemoveOldest()
    if _, ok := lru.Get("key1"); ok {
        t.Fatal("expected false, got", ok)
    }
    if _, ok := lru.Get("key2"); !ok {
        t.Fatal("expected true, got", ok)
    }
}

func TestOnEvicted(t *testing.T) {
    var out string
    lru := New(0, func(key string, value Value) {
        out += key
    })
    lru.Add("key1", String("1234"))
    lru.RemoveOldest()
    if out != "key1" {
        t.Fatalf("expected %s, got %s", "key1", out)
    }
}
