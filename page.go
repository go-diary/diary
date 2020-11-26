// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

// A private function used to parse a page instance from its JSON definition
var parsePage = func(data []byte) (func(Page), error) {
	return nil, nil
}

// An definition of the public functions for a page instance
type Page interface{
	Debug()
	Info()
	Notice()
	Warn()
	Error()
	Fatal()
}

// normally only used for troubleshooting
func (p page) Debug(key string, value interface{}) {
}

// normally inside of a loop
func (p page) Info(category string, meta M) {
}

// normally outside of a loop
func (p page) Notice(category string, meta M) {
}

// - category: (may be empty)
func (p page) Warning(category, message string, meta M) {
}

func (p page) Error(category, message string, meta M) {
}

// application will be force to exit
func (p page) Fatal(category, message string, meta M) {
}

func (p page) ToJSON() []byte {
	return nil
}

func (p page) Scope(category string) (func(p Page), error) {
}