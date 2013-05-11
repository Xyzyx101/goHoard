package config

import (
	"code.google.com/p/goconf/conf"
)

type Config map[string]map[string]string

func ParseFile(file string) (Config, error) {
	var config Config

	f, err := conf.ReadConfigFile(file)
	if err != nil {
		return config, err
	}

	sections := f.GetSections()
	config = make(map[string]map[string]string, len(sections))

	for _, section := range sections {
		options, err := f.GetOptions(section)
		if err != nil {
			return config, err
		}
		config[section] = make(map[string]string, len(options))
		for _, option := range options {
			config[section][option], err = f.GetString(section, option)
			if err != nil {
				return config, err
			}
		}
	}
	return config, err
}
