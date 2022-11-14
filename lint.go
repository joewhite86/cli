package cli

import "fmt"

func lint(cmd *Command) []error {
	errs := make([]error, 0)

	if cmd.Name == "" {
		errs = append(errs, fmt.Errorf("missing name on command %+v", cmd))
	}
	if root != cmd && cmd.Short == "" {
		errs = append(errs, fmt.Errorf("missing short description on command %s", cmd.Name))
	}
	if cmd.hasFlags() {
		for _, param := range cmd.Flags {
			if param.Name == "" {
				errs = append(errs, fmt.Errorf("missing name param %s in command %s", param.Name, cmd.Name))
			}
			if param.Description == "" {
				errs = append(errs, fmt.Errorf("missing description on param %s in command %s", param.Name, cmd.Name))
			}
		}
	}
	if cmd.hasSubCommands() {
		for _, sub := range cmd.Commands {
			s := sub
			errs = append(errs, lint(&s)...)
		}
	}

	return errs
}
