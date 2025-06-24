```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/renja-g/riotgo/clients"
)

func main() {

	riotClient := clients.NewRiotClient("XXX")

	acc, err := riotClient.GetAccountV1ByRiotID(
		context.Background(),
		clients.Europe,
		"Ayato",
		"11235",
	)
	if err != nil {
		log.Fatalf("Error getting account: %v", err)
	}

	fmt.Println(acc.GameName)
	fmt.Println(acc.TagLine)
	fmt.Println(acc.PUUID)
}

```