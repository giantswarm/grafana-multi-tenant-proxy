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
| fullnameOverride | string | `nil` | Overrides the chart's computed fullname |
| global.image.registry | string | `"gsoci.azurecr.io"` | Overrides the Docker registry globally for all images |
| nameOverride | string | `nil` | Overrides the chart's name |
| proxy.autoscaling.enabled | bool | `true` | Enable autoscaling for the multi-tenant proxy |
| proxy.autoscaling.maxReplicas | int | `4` | Maximum autoscaling replicas for the multi-tenant proxy |
| proxy.autoscaling.minReplicas | int | `2` | Minimum autoscaling replicas for the multi-tenant proxy |
| proxy.autoscaling.targetCPUUtilizationPercentage | int | `90` | Target CPU utilisation percentage for the multi-tenant proxy |
| proxy.autoscaling.targetMemoryUtilizationPercentage | string | `nil` | Target memory utilisation percentage for the multi-tenant proxy |
| proxy.configReloader.containerSecurityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"seccompProfile":{"type":"RuntimeDefault"}}` | Security context to apply to the config reloader containers. |
| proxy.configReloader.image.repository | string | `"giantswarm/configmap-reload"` | Repository to get config reloader image from. |
| proxy.configReloader.image.tag | string | `"v0.13.1"` | Tag of image to use for config reloading. |
| proxy.configReloader.resources | object | `{"requests":{"cpu":"1m","memory":"5Mi"}}` | Resource requests and limits to apply to the config reloader containers. |
| proxy.containerPort | int | `3501` |  |
| proxy.containerSecurityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"seccompProfile":{"type":"RuntimeDefault"}}` | The SecurityContext for Loki containers |
| proxy.credentials | string | `"users:\n  - username: Tenant1\n    password: 1tnaneT\n    orgid: tenant-1\n  - username: Tenant2\n    password: 2tnaneT\n    orgid: tenant-2"` |  |
| proxy.deployCredentials | bool | `true` |  |
| proxy.enabled | bool | `true` | Specifies whether the multi-tenant proxy should be enabled |
| proxy.image | object | `{"pullPolicy":"IfNotPresent","pullSecrets":[],"repository":"giantswarm/grafana-multi-tenant-proxy","tag":null}` | ref: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy |
| proxy.image.repository | string | `"giantswarm/grafana-multi-tenant-proxy"` | Repository to get multi-tenant proxy image from. |
| proxy.image.tag | string | `nil` | Overrides the image tag whose default is the chart's appVersion |
| proxy.ingress.enabled | bool | `false` | Specifies whether an ingress for the multi-tenant-proxy should be created |
| proxy.ingress.hosts | list | `[{"host":"multi-tenant-proxy.loki.example.com","paths":[{"path":"/"}]}]` | Hosts configuration for the multi-tenant-proxy ingress, passed through the `tpl` function to allow templating |
| proxy.ingress.ingressClassName | string | `""` | Ingress Class Name. MAY be required for Kubernetes versions >= 1.18 |
| proxy.ingress.tls | list | `[{"hosts":["write.multi-tenant-proxy.loki.example.com"],"secretName":"loki-multi-tenant-proxy-tls"}]` | TLS configuration for the gateway ingress. Hosts passed through the `tpl` function to allow templating |
| proxy.networkPolicy.enabled | bool | `true` |  |
| proxy.networkPolicy.flavor | string | `"cilium"` |  |
| proxy.podSecurityContext.fsGroup | int | `10001` |  |
| proxy.podSecurityContext.runAsGroup | int | `10001` |  |
| proxy.podSecurityContext.runAsNonRoot | bool | `true` |  |
| proxy.podSecurityContext.runAsUser | int | `10001` |  |
| proxy.podSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| proxy.replicas | int | `3` | Number of replicas for the multi-tenant proxy |
| proxy.resources | object | `{"limits":{"memory":"500Mi"},"requests":{"cpu":"50m","memory":"50Mi"}}` | Resource requests and limits for the write |
| proxy.service.port | int | `80` |  |
| proxy.targetServers | list | `[]` | List of target servers to proxy |

