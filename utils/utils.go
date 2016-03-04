package utils

import (
	"fmt"
	"time"

	math_rand "math/rand"

	"github.com/gorilla/feeds"
)

func init() {
	math_rand.Seed(time.Now().Unix())
}

func GenerateKey() string {
	u := feeds.NewUUID()
	return fmt.Sprintf("%x%x%x%x%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func GetNowSecond() int {
	return int(time.Now().Unix())
}

func GenerateUUID() string {
	u := feeds.NewUUID()
	return string(u.String())
}

func GetNowStringYMD() string {
	return time.Now().Format("2006-01-02")
}

func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}
