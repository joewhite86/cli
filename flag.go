package cli

import (
	"fmt"
)

// Flag that can be passed to commands. The name and short description should always
// be set.
type Flag struct {
	Short       string      // Name used as short parameter ("-") name
	Name        string      // Name used as long parameter ("--") name and as key in the params map
	HasValue    bool        // If true, the following argument will be taken as parameter value
	Description string      // Description text shown in help texts
	Parser      ParserFunc  // Parser function to use
	Required    bool        // If true, the execution will fail, if this flag is not set
	Default     interface{} // Default value
}

func (p *Flag) matches(arg string) bool {
	return arg == "-"+p.Short || arg == "--"+p.Name
}

func (p *Flag) parse(next string) interface{} {
	if p.HasValue {
		parser := StringParser
		if p.Parser != nil {
			parser = p.Parser
		}

		val, err := parser(next)
		if err != nil {
			fmt.Fprintln(Err, err.Error())
		}

		return val
	} else {
		return true
	}
}
