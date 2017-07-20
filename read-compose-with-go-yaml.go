package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

type Service struct {
  Labels map[string]string
}

type ComposeFile struct {
  Services map[string]Service
}

func main() {
  data, err := ioutil.ReadFile("docker-compose.yml")
  check(err)
  fmt.Print(string(data))

  compose := ComposeFile{}

  err = yaml.Unmarshal([]byte(data), &compose)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  fmt.Printf("services: %#+v\n", compose.Services)
}
