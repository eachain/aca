# ACA (Aho-Corasick automation)
Golang implementation of Aho-Corasick algorithm.

* Aho-Corasick Wikipedia : [Aho-Corasick algorithm wiki](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm)

Aho-Corasick automation，1975年产生于贝尔实验室，著名的多模匹配算法。

# Prerequisite

* Golang 1.7+

# Preview

Debug mode, you can print the hole trie (with color).

调试模式下，可以把整棵树(带颜色地)打印出来.

```go
a := aca.New()
a.Add("abcdefg")
a.Add("abcd")
a.Add("bcdefg")
a.Add("cde")
a.Add("cdefg")
a.Add("defg")
a.Add("efg")
a.Add("fg")
a.Add("g")
a.Add("abcdeg")
a.Build()
a.Debug() // not exported
```

Then, this is the trie(after optimised):

然后, 这是(优化后的)树:

```
.[0](fail->0)
├── a[1](fail->0)
|   └── b[2](fail->0)
|       └── c[3](fail->0)
|           └── d[4](fail->0) √
|               └── e[5](fail->17)
|                   ├── f[6](fail->0)
|                   |   └── g[7](fail->14) √
|                   └── g[8](fail->29) √
├── b[9](fail->0)
|   └── c[10](fail->0)
|       └── d[11](fail->0)
|           └── e[12](fail->17)
|               └── f[13](fail->0)
|                   └── g[14](fail->19) √
├── c[15](fail->0)
|   └── d[16](fail->0)
|       └── e[17](fail->0) √
|           └── f[18](fail->0)
|               └── g[19](fail->23) √
├── d[20](fail->0)
|   └── e[21](fail->0)
|       └── f[22](fail->0)
|           └── g[23](fail->26) √
├── e[24](fail->0)
|   └── f[25](fail->0)
|       └── g[26](fail->28) √
├── f[27](fail->0)
|   └── g[28](fail->29) √
└── g[29](fail->0) √
```

打印规则: "字符[ID]\(fail->ID)"，如果是完整单词，会在后面出现"√"。

优化过程: 一般fail指针指向的下个节点，仍然不是完整的单词，故一直fail下去，直到找到完整的单词，或到root。优化后可减少Find过程遍历的节点数。

# Usage

### 1. Find out all words appear in a sentence.(查找所有出现过的词)

***REMEBER: Your call sequence should always be: "Add/Del", "Build", and then "Find".***

***注意: 你应该永远以"Add/Del", "Build", "Find"的顺序调用。***

**If you call "Find" before "Build", it will panic.**

**如果你在"Build"之前调用了"Find"，你的程序会挂逼。**

```go
package main

import (
	"fmt"

	"github.com/eachain/aca"
)

func main() {
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

	words := a.Find("yasherhs")
	fmt.Println(words) // prints: [she he her]
}
```

### 2. You should write uniq function yourself.(你应该自己实现唯一函数。)

For example:

```go
func Uniq(words []string) []string {
    if len(words) == 0 {
        return words
    }
	sort.Strings(words)
	n := 0
	for i := 1; i < len(words); i++ {
		if words[n] != words[i] {
			n++
			words[n] = words[i]
		}
	}
	return words[:n+1]
}
```

### 3. If you want any concurrency for both read and write, encapsulate it yourself.(如果你想要读写并发安全，你应该自己封装实现。)

For example:

```go
type Filter struct {
	a   *aca.ACA
	mtx sync.RWMutex
}

func (f *Filter) AddAndBuild(word string) {
	f.mtx.Lock()
	if f.a == nil {
		f.a = aca.New()
	}
	f.a.Add(word)
	f.a.Build()
	f.mtx.Unlock()
}

func (f *Filter) DelAndBuild(word string) {
	f.mtx.Lock()
	f.a.Del(word)
	f.a.Build()
	f.mtx.Unlock()
}

func (f *Filter) Find(s string) []string {
	f.mtx.RLock()
	words := f.a.Find(s)
	f.mtx.RUnlock()
	return words
}
```

*But if you just build once, and only find after that, it is concurrently-read-secure. Like this:*

*如果你只会build一次，之后只find，那它本身是并发读安全的。*

```go
package main

import (
	"fmt"
	"time"

	"github.com/eachain/aca"
)

var globalAca *aca.ACA

func init() {
	a := aca.New()
	a.Add("say")
	a.Add("she")
	a.Add("shr")
	a.Add("he")
	a.Add("her")
	a.Build()

	globalAca = a
}

func main() {
	for i := 0; i < 100; i++ {
		go func() {
			words := globalAca.Find("yasherhs")
			fmt.Println(words) // prints: [she he her]
		}()
	}
	time.Sleep(time.Second)
}
```

# Debug

You can open debug mode by modify file `https://github.com/eachain/aca/blob/master/debug.go` line 99, 'debug' to 'Debug'.

你可以通过修改文件`https://github.com/eachain/aca/blob/master/debug.go`第99行，'debug'改为'Debug'，打开debug模式。

# Contribution

If you want to participate, you can create an issue or request a 'Pull Request'.

Welcome any and all suggestions.

