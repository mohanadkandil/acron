package main

import "testing"

func TestGreet(t *testing.T) {
	t.Run("normal name", func(t *testing.T) {
		got := Greet("mohanned")
		want := "hello, mohanned"
		if got != want {
			t.Fatalf("Greet() = %q, want %q", got, want)
		}
	})

	t.Run("blank name defaults", func(t *testing.T) {
		got := Greet("   ")
		want := "hello, world"
		if got != want {
			t.Fatalf("Greet() = %q, want %q", got, want)
		}
	})
}


