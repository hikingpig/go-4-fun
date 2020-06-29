package main

import (
	"fmt"
	"testing"
	"unicode"
)

func TestPalindrome(t *testing.T) {
	if !IsPalindrome("detartrated") {
		t.Error(`IsPalindrome("detartrated") = false`)
	}
	if !IsPalindrome("kayak") {
		t.Error(`IsPalindrome("kayak") = false`)
	}
}

func TestNonPalindrome(t *testing.T) {
	if IsPalindrome("palindrome") {
		t.Error(`IsPalindrome("palindrome") = true`)
	}
}

/*
- é is a non-ascii character. using string(byte) is not correct
- if we print, é is 195 for the 1st, 116 for the 2nd. and we can not printf with %s for string element
- when we turn it to rune, we have same value for both
- we must use rune
*/
func TestFrenchPalindrome(t *testing.T) {
	m := "été"
	fmt.Println("==== m[0] == m[2]", m[0] == m[2])
	fmt.Printf("==== m[0] %d, m[2] %d, m[0] %v, m[2] %v\n", m[0], m[2], m[0], m[2])
	if !IsPalindrome("été") {
		t.Error(`IsPalindrome("été") = false`)
	}
}

/*
we should ignore the white spaces and cases
*/
func TestCanalPalindrome(t *testing.T) {
	input := "a man, a plan, a canal: panama"
	if !IsPalindrome(input) {
		t.Errorf(`IsPalindrome(%q) = false`, input)
	}
}

func IsPalindrome(s string) bool {
	var letters []rune
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToLower(r))
		}
	}
	for i := range letters {
		if letters[i] != letters[len(letters)-1-i] {
			return false
		}
	}
	return true
}
