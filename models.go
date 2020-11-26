// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

import "time"

// A package shorthand for a map[string]interface
type M map[string]interface{}

// A public struct to encapsulate the service details for a log entry
type Service struct {
	Client string `json:"client"`
	Project string `json:"project"`
	Service string `json:"service"`
	Host string `json:"host"`
	HostIps []string `json:"hostIps"`
	ProcessId int `json:"pid"`
	ParentProcessId int `json:"ppid"`
	Meta M `json:"meta"`
}

// A public struct to encapsulate the commit details for a log entry
type Commit struct {
	Repository string `json:"repository"`
	Hash string `json:"hash"`
	Tags []string `json:"tags"`
	Meta M `json:"meta"`
}

// A public struct to encapsulate the auth details for a log entry
type Auth struct {
	Type string `json:"type"`
	Identifier string `json:"identifier"`
	Meta M `json:"meta"`
}

// A public struct to encapsulate the chain details for a log entry
type Chain struct {
	Id string `json:"id"`
	Meta M `json:"meta"`
	Auth Auth `json:"auth"`
}

// A public struct to encapsulate the log entry details
type Log struct {
	Service Service `json:"service"`
	Commit Commit `json:"commit"`
	Chain Chain `json:"chain"`
	Level string `json:"level"`
	Category string `json:"category"`
	Line string `json:"line"`
	Stack string `json:"stack"`
	Message string `json:"message"`
	Meta M `json:"meta"`
	Time time.Time `json:"time"`
}

// A private struct to encapsulate diary instance logic
type diary struct {
	Level int
	Handler func(log Log)
	Service Service
	Commit Commit
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
