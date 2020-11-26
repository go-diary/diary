package main

import "diary"

func main() {
	d := diary.Dear("client", "project", "service", diary.M{"type":"service"}, "repository", "hash", []string{}, diary.M{"type":"commit"}, diary.LevelTrace)
	if err := d.Page(-1, 1, false, "main", diary.M{}, "", "", nil, func(p diary.Page) {
		x := 100
		p.Debug("x", x)
		p.Info("info", nil)
		p.Notice("notice", nil)
		p.Warning("warning", "this is a warning", nil)
		p.Error("error", "this is an error", nil)
		p.Fatal("fatal", "this is a fatal error", 1, nil)
	}); err != nil {
		panic(err)
	}
}
