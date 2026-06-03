# gcp-service-learning

OpenAPI 定義から始める Go 注文 API と、GCP 実践学習教材を同居させたリポジトリです。  
このリポジトリの主目的はアプリ開発そのものではなく、GCP 案件に参画する PM/SE が主要サービスを構成・責務・運用観点まで含めて学べる教材土台を作ることです。

## この教材で学ぶこと

- OpenAPI を起点に API の責務を整理する方法
- Go のヘキサゴナルアーキテクチャで GCP 依存を adapter に閉じ込める方法
- Cloud Run / Spanner / Pub/Sub へ段階的に拡張する設計の考え方
- ネットワーク、IAM、運用監視を実装と一緒に読む方法
- PM/SE が設計レビューで見るべき論点

## 初期スコープ

- `POST /orders`
- `GET /orders/{orderId}`
- 保存先はインメモリ
- OpenAPI 定義あり
- Go
- ヘキサゴナルアーキテクチャ

## クイックスタート

1. OpenAPI 定義を確認する: [api/openapi.yaml](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/api/openapi.yaml)
2. サーバーを起動する: `go run ./cmd/server`
3. 注文を作成する:

```bash
curl -X POST http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{
    "customerId": "cust-001",
    "items": [
      { "productId": "book-001", "quantity": 2 }
    ]
  }'
```

4. 取得する:

```bash
curl http://localhost:8080/orders/<orderId>
```

## リポジトリ構成

```text
.
├── api/                         # OpenAPI 定義
├── cmd/server/                  # エントリポイント
├── internal/domain/             # ドメインモデルと port
├── internal/application/        # ユースケース
├── internal/adapters/http/      # HTTP adapter
├── internal/adapters/repository # 保存先 adapter
├── docs/                        # GCP 学習教材
└── .aof/                        # AOF の framing / governance / decisions
```

## システムアーキテクチャ図

```mermaid
flowchart LR
    user["API Client"] --> run["Order API on Cloud Run<br/>Current: local process"]
    run --> http["HTTP Adapter"]
    http --> app["Usecase Layer"]
    app --> repoport["OrderRepository Port"]
    repoport --> mem["In-Memory Repository"]
    repoport -. future .-> spanner["Spanner Repository Adapter"]
    app -. future .-> eventport["Event Publisher Port"]
    eventport -. future .-> pubsub["Pub/Sub Adapter"]
    run -. telemetry .-> logging["Cloud Logging"]
    run -. metrics .-> monitoring["Cloud Monitoring"]
```

## レイヤ構成図

```mermaid
flowchart TB
    subgraph outside["Outside World"]
        client["Client / Operator / CI"]
        gcp["GCP Services"]
    end

    subgraph adapter["Adapters"]
        httpa["HTTP Adapter"]
        repoa["Repository Adapter"]
        eventa["Event Adapter (future)"]
    end

    subgraph app["Application"]
        usecase["Order Usecase"]
    end

    subgraph domain["Domain"]
        entity["Order Entity"]
        port["Repository / Publisher Port"]
    end

    client --> httpa
    httpa --> usecase
    usecase --> entity
    usecase --> port
    repoa --> port
    eventa --> port
    gcp --> repoa
    gcp --> eventa
```

## データフロー図

```mermaid
flowchart LR
    req["Create Order Request"] --> validate["HTTP Validation"]
    validate --> uc["CreateOrder Usecase"]
    uc --> domain["Order Aggregate Creation"]
    domain --> save["OrderRepository.Save"]
    save --> store["In-Memory Store"]
    uc -. future .-> publish["Publish OrderCreated Event"]
    publish -. future .-> pubsub["Pub/Sub Topic"]
    store --> res["Order Response"]
```

## API シーケンス図

```mermaid
sequenceDiagram
    participant C as Client
    participant H as HTTP Adapter
    participant U as Order Usecase
    participant R as Repository

    C->>H: POST /orders
    H->>H: JSON decode / validate
    H->>U: CreateOrder(command)
    U->>U: Build Order entity
    U->>R: Save(order)
    R-->>U: saved
    U-->>H: OrderDTO
    H-->>C: 201 Created
```

## GCP リソース関連図

```mermaid
flowchart TB
    user["Developer / PM / SE"] --> github["GitHub Repository"]
    github --> cicd["CI/CD Pipeline"]
    cicd --> ar["Artifact Registry"]
    ar --> cr["Cloud Run Service"]
    cr --> ps["Pub/Sub Topic (future)"]
    cr --> sp["Spanner (future)"]
    cr --> gcs["GCS (future)"]
    ps --> wf["Workflows (future)"]
    cr --> log["Cloud Logging"]
    cr --> mon["Cloud Monitoring"]
    cr --> otel["OpenTelemetry Collector (future)"]
    ps --> bq["BigQuery (future analytics)"]
```

## IAM 関係図

```mermaid
flowchart LR
    dev["Developer / Student"] --> saDeployer["Deployer Service Account"]
    saDeployer --> cloudrunAdmin["roles/run.admin"]
    saDeployer --> artifactWriter["roles/artifactregistry.writer"]
    saDeployer --> iamUser["roles/iam.serviceAccountUser"]

    runtime["Cloud Run Runtime SA"] --> spannerRole["roles/spanner.databaseUser (future)"]
    runtime --> pubsubRole["roles/pubsub.publisher (future)"]
    runtime --> loggingRole["roles/logging.logWriter"]
    runtime --> metricRole["roles/monitoring.metricWriter"]

    analyst["Data Analyst"] --> bqRole["roles/bigquery.dataViewer (future)"]
    operator["Operator / SRE"] --> monitorRole["roles/monitoring.viewer"]
```

## 学習観点

### このサンプルで今すぐ学べること

- OpenAPI と実装の対応関係
- レイヤ分離の意味
- repository port による保存先差し替え準備
- 将来のイベント発行 port を先に意識した設計

### 将来追加で深掘ること

- Cloud Run へのデプロイ
- Pub/Sub による非同期イベント化
- Spanner への永続化
- Cloud Logging / Monitoring / OpenTelemetry による運用可視化
- Terraform による IaC
- GKE, PAM, VPC Service Controls を使った高度化

## 設計原則

- `domain/usecase` は GCP SDK に依存しない
- GCP SDK は adapter 層に限定する
- repository は interface 化する
- Spanner 実装へ差し替え可能にする
- Pub/Sub イベント発行へ拡張可能にする

## AOF の適用

このリポジトリでは AOF を「実装フレームワーク」ではなく「学習教材を成長させるための判断記録フレーム」として使います。

- Request Framing: [.aof/context/request-framing.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.aof/context/request-framing.md)
- Governance: [.aof/governance/governance.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.aof/governance/governance.md)
- Decision Record: [.aof/decisions/DEC-001-bootstrap-learning-repo.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.aof/decisions/DEC-001-bootstrap-learning-repo.md)

### AOF Runtime / CLI

この repo は AOF v1.10.0 の runtime / CLI を `tools/aof-runtime/` に同梱し、project-local な `.aof/` state を持つ構成にしています。

```bash
npm run aof -- run "AOF runtime でこの repo を運用したい" --project .
npm run aof -- answer --session ./.aof/sessions/<SESSION>.json --response "..."
npm run aof -- council --session ./.aof/sessions/<SESSION>.json --stage planning --project .
npm run aof -- council-exec --session ./.aof/sessions/<SESSION>.json --stage planning --project . --invoke-model --provider mock
```

### Visibility Service

live `.aof/` state から visibility 用 JSON を生成して viewer を起動できます。

```bash
npm run aof:visibility:build
npm run aof:visibility:serve
```

生成先:

- [.aof/artifacts/visibility/status-card.json](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.aof/artifacts/visibility/status-card.json)
- [.aof/artifacts/visibility/timeline-feed.json](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.aof/artifacts/visibility/timeline-feed.json)
- [.aof/artifacts/visibility/flow-snapshot.json](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/.aof/artifacts/visibility/flow-snapshot.json)

## 学習ロードマップ

```mermaid
flowchart LR
    s1["Step 1<br/>Local + In-Memory + OpenAPI"] --> s2["Step 2<br/>Cloud Run"]
    s2 --> s3["Step 3<br/>Pub/Sub Eventing"]
    s3 --> s4["Step 4<br/>Spanner Persistence"]
    s4 --> s5["Step 5<br/>Logging / Monitoring / OTel"]
    s5 --> s6["Step 6<br/>Terraform"]
    s6 --> s7["Step 7<br/>GKE / IAM / PAM / VPC SC"]
```

## docs ガイド

- 全体ガイド: [docs/README.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/README.md)
- ネットワークとセキュリティ観点: [docs/network-and-security.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/network-and-security.md)
- Cloud Run: [docs/services/cloud-run.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/services/cloud-run.md)
- Pub/Sub: [docs/services/pubsub.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/services/pubsub.md)
- GCS: [docs/services/gcs.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/services/gcs.md)
- BigQuery: [docs/services/bigquery.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/services/bigquery.md)
- Spanner: [docs/services/spanner.md](/Users/mn/Documents/Codex/2026-06-03/aof-go-openapi-api-github-aof/docs/services/spanner.md)