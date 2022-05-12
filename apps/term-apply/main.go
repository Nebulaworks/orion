package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nebulaworks/orion/apps/term-apply/pkg/server"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/version"
)

func main() {

	log.Print(version.BuildInfo())
	config := server.NewConfig()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s, err := server.NewServer(config)
	if err != nil {
		log.Printf("Cannot create server %v", err)
	}
	s.Start()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	s.Stop(ctx)

}
