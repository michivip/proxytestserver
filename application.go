package main

import (
	"github.com/michivip/proxytestserver/webserver"
	"log"
	"bufio"
	"os"
	"strings"
	"flag"
	"github.com/michivip/proxytestserver/config"
)

func main() {
	configurationFilePath := flag.String("config", "config.toml", "The path to your custom configuration file.")
	flag.Parse()
	var configLoader config.Loader
	configLoader = &config.TomlLoader{
		Filename: *configurationFilePath,
	}
	configuration, err := configLoader.Load()
	if err != nil {
		if os.IsNotExist(err) {
			if err = configLoader.Save(config.DefaultConfiguration); err != nil {
				panic(err)
			} else {
				log.Fatalln("Created the configuration file for the first time. Please adjust your values and restart the application.")
			}
		} else {
			panic(err)
		}
	}
	log.Printf("Starting proxytestserver on %v...", configuration.Address)
	server := webserver.StartWebserver(configuration)
	log.Println("Started webserver. Enter \"close\" to shutdown the webserver.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.EqualFold(scanner.Text(), "close") {
			break
		}
	}
	log.Println("Shutting down webserver...")
	server.Close()
	log.Println("Thank you for using proxytestserver. Bye!")
}
