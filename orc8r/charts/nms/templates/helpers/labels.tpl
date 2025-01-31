{{- define "nms.labels" -}}
{{- $envAll := index . 0 -}}
{{- $application := index . 1 -}}
{{- $component := index . 2 -}}
release_group: {{ $envAll.Values.release_group | default $envAll.Release.Name }}
app.kubernetes.io/name: {{ $application }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/instance: {{ $envAll.Release.Name }}
app.kubernetes.io/managed-by: helm
app.kubernetes.io/part-of: magma
{{- end -}}

{{/* Generate selector labels */}}
{{- define "nms.selector-labels" -}}
{{- $envAll := index . 0 -}}
{{- $application := index . 1 -}}
{{- $component := index . 2 -}}
release_group: {{ $envAll.Values.release_group | default $envAll.Release.Name }}
app.kubernetes.io/name: {{ $application }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/instance: {{ $envAll.Release.Name }}
{{- end -}}

{{/* Generate selector labels */}}
{{- define "magmalte-image-version-label" -}}
{{- end -}}

{{/* Generate selector labels */}}
{{- define "nginx-image-version-label" -}}
{{- end -}}
