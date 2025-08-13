# GitHub Actions Setup Guide

This directory contains GitHub Actions workflows for automatically deploying jplaw2epub-server to Google Cloud Run.

## Workflows

### deploy-server.yml
Deploys the server to Google Cloud Run when changes are merged to the master branch.

### test-server.yml
Runs tests and builds the Docker image on pull requests.

## Required GitHub Secrets

You need to configure the following secrets in your GitHub repository settings (Settings → Secrets and variables → Actions):

### Option 1: Workload Identity Federation (Recommended)

More secure as it doesn't require storing service account keys.

1. **GCP_PROJECT_ID**: Your Google Cloud project ID
2. **WIF_PROVIDER**: Workload Identity Federation provider 
   - Format: `projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_NAME/providers/PROVIDER_NAME`
3. **WIF_SERVICE_ACCOUNT**: Service account email for deployment
   - Format: `SERVICE_ACCOUNT_NAME@PROJECT_ID.iam.gserviceaccount.com`

### Option 2: Service Account Key (Simpler setup)

1. **GCP_PROJECT_ID**: Your Google Cloud project ID
2. **GCP_SA_KEY**: Service account JSON key with necessary permissions

To use this option, uncomment the alternative authentication section in `deploy-server.yml`.

## Setting up Workload Identity Federation

1. Enable required APIs:
```bash
gcloud services enable \
  iamcredentials.googleapis.com \
  cloudresourcemanager.googleapis.com \
  sts.googleapis.com
```

2. Create a Workload Identity Pool:
```bash
gcloud iam workload-identity-pools create "github-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="GitHub Actions Pool"
```

3. Create a Workload Identity Provider:
```bash
gcloud iam workload-identity-pools providers create-oidc "github-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-pool" \
  --display-name="GitHub Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"
```

4. Create a service account:
```bash
gcloud iam service-accounts create github-deploy \
  --display-name="GitHub Deploy Account"
```

5. Grant necessary permissions:
```bash
# Cloud Run Admin
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:github-deploy@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/run.admin"

# Service Account User
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:github-deploy@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# Storage Admin (for container images)
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
  --member="serviceAccount:github-deploy@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/storage.admin"
```

6. Allow impersonation from Workload Identity Pool:
```bash
gcloud iam service-accounts add-iam-policy-binding \
  "github-deploy@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-pool/attribute.repository/YOUR_GITHUB_USERNAME/jplaw2epub"
```

7. Get the provider and service account values:
```bash
# Get WIF_PROVIDER value
echo "projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-pool/providers/github-provider"

# Get WIF_SERVICE_ACCOUNT value
echo "github-deploy@${PROJECT_ID}.iam.gserviceaccount.com"
```

## Setting up with Service Account Key (Alternative)

1. Create a service account:
```bash
gcloud iam service-accounts create github-deploy \
  --display-name="GitHub Deploy Account"
```

2. Grant permissions (same as step 5 above)

3. Create and download the key:
```bash
gcloud iam service-accounts keys create key.json \
  --iam-account=github-deploy@${PROJECT_ID}.iam.gserviceaccount.com
```

4. Copy the contents of `key.json` and add it as the `GCP_SA_KEY` secret in GitHub

## Required Google Cloud APIs

Make sure these APIs are enabled in your project:
- Cloud Run API
- Container Registry API
- Cloud Build API (if using Cloud Build)

```bash
gcloud services enable \
  run.googleapis.com \
  containerregistry.googleapis.com \
  cloudbuild.googleapis.com
```

## Manual Deployment

If you need to deploy manually:

```bash
# From project root
gcloud builds submit \
  --tag gcr.io/${PROJECT_ID}/jplaw2epub-server \
  --file cmd/jplaw2epub-server/Dockerfile

gcloud run deploy jplaw2epub-server \
  --image gcr.io/${PROJECT_ID}/jplaw2epub-server \
  --region asia-northeast1 \
  --platform managed \
  --allow-unauthenticated
```