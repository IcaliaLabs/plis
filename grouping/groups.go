package grouping

import (
  "reflect"
  "io/ioutil"
  "strings"

  yaml "gopkg.in/yaml.v2"
)

// Service is a Docker Compose service entry with it's corresponding labels mapped
type Service struct {
  Labels map[string]string
}

type serviceRead struct {
  Labels interface{}
}

type composeDataRead struct {
  Services map[string]serviceRead
}

func convertLabelList(givenList []interface{}) (map[string]string, error) {
  convertedMap := make(map[string]string)
  for _, valueInterface := range givenList {
    keyVal := strings.SplitN(valueInterface.(string), "=", 2)
    keyVal = append(keyVal, "true") // in case it was a value-less label
    convertedMap[strings.TrimRight(keyVal[0], " ")] = strings.TrimLeft(keyVal[1], " ")
  }
  return convertedMap, nil
}

func convertLabelMap(givenMap map[interface {}]interface {}) (map[string]string, error) {
  convertedMap := make(map[string]string)
  for keyInterface, valueInterface := range givenMap {
    convertedMap[keyInterface.(string)] = valueInterface.(string)
  }
  return convertedMap, nil
}

func getServices(filename string) (map[string]Service, error) {
	bytes, err := ioutil.ReadFile(filename)
  if err != nil {
		return nil, err
	}

  var cfg composeDataRead
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

// = methods dealing with service groups: ==========================================================

func getServiceGrouping(services map[string]Service)(map[string][]string, error) {
  serviceGroups := make(map[string][]string)

  for serviceName, service := range services {
		for label, value := range service.Labels {
      if label != "com.icalialabs.plis.group" { continue }

      groups := strings.Split(value, " ")
      for _, groupNameUntrimmed := range groups {
        groupName := strings.Trim(groupNameUntrimmed, " ")
        serviceGroups[groupName] = append(serviceGroups[groupName], serviceName)
      }
    }
	}
  return serviceGroups, nil
}

// GetServiceGroupingFrom parses a Docker Compose file and figures out the service groups
// by looking for 'com.icalialabs.plis.group' labels.
func GetServiceGroupingFrom(filename string)(map[string][]string, error) {
  services, err := getServices(filename)
  if err != nil { return nil, err }

  serviceGroups, err := getServiceGrouping(services)
  if err != nil { return nil, err }

  return serviceGroups, nil
}
