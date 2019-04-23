#!/bin/bash
set -euo pipefail

TAG="${TAG:?specify TAG of the pushed image}"
PROJECT="${PROJECT:-ahmet-personal-api}"
IMAGE="${REPO:-gcr.io/$PROJECT/goodbye}"
REGION="${REGION:-us-central1}"
BERGLAS_SECRETS_BUCKET="${BERGLAS_SECRETS_BUCKET:-secrets-goodbye}"

echo >&2 "Deploying goodbye service."
set -x
gcloud alpha run deploy "goodbye" \
	--project="${PROJECT}" \
	--image="${IMAGE}:${TAG}" \
	--allow-unauthenticated \
	--region="${REGION}" \
	--memory=512Mi \
	--concurrency=1 \
	--set-env-vars="CONSUMER_KEY=berglas://${BERGLAS_SECRETS_BUCKET}/CONSUMER_KEY,\
CONSUMER_SECRET=berglas://${BERGLAS_SECRETS_BUCKET}/CONSUMER_SECRET,\
ACCESS_TOKEN=berglas://${BERGLAS_SECRETS_BUCKET}/ACCESS_TOKEN,\
ACCESS_TOKEN_SECRET=berglas://${BERGLAS_SECRETS_BUCKET}/ACCESS_TOKEN_SECRET"
set +x

echo >&2 "Done."
