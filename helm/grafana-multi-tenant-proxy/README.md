# grafana-multi-tenant-proxy

![Version: 0.5.0](https://img.shields.io/badge/Version-0.5.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.5.0](https://img.shields.io/badge/AppVersion-0.5.0-informational?style=flat-square)

Helm chart for Grafana Multi Tenant Proxy

**Homepage:** <https://github.com/giantswarm/grafana-multi-tenant-proxy>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| giantswarm/team-atlas | <team-atlas@giantswarm.io> |  |

## Source Code

* <https://github.com/giantswarm/grafana-multi-tenant-proxy>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| configReloader.containerSecurityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Security context to apply to the config reloader containers. |
| configReloader.image.repository | string | `"giantswarm/configmap-reload"` | Repository to get config reloader image from. |
| configReloader.image.tag | string | `"v0.13.1"` | Tag of image to use for config reloading. |
| configReloader.resources | object | `{"requests":{"cpu":"1m","memory":"5Mi"}}` | Resource requests and limits to apply to the config reloader containers. |
| fullnameOverride | string | `nil` | Overrides the chart's computed fullname |
| global.image.registry | string | `"gsoci.azurecr.io"` | Overrides the Docker registry globally for all images |
| ingress.annotations | object | `{}` | Annotations for the gateway ingress |
| ingress.enabled | bool | `false` | Specifies whether an ingress for the multi-tenant-proxy should be created |
| ingress.hosts | list | `[{"host":"multi-tenant-proxy.loki.example.com","paths":[{"path":"/"}]}]` | Hosts configuration for the multi-tenant-proxy ingress, passed through the `tpl` function to allow templating |
| ingress.ingressClassName | string | `""` | Ingress Class Name. MAY be required for Kubernetes versions >= 1.18 |
| ingress.labels | object | `{}` | Labels for the gateway ingress |
| ingress.tls | list | `[{"hosts":["write.multi-tenant-proxy.loki.example.com"],"secretName":"loki-multi-tenant-proxy-tls"}]` | TLS configuration for the gateway ingress. Hosts passed through the `tpl` function to allow templating |
| monitoring.enabled | bool | `true` |  |
| nameOverride | string | `nil` | Overrides the chart's name |
| networkPolicy.enabled | bool | `true` | Specifies whether the multi-tenant proxy should be deployed with a network policy |
| networkPolicy.flavor | string | `"cilium"` | Specifies the flavor of network policy to use |
| podSecurityContext | object | `{"fsGroup":10001,"runAsGroup":10001,"runAsNonRoot":true,"runAsUser":10001,"seccompProfile":{"type":"RuntimeDefault"}}` | The pod SecurityContext |
| proxy.autoscaling.enabled | bool | `true` | Enable autoscaling for the multi-tenant proxy |
| proxy.autoscaling.maxReplicas | int | `4` | Maximum autoscaling replicas for the multi-tenant proxy |
| proxy.autoscaling.minReplicas | int | `2` | Minimum autoscaling replicas for the multi-tenant proxy |
| proxy.autoscaling.targetCPUUtilizationPercentage | int | `90` | Target CPU utilisation percentage for the multi-tenant proxy |
| proxy.autoscaling.targetMemoryUtilizationPercentage | string | `nil` | Target memory utilisation percentage for the multi-tenant proxy |
| proxy.containerPort | int | `3501` | Default container port |
| proxy.containerSecurityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"seccompProfile":{"type":"RuntimeDefault"}}` | The container SecurityContext for the multi-tenant-proxy container |
| proxy.credentials | string | `"users:\n  - username: Tenant1\n    password: 1tnaneT\n    orgid: tenant-1\n  - username: Tenant2\n    password: 2tnaneT\n    orgid: tenant-2"` | The credentials for the multi-tenant-proxy |
| proxy.deployCredentials | bool | `false` |  |
| proxy.image.pullPolicy | string | `"IfNotPresent"` | Overrides the image pull policy whose default is 'IfNotPresent' |
| proxy.image.repository | string | `"giantswarm/grafana-multi-tenant-proxy"` | Repository to get multi-tenant proxy image from. |
| proxy.image.tag | string | `nil` | Overrides the image tag whose default is the chart's appVersion |
| proxy.replicas | int | `3` | Number of replicas for the multi-tenant proxy |
| proxy.resources | object | `{"limits":{"memory":"500Mi"},"requests":{"cpu":"50m","memory":"50Mi"}}` | Resource requests and limits |
| proxy.targetServers | list | `[]` | List of target servers to proxy |
| service.port | int | `80` |  |

