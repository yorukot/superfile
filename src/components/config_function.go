package components

import (
	"log"
	"os"
	"path/filepath"

	"github.com/barasher/go-exiftool"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

func initialConfig(dir string) (toggleDotFileBool bool, firstFilePanelDir string) {
	var err error

	logOutput, err = os.OpenFile(SuperFileCacheDir+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening superfile.log file: %v", err)
	}

	loadConfigFile()

	loadHotkeysFile()

	data, err := os.ReadFile(SuperFileMainDir + themeFolder + "/" + Config.Theme + ".toml")
	if err != nil {
		log.Fatalf("Theme file doesn't exist: %v", err)
	}

	err = toml.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error while decoding theme json( Your theme file may have errors ): %v", err)
	}
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
	viperConfig := viper.New()
	viperConfig.SetConfigName("config")
	viperConfig.SetConfigType("toml")
	viperConfig.SetConfigFile(SuperFileMainDir + configFile)
	err := viperConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}
	err = viperConfig.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("Error decoding config file( your config file may have misconfigured ): %v", err)
	}
}

func loadHotkeysFile() {
	
	viperHotkeys := viper.New()
	viperHotkeys.SetConfigName("hotkeys")
	viperHotkeys.SetConfigType("toml")
	viperHotkeys.SetConfigFile(SuperFileMainDir + hotkeysFile)

	err := viperHotkeys.ReadInConfig()
	if err != nil {
		log.Fatalf("Config file doesn't exist: %v", err)
	}
	err = viperHotkeys.Unmarshal(&hotkeys)
	if err != nil {
		log.Fatalf("Error decoding hotkeys file( your config file may have misconfigured ): %v", err)
	}
}
