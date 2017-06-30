package aca

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type view struct {
	remain rune
	last   rune
	a      *ACA
	depth  int
	nodeid map[*node]int
}

func sortedkeys(next map[rune]*node) []rune {
	if len(next) == 0 {
		return nil
	}
	ints := make([]int, 0, len(next))
	for r := range next {
		ints = append(ints, int(r))
	}
	sort.Ints(ints)
	runes := make([]rune, len(ints))
	for i, r := range ints {
		runes[i] = rune(r)
	}
	return runes
}

func newView(a *ACA, remain, last rune) *view {
	idx := 0
	nodeid := make(map[*node]int)
	maxDepth := 0

	var prescan func(*node, int)
	prescan = func(n *node, depth int) {
		nodeid[n] = idx
		idx++
		if depth > maxDepth {
			maxDepth = depth
		}

		for _, r := range sortedkeys(n.next) {
			prescan(n.next[r], depth+1)
		}
	}
	prescan(a.root, 1)

	return &view{
		remain: remain,
		last:   last,
		a:      a,
		depth:  maxDepth,
		nodeid: nodeid,
	}
}

func (v *view) nodestr(r rune, n *node) string {
	s := fmt.Sprintf("\033[36m%c\033[0m"+
		"[\033[32m%v\033[0m](fail->\033[33m%v\033[0m)",
		r, v.nodeid[n], v.nodeid[n.fail])
	if n.wordLength > 0 { // a complete word
		s += " \033[31m√\033[0m"
	}
	return s
}

func (v *view) print(w io.Writer, r rune, n *node, prefix []rune) {
	if len(prefix) > 0 {
		for i := 0; i < len(prefix)-1; i++ {
			if prefix[i] == v.remain {
				fmt.Fprintf(w, "%-4c", '|')
			} else {
				fmt.Fprintf(w, "%-4c", ' ')
			}
		}
		fmt.Fprintf(w, "%c── ", prefix[len(prefix)-1])
	}
	fmt.Fprintln(w, v.nodestr(r, n))

	runes := sortedkeys(n.next)
	for i, r := range runes {
		if i < len(runes)-1 {
			v.print(w, r, n.next[r], append(prefix, v.remain))
		} else {
			v.print(w, r, n.next[r], append(prefix, v.last))
		}
	}
}

func (v *view) show(w io.Writer, root rune) {
	v.print(w, root, v.a.root, make([]rune, 0, v.depth))
}

func (a *ACA) debug(w ...io.Writer) {
	var writer io.Writer
	if len(w) == 0 {
		writer = os.Stdout
	} else if len(w) == 1 {
		writer = w[0]
	} else {
		writer = io.MultiWriter(w...)
	}
	newView(a, '├', '└').show(writer, '.')
}
