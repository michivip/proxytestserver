package main

import (
	"github.com/michivip/proxytestserver/webserver"
	"log"
	"bufio"
	"os"
	"strings"
)

const (
	Address = "localhost:8070"
)

func main() {
	log.Printf("Starting proxytestserver on %v...", Address)
	server := webserver.StartWebserver(Address)
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
