package grouping

import (
  "fmt"
  "io/ioutil"
  "path/filepath"
  "reflect"
  "sort"
	"strings"

  // Import Docker CLI Compose stuff... good god they've already sorted out parsing of compose in go
  "github.com/docker/cli/cli/compose/schema"
  "github.com/docker/cli/cli/compose/template"
  "github.com/docker/cli/cli/compose/types"

  "github.com/mitchellh/mapstructure"
  "github.com/pkg/errors"
  yaml "gopkg.in/yaml.v2"
)

// Based on https://github.com/docker/cli/blob/79b6d376ce9d9d4ee24f44df25af56df3f981c9a/cli/command/stack/deploy_composefile.go

// = methods blatantly copied from github.com/docker/cli/cli/compose/loader ========================

func parseYAML(source []byte) (map[string]interface{}, error) {
	var cfg interface{}
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return nil, err
	}
	cfgMap, ok := cfg.(map[interface{}]interface{})
	if !ok {
		return nil, errors.Errorf("Top-level object must be a mapping")
	}
	converted, err := convertToStringKeysRecursive(cfgMap, "")
	if err != nil {
		return nil, err
	}
	return converted.(map[string]interface{}), nil
}

func getConfigDict(configDetails types.ConfigDetails) map[string]interface{} {
	return configDetails.ConfigFiles[0].Config
}

func getServices(configDict map[string]interface{}) map[string]interface{} {
	if services, ok := configDict["services"]; ok {
		if servicesDict, ok := services.(map[string]interface{}); ok {
			return servicesDict
		}
	}

	return map[string]interface{}{}
}

// keys needs to be converted to strings for jsonschema
func convertToStringKeysRecursive(value interface{}, keyPrefix string) (interface{}, error) {
	if mapping, ok := value.(map[interface{}]interface{}); ok {
		dict := make(map[string]interface{})
		for key, entry := range mapping {
			str, ok := key.(string)
			if !ok {
				return nil, formatInvalidKeyError(keyPrefix, key)
			}
			var newKeyPrefix string
			if keyPrefix == "" {
				newKeyPrefix = str
			} else {
				newKeyPrefix = fmt.Sprintf("%s.%s", keyPrefix, str)
			}
			convertedEntry, err := convertToStringKeysRecursive(entry, newKeyPrefix)
			if err != nil {
				return nil, err
			}
			dict[str] = convertedEntry
		}
		return dict, nil
	}
	if list, ok := value.([]interface{}); ok {
		var convertedList []interface{}
		for index, entry := range list {
			newKeyPrefix := fmt.Sprintf("%s[%d]", keyPrefix, index)
			convertedEntry, err := convertToStringKeysRecursive(entry, newKeyPrefix)
			if err != nil {
				return nil, err
			}
			convertedList = append(convertedList, convertedEntry)
		}
		return convertedList, nil
	}
	return value, nil
}

// ForbiddenPropertiesError is returned when there are properties in the Compose
// file that are forbidden.
type ForbiddenPropertiesError struct {
	Properties map[string]string
}

func (e *ForbiddenPropertiesError) Error() string {
	return "Configuration contains forbidden properties"
}

func getConfig(configDict map[string]interface{}, lookupEnv template.Mapping) (map[string]map[string]interface{}, error) {
	config := make(map[string]map[string]interface{})

	for _, key := range []string{"services", "networks", "volumes", "secrets", "configs"} {
    section, ok := configDict[key]
		if !ok {
			config[key] = make(map[string]interface{})
			continue
		}
		config[key] = section.(map[string]interface{})
	}
	return config, nil
}

func getProperties(services map[string]interface{}, propertyMap map[string]string) map[string]string {
	output := map[string]string{}

	for _, service := range services {
		if serviceDict, ok := service.(map[string]interface{}); ok {
			for property, description := range propertyMap {
				if _, isSet := serviceDict[property]; isSet {
					output[property] = description
				}
			}
		}
	}

	return output
}

// Load reads a ConfigDetails and returns a fully loaded configuration
func load(configDetails types.ConfigDetails) (*types.Config, error) {
	if len(configDetails.ConfigFiles) < 1 {
		return nil, errors.Errorf("No files specified")
	}
	if len(configDetails.ConfigFiles) > 1 {
		return nil, errors.Errorf("Multiple files are not yet supported")
	}

	configDict := getConfigDict(configDetails)

	if services, ok := configDict["services"]; ok {
		if servicesDict, ok := services.(map[string]interface{}); ok {
			forbidden := getProperties(servicesDict, types.ForbiddenProperties)

			if len(forbidden) > 0 {
				return nil, &ForbiddenPropertiesError{Properties: forbidden}
			}
		}
	}

	if err := schema.Validate(configDict, schema.Version(configDict)); err != nil {
		return nil, err
	}

	cfg := types.Config{}

	config, err := getConfig(configDict, configDetails.LookupEnv)
	if err != nil {
		return nil, err
	}

	cfg.Services, err = loadServices(config["services"], configDetails.WorkingDir, configDetails.LookupEnv)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string
	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n\n")
}

// loadServices produces a ServiceConfig map from a compose file Dict
// the servicesDict is not validated if directly used. Use Load() to enable validation
func loadServices(servicesDict map[string]interface{}, workingDir string, lookupEnv template.Mapping) ([]types.ServiceConfig, error) {
	var services []types.ServiceConfig

	for name, serviceDef := range servicesDict {
		serviceConfig, err := loadService(name, serviceDef.(map[string]interface{}), workingDir, lookupEnv)
		if err != nil {
			return nil, err
		}
		services = append(services, *serviceConfig)
	}

	return services, nil
}

// loadService produces a single ServiceConfig from a compose file Dict
// the serviceDict is not validated if directly used. Use Load() to enable validation
func loadService(name string, serviceDict map[string]interface{}, workingDir string, lookupEnv template.Mapping) (*types.ServiceConfig, error) {
	serviceConfig := &types.ServiceConfig{}
	if err := transform(serviceDict, serviceConfig); err != nil {
		return nil, err
	}
	serviceConfig.Name = name
	return serviceConfig, nil
}

func formatInvalidKeyError(keyPrefix string, key interface{}) error {
	var location string
	if keyPrefix == "" {
		location = "at top level"
	} else {
		location = fmt.Sprintf("in %s", keyPrefix)
	}
	return errors.Errorf("Non-string key %s: %#v", location, key)
}

func transform(source map[string]interface{}, target interface{}) error {
	data := mapstructure.Metadata{}
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			createTransformHook(),
			mapstructure.StringToTimeDurationHookFunc()),
		Result:   target,
		Metadata: &data,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(source)
}

func createTransformHook() mapstructure.DecodeHookFuncType {
	transforms := map[reflect.Type]func(interface{}) (interface{}, error){
		reflect.TypeOf(types.StringList{}):                       transformStringList,
		reflect.TypeOf(map[string]string{}):                      transformMapStringString,
		reflect.TypeOf(types.StringOrNumberList{}):               transformStringOrNumberList,
		reflect.TypeOf(types.MappingWithEquals{}):                transformMappingOrListFunc("=", true),
		reflect.TypeOf(types.Labels{}):                           transformMappingOrListFunc("=", false),
		reflect.TypeOf(types.MappingWithColon{}):                 transformMappingOrListFunc(":", false),
	}

	return func(_ reflect.Type, target reflect.Type, data interface{}) (interface{}, error) {
		transform, ok := transforms[target]
		if !ok {
			return data, nil
		}
		return transform(data)
	}
}

func transformStringList(data interface{}) (interface{}, error) {
	switch value := data.(type) {
	case string:
		return []string{value}, nil
	case []interface{}:
		return value, nil
	default:
		return data, errors.Errorf("invalid type %T for string list", value)
	}
}

func transformMapStringString(data interface{}) (interface{}, error) {
	switch value := data.(type) {
	case map[string]interface{}:
		return toMapStringString(value, false), nil
	case map[string]string:
		return value, nil
	default:
		return data, errors.Errorf("invalid type %T for map[string]string", value)
	}
}

func transformStringOrNumberList(value interface{}) (interface{}, error) {
	list := value.([]interface{})
	result := make([]string, len(list))
	for i, item := range list {
		result[i] = fmt.Sprint(item)
	}
	return result, nil
}

func transformMappingOrListFunc(sep string, allowNil bool) func(interface{}) (interface{}, error) {
	return func(data interface{}) (interface{}, error) {
		return transformMappingOrList(data, sep, allowNil), nil
	}
}

func toMapStringString(value map[string]interface{}, allowNil bool) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range value {
		output[key] = toString(value, allowNil)
	}
	return output
}

func transformMappingOrList(mappingOrList interface{}, sep string, allowNil bool) interface{} {
	switch value := mappingOrList.(type) {
	case map[string]interface{}:
		return toMapStringString(value, allowNil)
	case ([]interface{}):
		result := make(map[string]interface{})
		for _, value := range value {
			parts := strings.SplitN(value.(string), sep, 2)
			key := parts[0]
			switch {
			case len(parts) == 1 && allowNil:
				result[key] = nil
			case len(parts) == 1 && !allowNil:
				result[key] = ""
			default:
				result[key] = parts[1]
			}
		}
		return result
	}
	panic(errors.Errorf("expected a map or a list, got %T: %#v", mappingOrList, mappingOrList))
}

func toString(value interface{}, allowNil bool) interface{} {
	switch {
	case value != nil:
		return fmt.Sprint(value)
	case allowNil:
		return nil
	default:
		return ""
	}
}

// = methods blatantly copied from github.com/docker/cli/cli/command/stack =========================

func getConfigFile(filename string) (*types.ConfigFile, error) {
	bytes, err := ioutil.ReadFile(filename)
  if err != nil {
		return nil, err
	}

	config, err := parseYAML(bytes)
  if err != nil {
		return nil, err
	}

  return &types.ConfigFile{
		Filename: filename,
		Config:   config,
	}, nil
}

func getConfigDetails(composefile string) (types.ConfigDetails, error) {
	var details types.ConfigDetails

	absPath, err := filepath.Abs(composefile)
	if err != nil {
		return details, err
	}
	details.WorkingDir = filepath.Dir(absPath)

	configFile, err := getConfigFile(composefile)
	if err != nil {
		return details, err
	}

  details.ConfigFiles = []types.ConfigFile{*configFile}
  return details, err
}

// = methods dealing with service groups: ==========================================================

func getServiceGrouping(config *types.Config)(map[string][]string, error) {
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

// GetServiceGroupingFrom parses a Docker Compose file and figures out the service groups
// by looking for 'com.icalialabs.plis.group' labels.
func GetServiceGroupingFrom(filename string)(map[string][]string, error) {
  configDetails, err := getConfigDetails(filename)
  if err != nil { return nil, err }

  config, err := load(configDetails)
	if err != nil { return nil, err }

  serviceGroups, err := getServiceGrouping(config)
  if err != nil { return nil, err }

  return serviceGroups, nil
}
