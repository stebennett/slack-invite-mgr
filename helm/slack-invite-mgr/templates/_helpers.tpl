{{/*
Expand the name of the chart.
*/}}
{{- define "slack-invite-mgr.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "slack-invite-mgr.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "slack-invite-mgr.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "slack-invite-mgr.labels" -}}
helm.sh/chart: {{ include "slack-invite-mgr.chart" . }}
{{ include "slack-invite-mgr.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "slack-invite-mgr.selectorLabels" -}}
app.kubernetes.io/name: {{ include "slack-invite-mgr.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend labels
*/}}
{{- define "slack-invite-mgr.backend.labels" -}}
{{ include "slack-invite-mgr.labels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "slack-invite-mgr.backend.selectorLabels" -}}
{{ include "slack-invite-mgr.selectorLabels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Web labels
*/}}
{{- define "slack-invite-mgr.web.labels" -}}
{{ include "slack-invite-mgr.labels" . }}
app.kubernetes.io/component: web
{{- end }}

{{/*
Web selector labels
*/}}
{{- define "slack-invite-mgr.web.selectorLabels" -}}
{{ include "slack-invite-mgr.selectorLabels" . }}
app.kubernetes.io/component: web
{{- end }}

{{/*
Sheets labels
*/}}
{{- define "slack-invite-mgr.sheets.labels" -}}
{{ include "slack-invite-mgr.labels" . }}
app.kubernetes.io/component: sheets
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "slack-invite-mgr.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "slack-invite-mgr.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create image name with registry
*/}}
{{- define "slack-invite-mgr.image" -}}
{{- $registry := .global.imageRegistry -}}
{{- $username := .global.githubUsername -}}
{{- $repository := .image.repository -}}
{{- $tag := .image.tag -}}
{{- if $username }}
{{- printf "%s/%s/%s:%s" $registry $username $repository $tag }}
{{- else }}
{{- printf "%s/%s:%s" $registry $repository $tag }}
{{- end }}
{{- end }}

{{/*
Backend image
*/}}
{{- define "slack-invite-mgr.backend.image" -}}
{{- include "slack-invite-mgr.image" (dict "global" .Values.global "image" .Values.backend.image) }}
{{- end }}

{{/*
Web image
*/}}
{{- define "slack-invite-mgr.web.image" -}}
{{- include "slack-invite-mgr.image" (dict "global" .Values.global "image" .Values.web.image) }}
{{- end }}

{{/*
Sheets image
*/}}
{{- define "slack-invite-mgr.sheets.image" -}}
{{- include "slack-invite-mgr.image" (dict "global" .Values.global "image" .Values.sheets.image) }}
{{- end }}

{{/*
Image pull secrets
*/}}
{{- define "slack-invite-mgr.imagePullSecrets" -}}
{{- with .Values.global.imagePullSecrets }}
imagePullSecrets:
{{- toYaml . | nindent 2 }}
{{- end }}
{{- end }}

{{/*
Secret names
*/}}
{{- define "slack-invite-mgr.googleCredentialsSecretName" -}}
{{- printf "%s-google-credentials" (include "slack-invite-mgr.fullname" .) }}
{{- end }}

{{- define "slack-invite-mgr.appSecretsName" -}}
{{- printf "%s-app-secrets" (include "slack-invite-mgr.fullname" .) }}
{{- end }}

{{- define "slack-invite-mgr.smtpSecretsName" -}}
{{- printf "%s-smtp-secrets" (include "slack-invite-mgr.fullname" .) }}
{{- end }}
