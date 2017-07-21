package main

import (
  "fmt"
  "log"
  "io/ioutil"

  // "github.com/pkg/errors"
  yaml "gopkg.in/yaml.v2"
)

type Service struct {
  Labels interface{}
}

type ComposeData struct {
  Services map[string]Service
}

func getComposeData(filename string) (*ComposeData, error) {
	bytes, err := ioutil.ReadFile(filename)
  if err != nil {
		return nil, err
	}

  var cfg ComposeData
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

  return &cfg, nil
}

// GetServiceGroupingFrom parses a Docker Compose file and figures out the service groups
// by looking for 'com.icalialabs.plis.group' labels.
func main() {
  composeData, err := getComposeData("docker-compose.yml")
  if err != nil { log.Fatalf("error: %v", err) }

  fmt.Println(composeData)

  // config, err := load(configDetails)
	// if err != nil { return nil, err }
  //
  // serviceGroups, err := getServiceGrouping(config)
  // if err != nil { return nil, err }
  //
  // return serviceGroups, nil
}
