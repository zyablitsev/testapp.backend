package settings

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	instance *config
	once     sync.Once
)

func GetInstance() *config {
	once.Do(func() {
		var (
			goPath                   = os.Getenv("GOPATH")
			etcPath         string   = "/etc/testapp"
			yamlConfigPath  string   = filepath.Join(etcPath, "config.yaml")
			yamlConfigPaths []string = make([]string, 2)
		)

		instance = &config{}

		yamlConfigPaths[0] = yamlConfigPath
		if goPath != "" {
			yamlConfigPaths[1] = filepath.Join(goPath, yamlConfigPath)
		} else {
			yamlConfigPaths = yamlConfigPaths[:1]
		}

		instance.yamlConfigPaths = yamlConfigPaths
		instance.fromFile()
		instance.fromEnvDefault()
	})
	return instance
}
