package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	rand.Seed(time.Now().UnixNano())

	fmt.Println("It works")
}