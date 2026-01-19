{{- define "k8s-ecr-login-renew.jobPodTemplate" -}}
metadata:
{{- with .Values.podAnnotations }}
  annotations:
    {{- toYaml . | trim | nindent 4 }}
{{- end }}
{{- if .Values.forHelm }}
  labels:
    app.kubernetes.io/name: {{ .Chart.Name }}
    helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
    app.kubernetes.io/instance: {{ .Release.Name }}
    app.kubernetes.io/version: {{ .Chart.AppVersion }}
{{- end }}
spec:
  serviceAccountName: {{ required "Service account name is required" .Values.names.serviceAccount }}
  terminationGracePeriodSeconds: {{ required "Termination grace period is required" .Values.cronjob.terminationGracePeriodSeconds }}
  restartPolicy: Never
  containers:
  - name: k8s-ecr-login-renew
    imagePullPolicy: IfNotPresent
    image: {{ required "Docker image must be specficed" .Values.cronjob.dockerImage }}
    env:
{{- if .Values.aws }}
    - name: AWS_ACCESS_KEY_ID
      valueFrom:
        secretKeyRef:
          name: {{ required "AWS credentials secret name is required" .Values.aws.secretName }}
          key: {{ required "AWS credentials secret key acceess key id is required" .Values.aws.secretKeys.accessKeyId }}
    - name: AWS_SECRET_ACCESS_KEY
      valueFrom:
        secretKeyRef:
          name: {{ required "AWS credentials secret name is required" .Values.aws.secretName }}
          key: {{ required "AWS credentials secret key secret acceess key is required" .Values.aws.secretKeys.secretAccessKey }}
{{- end }}
    - name: AWS_REGION
      value: {{ required "AWS region must be specified" .Values.awsRegion }}
    - name: DOCKER_SECRET_NAME
      value: {{ required "Secret name for Docker credentials is required" .Values.dockerSecretName }}
    - name: TARGET_NAMESPACE
      value: {{ required "Target namespace is required" .Values.targetNamespace | quote }}
{{- if .Values.excludeNamespace }}
    - name: EXCLUDE_NAMESPACE
      value: {{ required "Target namespace is required" .Values.excludeNamespace | quote }}
{{- end }}
{{- if .Values.registries }}
    - name: DOCKER_REGISTRIES
      value: {{ .Values.registries }}
{{- end }}
{{- with .Values.tolerations }}
  tolerations:
    {{- toYaml . | nindent 4 }}
{{- end }}
{{- end -}}
