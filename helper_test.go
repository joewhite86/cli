package cli_test

func osArgs(args []string) []string {
	a := make([]string, 0, len(args)+1)
	a = append(a, "cmd")
	a = append(a, args...)
	return a
}
