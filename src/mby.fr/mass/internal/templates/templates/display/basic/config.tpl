---
labels: 
{{- range $key, $value := .Labels }}
  {{ $key }}: {{ $value }}
{{- end }}
tags: 
{{- range $key, $value := .Tags }}
  {{ $key }}: {{ $value }}
{{- end }}
environment: 
{{- range $key, $value := .Environment }}
  {{ $key }}: {{ $value }}
{{- end }}
buildArgs: 
{{- range $key, $value := .BuildArgs }}
  {{ $key }}: {{ $value }}
{{- end }}
runArgs: 
{{- range $key, $value := .RunArgs }}
  - {{ $value }}
{{- end }}
---
