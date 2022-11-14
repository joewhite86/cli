package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joewhite86/cli"
)

var Ls = cli.Command{
	Name:  "ls",
	Short: "Execute ls.",
	Run:   runLs,
}

func runLs(ctx context.Context, _ cli.Params) error {
	cmd := exec.CommandContext(ctx, "/bin/ls") // nolint:gosec
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

var Login = cli.Command{
	Name:  "login",
	Short: "Login to something.",
	Flags: []cli.Flag{
		{Short: "u", Name: "user", HasValue: true, Description: "User name"},
		{Short: "p", Name: "pass", HasValue: true, Description: "Password"}},
}

// nolint:gomnd
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	commands := []cli.Command{
		Ls,
		Login,
	}

	var c = cli.Command{
		Name:     "example-cli",
		Long:     "This is an example.",
		Short:    "This is an example.",
		Commands: commands,
	}

	if err := cli.Run(ctx, &c); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
