package main

import (
	"github.com/sharkbait0402/blog-aggregator/internal/config"
	"fmt"
)

func main() {

	cfg, err := config.Read()
	if err!=nil {
		fmt.Errorf("read unsuccessful")
	}

	cfg.SetUser("Casen")

	cfg, err = config.Read()
		if err!=nil {
			fmt.Errorf("read unsuccessful")
		}

	fmt.Printf("%+v\n", cfg)

}
