package doc

import (
	"bytes"
	"io"
	"strings"
	"whitefrog/cli"
)

func GenMarkdown(cmd *cli.Command, dir string) error {
	buf := strings.Builder{}
	buf.WriteString("#" + cmd.Name + "\n\n")
	if cmd.Long != "" {
		buf.WriteString(cmd.Long + "\n\n")
	}
	return nil
}

func genMarkdown(cmd *cli.Command, w io.Writer) error {
	buf := bytes.Buffer{}
	buf.WriteString("##" + cmd.Name + "\n\n")
	buf.WriteString(cmd.Short + "\n\n")
	if cmd.Long != "" {
		buf.WriteString("### Synopsis\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}
	if cmd.Runnable() {
		buf.WriteString("```bash\n" + cmd.Usage() + "\n```\n\n")
	}
	_, err := buf.WriteTo(w)
	return err
}
