package main

import (
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	Formatter := new(prefixed.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.ForceFormatting = true
	Formatter.ForceColors = true
	log.SetFormatter(Formatter)

	viper.SetConfigFile("./config.json")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		os.Exit(3)
	}
	// Confirm which config file is used
	log.Info("Using config: " + viper.ConfigFileUsed())
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed: ", e.Name)
	})
}

func main() {
	log.Println("Starting app")
	app := NewApp()
	app.serve()
}
