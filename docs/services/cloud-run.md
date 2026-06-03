# Cloud Run

## 何のために使うか

HTTP ベースの注文 API をサーバーレスに実行し、最小運用で本番相当のデプロイ体験を得るために使います。

## どのような構成で使うか

- コンテナ化した Go API を Cloud Run に配置する
- Artifact Registry にイメージを保存する
- 将来的に VPC Connector や Secret Manager と連携する
- 初期教材では `/healthz` を使って起動確認し、In-Memory 実装の限界も合わせて学ぶ

## このサンプル内でどこに登場するか

- 現在の `cmd/server` がそのまま Cloud Run へ載る前提です
- README のシステムアーキテクチャ図で実行基盤として登場します
- [Dockerfile](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/Dockerfile) でコンテナ化の入口を示しています
- [infra/cloud-run/service.yaml](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/infra/cloud-run/service.yaml) で Cloud Run Service の教材用テンプレートを示しています
- [docs/steps/step-02-cloud-run.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/steps/step-02-cloud-run.md) で Step 単位の読み解きを用意しています
- 今回の教材値は `Project ID: gcp-service-learning`, `Region: asia-northeast1`, `Artifact Registry repository: gcp-service-learning` です

## PM視点で何が難しいか

- スケール設定とコストのバランス
- 同期 API のタイムアウト設計
- 環境差分とリリース承認フロー
- 「Cloud Run に載ったら本番になる」と誤解せず、データ永続化や認証方式の未整備を別論点として管理すること

## 開発者視点で何が難しいか

- コンテナのビルド再現性
- リクエストタイムアウトとコールドスタート考慮
- ローカル実行との差分吸収
- In-Memory Repository がインスタンス再起動や水平スケールで壊れる前提を理解したうえで教材として使うこと
