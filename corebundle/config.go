package corebundle

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var Config SystemConfiguration

var configFileName string = "config.json"

func InitConfig() error {
	log.Println("reading configuration file...")
	flag.StringVar(&configFileName, "configFileName", configFileName, "system configurion file name")
	flag.Parse()

	configFileContent, err := os.ReadFile(configFileName)
	if err != nil {
		return fmt.Errorf("error opening configuration file: %s", configFileName)
	}

	config := SystemConfiguration{}
	if err := json.Unmarshal(configFileContent, &config); err != nil {
		return fmt.Errorf("unable to read configuration file: %s\n%s", configFileName, err)
	}

	Config = config

	log.Println("reading configuration file complete")

	return nil
}
