package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

type config struct {
	Port    uint
	LogPath string

	yamlConfigPaths []string
}

func (s *config) fromFile() {
	var (
		data []byte
		err  error
		ok   bool
	)

	for _, v := range s.yamlConfigPaths {
		if _, err = os.Stat(v); err != nil {
			continue
		}

		if data, err = ioutil.ReadFile(v); err != nil {
			continue
		}

		if err = yaml.Unmarshal(data, s); err != nil {
			continue
		}
		ok = true
		break
	}

	if !ok {
		log.Fatal("app config not found")
	}

	return
}

func (s *config) fromEnvDefault() {
	// app port
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		if v, err := strconv.ParseUint(envPort, 10, 64); err == nil {
			s.Port = uint(v)
		}
	}
	if s.Port == 0 {
		s.Port = 8701
	}

	// log path
	envLogPath := os.Getenv("LOG_PATH")
	if len(envLogPath) > 0 {
		s.LogPath = envLogPath
	}
	if len(s.LogPath) == 0 {
		s.LogPath = "/var/log/testapp/requests.log"
	}
}

var (
	instanceConfig *config
	onceConfig     sync.Once
)

func getConfig() *config {
	onceConfig.Do(func() {
		var (
			goPath                   = os.Getenv("GOPATH")
			etcPath         string   = "/etc/testapp"
			yamlConfigPath  string   = filepath.Join(etcPath, "config.yaml")
			yamlConfigPaths []string = make([]string, 2)
		)

		instanceConfig = &config{}

		yamlConfigPaths[0] = yamlConfigPath
		if goPath != "" {
			yamlConfigPaths[1] = filepath.Join(goPath, yamlConfigPath)
		} else {
			yamlConfigPaths = yamlConfigPaths[:1]
		}

		instanceConfig.yamlConfigPaths = yamlConfigPaths
		instanceConfig.fromFile()
		instanceConfig.fromEnvDefault()
	})
	return instanceConfig
}
