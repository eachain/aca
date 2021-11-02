// Aho-Corasick automation
// AC自动机
package aca

import (
	"unicode"
	"unicode/utf8"
)

type node struct {
	next map[rune]*node
	fail *node
	word string
	size int // len([]rune(word))
}

type ACA struct {
	root *node
	fold bool
}

// New returns an empty aca.
// equalFold likes strings.EqualFold
func New(equalFold ...bool) *ACA {
	fold := false
	if len(equalFold) > 0 {
		fold = equalFold[0]
	}
	return &ACA{root: &node{}, fold: fold}
}

// Add adds a new word to aca.
// After Add, and before Find,
// MUST Build.
func (a *ACA) Add(word string) {
	rs := []rune(word)
	n := a.root
	for _, r := range rs {
		if n.next == nil {
			n.next = make(map[rune]*node)
		}
		if a.fold {
			r = unicode.ToLower(r)
		}
		if n.next[r] == nil {
			n.next[r] = &node{}
		}
		n = n.next[r]
	}
	n.word = word
	n.size = len(rs)
}

// Del delete a word from aca.
// After Del, and before Find,
// MUST Build.
func (a *ACA) Del(word string) {
	a.del(a.root, []rune(word), 0)
}

func (a *ACA) del(n *node, rs []rune, i int) {
	if i >= len(rs) {
		n.fail = nil
		n.word = ""
		n.size = 0
		return
	}

	r := rs[i]
	if a.fold {
		r = unicode.ToLower(r)
	}

	t := n.next[r]
	if t == nil {
		return
	}

	a.del(t, rs, i+1)
	if t.word == "" && len(t.next) == 0 {
		delete(n.next, r)
	}
}

// Build builds the fail pointer.
// It MUST be called before Find.
func (a *ACA) Build() {
	// attention: must use queue, cannot build recusively
	q := []*node{a.root}

	for len(q) > 0 {
		n := q[0]
		q = q[1:]

		for r, c := range n.next {
			q = append(q, c)

			p := n.fail
			for p != nil {
				if p.next[r] != nil {
					c.fail = p.next[r]
					break
				}
				p = p.fail
			}
			if p == nil {
				c.fail = a.root
			}
		}
	}
}

type (
	MatchFunc   func(i, j int, match string)
	ExcludeFunc func(r rune) bool
)

// FindLiteral finds all the matched words contains in s, and skip the rune excluded.
func (a *ACA) FindLiteral(s string, match MatchFunc, exclude ExcludeFunc) {
	index := make([]int, 0, len(s))
	n := a.root
	for i, r := range s {
		if exclude != nil && exclude(r) {
			continue
		}
		index = append(index, i)

		if a.fold {
			r = unicode.ToLower(r)
		}

		for n.next[r] == nil && n != a.root {
			n = n.fail
		}
		n = n.next[r]
		if n == nil {
			n = a.root
			continue
		}

		rl := utf8.RuneLen(r)
		for t := n; t != a.root; t = t.fail {
			if t.word != "" {
				match(index[len(index)-t.size], i+rl, t.word)
			}
		}
	}
}

// FindExclude finds all the matched words contains in s, and skip the rune excluded.
// The results may duplicated.
// It is caller's responsibility to make results unique.
func (a *ACA) FindExclude(s string, exclude ExcludeFunc) (matches []string) {
	a.FindLiteral(s, func(_, _ int, match string) {
		matches = append(matches, match)
	}, exclude)
	return
}

// Find finds all the matched words contains in s.
// The results may duplicated.
// It is caller's responsibility to make results unique.
func (a *ACA) Find(s string) []string {
	return a.FindExclude(s, nil)
}

// Block records the low and high position that words appear, namely s[low:high].
type Block struct {
	Low     int
	High    int
	Literal string // equals to s[low:high]
	Match   string // match word, if set equalFold or exclude something, the Literal may not equal Match
}

// FindBlocks returns the blocks that words in aca appear.
func (a *ACA) FindBlocks(s string, exclude ExcludeFunc) (blocks []Block) {
	a.FindLiteral(s, func(i, j int, match string) {
		blocks = append(blocks, Block{Low: i, High: j, Literal: s[i:j], Match: match})
	}, exclude)
	return
}
