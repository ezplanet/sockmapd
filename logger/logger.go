package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(level string, output string) error {
	var caller bool = false
	// Set the log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch level {
	case "trace":
		caller = true
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		caller = true
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}
	// Set the log output
	var err error

	// to optimise disable console writer on prod JSON could be much easier to parse
	if output == "stdout" {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: time.RFC3339}
		Log = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
	} else {
		var logfile *os.File
		if output == "stderr" {
			logfile = os.Stderr
		} else {
			logfile, err = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		fileWriter := zerolog.ConsoleWriter{Out: logfile, NoColor: true, TimeFormat: time.RFC3339}
		if caller {
			Log = zerolog.New(fileWriter).With().Timestamp().Caller().Logger()
		} else {
			Log = zerolog.New(fileWriter).With().Timestamp().Logger()
		}
	}
	return err
}
