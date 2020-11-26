// Copyright 2020 The Uprate Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Capture useful application logs for troubleshooting, auditing, profiling and statistics

package diary

import (
	"fmt"
	"net"
	"os"
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
	Page(level string, sample int, catch bool, category string, authType, authIdentifier string, authMeta M, scope func(p Page)) error

	// Load a page using it JSON definition to chain multiple logs together when crossing micro-service boundaries
	Load()
}

// Dear returns a diary.Diary interface instance for consumption
//
// - client: The shorthand code used to identify which client the log belongs to
// - project: The shorthand code used to identify which client-project the log belongs to
// - service: The shorthand code used to identify which service the log belongs to
// - serviceMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
//
// - commitRepository: The URI of the source-code repository (may be empty)
// - commitHash: The shorthand hash of the given source commit that the service was built off of (may be empty)
// - commitTags: The tags associated with the given source commit that the service was built off of (may be empty)
// - commitMeta: Can contain any other additional data that you may require on logs for troubleshooting (may be empty) [WARNING: Don't ever log personal data without first encrypting or salt-hashing the data.]
// - level: The default level to log at [NOTE: Normally NOTICE for production services]
func Dear(client, project, service string, serviceMeta M, commitRepository, commitHash string, commitTags []string, commitMeta M, level int) Diary {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}

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

	fmt.Printf("host: %s\n", host)
	fmt.Printf("ips: %v\n", ips)
	fmt.Printf("pid: %d\n", os.Getpid())
	fmt.Printf("ppid: %d\n", os.Getppid())

	// todo: add a handler to push to agent, api, queue, systemctl, etc. (do something with log instead of just outputting it to console)
	// todo: add flag to disable human readable output to console (in case we want to send service console output directly to systemctl logs for legacy systems)

	return nil
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
func (d diary) Page(level, sample int, catch bool, category string, authType, authIdentifier string, authMeta M, scope func(p Page)) error {
	// todo: automatically log trace
	// todo: automatically catch and return panic as error
	return nil
}

func (d diary) Load(data []byte, category string) (func(p Page), error) {
	return parsePage(data)
}
