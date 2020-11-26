// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

import (
	"net"
	"os"
	"strings"
)

// An definition of the public functions for a diary instance
type Diary interface{
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
	Page(level int, sample int, catch bool, category string, pageMeta M, authType, authIdentifier string, authMeta M, scope func(p Page)) error

	// Load a page using it JSON definition to chain multiple logs together when crossing micro-service boundaries
	Load(data []byte, category string, scope func(p Page)) error
}

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
func Dear(client, project, service string, serviceMeta M, repository, commitHash string, commitTags []string, commitMeta M, level int) Diary {
	if level < LevelTrace || level > LevelFatal {
		panic("level must be a value between 0 - 6")
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

	// todo: add a handler to push to agent, api, queue, systemctl, etc. (do something with log instead of just outputting it to console)
	// todo: add flag to disable human readable output to console (in case we want to send service console output directly to systemctl logs for legacy systems)

	return &diary{
		Level: level,
		Service: Service{
			Client: client,
			Project: project,
			Service: service,
			Meta: serviceMeta,

			Host: host,
			HostIps: ips,

			ParentProcessId: os.Getppid(),
			ProcessId: os.Getpid(),
		},
		Commit: Commit{
			Repository: repository,
			Hash: commitHash,
			Tags: commitTags,
			Meta: commitMeta,
		},
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
func (d diary) Page(level int, sample int, catch bool, category string, pageMeta M, authType, authIdentifier string, authMeta M, scope func(p Page)) (response error) {
	if level == -1 {
		level = d.Level
	}
	if level < LevelTrace || level > LevelFatal {
		panic("level must be a value between 0 - 6 or -1 to use default level")
	}

	p := page{
		Diary: d,
		Chain: Chain{
			Id: "abc...xyz",
			Meta: pageMeta,
			Auth: Auth{
				Type: authType,
				Identifier: authIdentifier,
				Meta: authMeta,
			},
		},
		Sample: sample,
		Level: level,
		Catch: catch,
		Category: category,
	}

	func() {
		if catch {
			// defer func() {
			//     response = catch(category)()
			// }()
		}
		// defer trace(category)()
		scope(p)
	}()
	return nil
}

func (d diary) Load(data []byte, category string, scope func(p Page)) error {
	if len(strings.TrimSpace(category)) == 0 {
		panic("category may not be empty")
	}
	if scope == nil {
		panic("scope must be defined")
	}

	p, err := parsePage(data)
	if err != nil {
		return err
	}
	func() {
		// defer func() {
		//     response = catch(category)()
		// }()
		// defer trace(category)()
		scope(p)
	}()

	return nil
}
