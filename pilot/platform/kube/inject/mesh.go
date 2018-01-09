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

// type Init struct {
// 	Image      string
// 	Args       []string
// 	PullPolicy string
// }

// type Container struct {
// 	Image      string
// 	Args       []string
// 	PullPolicy string
// 	Env        []string
// }

// type Config struct {
// 	Init
// 	Container
// 	Mesh meshconfig.MeshConfig
// }

var (
	productionTemplate = `
initContainers:
- name: istio-init
  image: {{ printf "%s" .MConfig.InitImage }}
  args: ["-p", "15001", "-u", "1337"]
  {{ if eq .MConfig.ImagePullPolicy "" -}}
  imagePullPolicy: {{ "IfNotPresent" }}
  {{ else -}}
  imagePullPolicy: {{ printf "%s" .MConfig.ImagePullPolicy }}
  {{ end -}}
  securityContext:
    capabilities:
      add:
      - NET_ADMIN
    priviledged: true
  restartPolicy: Always
containers:
- name: istio-proxy
  image: {{ printf "%s" .MConfig.ProxyImage }}
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
  - NONE
  - --serviceCluster
  {{ if eq .ServiceCluster "" -}}
  - {{ printf "%s" .ServiceCluster }}
  {{ else -}}
  - {{ printf "%s" .ServiceCluster }}
  {{ end -}}
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
  {{ if eq .MConfig.ImagePullPolicy "" -}}
  imagePullPolicy: {{ "IfNotPresent" }}
  {{ else -}}
  imagePullPolicy: {{ printf "%s" .MConfig.ImagePullPolicy }}
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
    optional: true
    {{ if eq .Spec.ServiceAccountName "" -}}
    secretName: {{ "istio.default" }}
    {{ else -}}
    secretName: {{ printf "istio.%s" .Spec.ServiceAccountName }}
    {{ end -}}`
)

// func CreateProductionTemplate(p *Params) error {
// 	t := template.Must(template.New("inject").Parse(productionTemplate))
// 	if err := t.Execute(os.Stdout, mc); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func CreateDebugTemplate(mc *Config) (*template.Template, error) {
// 	t := template.Must(template.New("inject").Parse(meshTemplate))
// 	if err := t.Execute(os.Stdout, mc); err != nil {
// 		return nil, err
// 	}
// 	return t, nil
// }

// func escape(in string) string {
// 	return strings.TrimSpace("{{`" + in + "`}}")
// }
