# Cloud Run 配備テンプレート

このディレクトリは、Step 2 で学ぶ Cloud Run 配備の最小テンプレートです。

## 含まれるもの

- [service.yaml](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/infra/cloud-run/service.yaml)
  - Cloud Run Service 定義の教材用テンプレート

## この repo で使う値

- `Project ID`: `gcp-service-learning-498313`
- `Region`: `asia-northeast1`
- `Artifact Registry repository`: `gcp-service-learning`
- 想定 image URI:
  - `asia-northeast1-docker.pkg.dev/gcp-service-learning-498313/gcp-service-learning/order-api:latest`

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
- 実行用 Service Account を `order-api-runtime@gcp-service-learning-498313.iam.gserviceaccount.com` として用意するか

## 初回デプロイ後の確認

1. `Cloud Run > order-api`
   - URL が払い出されているか
2. `/healthz`
   - `{"status":"ok","service":"order-api"}` を返すか
3. `/orders`
   - POST が `201` で返るか
4. `ログ エクスプローラ`
   - Cloud Run のリクエストログが出ているか

## 補足

今回は `Project ID` が `gcp-service-learning-498313`、`Artifact Registry repository` が `gcp-service-learning` です。  
同じ名前でそろっていなくても問題ありません。`Project` と `Repository` は別リソースです。
