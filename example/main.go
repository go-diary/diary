package main

import (
	"github.com/go-diary/diary"
	"sync"
)

var d diary.Diary

func init() {
	d = diary.Dear("client", "project", "service", diary.M{"type":"service"}, "repository", "hash", []string{}, diary.M{"type":"commit"}, diary.LevelTrace, diary.HumanReadableHandler)
}

func main() {
	group := sync.WaitGroup{}
	channel := make(chan []byte)

	go func() {
		group.Add(1)
		defer group.Done()

		select {
		case data := <-channel:
			d.Load(data, "channel", func(p diary.Page) {
				p.Debug("x", true)
			})
			break
		}
	}()

	d.Page(-1, 1000, true, "main", diary.M{}, "", "", nil, func(p diary.Page) {
		p.Debug("x", true)
		channel <- p.ToJSON()
		panic("test")
	})

	group.Wait()
}
