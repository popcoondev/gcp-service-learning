# Request Framing

## Need
- GCP案件に参画する PM/SE が、主要な GCP サービスを実際に動かしながら学べる教材リポジトリを用意する。

## Intent
- OpenAPI 駆動の Go サンプルを足場にしつつ、今後 Cloud Run / Spanner / Pub/Sub などへ拡張できる教材の土台を作る。

## Context
- 初期スコープは `POST /orders` と `GET /orders/{orderId}` のみ。
- 保存先はインメモリ。
- `domain/usecase` は GCP SDK 非依存。
- GCP 依存は adapter 層に閉じ込める。
- 将来的な Spanner 実装差し替えと Pub/Sub イベント拡張を可能にする。
- README は Mermaid 中心の学習教材として構成する。
