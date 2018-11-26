package main

import (
	"log"
	"net/http"
	"time"
)

const (
	numPollers     = 2
	pollInterval   = 60 * time.Second
	statusInterval = 10 * time.Second
	errTimeout     = 10 * time.Second
)

var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
}

// State : Stateは最後に確認したURLの状態を保持
type State struct {
	url    string
	status string
}

// Resource :
// ポーリングしたURLのの状態を保持
// 起動するとURLごとに１つのりResourceが割り当てられて、
// メインgoroutineとPoller goroutineがチャネル上でお互いに、
// それぞれのリソースを送り合う
type Resource struct {
	url      string
	errCount int
}

func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(updateInterval)

	// selectによって受付をブロックしている
	// urlStatusを持っているので、パラレル読み書きを制限し、
	// シリアルなアクセスにしているらしい
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates
}

func logState(s map[string]string) {
	log.Println("Current State:")
	for k, v := range s {
		log.Printf("%s %s", k, v)
	}
}

// Poller :
// チャネルからResourceポインタをもらってPollする。
// 完了したら、StateMonitorにStateを伝えて、Resourceポインタを返却する
// ポインタ渡しによって同一データに2つのPoller goroutineが操作することがなくなる。
// すなわちSync Packageを利用してロックをする必要がなくなる
func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	for r := range in {
		s := r.Poll()
		status <- State{r.url, s}
		out <- r
	}
}

// Poll :
// HTTP HEADリクエストをすることで対象URLのステータスコードを取得する
// ミスするとエラーメッセージを返却しResourceのエラー回数を増やす
func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval + errTimeout*time.Duration(r.errCount))
	done <- r
}

func ResourceQueing(in chan *Resource) {
	for _, url := range urls {
		in <- &Resource{url: url}
	}
}

func access_check() {
	// 入力チャネルと出力チャネルの作成
	pending, complete := make(chan *Resource), make(chan *Resource)

	// StateMonitorを起動
	status := StateMonitor(statusInterval)

	// goroutineでPollerを起動
	for i := 0; i < numPollers; i++ {
		go Poller(pending, complete, status)
	}

	// Resourceをキューに挿入する
	go ResourceQueing(pending)

	for r := range complete {
		go r.Sleep(pending)
	}
}
