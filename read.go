package main

import (
  "fmt"
  "log"
  "reflect"
  "io/ioutil"
  "strings"

  yaml "gopkg.in/yaml.v2"
)

type Service struct {
  Labels map[string]interface{}
}

type ServiceRead struct {
  Labels interface{}
}

type ComposeDataRead struct {
  Services map[string]ServiceRead
}

func convertLabelList(givenList []interface{}) (map[string]interface{}, error) {
  convertedMap := make(map[string]interface{})
  for _, valueInterface := range givenList {
    keyVal := strings.SplitN(valueInterface.(string), "=", 2)
    key := strings.TrimRight(keyVal[0], " ")

    var value interface{}
    if len(keyVal) > 1 {
      value = strings.TrimLeft(keyVal[1], " ")
    } else {
      value = true // It was a 'value-less' label...
    }
    convertedMap[key] = value
  }
  return convertedMap, nil
}

func convertLabelMap(givenMap map[interface {}]interface {}) (map[string]interface{}, error) {
  convertedMap := make(map[string]interface{})
  for keyInterface, valueInterface := range givenMap {
    convertedMap[keyInterface.(string)] = valueInterface
  }
  return convertedMap, nil
}

func getServices(filename string) (map[string]Service, error) {
	bytes, err := ioutil.ReadFile(filename)
  if err != nil {
		return nil, err
	}

  var cfg ComposeDataRead
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

  serviceMap := make(map[string]Service)

  for serviceName, serviceRead := range cfg.Services {
    service := Service{}
    labelsType := reflect.TypeOf(serviceRead.Labels)
    labelsKind := labelsType.Kind()

    var err error
    if labelsKind == reflect.Map {
      service.Labels, err = convertLabelMap(serviceRead.Labels.(map[interface {}]interface {}))
    } else if labelsKind == reflect.Slice {
      service.Labels, err = convertLabelList(serviceRead.Labels.([]interface{}))
    }
    if err != nil { return nil, err }
    serviceMap[serviceName] = service
  }

  return serviceMap, nil
}

// GetServiceGroupingFrom parses a Docker Compose file and figures out the service groups
// by looking for 'com.icalialabs.plis.group' labels.
func main() {
  services, err := getServices("docker-compose.yml")
  if err != nil { log.Fatalf("error: %v", err) }

  for serviceName, service := range services {
    fmt.Print("Service: ")
    fmt.Println(serviceName)
    for labelName, labelValue := range service.Labels {
      fmt.Print(" - label: '")
      fmt.Print(labelName)
      fmt.Println("'")
      fmt.Print("   value: '")
      fmt.Print(labelValue)
      fmt.Println("'")
    }
  }
}
