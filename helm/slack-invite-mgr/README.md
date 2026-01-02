# Slack Invite Manager Helm Chart

A Helm chart for deploying Slack Invite Manager to Kubernetes via ArgoCD.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Traefik Ingress Controller
- External Secrets Operator (ESO) installed
- A configured SecretStore or ClusterSecretStore
- ArgoCD (for GitOps deployment)

## Components

This chart deploys the following components:

| Component | Type | Description |
|-----------|------|-------------|
| Backend API | Deployment + Service | Go REST API server |
| Frontend Web | Deployment + Service | React SPA served by nginx |
| Sheets Sync | CronJob | Scheduled Google Sheets synchronization |
| Ingress | Traefik IngressRoute | TLS-enabled ingress with path routing |
| External Secrets | ExternalSecret | Secrets fetched from external store |

## Installation

### Using Helm directly

```bash
# Add your values
helm install slack-invite-mgr ./helm/slack-invite-mgr \
  --namespace slack-invite-mgr \
  --create-namespace \
  -f values-prod.yaml \
  --set global.githubUsername=your-username \
  --set externalSecrets.secretStoreName=your-secret-store
```

### Using ArgoCD

1. Update `helm/argocd/application.yaml` with your repository URL
2. Apply the ArgoCD Application:

```bash
kubectl apply -f helm/argocd/application.yaml
```

## Configuration

### Required Values

| Parameter | Description |
|-----------|-------------|
| `global.githubUsername` | GitHub username for container registry |
| `externalSecrets.secretStoreName` | Name of your SecretStore/ClusterSecretStore |
| `ingress.host` | Hostname for the ingress |
| `web.env.apiUrl` | Browser-accessible API URL |

### External Secrets Setup

Before deploying, ensure your external secret store contains:

1. **Google Credentials** (`google-credentials`):
   - `credentials.json`: Google service account JSON

2. **App Secrets** (`slack-invite/app`):
   - `spreadsheetId`: Google Spreadsheet ID
   - `sheetName`: Sheet name within the spreadsheet

3. **SMTP Secrets** (`slack-invite/smtp`) - only needed if sheets is enabled:
   - `emailRecipient`: Notification email address
   - `fromEmail`: Sender email address
   - `username`: SMTP2Go username
   - `password`: SMTP2Go API key
   - `dashboardUrl`: Dashboard URL for email templates

### Full Configuration Reference

See `values.yaml` for all available configuration options.

#### Global Settings

```yaml
global:
  imageRegistry: ghcr.io
  githubUsername: "your-username"
  imagePullSecrets: []
```

#### Backend Configuration

```yaml
backend:
  enabled: true
  replicas: 1
  port: 8080
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 256Mi
```

#### Frontend Configuration

```yaml
web:
  enabled: true
  replicas: 1
  port: 80
  env:
    apiUrl: "https://slack-invite.example.com/api"
    publicUrl: "/slack-invite"
```

#### Sheets CronJob Configuration

```yaml
sheets:
  enabled: true
  schedule: "0 * * * *"  # Every hour
  concurrencyPolicy: Forbid
```

#### Ingress Configuration

```yaml
ingress:
  enabled: true
  host: slack-invite.example.com
  path: /slack-invite
  entryPoints:
    - websecure
  tls:
    enabled: true
    secretName: slack-invite-tls
```

#### External Secrets Configuration

```yaml
externalSecrets:
  enabled: true
  secretStoreName: "your-secret-store"
  secretStoreKind: ClusterSecretStore
  refreshInterval: 1h
  googleCredentials:
    remoteRef: google-credentials
  appSecrets:
    remoteRef: slack-invite/app
  smtpSecrets:
    remoteRef: slack-invite/smtp
```

## Upgrading

```bash
helm upgrade slack-invite-mgr ./helm/slack-invite-mgr \
  --namespace slack-invite-mgr \
  -f values-prod.yaml
```

Or with ArgoCD, simply push changes to your Git repository.

## Uninstalling

```bash
helm uninstall slack-invite-mgr --namespace slack-invite-mgr
```

Or delete the ArgoCD Application:

```bash
kubectl delete application slack-invite-mgr -n argocd
```

## Troubleshooting

### Pods not starting

Check if secrets are being created:

```bash
kubectl get externalsecrets -n slack-invite-mgr
kubectl get secrets -n slack-invite-mgr
```

### API connection issues

Ensure `web.env.apiUrl` is set to a browser-accessible URL, not an internal service URL.

### CronJob not running

Check the CronJob status:

```bash
kubectl get cronjobs -n slack-invite-mgr
kubectl get jobs -n slack-invite-mgr
```

## Architecture

```
                    ┌─────────────────┐
                    │    Traefik      │
                    │   IngressRoute  │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              │
    ┌─────────────┐  ┌─────────────┐       │
    │   Web (80)  │  │ Backend     │       │
    │   nginx     │  │ API (8080)  │       │
    └─────────────┘  └──────┬──────┘       │
                            │              │
                            ▼              │
                    ┌─────────────┐        │
                    │   Sheets    │◄───────┘
                    │   CronJob   │  (scheduled)
                    └─────────────┘
```

## License

See LICENSE file in the root of the repository.
