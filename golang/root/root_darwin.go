package root

import (
	"fmt"
	"log"
	"os/exec"
)

type macExecutor struct{}

func NewExecutor() Executor {
	return new(macExecutor)
}

func (e *macExecutor) Exec(s, msg string) string {
	c := fmt.Sprintf("do shell script \"%s\" with administrator privileges", s)
	b, err := exec.Command("osascript", "-e", c).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
