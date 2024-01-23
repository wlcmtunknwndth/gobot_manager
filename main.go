package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/wlcmtunknwndth/gobot_manager/clients/telegram"
	event_consumer "github.com/wlcmtunknwndth/gobot_manager/consumer/event-consumer"
	"github.com/wlcmtunknwndth/gobot_manager/events/telegram"
	"github.com/wlcmtunknwndth/gobot_manager/storage/sqlite"
)

const (
	batchSize         = 100
	sqliteStoragePath = "data/sqlite/storage.db"
	host              = "api.telegram.org"
)

func main() {
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't coonect to storage: %w", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: %w", err)
	}

	eventProcessor := telegram.New(
		tgClient.New(host, mustToken()),
		s,
		0,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(context.TODO()); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {

	token := flag.String(
		"token",
		"",
		"token for access to tg-bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is not specified")
	}
	return *token
}

// func hostFlag() *string {
// 	host := flag.String(
// 		"host",
// 		"",
// 		"host for tg-bot",
// 	)
// 	// flag.Parse()
// 	// if *host == "" {
// 	// 	log.Fatal("Host is not specified")
// 	// }
// 	return host
// }

// func initialCondition() (string, string) {
// 	host := hostFlag()
// 	token := tokenFlag()

// 	flag.Parse()
// 	if *host == "" || *token == "" {
// 		log.Fatal("Token or Host is not specified")
// 	}

// 	return *host, *token
// }

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
