package translation

import (
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Exec(command []string) {
	color.Cyan(strings.Join(command, " "))

	binary, lookErr := exec.LookPath(command[0])
	if lookErr != nil {
		panic(lookErr)
	}

	execErr := syscall.Exec(binary, command, os.Environ())
	if execErr != nil {
		panic(execErr)
	}
}
