# Spanner

## 何のために使うか

インメモリ保存を本番向け永続化へ置き換え、可用性と整合性を高めるために使います。

## どのような構成で使うか

- `domain.OrderRepository` を実装する Spanner adapter を追加する
- `application` と `domain` は変更せず、adapter 差し替えで導入する

## このサンプル内でどこに登場するか

- `internal/domain/order_repository.go` が差し替えポイントです
- README に future adapter として記載しています

## PM視点で何が難しいか

- なぜ RDB ではなく Spanner なのかの説明
- 可用性要求とコストの妥当性
- スキーマ変更計画

## 開発者視点で何が難しいか

- トランザクションモデル理解
- 主キー設計
- エミュレータと本番差分
