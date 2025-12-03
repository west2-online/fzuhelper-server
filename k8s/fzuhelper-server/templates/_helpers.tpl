
{{/*镜像地址（默认tag为服务名，可以自定义imageOverride）*/}}

{{- define "fzuhelper.image" -}}
{{- $service := .service -}}
{{- $values := .values -}}
{{- $svc := index $values.services $service -}}
{{- if $svc.imageOverride }}
{{- $svc.imageOverride }}
{{- else -}}
{{- printf "%s/%s/%s:%s" $values.global.image.registry $values.global.image.project $values.global.image.baseName (default $service $svc.imageTag) -}}
{{- end -}}
{{- end -}}


{{/*环境变量合并函数，合并来自全局和服务本身的环境变量*/}}

{{- define "mergeEnvs" -}}
{{- $globalEnvs := .Values.global.env -}}
{{- $serviceEnvs := .env -}}
{{- $mergedEnvs := concat $globalEnvs $serviceEnvs -}}
{{- if $mergedEnvs -}}
{{- toYaml $mergedEnvs | nindent 12 }}
{{- end -}}
{{- end -}}
