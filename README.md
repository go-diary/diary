# Diary

**Note:** This project is still in development, has not even started pre-alpha yet.

## Logging Checklist [ [pdf](https://github.com/go-diary/diary/raw/main/Logging%20Checklist.pdf) ]
### 10 Commandments of Logging “Masterzen”

- [X] **1. Thou shalt not write log by yourself**
  - Make use of syslog or similar systems as they handle auto-ration and more.

- [X] **2. Thou shalt log at the proper level**
  - **TRACE** track code routines.
  - **DEBUG** activated during troubleshooting for debugging.
  - **INFO** user-driven or system specific actions.
  - **NOTICE** notable events that are not errors (default level for prod), always add context.
  - **WARN** events that could result in error, e.g low disk space, always add context.
  - **ERROR** for all errors, always add context.
  - **FATAL** signifies the end of a program (exit program), always add context. 

- [X] **3. Honor thy log category**
  - The category allows classification of the log message. E.g. my.service.api.<apitoken>
  - Log categories are hierarchical to allow inherited behaviours, rules and filters.

- [X] **4. Thou shalt write meaningful logs**
  - Treat logging as if there is no access to the program source-code.
  - Log should not depend on previous log for context as it may not be there (async).

- [X] **5. Thy log shalt be written in English**
  - English is an internationally recognized language.
  - English is a better suited technical language.
  - English can be written in ASCII characters (stays the same over UTF-8, Unicode, etc).

- [X] **6. Thou shalt log with context**
  - Treat logging as if there is no access to the program source-code.
  - Log messages without context are just noise.
  - Ensure that at least the local scope of variables are included in the log context.

- [X] **7. Thou shalt log in machine parsable format**
  - Text logs are good for humans but very poor for machines.
  - Use a simple international standard format, like JSON.
  - Simplifies automation processing for alerting/auditing.
  - Log parsers require less processing of messages.
  - Log search engine indexing becomes straightforward.

- [X] **8. Thou shalt not log too much or too little**
  - Too much logging generates too much clutter (time consuming to browse).
  - Too little logging leads to troubleshooting problems (not enough context).
  - Log more than enough in dev, then tighten it up before shipping to prod.

- [X] **9. Thou shalt think to the reader**
  - The whole purpose of logging is so that someone will read it one day.
  - Readers deserve consistency, so stick to standards like RFC3339 (datetime).
  - Log parsers are readers too, so ensure consistency of log format/dictionary.

- [X] **10. Thou shalt not log only for troubleshooting**
  - **Auditing** management/legal events, describe what users of the system are doing.
  - **Profiling** logs are time stamped this allows us to infer performance metrics.
  - **Statistics** compute user behaviours, e.g. alert when too many errors detected in a row.

### References
* http://www.masterzen.fr/2013/01/13/the-10-commandments-of-logging/

## Basic Usage

### Installation
```
go get github.com/go-diary/diary
```

### Create page
```
package main

import (
    "github.com/go-diary/diary"
)

func main() {
    instance := diary.Dear("client", "project", "service", diary.M{}, "repository", "hash", []string{}, diary.M{}, diary.LevelTrace, diary.HumanReadableHandler)
    instance.Page(-1, 1000, true, "main", diary.M{}, "", "", nil, func(p diary.Page) {
        x := 100
        p.Debug("x", x)
    })
}
```

### Load Page
```
package main

import (
	"diary"
	"sync"
)

var channel = make(chan []byte)
var instance diary.Diary

func main() {
	group := sync.WaitGroup{}
	go func() {
		group.Add(1)
		defer group.Done()

		routine()
	}()

	instance = diary.Dear("uprate", "go-diary", "diary", diary.M{}, "git@github.com:go-diary/diary.git", "084c59f", []string{}, diary.M{}, diary.LevelTrace, diary.HumanReadableHandler)
	instance.Page(-1, 1000, true, "main", diary.M{}, "", "", nil, func(p diary.Page) {
		x := 100
		p.Debug("x", x)
		channel <- p.ToJSON()
	})

	group.Wait()
}

func routine() {
	select {
		case data := <-channel:
			instance.Load(data, "routine", func(p diary.Page) {
				y := 200
				p.Debug("y", y)
			})
			break
	}
}
```
#### NOTES
- Use `git config --get remote.origin.url` to get repository URL using script.
- Use `git rev-parse --short HEAD` to get short commit has using script.

### License
This project is licensed under the MIT license. See the [LICENSE](https://github.com/go-diary/diary/blob/main/LICENSE) file for more info.