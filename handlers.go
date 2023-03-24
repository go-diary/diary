package diary

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func DefaultHandler(log Log) {
	data, err := json.Marshal(log)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func HumanReadableHandler(log Log) {
	fmt.Printf("[%s] %-8s %-90s %s %s\n", log.Time.Format(time.RFC3339), strings.ToUpper(log.Level), log.Category, log.Chain.Id, log.Line)
	if log.Meta != nil && len(log.Meta) > 0 {
		out, err := json.MarshalIndent(log.Meta, "", "    ")
		if err == nil && out != nil {
			fmt.Println(string(out))
		}
	}
	if len(strings.TrimSpace(log.Message)) > 0 {
		fmt.Printf("**> %s\n", log.Message)
	}
	if len(strings.TrimSpace(log.Stack)) > 0 {
		fmt.Println(strings.TrimSuffix(log.Stack, "\n"))
	}
}