package main

import (
	"context"
	"fmt"
	"os"

	"github.com/renja-g/riotgo/clients"
	"github.com/renja-g/riotgo/middleware"
)

func main() {

	riotClient := clients.NewRiotClient(
		os.Getenv("RIOT_API_KEY"),
		clients.WithMiddleware(middleware.Logging()),
	)

	// Will use default context.Background()
	acc, _ := riotClient.GetAccountV1ByRiotID(
		clients.Europe,
		"Ayato",
		"11235",
	)
	fmt.Printf("%s#%s\n", acc.GameName, acc.TagLine)

	// WithContext custom context
	acc, _ = riotClient.WithContext(context.Background()).GetAccountV1ByRiotID(
		clients.Europe,
		"Ayato",
		"11235",
	)
	fmt.Printf("%s#%s\n", acc.GameName, acc.TagLine)
}
