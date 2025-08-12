# jplaw2epub-server

Web API server that converts Japanese Law XML to EPUB format and provides GraphQL API for law data queries

## Features

- POST `/convert` - Accepts XML data and returns EPUB file
- GET `/epubs/{law_id}` - Get EPUB file for a specific law ID
- POST `/graphql` - GraphQL API endpoint for law queries
- GET `/graphiql` - Interactive GraphQL playground
- GET `/health` - Health check endpoint

## Local Development

```bash
# Default (automatically finds available port)
go run main.go

# Specify port
go run main.go -port 8080

# Using environment variable
PORT=3000 go run main.go
```

## API Usage

### Convert XML to EPUB
```bash
# Convert XML file to EPUB
curl -X POST \
  -H "Content-Type: application/xml" \
  --data-binary @law.xml \
  http://localhost:8090/convert \
  -o output.epub
```

### Get EPUB by Law ID
```bash
# Download EPUB for a specific law (e.g., Radio Act)
curl http://localhost:8090/epubs/325AC0000000131 -o radio_act.epub
```

### GraphQL Queries
```bash
# Search laws by title
curl -X POST http://localhost:8090/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ laws(lawTitle: \"電波\", limit: 5) { totalCount laws { lawInfo { lawId lawNum } revisionInfo { lawTitle } } } }"}'

# Get law revisions
curl -X POST http://localhost:8090/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ revisions(lawId: \"325AC0000000131\") { lawInfo { lawNum } revisions { amendmentLawTitle amendmentEnforcementDate } } }"}'

# Keyword search
curl -X POST http://localhost:8090/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ keyword(keyword: \"無線\", limit: 3) { totalCount items { lawInfo { lawId } revisionInfo { lawTitle } sentences { text } } } }"}'
```

### Health check
```bash
curl http://localhost:8090/health
```

## Docker Build

```bash
# Build from project root
docker build -t jplaw2epub-server -f cmd/jplaw2epub-server/Dockerfile .

# Run
docker run -p 8080:8080 jplaw2epub-server
```

## Deploy to Google Cloud Run

### Prerequisites

- Google Cloud SDK installed
- Project configured
- Cloud Run API enabled

### Manual Deployment

```bash
# Execute from project root

# 1. Build container image
gcloud builds submit \
  --tag gcr.io/YOUR_PROJECT_ID/jplaw2epub-server \
  --file cmd/jplaw2epub-server/Dockerfile

# 2. Deploy to Cloud Run
gcloud run deploy jplaw2epub-server \
  --image gcr.io/YOUR_PROJECT_ID/jplaw2epub-server \
  --region asia-northeast1 \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --memory 512Mi \
  --max-instances 10 \
  --min-instances 0 \
  --timeout 60
```

### Automated Deployment with Cloud Build

```bash
# Execute from project root
gcloud builds submit \
  --config cmd/jplaw2epub-server/cloudbuild.yaml \
  --substitutions=_REGION=asia-northeast1
```

### Continuous Deployment with GitHub

1. Create Cloud Build trigger
```bash
gcloud builds triggers create github \
  --repo-name=jplaw2epub \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cmd/jplaw2epub-server/cloudbuild.yaml
```

2. Automatic deployment will run on push to main branch

## Environment Variables

- `PORT` - Server listening port (default: auto-select)

## Recommended Cloud Run Settings

- **Region**: asia-northeast1 (Tokyo)
- **Memory**: 512Mi
- **Max instances**: 10
- **Min instances**: 0 (allows cold start)
- **Timeout**: 60 seconds
- **Concurrency**: 1000 (default)

## Troubleshooting

### Out of Memory Error
For large XML files, increase Cloud Run memory:
```bash
gcloud run services update jplaw2epub-server --memory 1Gi
```

### Timeout Error
For longer processing times, extend timeout:
```bash
gcloud run services update jplaw2epub-server --timeout 300
```