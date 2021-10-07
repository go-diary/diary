// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

// A private function used to parse a page instance from its JSON definition
var parsePage = func(data []byte, d diary) (page, error) {
	var p page
	if err := json.Unmarshal(data, &p); err != nil {
		return page{}, err
	}
	p.Diary = d
	return p, nil
}

// A private struct to encapsulate page instance logic
type page struct {
	Diary diary
	Chain Chain
	Category string
	Sample int
	Level int
	Catch bool
}

// return parent diary
func (p page) Parent() IDiary {
	return p.Diary
}

// normally only used for troubleshooting
func (p page) Debug(key string, value interface{}) {
	if p.Level > LevelDebug {
		return
	}

	cat := key
	if p.Category != "" {
		cat = fmt.Sprintf("%s.%s", p.Category, key)
	}

	_, file, line, _ := runtime.Caller(1)
	log := Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: TextLevelDebug,
		Category: cat,
		Line: fmt.Sprintf("%s:%d", file, line),
		Stack: "",
		Message: "",
		Meta: M{
			key: value,
		},
		Time: time.Now(),
	}
	if p.Diary.Handler != nil {
		p.Diary.Handler(log)
	} else {
		DefaultHandler(log)
	}
}

// normally inside of a loop
func (p page) Info(category string, meta M) {
	if p.Level > LevelInfo {
		return
	}

	cat := category
	if p.Category != "" {
		cat = fmt.Sprintf("%s.%s", p.Category, category)
	}

	_, file, line, _ := runtime.Caller(1)
	if meta == nil {
		meta = M{}
	}
	log := Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: TextLevelInfo,
		Category: cat,
		Line: fmt.Sprintf("%s:%d", file, line),
		Stack: "",
		Message: "",
		Meta: meta,
		Time: time.Now(),
	}
	if p.Diary.Handler != nil {
		p.Diary.Handler(log)
	} else {
		DefaultHandler(log)
	}
}

// normally outside of a loop
func (p page) Notice(category string, meta M) {
	if p.Level > LevelNotice {
		return
	}

	cat := category
	if p.Category != "" {
		cat = fmt.Sprintf("%s.%s", p.Category, category)
	}

	_, file, line, _ := runtime.Caller(1)
	if meta == nil {
		meta = M{}
	}
	log := Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: TextLevelNotice,
		Category: cat,
		Line: fmt.Sprintf("%s:%d", file, line),
		Stack: "",
		Message: "",
		Meta: meta,
		Time: time.Now(),
	}
	if p.Diary.Handler != nil {
		p.Diary.Handler(log)
	} else {
		DefaultHandler(log)
	}
}

// - category: (may be empty)
func (p page) Warning(category, message string, meta M) {
	if p.Level > LevelWarning {
		return
	}

	cat := category
	if p.Category != "" {
		cat = fmt.Sprintf("%s.%s", p.Category, category)
	}

	_, file, line, _ := runtime.Caller(1)
	if meta == nil {
		meta = M{}
	}
	log := Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: TextLevelWarning,
		Category: cat,
		Line: fmt.Sprintf("%s:%d", file, line),
		Stack: "",
		Message: message,
		Meta: meta,
		Time: time.Now(),
	}
	if p.Diary.Handler != nil {
		p.Diary.Handler(log)
	} else {
		DefaultHandler(log)
	}
}

func (p page) Error(category, message string, meta M) {
	if p.Level > LevelError {
		return
	}

	cat := category
	if p.Category != "" {
		cat = fmt.Sprintf("%s.%s", p.Category, category)
	}

	_, file, line, _ := runtime.Caller(1)
	if meta == nil {
		meta = M{}
	}
	log := Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: TextLevelError,
		Category: cat,
		Line: fmt.Sprintf("%s:%d", file, line),
		Stack: string(debug.Stack()),
		Message: message,
		Meta: meta,
		Time: time.Now(),
	}
	if p.Diary.Handler != nil {
		p.Diary.Handler(log)
	} else {
		DefaultHandler(log)
	}
}

// application will be force to exit
func (p page) Fatal(category, message string, code int, meta M) {
	cat := category
	if p.Category != "" {
		cat = fmt.Sprintf("%s.%s", p.Category, category)
	}

	_, file, line, _ := runtime.Caller(1)
	if meta == nil {
		meta = M{}
	}
	log := Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: TextLevelFatal,
		Category: cat,
		Line: fmt.Sprintf("%s:%d", file, line),
		Stack: string(debug.Stack()),
		Message: message,
		Meta: meta,
		Time: time.Now(),
	}
	if p.Diary.Handler != nil {
		p.Diary.Handler(log)
	} else {
		DefaultHandler(log)
	}
	os.Exit(code)
}

func (p page) Scope(category string, scope S) error {
	return p.Diary.LoadX(p.ToJson(), category, scope)
}

func (p page) ToJson() []byte {
	data, err := json.Marshal(struct {
		Service Service `json:"service"`
		Commit Commit `json:"commit"`
		Chain Chain `json:"chain"`
		Category string `json:"category"`
		Sample int `json:"sample"`
		Level int `json:"level"`
		Catch bool `json:"catch"`
	}{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Category: p.Category,
		Sample: p.Sample,
		Level: p.Level,
		Catch: p.Catch,
	})
	if err != nil {
		panic(err)
	}
	return data
}