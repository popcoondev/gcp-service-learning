# Cloud Run 配備テンプレート

このディレクトリは、Step 2 で学ぶ Cloud Run 配備の最小テンプレートです。

## 含まれるもの

- [service.yaml](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/infra/cloud-run/service.yaml)
  - Cloud Run Service 定義の教材用テンプレート

## いまの段階で見るポイント

- `serviceAccountName`
  - 実行主体を誰にするか
- `containerConcurrency`
  - 同時実行数をどこまで許容するか
- `timeoutSeconds`
  - 同期 API のタイムアウト方針
- `startupProbe` / `livenessProbe`
  - `/healthz` をどう運用に使うか

## 置換が必要な値

- `PROJECT_ID`
- `REGION`
- `REPOSITORY`

これらは実際の GCP 環境に合わせてユーザーと確認してから埋める前提です。
