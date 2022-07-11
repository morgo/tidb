package main

import (
	"fmt"

	"github.com/pingcap/tidb/parser/auth"
)

func isAuditable(user *auth.UserIdentity) bool {
	return true
}

func log(user *auth.UserIdentity, eventType string, info map[string]string) {
	fmt.Printf("[%s] [%s]: %#v\n", eventType, user.Username, info)
}
