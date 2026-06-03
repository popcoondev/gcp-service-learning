# Decision Record: DEC-001

## Scope
- Record Format Version: 1.0.0
- Created At: 2026-06-03T10:00:00+09:00
- Canonical Markdown Path: .aof/decisions/DEC-001-bootstrap-learning-repo.md
- Scope: GCP実践学習教材リポジトリ初期化
- Stage: Requirements and design bootstrap
- Organization: GCP Service Learning Team

## Input
- Request: AOF を用い、OpenAPI 定義から Go の注文 API サンプルと GCP 学習教材リポジトリを作る
- Need: PM/SE がサービス全体像を学べる構造化教材が必要
- Intent: コード実装よりも、将来の GCP 学習導線を中心にした拡張可能な教材土台を整備する
- Context: 初期スコープはインメモリ保存、将来拡張は Cloud Run / Spanner / Pub/Sub ほか
- Existing Artifacts Reviewed: AOF v1.10.0 README, quickstart, governance template, decision record template
- Background or Prior Decisions: なし
- Clarifications or Assumptions: GitHub リポジトリ名は `gcp-service-learning`。公開範囲指定がないため public 前提で進める

## Options Considered
- Option A: コード中心の小さな API repo を先に作る
- Option B: 教材ドキュメント中心でコードは最小限にする
- Option C: コードと教材を同時に整備し、拡張ロードマップを最初から埋め込む

## Decision
- Selected Option: Option C
- Decision Summary: 初期価値として動く API と教材導線を同時に成立させる

## Governance
- Governance Model: Council of Three
- Decision Makers: Visionary, Builder, Guardian
- Governance Rule Applied: Majority vote with Guardian veto on architectural boundary violations
- Veto Used: No

## Rationale
- Why this option: 学習教材の主目的に沿いながら、将来の技術拡張を具体的に説明できる
- Why other options were not selected: A は教材価値が弱く、B は実装の足場が弱い
- Policy priorities applied: Learning Value > Architectural Extendability > Simplicity > Delivery Speed
- Policy tradeoffs accepted: 初期実装はインメモリに限定し、GCP 実接続は後続ロードマップへ送る

## Execution
- Actions: OpenAPI 作成、Go API 実装、README 図解、GCP 学習 docs 整備、GitHub 公開
- Expected Artifact: 学習用サンプルリポジトリ一式
- Expected Outcome: PM/SE がコードと構成図を往復しながら GCP 導入論点を学べる
- Completion Criteria: API 実装、README、docs、AOF 記録、GitHub push が揃う
- Success Criteria: 今後の Cloud Run / Spanner / Pub/Sub 導入演習にそのまま使える教材土台になっている
- Completion Approval Scope: Repository bootstrap approval
- Success Evaluation Scope: 学習教材レビュー

## Review
- Change Trigger: GCP 実サービス導入時
- Review Trigger: Cloud Run または Spanner を導入するタイミング
- Review Date or Condition: 次の value slice 着手時
- Re-open Conditions: ドメイン層に GCP SDK 依存が混入した場合、または教材としての導線が不足した場合
