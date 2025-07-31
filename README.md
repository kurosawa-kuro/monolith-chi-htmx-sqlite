# monolith-chi-htmx-sqlite

シンプルなToDo管理アプリケーション（カテゴリ機能付き）

## 技術スタック

- **バックエンド**: Go + [chi](https://github.com/go-chi/chi)
- **フロントエンド**: HTML + [htmx](https://htmx.org/) + [Tailwind CSS](https://tailwindcss.com/)
- **データベース**: SQLite
- **テンプレート**: Go標準 `html/template`

## 機能

### ✅ Todo機能
- 一覧表示（カテゴリと紐付け含む）
- 新規作成／編集／削除
- ステータス更新（未完了 ↔ 完了）
- ページネーション機能（ページサイズ選択可能）

### ✅ カテゴリ機能
- カテゴリの追加／編集／削除
- カテゴリによる絞り込み表示

### ✅ UX向上
- htmx による非同期更新（部分的な再描画）
- Tailwind CSS で簡易UIスタイル

## 起動方法

### 方法1: 直接実行
```bash
# 依存関係のインストール
go mod tidy

# アプリケーションの実行
go run src/main.go
```

### 方法2: ビルドして実行
```bash
# 依存関係のインストール
go mod tidy

# アプリケーションのビルド
make build

# アプリケーションの実行
./bin/todo-app
```

### 方法3: Makefileを使用
```bash
# 開発モード（ビルド + 実行）
make dev

# または直接実行
make run
```

4. ブラウザで http://localhost:8080 にアクセス

## ディレクトリ構成

```
/
├── src/
│   ├── main.go             # アプリケーションエントリーポイント
│   ├── templates/          # HTMLテンプレート
│   │   ├── layout.html
│   │   └── index.html
│   ├── static/             # 静的ファイル
│   │   └── style.css      # Tailwind CSS
│   ├── models/             # データモデル
│   │   └── models.go
│   ├── handlers/           # HTTPハンドラー
│   │   └── todo.go
│   └── db/                 # データベース関連
│       ├── schema.sql      # データベーススキーマ
│       └── todo.db         # SQLiteデータベースファイル
├── bin/                    # ビルド成果物
│   └── todo-app           # 実行ファイル
├── go.mod                  # Go modules設定
├── go.sum                  # 依存関係チェックサム
├── Makefile                # ビルド・実行スクリプト
└── docs/                   # 仕様書
    └── 仕様書.md
```

## データベーススキーマ

```text
+---------+          +---------------+          +--------------+
|  todos  |          | todo_category |          |  categories  |
+---------+          +---------------+          +--------------+
| id      |<--+    +--| todo_id       |    +--->| id           |
| title   |   |    |  | category_id   |<---+     | title        |
| status  |   |    +---------------+            +--------------+
+---------+
```

- `todos`: タスク本体（ID、タイトル、ステータス）
- `categories`: カテゴリ（ID、タイトル）
- `todo_category`: 多対多リレーション（中間テーブル）

## 開発

### Makefileコマンド
```bash
make help      # 利用可能なコマンドを表示
make build     # アプリケーションをビルド
make run       # アプリケーションを実行
make dev       # ビルドして実行
make clean     # ビルド成果物を削除
make fmt       # コードをフォーマット
make test      # テストを実行
make db-init   # データベースを初期化
make seed-data # サンプルデータを追加（ページネーション機能のテスト用）
```

### CSS開発（Tailwind）
```bash
# 現在は事前ビルドされたCSSを使用
# 必要に応じてTailwind CLIで再ビルド
```