package aca_test

import (
	"testing"

	"github.com/eachain/aca"
)

func TestEmpty(t *testing.T) {
	a := aca.New()
	a.Build()
	s := "test text"
	if len(a.Find(s)) != 0 {
		t.Errorf("empty aca find result: %v", a.Find(s))
	}
}

func TestDelLeaf(t *testing.T) {
	a := aca.New()
	a.Add("abcd")
	a.Add("abcdefg")
	a.Add("abcdexy")
	a.Del("abcdefg")
	a.Build()

	s := "abcdefgabcdexy"
	if len(a.Find(s)) != 3 { // [abcd, abcd, abcdexy]
		t.Errorf("aca find after del result: %v", a.Find(s))
	}
}

func TestDelMid(t *testing.T) {
	a := aca.New()
	a.Add("abcd")
	a.Add("abcdefg")
	a.Add("abcdxyz")
	a.Add("abcdemn")
	a.Del("abcd")
	a.Build()

	s := "abcdefg"
	if len(a.Find(s)) != 1 { // [abcdefg]
		t.Errorf("aca find after del result: %v", a.Find(s))
	}
}

func TestFind(t *testing.T) {
	a := aca.New()
	a.Add("say")
	a.Add("erh")
	a.Add("she")
	a.Add("shr")
	a.Del("erh")
	a.Add("he")
	a.Del("shr")
	a.Add("her")
	a.Build()

	words := a.Find("yasherhs") // [she, he, her]
	if len(words) != 3 {
		t.Errorf("aca find words count: %v", len(words))
	}
	results := []string{"she", "he", "her"}
	for i := range words {
		if words[i] != results[i] {
			t.Errorf("aca find word[%v]: %v", i, words[i])
		}
	}
}
