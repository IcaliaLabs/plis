package translation

func BypassToCompose(cmd string, args []string) {
  Exec(append([]string{"docker-compose", cmd}, args...))
}
