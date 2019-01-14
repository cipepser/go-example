# go-example

Goのサンプル集。

## Functional Options

```go
package functionalOptions

type RequestOptions struct {
	page    int
	perPage int
	sort    string
}

type Option func(request *RequestOptions)

func Page(p int) Option {
	return func(r *RequestOptions) {
		if r != nil {
			r.page = p
		}
	}
}

func PerPage(pp int) Option {
	return func(r *RequestOptions) {
		if r != nil {
			r.perPage = pp
		}
	}
}

func Sort(s string) Option {
	return func(r *RequestOptions) {
		if r != nil {
			r.sort = s
		}
	}
}

func NewRequest(opts ...Option) *RequestOptions {
	// set default values
	r := &RequestOptions{
		page:    1,
		perPage: 30,
		sort:    "desc",
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}
```

実行コード

```go
package main

import (
	"fmt"

	"github.com/cipepser/go-example/functionalOptions"
)

func main() {
	r := functionalOptions.NewRequest()
	fmt.Println(r) // &{1 30 desc}

	r = functionalOptions.NewRequest(functionalOptions.Page(10))
	fmt.Println(r) // &{10 30 desc}

	r = functionalOptions.NewRequest(
		functionalOptions.Page(10),
		functionalOptions.PerPage(2),
		functionalOptions.Sort("ast"),
	)
	fmt.Println(r) // &{10 2 ast}
}
```

やること自体はシンプルだけど、自分で書かないといけないコードが多いので、
wireを使った[こっち](https://github.com/cipepser/go-DI-sample)のほうが保守性も込みでよいと思う。  
簡単な使い捨てコードを書くときはこの`Functional Options`を使ってもいいけど、
長期的にはあんまり使わないかなぁという所感。

## Observer Pattern

※レシーバの変数名が、元記事の[Go Patternsで学ぶGo \- Qiita](https://qiita.com/mnuma/items/109458d90ce9dbdde426)では、
型名と一致していなくてわかりづらいので書き換えている。  
同様にコンストラクタも追加で実装している（実際に利用するときはパッケージ外から利用することを想定）。

```go
package observerPattern

import "fmt"

type (
	Event struct {
		Date int64
	}

	Observer interface {
		OnNotify(Event)
	}

	Notifier interface {
		Register(Observer)
		Deregister(Observer)
		Notify(Event)
	}
)

type (
	EventObserver struct {
		id int
	}

	EventNotifier struct {
		observers map[Observer]struct{}
	}
)

func NewEventObserver(id int) *EventObserver {
	return &EventObserver{id: id}
}

func NewEventNotifier() EventNotifier {
	return EventNotifier{
		observers: make(map[Observer]struct{}),
	}
}

func (eo *EventObserver) OnNotify(e Event) {
	fmt.Printf("*** Observer %d received: %d\n", eo.id, e.Date)
}

func (en *EventNotifier) Register(o Observer) {
	en.observers[o] = struct{}{}
}

func (en *EventNotifier) Deregister(o Observer) {
	delete(en.observers, o)
}

func (en *EventNotifier) Notify(e Event) {
	for o := range en.observers {
		o.OnNotify(e)
	}
}
```

利用例

```go
package main

import (
	"time"

	"github.com/cipepser/go-example/observerPattern"
)

func main() {
	en := observerPattern.NewEventNotifier()

	en.Register(observerPattern.NewEventObserver(1))
	en.Register(observerPattern.NewEventObserver(2))

	stop := time.NewTimer(10 * time.Second).C
	tick := time.NewTicker(time.Second).C

	for {
		select {
		case <-stop:
			return
		case t := <-tick:
			en.Notify(observerPattern.Event{
				Date: t.UnixNano(),
			})
		}
	}
}
```

実行結果（以下を10秒間で表示したのちexit 0になる）

```
*** Observer 1 received: 1546421712877318000
*** Observer 2 received: 1546421712877318000
*** Observer 1 received: 1546421713881677000
*** Observer 2 received: 1546421713881677000
*** Observer 1 received: 1546421714882143000
*** Observer 2 received: 1546421714882143000
*** Observer 1 received: 1546421715882261000
*** Observer 2 received: 1546421715882261000
*** Observer 1 received: 1546421716882365000
*** Observer 2 received: 1546421716882365000
*** Observer 1 received: 1546421717877488000
*** Observer 2 received: 1546421717877488000
*** Observer 1 received: 1546421718882481000
*** Observer 2 received: 1546421718882481000
*** Observer 1 received: 1546421719881285000
*** Observer 2 received: 1546421719881285000
*** Observer 1 received: 1546421720877872000
*** Observer 2 received: 1546421720877872000
```

## Semaphore

```go
package semaphore

import (
	"errors"
	"time"
)

var (
	ErrNoTickets      = errors.New("semaphore: could not acquire semaphore")
	ErrIllegalRelease = errors.New("semaphore: can't release the semaphore without acquiring it first")
)

type Interface interface {
	Acquire() error
	Release() error
}

type Implementation struct {
	sem     chan struct{}
	timeout time.Duration
}

func (s *Implementation) Acquire() error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-time.After(s.timeout):
		return ErrNoTickets
	}
}

func (s *Implementation) Release() error {
	select {
	case <-s.sem:
		return nil
	case <-time.After(s.timeout):
		return ErrIllegalRelease
	}
}

func NewInterface(tickets int, timeout time.Duration) Interface {
	return &Implementation{
		sem:     make(chan struct{}, tickets),
		timeout: timeout,
	}
}
```

使い方

```go
package main

import (
	"fmt"
	"time"

	"github.com/cipepser/go-example/semaphore"
)

func main() {
	tickets, timeout := 10, 6*time.Second
	s := semaphore.NewInterface(tickets, timeout)

	for i := 0; i <= 100; i++ {
		if err := s.Acquire(); err != nil {
			panic(err)
		}

		go func(i int) {
			doHeavyProcess(i)

			if err := s.Release(); err != nil {
				panic(err)
			}
		}(i)
	}
}

func doHeavyProcess(i int) {
	fmt.Printf("process[%d] starts\n", i)
	time.Sleep(5 * time.Second)
	fmt.Printf("process[%d] ends\n", i)
}
```

結果

```
process[9] starts
process[4] starts
process[3] starts
process[5] starts
process[6] starts
process[7] starts
process[8] starts
process[1] starts
process[2] starts
process[0] starts
process[2] ends
process[0] ends
process[5] ends
process[8] ends
process[10] starts
process[1] ends
process[3] ends
process[4] ends
process[9] ends
process[17] starts
process[6] ends
process[11] starts
process[12] starts
process[13] starts
process[14] starts
process[15] starts
process[16] starts
process[7] ends
process[19] starts
process[18] starts
process[13] ends
process[20] starts
process[17] ends
process[21] starts
process[11] ends
process[22] starts
process[10] ends
process[23] starts
process[15] ends
process[14] ends
process[24] starts
process[25] starts
process[16] ends
process[19] ends
process[26] starts
process[18] ends
process[27] starts
process[28] starts
process[12] ends
process[29] starts
process[29] ends
process[25] ends
process[31] starts
process[21] ends
process[32] starts
process[20] ends
process[27] ends
process[34] starts
process[30] starts
process[23] ends
process[26] ends
process[35] starts
process[24] ends
process[28] ends
process[33] starts
process[37] starts
process[38] starts
process[36] starts
process[22] ends
process[39] starts
process[34] ends
process[39] ends
process[40] starts
process[38] ends
process[42] starts
process[30] ends
process[43] starts
process[32] ends
process[31] ends
process[41] starts
process[35] ends
process[36] ends
process[33] ends
process[37] ends
process[48] starts
process[46] starts
process[47] starts
process[45] starts
process[44] starts
process[49] starts
process[48] ends
process[46] ends
process[49] ends
process[43] ends
process[40] ends
process[53] starts
process[45] ends
process[54] starts
process[44] ends
process[55] starts
process[51] starts
process[47] ends
process[56] starts
process[41] ends
process[52] starts
process[57] starts
process[50] starts
process[42] ends
process[58] starts
process[59] starts
process[52] ends
process[51] ends
process[53] ends
process[56] ends
process[61] starts
process[60] starts
process[54] ends
process[57] ends
process[62] starts
process[58] ends
process[63] starts
process[64] starts
process[50] ends
process[65] starts
process[66] starts
process[55] ends
process[68] starts
process[59] ends
process[69] starts
process[67] starts
process[66] ends
process[70] starts
process[60] ends
process[64] ends
process[72] starts
process[68] ends
process[61] ends
process[74] starts
process[67] ends
process[75] starts
process[71] starts
process[69] ends
process[76] starts
process[63] ends
process[65] ends
process[78] starts
process[73] starts
process[77] starts
process[62] ends
process[79] starts
process[72] ends
process[70] ends
process[76] ends
process[71] ends
process[73] ends
process[82] starts
process[83] starts
process[78] ends
process[85] starts
process[77] ends
process[86] starts
process[75] ends
process[87] starts
process[79] ends
process[88] starts
process[84] starts
process[74] ends
process[89] starts
process[81] starts
process[80] starts
process[84] ends
process[90] starts
process[88] ends
process[91] starts
process[86] ends
process[92] starts
process[87] ends
process[81] ends
process[83] ends
process[80] ends
process[85] ends
process[82] ends
process[89] ends
process[93] starts
process[94] starts
process[99] starts
process[96] starts
process[95] starts
process[98] starts
process[97] starts
process[97] ends
process[96] ends
process[98] ends
process[94] ends
process[93] ends
process[99] ends
process[95] ends
process[92] ends
process[91] ends
```

`doHeavyProcess`内のタイムアウトを7秒（セマフォのタイムアウトより長い）にすると
以下のように次のプロセグが`Acquire`できない。

```
process[9] starts
process[5] starts
process[6] starts
process[4] starts
process[7] starts
process[8] starts
process[1] starts
process[2] starts
process[0] starts
process[3] starts
panic: semaphore: could not acquire semaphore
```


## References
* [Go Patternsで学ぶGo \- Qiita](https://qiita.com/mnuma/items/109458d90ce9dbdde426)
* [Go Patterns · GitBook](http://tmrts.com/go-patterns/)