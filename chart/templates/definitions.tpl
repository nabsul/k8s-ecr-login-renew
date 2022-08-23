{{/* vim: set filetype=mustache: */}}
{{/*
Define resource names
*/}}
{{- define "k8s-ecr-login-renew.namespace" }}
{{- default (printf "%s-ns" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.serviceAccount" }}
{{- default (printf "%s-account" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.role" }}
{{- default (printf "%s-role" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.roleBinding" }}
{{- default (printf "%s-rolebinding" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.job" }}
{{- default (printf "%s-job" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.cronJob" }}
{{- default (printf "%s-cron" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.secret" }}
{{- .Values.ecr.auth.existingSecret | default (printf "%s-secret" .Release.Name) -}}
{{- end }}

{{- define "k8s-ecr-login-renew.targetNamespace" }}
{{- default .Release.Namespace .Values.targetNamespace -}}
{{- end }}
