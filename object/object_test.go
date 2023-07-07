package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := NewString("Hello World")
	hello2 := NewString("Hello World")
	diff1 := NewString("My name is johnny")
	diff2 := NewString("My name is johnny")

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
