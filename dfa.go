package aca

import (
	"unicode"
)

type DFA struct {
	root *node
	fold bool
}

func NewDFA(equalFold ...bool) *DFA {
	fold := false
	if len(equalFold) > 0 {
		fold = equalFold[0]
	}
	return &DFA{root: &node{}, fold: fold}
}

func (a *DFA) Add(word string) {
	n := a.root
	for _, r := range word {
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
}

func (a *DFA) Del(word string) {
	a.del(a.root, []rune(word), 0)
}

func (a *DFA) del(n *node, rs []rune, i int) {
	if i >= len(rs) {
		n.word = ""
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

func (a *DFA) Build() {}

func (a *DFA) FindLiteral(s string, match MatchFunc, exclude ExcludeFunc) {
	index := make([]int, 0, len(s)+1)
	rs := make([]rune, 0, len(s))
	for i, r := range s {
		if exclude != nil && exclude(r) {
			continue
		}
		index = append(index, i)
		if a.fold {
			r = unicode.ToLower(r)
		}
		rs = append(rs, r)
	}
	index = append(index, len(s))

	for i := 0; i < len(rs); i++ {
		n := a.root
		for j := i; j < len(rs); j++ {
			n = n.next[rs[j]]
			if n == nil {
				break
			}
			if n.word != "" {
				match(index[i], index[j+1], n.word)
			}
		}
	}
}

func (a *DFA) FindExclude(s string, exclude ExcludeFunc) (matches []string) {
	a.FindLiteral(s, func(_, _ int, match string) {
		matches = append(matches, match)
	}, exclude)
	return
}

func (a *DFA) Find(s string) []string {
	return a.FindExclude(s, nil)
}

func (a *DFA) FindBlocks(s string, exclude ExcludeFunc) (blocks []Block) {
	a.FindLiteral(s, func(i, j int, match string) {
		blocks = append(blocks, Block{Low: i, High: j, Literal: s[i:j], Match: match})
	}, exclude)
	return
}
