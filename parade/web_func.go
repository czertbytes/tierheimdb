package main

import (
	"fmt"
	"time"
)

const timeFormat = "Jan 2, 2006 at 15:04pm"

func DateFormatter(args ...interface{}) string {
	if s, ok := args[0].(string); ok == true {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		return t.Format(timeFormat)
	}

	return ""
}
