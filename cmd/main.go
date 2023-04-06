package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mttchrry/oxio-phone-lookup/pkg/app"
	httppkg "github.com/mttchrry/oxio-phone-lookup/pkg/http"
	"github.com/mttchrry/oxio-phone-lookup/pkg/keyDataStore"
)

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	inputFile := os.Args[1]

	k, err := keyDataStore.New(inputFile)
	if err != nil {
		fmt.Printf("error starting server: ", err)
		return nil, err
	}

	h, err := httppkg.New(k, "8081")
	if err != nil {
		return nil, err
	}

	// Start listening for HTTP requests
	return []app.Listener{
		h,
	}, nil
}
