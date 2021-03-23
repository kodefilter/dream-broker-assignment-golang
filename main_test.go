package main

import "testing"

func TestGetWordCount(t *testing.T) {
	if GetWordCount("This is three") != 3 {
		t.Error("Expected 3 as word count")
	}
}
