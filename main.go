package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

func main() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	logFileName := fmt.Sprintf("%s/%s", dirname, "tv-source-switch.log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("%s/.config/tvsourceswitch", dirname))
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	source := viper.GetString("source")
	token := viper.GetString("smartthings_token")
	deviceID := viper.GetString("smartthings_device_id")

	client := NewSmartThingsTVClient(token, deviceID)
	defer client.Close()
	status, err := client.GetStatus()

	if err != nil {
		log.Fatal(err)
	}

	if status.State != "ONLINE" {
		log.Fatal("Device is not online")
	}

	result, err := client.SetPower("on")

	if err != nil {
		log.Fatal(err)
	}

	if result.Status != "ACCEPTED" {
		log.Fatal("Error switching on power")
	}

	result, err = client.SetSource(source)

	if err != nil {
		log.Fatal(err)
	}

	if result.Status != "ACCEPTED" {
		log.Fatal("Error switching source")
	}

}
