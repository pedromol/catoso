package main

import (
	"github.com/pedromol/catoso/pkg/catoso"
	"github.com/pedromol/catoso/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	catoso, err := catoso.NewCatoso(cfg)
	if err != nil {
		panic(err)
	}

	catoso.Start()
}
