package cli_test

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/joewhite86/cli"
)

var ctx = context.Background()

func TestRun_Version(t *testing.T) {
	type args struct {
		cmd     *cli.Command
		cliArgs []string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{
		{
			name: "print_default_version",
			args: args{
				cmd:     &cli.Command{},
				cliArgs: []string{"-v"},
			},
			wantOutput: "0.1.0",
		},
		{
			name: "print_root_version",
			args: args{
				cmd:     &cli.Command{Version: "1.0"},
				cliArgs: []string{"-v"},
			},
			wantOutput: "1.0",
		},
		{
			name: "print_root_version_long",
			args: args{
				cmd:     &cli.Command{Version: "1.0"},
				cliArgs: []string{"--version"},
			},
			wantOutput: "1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := bytes.Buffer{}
			cli.Out = &out
			os.Args = osArgs(tt.args.cliArgs)
			if err := cli.Run(ctx, tt.args.cmd); err != nil {
				t.Errorf("Run() returned an error: %v", err)
			}
			if out.String()[:out.Len()-1] != tt.wantOutput {
				t.Errorf("Run() output = %v, wantOutput %v", out.String(), tt.wantOutput)
			}
		})
	}
}

func TestRun_ShouldParseFlag(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		flags      []cli.Flag
		arguments  []cli.Arg
		expected   cli.Params
		shouldFail bool
	}{{
		name:       "NotFound",
		args:       []string{"-s"},
		shouldFail: true,
	}, {
		name:     "StringValue",
		args:     []string{"-s", "test"},
		flags:    []cli.Flag{{Name: "string", Short: "s", HasValue: true}},
		expected: map[string]interface{}{"string": "test"},
	}, {
		name:     "Int32Value",
		args:     []string{"-i", "32"},
		flags:    []cli.Flag{{Name: "int32", Short: "i", HasValue: true, Parser: cli.Int32Parser}},
		expected: map[string]interface{}{"int32": int32(32)},
	}, {
		name:     "TwoFlags",
		args:     []string{"-i", "-s"},
		flags:    []cli.Flag{{Name: "i", Short: "i"}, {Name: "s", Short: "s"}},
		expected: map[string]interface{}{"i": true, "s": true},
	}, {
		name:     "FlagUnset",
		args:     []string{},
		flags:    []cli.Flag{{Name: "i", Short: "i"}},
		expected: map[string]interface{}{},
	}, {
		name:       "RequiredMissing",
		args:       []string{},
		flags:      []cli.Flag{{Name: "i", Short: "i", Required: true}},
		shouldFail: true,
	}, {
		name:      "StringValue",
		args:      []string{"test"},
		arguments: []cli.Arg{{Name: "string"}},
		expected:  cli.Params{"string": "test"},
	}, {
		name:      "Int32Value",
		args:      []string{"32"},
		arguments: []cli.Arg{{Name: "int32", Parser: cli.Int32Parser}},
		expected:  cli.Params{"int32": int32(32)},
	}, {
		name:      "Int32InvalidValue",
		args:      []string{"test"},
		arguments: []cli.Arg{{Name: "int32", Parser: cli.Int32Parser}},
		expected:  cli.Params{"int32": int32(-1)},
	}, {
		name:      "TwoArgs",
		args:      []string{"test1", "test2"},
		arguments: []cli.Arg{{Name: "i"}, {Name: "s"}},
		expected:  cli.Params{"i": "test1", "s": "test2"},
	}, {
		name:      "ArgUnset",
		args:      []string{},
		arguments: []cli.Arg{{Name: "i"}},
		expected:  cli.Params{},
	}, {
		name:       "RequiredMissing",
		args:       []string{},
		arguments:  []cli.Arg{{Name: "i", Required: true}},
		shouldFail: true,
	}, {
		name:      "RequiredAfterParam",
		args:      []string{"-i", "test"},
		arguments: []cli.Arg{{Name: "required", Required: true}},
		flags:     []cli.Flag{{Name: "i", Short: "i"}},
		expected:  cli.Params{"required": "test", "i": true},
	}, {
		name:      "Varargs",
		args:      []string{"test1", "test2"},
		arguments: []cli.Arg{{Name: "required", Required: true, Vararg: true}},
		expected:  cli.Params{"required": []string{"test1", "test2"}},
	}, {
		name:      "VarargsAndFlag",
		args:      []string{"test1", "test2", "-p", "test"},
		arguments: []cli.Arg{{Name: "required", Required: true, Vararg: true}},
		flags:     []cli.Flag{{Name: "p", Short: "p", HasValue: true}},
		expected:  cli.Params{"required": []string{"test1", "test2"}, "p": "test"},
	}, {
		name:      "FlagAndVarargs",
		args:      []string{"-p", "test", "test1", "test2"},
		arguments: []cli.Arg{{Name: "required", Required: true, Vararg: true}},
		flags:     []cli.Flag{{Name: "p", Short: "p", HasValue: true}},
		expected:  cli.Params{"required": []string{"test1", "test2"}, "p": "test"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = osArgs(tt.args)
			run := func(_ context.Context, params cli.Params) error {
				delete(params, "_args")
				if equal := reflect.DeepEqual(tt.expected, params); !equal {
					t.Errorf("expected = %+v, got = %+v", tt.expected, params)
				}
				return nil
			}
			cmd := &cli.Command{Flags: tt.flags, Args: tt.arguments, Run: run}
			if err := cli.Run(ctx, cmd); (err != nil) != tt.shouldFail {
				t.Errorf("unexpected error %v", err)
			}
			os.Args = append([]string{"cmd"}, osArgs(tt.args)...)
			subCmd := cli.Command{Name: "cmd", Flags: tt.flags, Args: tt.arguments, Run: run}
			cmd = &cli.Command{Commands: []cli.Command{subCmd}}
			if err := cli.Run(ctx, cmd); (err != nil) != tt.shouldFail {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

func TestRun_ShouldLint(t *testing.T) {
	cmd := cli.Command{}
	os.Args = []string{"cmd", "lint"}
	buf := bytes.Buffer{}
	cli.Err = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("Output doesn't contain a warning")
	}
}

func TestRun_ShouldPrintHelpUsage(t *testing.T) {
	cmd := cli.Command{Name: "cmd"}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "Usage:\n  cmd") {
		t.Errorf("Output doesn't contain a usage line")
	}
}

func TestRun_ShouldPrintHelpSubCommand(t *testing.T) {
	cmd := cli.Command{Name: "cmd", Commands: []cli.Command{{
		Group: "group-name", Name: "sub-cmd", Short: "description text",
	}}}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "sub-cmd") {
		t.Errorf("Output doesn't contain sub command name")
	}
	if !strings.Contains(buf.String(), "description text") {
		t.Errorf("Output doesn't contain sub command short-description text")
	}
}

func TestRun_ShouldPrintHelpCommandGroup(t *testing.T) {
	cmd := cli.Command{Name: "cmd", Commands: []cli.Command{{Group: "group-name", Name: "sub-cmd"}}}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "group-name:") {
		t.Errorf("Output doesn't contain group name line")
	}
}

func TestRun_ShouldPrintHelpArguments(t *testing.T) {
	cmd := cli.Command{Name: "cmd", Args: []cli.Arg{{Name: "arg1"}}}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "Arguments:") {
		t.Errorf("Output doesn't contain arguments line")
	}
}

func TestRun_ShouldPrintHelpArgumentLine(t *testing.T) {
	cmd := cli.Command{Name: "cmd", Args: []cli.Arg{{Name: "arg1", Description: "desc1"}}}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "arg1") {
		t.Errorf("Output doesn't contain argument name")
	}
	if !strings.Contains(buf.String(), "desc1") {
		t.Errorf("Output doesn't contain argument description")
	}
}

func TestRun_ShouldPrintHelpFlags(t *testing.T) {
	cmd := cli.Command{Name: "cmd", Args: []cli.Arg{{Name: "arg1"}}}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "Flags:") {
		t.Errorf("Output doesn't contain arguments line")
	}
}

func TestRun_ShouldPrintHelpFlagLine(t *testing.T) {
	cmd := cli.Command{Name: "cmd", Flags: []cli.Flag{{Name: "flag1", Description: "desc1"}}}
	os.Args = []string{"cmd", "help"}
	buf := bytes.Buffer{}
	cli.Out = &buf
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !strings.Contains(buf.String(), "flag1") {
		t.Errorf("Output doesn't contain flag name")
	}
	if !strings.Contains(buf.String(), "desc1") {
		t.Errorf("Output doesn't contain flag description")
	}
}

func TestRun_ShouldExecCommand(t *testing.T) {
	ran := false
	res := cli.Command{
		Name: "exec",
		Run: func(_ context.Context, _ cli.Params) error {
			ran = true
			return nil
		},
	}
	cmd := cli.Command{Commands: []cli.Command{res}}
	os.Args = []string{"cmd", "exec"}
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !ran {
		t.Error("Handler not executed")
	}
}

func TestRun_ShouldExecSecondCommand(t *testing.T) {
	ran := false
	first := cli.Command{
		Name: "other",
	}
	res := cli.Command{
		Name: "exec",
		Run: func(_ context.Context, _ cli.Params) error {
			ran = true
			return nil
		},
	}
	cmd := cli.Command{Commands: []cli.Command{first, res}}
	os.Args = []string{"cmd", "exec"}
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !ran {
		t.Error("Handler not executed")
	}
}

func TestRun_ShouldExecSubCommand(t *testing.T) {
	ran := false
	subres := cli.Command{
		Name: "sub",
		Run: func(_ context.Context, _ cli.Params) error {
			ran = true
			return nil
		},
	}
	res := cli.Command{Name: "exec", Commands: []cli.Command{subres}}
	cmd := cli.Command{Commands: []cli.Command{res}}
	os.Args = []string{"cmd", "exec", "sub"}
	if err := cli.Run(ctx, &cmd); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !ran {
		t.Error("Handler not executed")
	}
}
