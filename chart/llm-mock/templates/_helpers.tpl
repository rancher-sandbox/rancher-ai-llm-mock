{{/* Generate a fullname */}}
{{- define "llm-mock.fullname" -}}
{{- printf "%s" .Release.Name -}}
{{- end -}}

{{/* Generate common labels */}}
{{- define "llm-mock.labels" -}}
app.kubernetes.io/name: {{ include "llm-mock.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/* Chart name */}}
{{- define "llm-mock.name" -}}
llm-mock
{{- end -}}
