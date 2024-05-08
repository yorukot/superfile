package internal

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/barasher/go-exiftool"
	"github.com/pelletier/go-toml/v2"
)

func initialConfig(dir string) (toggleDotFileBool bool, firstFilePanelDir string) {
	var err error

	logOutput, err = os.OpenFile(SuperFileStateDir+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening superfile.log file: %v", err)
	}

	loadConfigFile()

	loadHotkeysFile()

	loadThemeFile()

	toggleDotFileData, err := os.ReadFile(SuperFileDataDir + toggleDotFile)
	if err != nil {
		outPutLog("Error while reading toggleDotFile data error:", err)
	}
	if string(toggleDotFileData) == "true" {
		toggleDotFileBool = true
	} else if string(toggleDotFileData) == "false" {
		toggleDotFileBool = false
	}
	LoadThemeConfig()

	if Config.Metadata {
		et, err = exiftool.NewExiftool()
		if err != nil {
			outPutLog("Initial model function init exiftool error", err)
		}
	}

	firstFilePanelDir = HomeDir
	if dir != "" {
		firstFilePanelDir, err = filepath.Abs(dir)
		if err != nil {
			firstFilePanelDir = HomeDir
		}
	}
	return toggleDotFileBool, firstFilePanelDir
}

func loadConfigFile() {

	_ = toml.Unmarshal([]byte(ConfigTomlString), &Config)
	tempForCheckMissingConfig := ConfigType{}

	data, err := os.ReadFile(SuperFileMainDir + configFile)
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}

	_ = toml.Unmarshal(data, &tempForCheckMissingConfig)
	err = toml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error decoding config file ( your config file may have misconfigured ): %v", err)
	}

	if !reflect.DeepEqual(Config, tempForCheckMissingConfig) {
		tomlData, err := toml.Marshal(Config)
		if err != nil {
			log.Fatalf("Error encoding config: %v", err)
		}

		err = os.WriteFile(SuperFileMainDir+configFile, tomlData, 0644)
		if err != nil {
			log.Fatalf("Error writing config file: %v", err)
		}
	}
}

func loadHotkeysFile() {

	_ = toml.Unmarshal([]byte(HotkeysTomlString), &hotkeys)
	tempForCheckMissingConfig := HotkeysType{}
	data, err := os.ReadFile(SuperFileMainDir + hotkeysFile)

	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}

	_ = toml.Unmarshal(data, &tempForCheckMissingConfig)
	err = toml.Unmarshal(data, &hotkeys)
	if err != nil {
		log.Fatalf("Error decoding hotkeys file ( your config file may have misconfigured ): %v", err)
	}

	if !reflect.DeepEqual(hotkeys, tempForCheckMissingConfig) {
		tomlData, err := toml.Marshal(hotkeys)
		if err != nil {
			log.Fatalf("Error encoding hotkeys: %v", err)
		}

		err = os.WriteFile(SuperFileMainDir+hotkeysFile, tomlData, 0644)
		if err != nil {
			log.Fatalf("Error writing hotkeys file: %v", err)
		}
	}

}

func loadThemeFile() {
	data, err := os.ReadFile(SuperFileMainDir + themeFolder + "/" + Config.Theme + ".toml")
	if err != nil {
		log.Fatalf("Theme file doesn't exist: %v", err)
	}

	err = toml.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error while decoding theme file( Your theme file may have errors ): %v", err)
	}
}
