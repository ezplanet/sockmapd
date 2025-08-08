package base

import (
	"os"
	"os/signal"
	"sockmapd/logger"
	"syscall"
)

// ProcessSignal upon SIGHUP or SIGUSR1 signals it re-loads the configuration and re-initializes the database
func ProcessSignal() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGUSR1, syscall.SIGHUP)
	go func() {
		for {
			signalType := <-sigch
			logger.Log.Info().Msgf("Received signal: %v - re-initializing...", signalType)
			err := InitializeConfiguration()
			if err != nil {
				logger.Log.Fatal().Msgf("%s %s", err, StrTerminated)
			}
			err = InitializeDatabase()
			if err != nil {
				logger.Log.Fatal().Msgf("%s %s", err, StrTerminated)
			}
		}
	}()
}
