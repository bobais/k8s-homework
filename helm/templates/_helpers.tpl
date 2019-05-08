{{/*
Computed definitions
*/}}

{{- define "prometheus.name" -}}
{{ .Chart.Name }}-{{ .Release.Name }}-prometheus
{{- end -}}

{{- define "prometheus.url" -}}
http://{{ include "prometheus.name" . }}/
{{- end -}}

{{- define "reporter.name" -}}
{{ .Chart.Name }}-{{ .Release.Name }}-reporter
{{- end -}}

{{- define "watcher.name" -}}
{{ .Chart.Name }}-{{ .Release.Name }}-watcher
{{- end -}}

{{- define "eventstats.common.labels" -}}
chart: {{ .Chart.Name }}-{{ .Chart.Version }}
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
{{- end -}}

{{- define "eventstats.watcher.labels" -}}
component: watcher
{{ include "eventstats.common.labels" . }}
{{- end -}}

{{- define "eventstats.reporter.labels" -}}
component: reporter
{{ include "eventstats.common.labels" . }}
{{- end -}}

{{- define "eventstats.prometheus.labels" -}}
component: prometheus
{{ include "eventstats.common.labels" . }}
{{- end -}}


