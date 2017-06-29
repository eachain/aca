// Aho-Corasick automation
// AC自动机
package aca

import "unicode/utf8"

type node struct {
	next       map[rune]*node
	fail       *node
	wordLength int
}

type ACA struct {
	root      *node
	nodeCount int
}

// New returns an empty aca.
func New() *ACA {
	return &ACA{root: &node{}}
}

// Add adds a new word to aca,
// after Add, and before Find,
// MUST Build.
func (a *ACA) Add(word string) {
	n := a.root
	for _, r := range word {
		if n.next == nil {
			n.next = make(map[rune]*node)
		}
		if n.next[r] == nil {
			n.next[r] = &node{}
			a.nodeCount++
		}
		n = n.next[r]
	}
	n.wordLength = len(word)
}

// Build builds the fail pointer.
// It MUST be called before Find.
func (a *ACA) Build() {
	// allocate enough memory as a queue
	q := make([]*node, 0, a.nodeCount)
	q = append(q, a.root)

	for len(q) > 0 {
		n := q[0]
		q = q[1:]

		for r, c := range n.next {
			q = append(q, c)

			p := n.fail
			for p != nil {
				// ATTENTION: nil map cannot be writen, but CAN BE READ!!!
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

func (a *ACA) find(s string, cb func(start, end int)) {
	n := a.root
	for i, r := range s {
		for n.next[r] == nil && n != a.root {
			n = n.fail
		}
		n = n.next[r]
		if n == nil {
			n = a.root
			continue
		}

		end := i + utf8.RuneLen(r)
		for t := n; t != a.root; t = t.fail {
			if t.wordLength > 0 {
				cb(end-t.wordLength, end)
			}
		}
	}
}

// Find finds all the words contains in s.
// The results may duplicated.
// It is caller's responsibility to make results unique.
func (a *ACA) Find(s string) (words []string) {
	a.find(s, func(start, end int) {
		words = append(words, s[start:end])
	})
	return
}

// Block records the start and end position
// that words appear, namely s[start:end].
type Block struct {
	Start, End int
}

// Blocks returns the blocks that words in aca appear.
func (a *ACA) Blocks(s string) (blocks []Block) {
	a.find(s, func(start, end int) {
		blocks = append(blocks, Block{Start: start, End: end})
	})
	return
}
