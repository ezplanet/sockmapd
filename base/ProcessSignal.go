package base

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// ProcessSignal: upon SIGHUP or SIGUSR1 signals it re-loads the configuration and re-initializes the database
func ProcessSignal() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGUSR1, syscall.SIGHUP)
	go func() {
		for {
			signalType := <-sigch
			log.Printf("Received signal: %v - re-initializing...", signalType)
			err := InitializeConfiguration()
			if err != nil {
				log.Fatalln(err, StrTerminated)
			}
			err = InitializeDatabase()
			if err != nil {
				log.Fatalln(err, StrTerminated)
			}
		}
	}()
}
