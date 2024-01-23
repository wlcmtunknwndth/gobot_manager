package main

import (
	"flag"
	"log"

	tgClient "github.com/wlcmtunknwndth/gobot_manager/clients/telegram"
	event_consumer "github.com/wlcmtunknwndth/gobot_manager/consumer/event-consumer"
	"github.com/wlcmtunknwndth/gobot_manager/events/telegram"
	"github.com/wlcmtunknwndth/gobot_manager/storage/files"
)

const (
	batchSize   = 100
	storagePath = "storage"
)

func main() {

	eventProcessor := telegram.New(
		tgClient.New(initialCondition()),
		files.New(storagePath),
		0,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() *string {

	token := flag.String(
		"token",
		"",
		"token for access to tg-bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}
	return token
}

func mustHost() *string {
	host := flag.String(
		"host",
		"",
		"host for tg-bot",
	)
	// flag.Parse()
	// if *host == "" {
	// 	log.Fatal("Host is not specified")
	// }
	return host
}

func initialCondition() (string, string) {
	host := mustHost()
	token := mustToken()

	flag.Parse()
	if *host == "" || *token == "" {
		log.Fatal("Token or Host is not specified")
	}

	return *host, *token
}

// func mustStoragePath() string {
// 	path := flag.String(
// 		"storage-path",
// 		"",
// 		"path to Storage",
// 	)

// 	flag.Parse()
// 	if *path == "" {
// 		log.Fatal("Storage path is not spicified")
// 	}
// 	return *path
// }
