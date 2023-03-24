// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"net"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var counter int = 0
var flux int = 0

// Dear returns a diary.Diary interface instance for consumption
//
// - client: The shorthand code used to identify which client the log belongs to
// - project: The shorthand code used to identify which client-project the log belongs to
// - service: The shorthand code used to identify which service the log belongs to
// - serviceMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
//
// - repository: The URI of the source-code repository (may be empty)
// - commitHash: The shorthand hash of the given source commit that the service was built off of (may be empty)
// - commitTags: The tags associated with the given source commit that the service was built off of (may be empty)
// - commitMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
//
// - level: The default level to log at [NOTE: Normally NOTICE for production services]
// - handler: A routine to handle log entries  [NOTE: ]
func Dear(client, project, service string, serviceMeta M, repository, commitHash string, commitTags []string, commitMeta M, level int, handler H) IDiary {
	if level < LevelTrace || level > LevelAudit {
		panic("level must be a value between 0 - 7")
	}

	// get the hostname of the server that the service is running on
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// get the host ips of the server that the service is running on
	ips := make([]string, 0)
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range networkInterfaces {
		addresses, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		for _, address := range addresses {
			if ip, ok := address.(*net.IPNet); ok {
				if ip != nil && ip.IP != nil && ip.IP.IsGlobalUnicast() {
					ips = append(ips, ip.IP.String())
				}
			}
		}
	}

	return &diary{
		Level:   level,
		Handler: handler,
		Service: Service{
			Client:  client,
			Project: project,
			Service: service,
			Meta:    serviceMeta,

			Host:    host,
			HostIps: ips,

			ParentProcessId: os.Getppid(),
			ProcessId:       os.Getpid(),
		},
		Commit: Commit{
			Repository: repository,
			Hash:       commitHash,
			Tags:       commitTags,
			Meta:       commitMeta,
		},
	}
}

// A private struct to encapsulate diary instance logic
type diary struct {
	Level   int
	Handler H
	Service Service
	Commit  Commit
}

// Page issues a diary.Page interface instance for consumption
// In a page scope all logs will be linked to the same page identifier
// This allows us to trace the entire page chain easily for troubleshooting and profiling
//
// - level: The default level to log at [NOTE: If -1 will inherit from diary instance]
// - sample: A per second count indicating how frequently traces should be sampled, only applicable if level is higher than DEBUG [NOTE: if value is less than zero then will sample all traces]
// - catch: A flag indicating if the scope should automatically catch and return errors. [NOTE: If set to true panics will be caught and returned as an error.]
// - category: The shorthand code used to identify the given workflow category [NOTE: Categories will be concatenated by dot-nation: "main.sub1.sub2.sub3".]
// - authType: The shorthand code for the type of auth account (may be empty)
// - authIdentifier: The identifier, which can be anything, used to identify the given auth account (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
// - authMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
func (d diary) Page(level int, sample int, catch bool, category string, pageMeta M, authType, authIdentifier string, authMeta M, scope S) {
	if err := d.PageX(level, sample, catch, category, pageMeta, authType, authIdentifier, authMeta, scope); err != nil && !catch {
		panic(err)
	}
}

// Page returns a diary.Page interface instance for consumption
// In a page scope all logs will be linked to the same page identifier
// This allows us to trace the entire page chain easily for troubleshooting and profiling
//
// - level: The default level to log at [NOTE: If -1 will inherit from diary instance]
// - sample: A per second count indicating how frequently traces should be sampled, only applicable if level is higher than DEBUG [NOTE: if value is less than zero then will sample all traces]
// - catch: A flag indicating if the scope should automatically catch and return errors. [NOTE: If set to true panics will be caught and returned as an error.]
// - category: The shorthand code used to identify the given workflow category [NOTE: Categories will be concatenated by dot-nation: "main.sub1.sub2.sub3".]
// - authType: The shorthand code for the type of auth account (may be empty)
// - authIdentifier: The identifier, which can be anything, used to identify the given auth account (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
// - authMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
func (d diary) PageX(level int, sample int, catch bool, category string, pageMeta M, authType, authIdentifier string, authMeta M, scope S) (response error) {
	if level == -1 {
		level = d.Level
	}
	if level < LevelTrace || level > LevelAudit {
		panic("level must be a value between 0 - 7 or -1 to use default level")
	}

	if authMeta == nil {
		authMeta = M{}
	}
	if pageMeta == nil {
		pageMeta = M{}
	}
	if sample < 0 {
		sample = 0
	}
	fluxRate := int(float64(sample) * 0.05) // add 5% flux to ensure that a different trace is sampled each time
	if fluxRate > 0 {
		flux = rand.Intn(fluxRate)
	}

	p := page{
		Diary: d,

		Chain: Chain{
			Id:   primitive.NewObjectID().Hex(),
			Meta: pageMeta,
			Auth: Auth{
				Type:       authType,
				Identifier: authIdentifier,
				Meta:       authMeta,
			},
		},
		Sample:   sample,
		Level:    level,
		Catch:    catch,
		Category: strings.TrimPrefix(strings.TrimPrefix(category, d.Service.Service), "."),
	}

	return pageScope(p, scope)
}

func (d diary) Load(data []byte, category string, scope S) {
	if err := d.LoadX(data, category, scope); err != nil {
		panic(err)
	}
}

func (d diary) LoadX(data []byte, category string, scope S) error {
	if len(strings.TrimSpace(category)) == 0 {
		panic("category may not be empty")
	}
	if scope == nil {
		panic("scope must be defined")
	}

	p, err := parsePage(data, d)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(p.Category, category) {
		if p.Category != "" {
			p.Category = fmt.Sprintf("%s.%s", p.Category, category)
		} else {
			p.Category = category
		}
	}

	return pageScope(p, scope)
}

func pageScope(p page, scope S) (response error) {
	cat := p.Category
	if cat == "" {
		cat = p.Diary.Service.Service
	}

	if p.Catch {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				if assertedErr, ok := r.(error); ok {
					err = assertedErr
				}
				response = err

				if p.Level > LevelError {
					return
				}

				_, file, line, _ := runtime.Caller(2)
				log := Log{
					Service:  p.Diary.Service,
					Commit:   p.Diary.Commit,
					Chain:    p.Chain,
					Level:    TextLevelError,
					Category: cat,
					Line:     fmt.Sprintf("%s:%d", file, line),
					Stack:    string(debug.Stack()),
					Message:  fmt.Sprint(response),
					Meta:     M{},
					Time:     time.Now(),
				}
				if p.Diary.Handler != nil {
					p.Diary.Handler(log)
				} else {
					DefaultHandler(log)
				}
			}
		}()
	}

	trace := true
	if p.Level > LevelTrace {
		trace = false
		counter++
		if counter > p.Sample-flux {
			trace = true
			counter = 0
			fluxRate := int(float64(p.Sample) * 0.05) // add 5% flux to ensure that a different trace is sampled each time
			if fluxRate > 0 {
				flux = rand.Intn(fluxRate)
			}
		}
	}

	if trace {
		defer func() func() {
			_, file, line, _ := runtime.Caller(3)
			enter := time.Now()
			log := Log{
				Service:  p.Diary.Service,
				Commit:   p.Diary.Commit,
				Chain:    p.Chain,
				Level:    TextLevelTraceEnter,
				Category: cat,
				Line:     fmt.Sprintf("%s:%d", file, line),
				Stack:    "",
				Message:  "",
				Time:     time.Now(),
			}
			if p.Diary.Handler != nil {
				p.Diary.Handler(log)
			} else {
				DefaultHandler(log)
			}
			return func() {
				exit := time.Now()
				var minutes = exit.Sub(enter).Minutes()
				var seconds = exit.Sub(enter).Seconds()
				var milliSeconds = exit.Sub(enter).Milliseconds()
				var microSeconds = exit.Sub(enter).Microseconds()
				if milliSeconds < 1 {
					milliSeconds = 0
					seconds = 0
					minutes = 0
				} else if seconds < 1 {
					seconds = 0
					minutes = 0
				} else if minutes < 1 {
					minutes = 0
				}

				log := Log{
					Service:  p.Diary.Service,
					Commit:   p.Diary.Commit,
					Chain:    p.Chain,
					Level:    TextLevelTraceExit,
					Category: cat,
					Line:     fmt.Sprintf("%s:%d", file, line),
					Stack:    "",
					Message:  "",
					Meta: M{
						"enter":    enter,
						"exit":     exit,
						"durationMinutes": minutes,
						"durationSeconds": seconds,
						"durationMilliSeconds": milliSeconds,
						"durationMicroSeconds": microSeconds,
					},
					Time: time.Now(),
				}
				if p.Diary.Handler != nil {
					p.Diary.Handler(log)
				} else {
					DefaultHandler(log)
				}
			}
		}()()
	}

	scope(p)

	return nil
}
