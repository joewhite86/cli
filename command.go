package cli

import (
	"context"
	"fmt"
	"io"
	"os"
)

var defaultVersion = "0.1.0"
var versionFlag = Flag{Short: "v", Name: "version", Description: "Print the version."}
var Out io.Writer = os.Stdout
var Err io.Writer = os.Stderr

// A Runner function can be defined as Command.Run or Command.DryRun function.
// It will be executed when the command is resolved by the provided arguments.
// The params map will contain a field "_args" with the original argument list.
type Runner func(ctx context.Context, params Params) error

func rootRunner(_ context.Context, params Params) error {
	if params[versionFlag.Name] != nil {
		if root.Version != "" {
			fmt.Fprintln(Out, root.Version)
		} else {
			fmt.Fprintln(Out, defaultVersion)
		}
	}

	return nil
}

// A Command defines a command, or sub-command that can be run by the user.
type Command struct {
	Name     string    // Command name used in help text and the params map
	Group    string    // Group name used to group commands inside help
	Short    string    // Short description, shown in the containing commands help
	Long     string    // Long description, shown in the help text
	Args     []Arg     // Positional arguments
	Flags    []Flag    // Flags for this command
	Commands []Command // Contains a list of sub-commands
	Run      Runner    // The command handler to execute
	Version  string    // Version used in the root command to print the cli's version

	showHelp bool
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) hasFlags() bool {
	return c.Flags != nil && len(c.Flags) > 0
}

func (c *Command) hasArgs() bool {
	return c.Args != nil && len(c.Args) > 0
}

func (c *Command) hasSubCommands() bool {
	return c.Commands != nil && len(c.Commands) > 0
}
