package main

import (
    "testing"
)

func TestStartOfLine(t *testing.T) {
    pos := startOfLine("hej\nfoobarbaz\nturtles", 7)
    if pos != 4 {
        t.Errorf("wrong: %d"  , pos)
    }
}


func TestSplitPos(t *testing.T) {
    split := splitPos("hej\nfoobarbaz\nturtles", 7)
    if split[0] != "hej\nfoo" {
        t.Errorf("wrong: %s"  , split[0])
    }
    if split[1] != "barbaz\nturtles" {
        t.Errorf("wrong: %s"  , split[1])
    }
}
