package utils

import (
	"fmt"
	"regexp"
	"time"
)

func RegexpUtils(str, s string) []string {
	//定义正则表达式
	regexpCompile := regexp.MustCompile(str)
	//使用正则表达式找与之相匹配的字符串，返回一个数组包含子表达式匹配的字符串
	return regexpCompile.FindStringSubmatch(s)
}
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
func TimeFormat20060102(t time.Time) string {
	return t.Format("20060102150405")
}
func TimeParse(t string) string {
	t0, _ := time.Parse("20060102150405", t)
	return TimeFormat(t0)
}

// TimeExpire 是否在当前时间点的有效时间段内
func TimeExpire(unix int64, t time.Duration) bool {

	startTime := time.Unix(unix, 0)
	sub := time.Now().Sub(startTime)
	fmt.Println(sub.Minutes(),t.Minutes())

	if sub.Minutes() < t.Minutes() {
		return true
	}

	return false
}
