package main

import (
	"fmt"

	"github.com/S-a-b-r/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
	// TODO: init config: cleanenv
	// TODO: init logger: slog
	// TODO: init storage: sqlight
	// TODO: init router : chi
	// TODO: run server
}
