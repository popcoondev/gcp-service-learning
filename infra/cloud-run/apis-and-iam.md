# APIs と IAM の初期セット

このファイルは、Cloud Run 配備前に最低限そろえる API と IAM をまとめたものです。

## まず有効化する API

- `run.googleapis.com`
- `artifactregistry.googleapis.com`
- `cloudbuild.googleapis.com`
- `iam.googleapis.com`
- `logging.googleapis.com`
- `monitoring.googleapis.com`

## なぜ必要か

- `Cloud Run API`
  - サービス配備そのものに必要
- `Artifact Registry API`
  - Docker image の保存先に必要
- `Cloud Build API`
  - `gcloud builds submit` を使う場合に必要
- `IAM API`
  - Service Account と権限付与に必要
- `Logging / Monitoring API`
  - 運用可視化の初期土台に必要

## 最低限の Service Account

### 1. 実行用 SA

- `order-api-runtime@gcp-service-learning.iam.gserviceaccount.com`

想定 role:
- `roles/logging.logWriter`
- `roles/monitoring.metricWriter`

将来追加:
- `roles/pubsub.publisher`
- `roles/spanner.databaseUser`

UI での作業:

1. `IAM と管理 > サービス アカウント`
2. `サービス アカウントを作成`
3. 名前: `order-api-runtime`
4. ロール付与:
   - `Logs Writer`
   - `Monitoring Metric Writer`

### 2. デプロイ実行主体

人間ユーザーか CI/CD のどちらかです。

想定 role:
- `roles/run.admin`
- `roles/artifactregistry.writer`
- `roles/iam.serviceAccountUser`
- 必要に応じて `roles/cloudbuild.builds.editor`

UI での作業:

1. `IAM と管理 > IAM`
2. デプロイする人か CI 用 SA を選ぶ
3. 次のロールを付与する
   - `Cloud Run Admin`
   - `Artifact Registry Writer`
   - `Service Account User`
   - 必要に応じて `Cloud Build Editor`

## CLI で付与する場合の例

```bash
gcloud projects add-iam-policy-binding gcp-service-learning \
  --member="serviceAccount:order-api-runtime@gcp-service-learning.iam.gserviceaccount.com" \
  --role="roles/logging.logWriter"

gcloud projects add-iam-policy-binding gcp-service-learning \
  --member="serviceAccount:order-api-runtime@gcp-service-learning.iam.gserviceaccount.com" \
  --role="roles/monitoring.metricWriter"
```

デプロイ主体への付与例:

```bash
gcloud projects add-iam-policy-binding gcp-service-learning \
  --member="user:YOUR_EMAIL@example.com" \
  --role="roles/run.admin"
```

## PM 観点での注意

- 「誰が deploy できるか」と「アプリが何にアクセスできるか」は分けて考える
- 最初は広めの権限で立ち上げても、教材ではあとで絞る前提を明記する

## 開発者観点での注意

- `allow unauthenticated` は学習用に分かりやすいが、本番要件とは分けて考える
- 権限不足時のエラーは Cloud Run, Cloud Build, Artifact Registry のどこで落ちたか切り分ける
- `roles/editor` でまとめて逃げず、最低限ロールを積み上げて理解するほうが教材価値が高い
