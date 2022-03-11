# ACME webhook for Active24 DNS API

This repository contains code and supporting files for ACME webhook that interacts with [active24.cz](https://customer.active24.com/user/api) DNS API.

## Installation

### Requirements

- [cert-manager](https://cert-manager.io/docs/installation/)

- [API token](https://customer.active24.com/user/api) to access your domain

Create secret with API token

```
kubectl create secret generic active24-apikey --namespace cert-manager \
	--from-literal='apiKey=abcd1234567890'
```

Create ClusterIssuer


Apply following manifest into cluster

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    # The ACME server URL
    server: https://acme-v02.api.letsencrypt.org/directory
    # Email address used for ACME registration
    email: admin@somegreatdomain.tld
    # Name of a secret used to store the ACME account private key
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - selector:
        dnsZones:
          - somegreatdomain.tld
      dns01:
        webhook:
          groupName: acme.yourdomain.tld
          solverName: active24
          config:
            apiKeySecretRef:
              name: 'active24-apikey'
              key: 'apiKey'
            domain: somegreatdomain.tld
```

Replace `somegreatdomain.tld` with actual domain managed by Active24

Install using helm

```
helm upgrade --install ac24 ./deploy/chart --namespace cert-manager
```

Create certificate

```yaml
kind: Certificate
apiVersion: cert-manager.io/v1
metadata:
  name: my-certificate
spec:
  commonName: somegreatdomain.tld
  dnsNames:
    - somegreatdomain.tld
    - '*.somegreatdomain.tld'
  issuerRef:
    kind: ClusterIssuer
    name: letsencrypt-prod
  secretName: somegreatdomain.tld-tls
```