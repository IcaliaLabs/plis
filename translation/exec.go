package translation

import (
  "strings"
  "syscall"
  "os"
  "os/exec"
  "github.com/fatih/color"
)

func Exec(command []string) {
  color.Cyan(strings.Join(command, " "))

  binary, lookErr := exec.LookPath(command[0])
  if lookErr != nil { panic(lookErr) }

  execErr := syscall.Exec(binary, command, os.Environ())
  if execErr != nil { panic(execErr) }
}
