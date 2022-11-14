{{ if .Long -}}
{{ .Long }}
{{ else if .Short -}}
{{ .Short }}
{{ end -}}

{{ range $group, $commands := groups }}
{{ if $group }}{{ $group -}}{{ else }}Available Commands{{ end }}:
  {{- range $commands -}}
    {{ .Name | formatSubCommand -}}{{ .Short }}
  {{ end -}}
{{ end }}

{{- if .Args }}
Arguments:
{{- range .Args }}
  {{ .Name | formatArg }}{{ .Description }}{{ if not .Required }} {{ if .Default }}(Default: {{ .Default }}){{ else }}(Optional){{ end }}{{ end -}}
{{ end }}
{{ end }}

{{- if .Flags }}
Flags:
{{- range .Flags }}
  {{ if .Short }}-{{ .Short }},{{ end }}--{{ .Name }}{{ if .Default }}={{ .Default }}{{ end }}:     {{ .Description }}
{{ end -}}
{{ end }}
Usage:
  {{ usage }}
