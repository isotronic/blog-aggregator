package main

import (
	"fmt"

	"github.com/isotronic/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	cfg.SetUser("isotronic")

	newCfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", newCfg)
}