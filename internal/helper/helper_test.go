package helper

import (
	"testing"
)

func TestIsCharAllowedInKeyOrVar(t *testing.T) {
	t.Run("it should return false", func(t *testing.T) {
		chars := []byte{ '2', '#', '+', '-', '~' }
		
		for i, ch := range chars {
			result := IsCharAllowedInKeyOrVar(ch)

			if result {
				t.Fatalf(
					"[test #%d]: Expected IsCharAllowedInKeyOrVar('%c') to return false, but got true",
					i, ch,
				)
			}
		}
	})

	t.Run("it should return true", func(t *testing.T) {
		chars := []byte{'l', 'c', '_', 'a'}
		
		for i, ch := range chars {
			result := IsCharAllowedInKeyOrVar(ch)

			if !result {
				t.Fatalf(
					"[test #%d]: Expected IsCharAllowedInKeyOrVar('%c') to return true, but got false",
					i, ch,
				)
			}
		}
	})
}

func TestIsDigit(t *testing.T) {

	t.Run("it should return false", func(t *testing.T) {
		chars := []byte{ 'e', 'f', 'g', 'h', 'i' }

		for i, ch := range chars {
			result := IsDigit(ch)

			if result {
				t.Fatalf(
					"[test #%d]: Expected IsDigit('%c') to return false, but got true",
					i, ch,
				)
			}
		}
	})

	t.Run("it should return true", func(t *testing.T) {
		for i := range 10 {
			result := IsDigit(byte(48 + i))

			if !result {
				t.Fatalf(
					"[test #%d]: Expected IsDigit('%c') to return true, but got false",
					i, 48 + i,
				)
			}
		}
	})
}

func TestFlipMap(t *testing.T) {
	
	m := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	fm := FlipMap(m)
	i := 1

	for k, v := range fm {

		if k != m[i] {
			t.Fatalf(
				"Expecting %q, got %q\n",
				k, m[i],
			)
		}

		if v != i {
			t.Fatalf(
				"Expecting %d, got %d\n",
				i, v,
			)
		}

		i++
	}

}
