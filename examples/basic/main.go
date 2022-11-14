package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joewhite86/cli"
)

// nolint:gomnd
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var login = cli.Command{
		Name: "print",
		Flags: []cli.Flag{
			{Name: "user", Short: "u", Required: true, Description: "User to run with.", HasValue: true},
		},
		Run: func(ctx context.Context, params cli.Params) error {
			fmt.Printf("The passed user was: %s.\n", params["user"].(string))
			return nil
		},
	}

	var c = cli.Command{
		Name:     "example-cli",
		Long:     "This is an example",
		Commands: []cli.Command{login},
	}

	if err := cli.Run(ctx, &c); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
