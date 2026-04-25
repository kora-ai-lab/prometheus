package main

import "testing"

func TestModelQwen(t *testing.T) {
	if !TestModel("qwen2.5:0.5b") {
		t.Fail()
	}
}

func TestModelPhi3(t *testing.T) {
	if !TestModel("phi3:mini") {
		t.Fail()
	}
}

func TestModelGemma(t *testing.T) {
	if !TestModel("gemma4:e2b") {
		t.Fail()
	}
}
