package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ahatojli4/ghelper/internal/entities"
)

func LoadConfig() *entities.Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}
	configPath := path.Join(homeDir, ".config/ghelper/config.json")
	fileStat, _ := os.Stat(configPath)
	isConfigFileExist := fileStat != nil
	if isConfigFileExist {
		file, err := os.OpenFile(configPath, os.O_RDWR, 0755)
		if err != nil {
			panic(fmt.Errorf("can't file file: %s", err))
		}
		defer func() {
			_ = file.Close()
		}()

		rw := bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file))
		data, err := io.ReadAll(rw)
		if err != nil {
			panic(fmt.Errorf("can't read data: %s", err))
		}
		cfg := &entities.Config{}
		err = json.Unmarshal(data, &cfg)
		if err != nil {
			panic(fmt.Errorf("can't unmarshal data: %s\n\n%s", err, "try set config again. To do it run this command:\n\n ghelper -set-config"))
		}

		return cfg
	}

	cfg, err := SetCfg()
	if err != nil {
		panic(fmt.Errorf("can't set config: %s", err))
	}

	return cfg
}
