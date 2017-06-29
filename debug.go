package aca

// for debug

func (n *node) show(s []rune) {
	if n.wordLength > 0 {
		println(string(s))
	}
	for r, c := range n.next {
		c.show(append(s, r))
	}
}

func (a *ACA) show() {
	println("node count:", a.nodeCount)
	a.root.show(make([]rune, 0, 32))
}
