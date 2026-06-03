# Cloud Run 配備テンプレート

このディレクトリは、Step 2 で学ぶ Cloud Run 配備の最小テンプレートです。

## 含まれるもの

- [service.yaml](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/infra/cloud-run/service.yaml)
  - Cloud Run Service 定義の教材用テンプレート

## この repo で使う値

- `Project ID`: `gcp-service-learning`
- `Region`: `asia-northeast1`
- `Artifact Registry repository`: `gcp-service-learning`
- 想定 image URI:
  - `asia-northeast1-docker.pkg.dev/gcp-service-learning/gcp-service-learning/order-api:latest`

## いまの段階で見るポイント

- `serviceAccountName`
  - 実行主体を誰にするか
- `containerConcurrency`
  - 同時実行数をどこまで許容するか
- `timeoutSeconds`
  - 同期 API のタイムアウト方針
- `startupProbe` / `livenessProbe`
  - `/healthz` をどう運用に使うか

## 次に UI で確認すること

- Cloud Run と Artifact Registry がどちらも `asia-northeast1` になっているか
- `gcp-service-learning` repository が Docker 形式で作成されているか
- 実行用 Service Account を `order-api-runtime@gcp-service-learning.iam.gserviceaccount.com` として用意するか

## 補足

`Project ID` と `Artifact Registry repository` がどちらも `gcp-service-learning` でも問題ありません。  
教材としては、むしろ命名ルールが単純で追いやすい構成です。
