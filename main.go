package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	Formatter := new(prefixed.TextFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.ForceFormatting = true
	Formatter.ForceColors = true
	log.SetFormatter(Formatter)

	errenv := godotenv.Load()
	if errenv != nil {
		log.Println("Error loading .env file")
	}

}

func main() {
	log.Println("Starting app")
	app := NewApp()
	app.serve()
}
