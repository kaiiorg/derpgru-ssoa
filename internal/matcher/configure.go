package matcher

import (
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func (matcher *Matcher) configure(cmd *cobra.Command, args []string) error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	matcher.configureLogLevel()
	rand.Seed(time.Now().UnixNano())

	log.Trace().Str("cmd", cmd.Use).Msg("configure executed")
	return nil
}

func (matcher *Matcher) configureLogLevel() {
	logLevel, err := zerolog.ParseLevel(matcher.logLevel)
	if err != nil || logLevel == zerolog.NoLevel {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}
