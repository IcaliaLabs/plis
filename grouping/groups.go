package grouping

import (
  "fmt"
  "io/ioutil"
  "path/filepath"
  "sort"
	"strings"

  // Import Docker CLI Compose stuff... good god they've already sorted out parsing of compose in go
  "github.com/docker/cli/cli/compose/loader"
	composetypes "github.com/docker/cli/cli/compose/types"
)

// Based on https://github.com/docker/cli/blob/79b6d376ce9d9d4ee24f44df25af56df3f981c9a/cli/command/stack/deploy_composefile.go

func getConfigFile(filename string) (*composetypes.ConfigFile, error) {
	bytes, err := ioutil.ReadFile(filename)
  if err != nil {
		return nil, err
	}

	config, err := loader.ParseYAML(bytes)
  if err != nil {
		return nil, err
	}

  return &composetypes.ConfigFile{
		Filename: filename,
		Config:   config,
	}, nil
}

func getConfigDetails(composefile string) (composetypes.ConfigDetails, error) {
	var details composetypes.ConfigDetails

	absPath, err := filepath.Abs(composefile)
	if err != nil {
		return details, err
	}
	details.WorkingDir = filepath.Dir(absPath)

	configFile, err := getConfigFile(composefile)
	if err != nil {
		return details, err
	}

  details.ConfigFiles = []composetypes.ConfigFile{*configFile}
  return details, err
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string
	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n\n")
}

func getServiceGrouping(config *composetypes.Config)(map[string][]string, error) {
  serviceGroups := map[string][]string{}
  for _, service := range config.Services {
		for label, value := range service.Labels {
      if label != "com.icalialabs.plis.group" { continue }

      groups := strings.Split(value, " ")
      for _, group := range groups {
        serviceGroups[group] = append(serviceGroups[group], service.Name)
      }
    }
	}
  return serviceGroups, nil
}

func GetServiceGroupingFrom(filename string)(map[string][]string, error) {
  configDetails, err := getConfigDetails(filename)
  if err != nil { return nil, err }

  config, err := loader.Load(configDetails)
	if err != nil { return nil, err }

  serviceGroups, err := getServiceGrouping(config)
  if err != nil { return nil, err }

  return serviceGroups, nil
}
