package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// ParserFunc defines a specific parameter parsing.
type ParserFunc func(val string) (interface{}, error)

// StringParser just returns the given val.
func StringParser(val string) (interface{}, error) {
	return val, nil
}

// Int32Parser returns the int32 representation of the given val.
func Int32Parser(val string) (interface{}, error) {
	i64, err := strconv.ParseInt(val, 10, 32) // nolint:gomnd
	if err != nil {
		return int32(-1), err
	}
	return int32(i64), nil
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--")
}

func parseArgs(cmd *Command, args *[]string) Params {
	params := make(Params)
	params["_args"] = *args

	skip := 0
	notFound := []string{}
	for i, arg := range *args {
		if skip > 0 {
			skip--
			continue
		}
		var found bool
		if !isFlag(arg) {
			skip, found = parseArguments(cmd, params, i, args)
			if !found {
				notFound = append(notFound, arg)
			}
			continue
		}
		skip, found = parseFlags(cmd, params, i, args)
		if !found {
			notFound = append(notFound, arg)
		}
	}
	*args = notFound

	return params
}

func parseFlags(cmd *Command, params Params, index int, args *[]string) (skip int, found bool) {
	for _, param := range cmd.Flags {
		if params[param.Name] != nil || !param.matches((*args)[index]) {
			continue
		}
		next := ""
		if param.HasValue && len(*args) > index+1 {
			next = (*args)[index+1]
		}
		params[param.Name] = param.parse(next)
		found = true
		if param.HasValue {
			skip++
		}
		break
	}
	return
}

func parseArguments(cmd *Command, params Params, index int, args *[]string) (skip int, found bool) {
	for _, param := range cmd.Args {
		if params[param.Name] != nil {
			continue
		}
		if param.Vararg {
			skip, found = parseVargs(param, params, index, args)
			break
		}
		params[param.Name] = param.parse((*args)[index])
		found = true
		break
	}
	return
}

func parseVargs(param Arg, params Params, index int, args *[]string) (skip int, found bool) {
	params[param.Name] = []string{}
	for len(*args) > index+skip {
		cur := (*args)[index+skip]
		if isFlag(cur) {
			skip--
			break
		}
		params[param.Name] = append(params[param.Name].([]string), param.parse(cur).(string))
		found = true
		skip++
	}
	return
}

func checkRequiredParams(cmd *Command, params map[string]interface{}) error {
	for _, arg := range cmd.Args {
		if !arg.Required {
			continue
		}
		if _, ok := params[arg.Name]; !ok {
			return fmt.Errorf("required argument <%s> not set", arg.Name)
		}
	}

	for _, flag := range cmd.Flags {
		if !flag.Required {
			continue
		}
		if _, ok := params[flag.Name]; !ok {
			return fmt.Errorf("required flag [%s] not set", flag.Name)
		}
	}

	return nil
}
