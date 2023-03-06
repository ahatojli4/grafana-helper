package cmd

import (
	"encoding/json"
	"fmt"
	"grafana-helper/internal/entities"
	"os"
	"path"
)

func SetCfg() (*entities.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}
	err = os.MkdirAll(path.Join(homeDir, ".config/ghelper"), 0755)
	if err != nil {
		panic(fmt.Errorf("can't create directory: %s", err))
	}
	file, err := os.OpenFile(path.Join(homeDir, ".config/ghelper/config.json"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(fmt.Errorf("can't open file: %s", err))
	}
	defer func() {
		_ = file.Close()
	}()

	cfg := &entities.Config{}
	cfg.Grafana.Host = askHost()
	cfg.Grafana.Auth.Username = askUsername()
	cfg.Grafana.Auth.Password = askPassword()

	if checkCredentials(cfg.Grafana.Host, cfg.Grafana.Auth.BasicHeader()) {
		data, err := json.Marshal(cfg)
		if err != nil {
			panic(fmt.Errorf("can't marshal data: %s", err))
		}
		_, err = file.Write(data)
		if err != nil {
			panic(fmt.Errorf("can't write data: %s", err))
		}

		return cfg, nil
	}

	return nil, fmt.Errorf("can't set config")
}
