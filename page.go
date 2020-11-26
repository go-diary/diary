// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
)

// A private function used to parse a page instance from its JSON definition
var parsePage = func(data []byte) (Page, error) {
	return nil, nil
}

// An definition of the public functions for a page instance
type Page interface{
	Debug(key string, value interface{})
	Info(category string, meta M)
	Notice(category string, meta M)
	Warning(category, message string, meta M)
	Error(category, message string, meta M)
	Fatal(category, message string, code int, meta M)
}

// normally only used for troubleshooting
func (p page) Debug(key string, value interface{}) {
	data, err := json.Marshal(Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: levelToText(p.Level),
		Category: p.Category,
		Line: "",
		Stack: "",
		Message: "",
		Meta: M{
			key: value,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// normally inside of a loop
func (p page) Info(category string, meta M) {
	if meta == nil {
		meta = M{}
	}
	data, err := json.Marshal(Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: levelToText(p.Level),
		Category: fmt.Sprintf("%s.%s", p.Category, category),
		Line: "",
		Stack: "",
		Message: "",
		Meta: meta,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// normally outside of a loop
func (p page) Notice(category string, meta M) {
	if meta == nil {
		meta = M{}
	}
	data, err := json.Marshal(Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: levelToText(p.Level),
		Category: fmt.Sprintf("%s.%s", p.Category, category),
		Line: "",
		Stack: "",
		Message: "",
		Meta: meta,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// - category: (may be empty)
func (p page) Warning(category, message string, meta M) {
	if meta == nil {
		meta = M{}
	}
	data, err := json.Marshal(Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: levelToText(p.Level),
		Category: fmt.Sprintf("%s.%s", p.Category, category),
		Line: "",
		Stack: "",
		Message: message,
		Meta: meta,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func (p page) Error(category, message string, meta M) {
	if meta == nil {
		meta = M{}
	}
	data, err := json.Marshal(Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: levelToText(p.Level),
		Category: fmt.Sprintf("%s.%s", p.Category, category),
		Line: "",
		Stack: string(debug.Stack()),
		Message: message,
		Meta: meta,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

// application will be force to exit
func (p page) Fatal(category, message string, code int, meta M) {
	if meta == nil {
		meta = M{}
	}
	data, err := json.Marshal(Log{
		Service: p.Diary.Service,
		Commit: p.Diary.Commit,
		Chain: p.Chain,
		Level: levelToText(p.Level),
		Category: fmt.Sprintf("%s.%s", p.Category, category),
		Line: "",
		Stack: string(debug.Stack()),
		Message: message,
		Meta: meta,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	os.Exit(code)
}

func (p page) Scope(category string, scope func(p Page)) error {
	return nil
}

func (p page) ToJSON() []byte {
	return nil
}