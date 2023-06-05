# 🌱 Mikochi: a minimalist remote file browser

[GitHub](https://github.com/zer0tonin/Mikochi)

## Values


| Key                           | Description                        | Default                    |
|-------------------------------|------------------------------------|----------------------------|
| mikochi.username              | Username for credentials           | root                       |
| mikochi.password              | Password for credentials           | pass                       |
| image.repository              | Docker image repo                  | zer0tonin/mikochi          |
| image.tag                     | Docker image tag                   | 1.2.0                      |
| image.pullPolicy              | Image pull policy                  | IfNotPresent               |
| imagePullSecrets              | Secrets for private docker repos   | {}                         |
| service.type                  | Service type                       | ClusterIP                  |
| service.port                  | Service inbound/outbound port      | 8080                       |
| ingress.enabled               | Automatically create ingress object| false                      |
| ingress.className             | Ingress class name                 |                            |
| ingress.annotations           | Custom ingress annotations         | {}                         |
| ingress.hosts                 | Hosts and paths                    | []                         |
| ingress.tls                   | SSL/TLS configuration              | []                         |
| persistence.enabled           | Use PVC                            | false                      |
| persistence.existingClaim     | Use existing PVC                   | false                      |
| persistence.annotations       | PVC annotations                    | []                         |
| persistence.skipuninstall     | Skip PVC uninstall                 | false                      |
| persistence.accessModes       | PVC access modes                   | [ReadWriteOnce]            |
| persistence.size              | PVC requested size                 | 1Gi                        |
| persistence.storageClassName  | PVC storage class                  | local-path                 |
| resources                     | CPU/RAM requests/limits            | {}                         |
| nodeSelector                  | Custom node selectors              | {}                         |
| tolerations                   | Custom node tolerations            | []                         |
| nameOverride                  | App name override                  |                            |
| fullnameOverride              | Chart name override                |                            |
| podAnnotations                | Custom pod annotations             | {}                         |
| podSecurityContext            | Custom pod security context        | {}                         |
| securityContext               | Custom security context            | {}                         |
