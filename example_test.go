package aca_test

import (
	"bytes"
	"fmt"
	"sort"
	"unicode/utf8"

	"github.com/eachain/aca"
)

type byPos []aca.Block

func (bs byPos) Len() int { return len(bs) }

func (bs byPos) Swap(i, j int) { bs[i], bs[j] = bs[j], bs[i] }

func (bs byPos) Less(i, j int) bool {
	if bs[i].Low < bs[j].Low {
		return true
	}
	if bs[i].Low == bs[j].Low {
		return bs[i].High < bs[j].High
	}
	return false
}

func UnionBlocks(blocks []aca.Block) []aca.Block {
	if len(blocks) == 0 {
		return blocks
	}

	sort.Sort(byPos(blocks))
	n := 0
	for i := 1; i < len(blocks); i++ {
		if blocks[i].Low <= blocks[n].Low {
			if blocks[i].High > blocks[n].High {
				blocks[n].High = blocks[i].High
			}
		} else {
			n++
			blocks[n] = blocks[i]
		}
	}
	return blocks[:n+1]
}

func ReplaceAll(a *aca.ACA, s string, new rune) string {
	tmp := make([]rune, utf8.RuneCountInString(s))
	for i := range tmp {
		tmp[i] = new
	}

	now := 0
	buf := &bytes.Buffer{}
	for _, b := range UnionBlocks(a.FindBlocks(s, nil)) {
		buf.WriteString(s[now:b.Low])
		cnt := utf8.RuneCountInString(s[b.Low:b.High])
		buf.WriteString(string(tmp[:cnt]))
		now = b.High
	}
	if now < len(s) {
		buf.WriteString(s[now:])
	}
	return buf.String()
}

func ExampleSensitives() {
	a := aca.New()
	a.Add("fuck")
	a.Add("shit")
	a.Add("bitch")
	a.Add("艹")
	a.Add("就是")
	a.Add("傻X")
	a.Add("他奶奶的")
	a.Del("就是")
	a.Build()

	s := "我fuck你shit up, 艹他奶奶的个球嘞, you这个bitch，就是个傻X!"
	fmt.Println(a.Find(s))
	fmt.Println(a.FindBlocks(s, nil))
	fmt.Println(UnionBlocks(a.FindBlocks(s, nil)))
	fmt.Println(ReplaceAll(a, s, '*'))
	// Output:
	// [fuck shit 艹 他奶奶的 bitch 傻X]
	// [{3 7 fuck fuck} {10 14 shit shit} {19 22 艹 艹} {22 34 他奶奶的 他奶奶的} {54 59 bitch bitch} {71 75 傻X 傻X}]
	// [{3 7 fuck fuck} {10 14 shit shit} {19 22 艹 艹} {22 34 他奶奶的 他奶奶的} {54 59 bitch bitch} {71 75 傻X 傻X}]
	// 我****你**** up, *****个球嘞, you这个*****，就是个**!
}
