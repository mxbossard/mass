Labels: 
{{ range $key, $value := .Labels }}
	{{ $key }}: {{ $value }}
{{ end -}}
Tags: 
{{ range $key, $value := .Tags }}
	{{ $key }}: {{ $value }}
{{ end -}}
Environments: 
{{ range $key, $value := .Environment }}
	{{ $key }}: {{ $value }}
{{ end -}}
