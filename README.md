# ACA (Aho-Corasick automation)
Golang implementation of Aho-Corasick algorithm.

* Aho-Corasick Wikipedia : [Aho-Corasick algorithm wiki](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm)

# Prerequisite

* Golang 1.7+

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

# Contribution

If you want to participate, you can create an issue or request a 'Pull Request'.

Welcome any and all suggestions.

# License

MIT License
