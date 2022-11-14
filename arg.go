package cli

import "fmt"

// Arg is an positional argument passed to a command.
type Arg struct {
	Name        string      // Name used in help texts and as key in the params map
	Description string      // Description shown in hep texts
	Parser      ParserFunc  // Parser function to use (Default: StringParser)
	Required    bool        // If true, the execution will fail if the argument is not passed
	Default     interface{} // Default value
	Vararg      bool        // If true, the argument has an undefined length and is of type []string
}

func (a *Arg) parse(val string) interface{} {
	parser := StringParser
	if a.Parser != nil {
		parser = a.Parser
	}

	parsed, err := parser(val)
	if err != nil {
		fmt.Fprintln(Out, err.Error())
	}

	return parsed
}
