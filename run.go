package cli

import (
	"context"
	"fmt"
	"os"
)

var root *Command

// Map of parsed command line parameters.
type Params map[string]interface{}

func (p Params) AdditionalArguments() []string {
	return p["_args"].([]string)
}

func Run(ctx context.Context, cmd *Command) error {
	root = cmd
	cmd.Flags = append(cmd.Flags, versionFlag)
	if cmd.Run == nil {
		cmd.Run = rootRunner
	}

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "lint" {
		errs := lint(cmd)
		for _, err := range errs {
			fmt.Fprintln(Err, "[WARN] "+err.Error())
		}
		return nil
	}

	params := Params{}
	cmd, err := resolve(cmd, args, params)
	if err != nil {
		return err
	}

	if len(args) == 0 || cmd.showHelp || !cmd.Runnable() {
		cmd.printHelp()
		return nil
	}

	return cmd.Run(ctx, params)
}

func resolve(cmd *Command, args []string, params Params) (*Command, error) {
	if len(args) == 1 && isHelp(args[0]) {
		cmd.showHelp = true
		return cmd, nil
	}
	p := parseArgs(cmd, &args)
	if err := checkRequiredParams(cmd, p); err != nil {
		return nil, err
	}
	appendParams(params, p)

	if len(args) > 0 && cmd.hasSubCommands() {
		for _, res := range cmd.Commands {
			if args[0] != res.Name {
				continue
			}
			clone := res
			return resolve(&clone, args[1:], params)
		}
	}
	// if len(args) > 0 {
	// 	return nil, fmt.Errorf("invalid arguments: %s", args)
	// }

	return cmd, nil
}

func appendParams(params Params, append Params) {
	for key, val := range append {
		params[key] = val
	}
}
