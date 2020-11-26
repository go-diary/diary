package main

import "diary"

func main() {
	d := diary.Dear("client", "project", "service", diary.M{}, "repository", "hash", []string{}, diary.M{}, diary.LevelTrace)
}
