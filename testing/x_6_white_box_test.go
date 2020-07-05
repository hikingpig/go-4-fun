package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"testing"
)

func bytesInUse(username string) int64 {
	return usage[username]
}

var usage = make(map[string]int64)

const sender = "notifications@example.com"
const password = "correcthorsebatterystaple"
const hostname = "smtp.example.com"

/*
- we must use %% to get % in the string
*/
const template = `Warning: you are using %d bytes of storage, %d%% of your quota.`

var notifyUser = func(username, msg string) {
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender, []string{username}, []byte(msg))
	if err != nil {
		log.Printf("smtp.SendEmail(%s) failed: %s", username, err)
	}
}

func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 1000000000 // 1GB
	percent := 100 * used / quota
	if percent < 90 {
		return
	}
	msg := fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)
}

/*
this test is to check if:
	- CheckQuota is called with correct input
	- CheckQuota produces correct error msg
*/
func TestCheckQuotaNotifiesUser(t *testing.T) {
	const user = "joe@example.org"
	var notifiedUser, notifiedMsg string
	saved := notifyUser
	savedUsage := usage[user]
	defer func() {
		notifyUser = saved
		usage[user] = savedUsage
	}()
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}
	usage[user] = 980000000
	CheckQuota(user)
	if notifiedUser == "" && notifiedMsg == "" {
		t.Fatalf("notifyUser not called")
	}
	if notifiedUser != user {
		t.Errorf("wrong user (%s) notified, want %s", notifiedUser, user)
	}
	const wantSubstring = "98% of your quota"
	if !strings.Contains(notifiedMsg, wantSubstring) {
		t.Errorf("unexpected notification message <<%s>>, "+"want substring %q", notifiedMsg, wantSubstring)
	}
}
