package cli

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

var defaultHelpFlags = []string{"-h", "--help", "help"}

//go:embed help.tpl
var helpTemplate string

// Print the help text for a command.
// This includes the description (if any), the arguments, the parameters and
// available sub-commands.
func (c *Command) printHelp() {
	funcMap := template.FuncMap{
		"usage":            c.Usage,
		"formatArg":        templateFormatArg(c),
		"formatSubCommand": templateFormatSubCommand(c),
		"groups":           templateCommandGroups(c),
	}

	tmpl, err := template.New("help").Funcs(funcMap).Parse(helpTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(Out, c)
	if err != nil {
		panic(err)
	}
}

func (c *Command) Usage() string {
	buf := strings.Builder{}
	fmt.Fprintf(&buf, "%s ", c.Name)

	if c.hasFlags() {
		buf.WriteString("[flags] ")
	}

	if c.hasSubCommands() {
		buf.WriteString("<command> ")
	}

	if c.hasArgs() {
		str := make([]string, 0, len(c.Args))
		for _, flag := range c.Args {
			str = append(str, "<"+flag.Name+">")
		}
		fmt.Fprint(&buf, strings.Join(str, " "))
	}

	return buf.String()
}

func templateFormatSubCommand(c *Command) func(arg string) string {
	nameWidth := subCommandMaxLen(c)
	space := 4
	spaceString := strconv.Itoa(nameWidth + space)
	return func(name string) string {
		return fmt.Sprintf("%-"+spaceString+"s", name)
	}
}

func templateFormatArg(c *Command) func(arg string) string {
	argWidth := argMaxLen(c)
	space := 6
	spaceString := strconv.Itoa(argWidth + space)
	return func(arg string) string {
		return fmt.Sprintf("%-"+spaceString+"s", "<"+arg+">:")
	}
}

func templateCommandGroups(c *Command) func() map[string][]Command {
	groups := make(map[string][]Command)

	for _, cmd := range c.Commands {
		grp := cmd.Group
		if groups[grp] == nil {
			groups[grp] = make([]Command, 1)
		}
		groups[grp] = append(groups[grp], cmd)
	}
	return func() map[string][]Command {
		return groups
	}
}

func isHelp(arg string) bool {
	for _, flag := range defaultHelpFlags {
		if arg == flag {
			return true
		}
	}
	return false
}

func subCommandMaxLen(c *Command) int {
	nameWidth := 3
	for _, cmd := range c.Commands {
		if len(cmd.Name) > nameWidth {
			nameWidth = len(cmd.Name)
		}
	}
	return nameWidth
}

func argMaxLen(c *Command) int {
	nameWidth := 3
	for _, cmd := range c.Args {
		if len(cmd.Name) > nameWidth {
			nameWidth = len(cmd.Name)
		}
	}
	return nameWidth
}
