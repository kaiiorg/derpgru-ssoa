package main

import (
	"github.com/kaiiorg/derpgru-ssoa/internal/matcher"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := matcher.New().CobraCommand().Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
