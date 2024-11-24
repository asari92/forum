package main

import (
	"fmt"
	"regexp"
)

var (
	EmailRX    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	TextRX     = regexp.MustCompile(`^[а-яА-ЯёЁa-zA-Z0-9.,:;!?'"()\-–—\[\]{}<>/|@#$%^&*+=_~\s]+$`)
	UsernameRX = regexp.MustCompile(`^[^._ ](?:[\w-]|\.[\w-])+[^._ ]$`)
)

func main() {
	fmt.Println(Matches(`а-яА-ЯёЁa-zA-Z0-9.,:;!?'"()-–—[]{}<>/|#$%^&*+=_~]+$`, TextRX))
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
