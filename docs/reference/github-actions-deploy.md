# GitHub Actions Deploy ガイド

このファイルは、[.github/workflows/cloud-run-deploy.yml](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.github/workflows/cloud-run-deploy.yml) を実際に使うときの前提をまとめたものです。

## 何をする workflow か

- Cloud Build で image を build
- Artifact Registry に push
- Cloud Run に deploy

## 事前に必要な secrets

- `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `GCP_DEPLOY_SERVICE_ACCOUNT`

## それぞれの意味

- `GCP_WORKLOAD_IDENTITY_PROVIDER`
  - GitHub Actions から Google Cloud へ federated login するための provider
- `GCP_DEPLOY_SERVICE_ACCOUNT`
  - GitHub Actions が impersonate する deploy 用 SA

## 教材としての意図

- 人間のローカル端末だけに依存せず、CI/CD からも配備できる形を見せる
- deploy 権限と runtime 権限が違うことを学べる

## 初心者のつまづきポイント

- `order-api-runtime@gcp-service-learning-498313.iam.gserviceaccount.com` は Cloud Run 実行用であり、GitHub Actions deploy 用ではない
- GitHub secrets が未設定だと workflow は認証前に失敗する

## ベテランからのアドバイス

- 最初はローカル deploy を通してから Actions に移ると切り分けが楽
- Workload Identity Federation を使うと、長期鍵を置かずに済む
