# Pulumi Go GCP Cloud SQL

This Pulumi project provisions a Google Cloud SQL instance (MySQL or PostgreSQL) with a database and user.

## Resources Created

- **Cloud SQL Instance**: MySQL 8.0 or PostgreSQL 14/15 instance
- **Database**: Initial database within the instance
- **User**: Database user with generated password
- **Automated Backups**: Enabled by default

## Prerequisites

- Pulumi CLI installed
- Go 1.21+
- GCP credentials configured
- A GCP project with billing enabled

## Configuration

Copy `Pulumi.dev.yaml.example` to `Pulumi.dev.yaml` and update values:

```yaml
config:
  gcp:project: your-project-id
  gcp:region: us-central1
  cloudSql:instanceName: my-sql-instance
  cloudSql:databaseName: mydatabase
  cloudSql:databaseUser: myuser
  cloudSql:databaseVersion: MYSQL_8_0
  cloudSql:tier: db-f1-micro
```

### Config Options

| Key | Description | Default |
|-----|-------------|---------|
| `gcp:project` | GCP Project ID | **Required** |
| `gcp:region` | GCP Region (Free Tier: us-west1, us-central1, us-east1) | `us-central1` |
| `cloudSql:instanceName` | Cloud SQL instance name | **Required** |
| `cloudSql:databaseName` | Initial database name | **Required** |
| `cloudSql:databaseUser` | Initial database user | **Required** |
| `cloudSql:databaseVersion` | Engine version (MYSQL_8_0, POSTGRES_14, POSTGRES_15) | `MYSQL_8_0` |
| `cloudSql:tier` | Machine type (Free Tier: db-f1-micro) | `db-f1-micro` |

## Deployment

```bash
# Install dependencies
go mod tidy

# Preview changes
pulumi preview

# Deploy
pulumi up
```

## Outputs

| Output | Description |
|--------|-------------|
| `instanceName` | The name of the Cloud SQL instance |
| `connectionName` | The connection name (project:region:instance) |
| `publicIpAddress` | The public IP address of the instance |
| `databaseName` | The name of the initial database |
| `databaseUser` | The name of the initial database user |
| `databasePassword` | The generated password (sensitive) |

## Free Tier Notes

Cloud SQL Free Tier includes:
- 1 db-f1-micro instance (shared-core, 0.6 GB RAM)
- 30 GB HDD storage
- 1 GB network egress/month (within Americas)
- Available in: us-west1, us-central1, us-east1

See [Cloud SQL Pricing](https://cloud.google.com/sql/pricing) for details.

## Connecting

```bash
# Using Cloud SQL Proxy (recommended)
cloud-sql-proxy <connection-name>

# Or direct connection (requires authorized networks)
mysql -h <public-ip> -u <user> -p
```

## Cleanup

```bash
pulumi destroy
```