// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inject

import (
	"os"
	"strings"
	"text/template"

	meshconfig "istio.io/api/mesh/v1alpha1"
)

type Init struct {
	Image      string
	Args       []string
	PullPolicy string
}

type Container struct {
	Image      string
	Args       []string
	PullPolicy string
	Env        []string
}

type Config struct {
	Init
	Container
	Mesh meshconfig.MeshConfig
}

var (
	serviceDefault = `
- --serviceCluster
{{ if eq .ServiceCluster "" }}
- {{ "istio-proxy" }}
{{ else }}
- {{ printf "%s" .ServiceCluster }}
{{ end }}`
	meshTemplate = `
initContainers:
- name: istio-init
  {{ if eq .Init.Image "" -}}
  image: {{ "gcr.io/istio-testing/proxy_init:3101ea9d82a5f83b699c2d3245b371a19fa6bef4" }}
  {{ else -}}
  image: {{ printf "%s" .Init.Image }}
  {{ end -}}
  {{ $length := len .Init.Args }} {{ if eq $length 0 -}}
  args: ["-p", "15001", "-u", "1337"]
  {{ else -}}
  args: {{ printf "%s" .Init.Args }}
  {{ end -}}
  {{ if eq .Init.PullPolicy "" -}}
  imagePullPolicy: {{ "IfNotPresent" }}
  {{ else -}}
  imagePullPolicy: {{ printf "%s" .Init.PullPolicy }}
  {{ end -}}
  securityContext:
    capabilities:
      add:
      - NET_ADMIN
    priviledged: true
  restartPolicy: Always
containers:
- name: istio-proxy
  {{ if eq .Container.Image "" -}}
  image: {{ "gcr.io/istio-testing/proxy:3101ea9d82a5f83b699c2d3245b371a19fa6bef4" }}
  {{ else -}}
  image: {{ printf "%s" .Container.Image }}
  {{ end -}}
  {{ $length := len .Container.Args }} {{ if eq $length 0 }}
  args:
  - proxy
  - sidecar
  - -v
  - "2"
  - --configPath
  - {{ printf "%s" .Mesh.DefaultConfig.ConfigPath }}
  - --binaryPath
  - {{ printf "%s" .Mesh.DefaultConfig.BinaryPath }}
  - --drainDuration
  - {{ printf "%v" .Mesh.DefaultConfig.DrainDuration.Seconds }}
  - --parentShutdownDuration
  - {{ printf "%v" .Mesh.DefaultConfig.ParentShutdownDuration.Seconds }}
  - --discoveryAddress
  - {{ printf "%s" .Mesh.DefaultConfig.DiscoveryAddress }}
  - --discoveryRefreshDelay
  - {{ printf "%v" .Mesh.DefaultConfig.DiscoveryRefreshDelay.Seconds }}
  - --zipkinAddress
  - zipkin.istio-system:9411
  - --connectTimeout
  - {{ printf "%v" .Mesh.DefaultConfig.ConnectTimeout.Seconds }}
  - --statsdUdpAddress
  - istio-mixer.istio-system:9125
  - --proxyAdminPort
  - {{ printf "%v" .Mesh.DefaultConfig.ProxyAdminPort }}
  - --controlPlaneAuthPolicy
  - NONE` + escape(`
  - --serviceCluster
  {{ if eq .ServiceCluster "" -}}
  - {{ printf "%s" .ServiceCluster }}
  {{ else -}}
  - {{ printf "%s" .ServiceCluster }}
  {{ end -}}`) + `
  {{ else -}}
  args: {{ printf "%s" .Container.Args }}
  {{ end -}}
  {{ $length := len .Container.Env }} {{ if eq $length 0 }}
  env:
  - name: POD_NAME
    valueFrom:
      fieldRef:
        apiVersion: v1
        fieldPath: metadata.name
  - name: POD_NAMESPACE
    valueFrom:
      fieldRef:
        apiVersion: v1
        fieldPath: metadata.namespace
  - name: INSTANCE_IP
    valueFrom:
      fieldRef:
        apiVersion: v1
        fieldPath: status.podIP
  {{ else -}}
  env: {{ printf "%s" .Container.Env }}
  {{ end -}}
  {{ if eq .Container.PullPolicy "" -}}
  imagePullPolicy: {{ "IfNotPresent" }}
  {{ else -}}
  imagePullPolicy: {{ printf "%s" .Container.PullPolicy }}
  {{ end -}}
  securityContext:
    capabilities:
      add:
      - NET_ADMIN
    priviledged: true
  restartPolicy: Always
  volumeMounts:
  - mountPath: /etc/istio/proxy
    name: istio-envoy
  - mountPath: /etc/certs/
    name: istio-certs
    readOnly: true
volumes:
- emptyDir:
    medium: Memory
  name: istio-envoy
- name: istio-certs
  secret:
    defaultMode: 420
    optional: true` + escape(`
    {{ if eq .Spec.ServiceAccountName "" -}}
    secretName: {{ "istio.default" }}
    {{ else -}}
    secretName: {{ printf "istio.%s" .Spec.ServiceAccountName }}
    {{ end -}}`)
)

func DefaultMeshTemplate(mc *Config) error {
	t := template.Must(template.New("inject").Parse(meshTemplate))
	if err := t.Execute(os.Stdout, mc); err != nil {
		return err
	}
	return nil
}

func CreateMeshTemplate(mc *Config) (*template.Template, error) {
	t := template.Must(template.New("inject").Parse(meshTemplate))
	if err := t.Execute(os.Stdout, mc); err != nil {
		return nil, err
	}
	return t, nil
}

func escape(in string) string {
	return strings.TrimSpace("{{`" + in + "`}}")
}
