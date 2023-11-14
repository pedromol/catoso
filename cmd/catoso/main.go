package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pedromol/catoso/pkg/catoso"
	"github.com/pedromol/catoso/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	conf, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		panic(err)
	}

	log.Printf("starting with config: \n%s\n", conf)

	catoso, err := catoso.NewCatoso(cfg)
	if err != nil {
		panic(err)
	}

	var exitAfter int64
	if cfg.ExitAfterMin != "" {
		exitAfter, _ = strconv.ParseInt(cfg.ExitAfterMin, 10, 64)
	}

	if exitAfter == 0 {
		catoso.Start()
	} else {
		go catoso.Start()
		time.Sleep(time.Duration(exitAfter) * time.Minute)
		log.Printf("exiting after %d minutes", exitAfter)
		os.Exit(0)
	}
}
