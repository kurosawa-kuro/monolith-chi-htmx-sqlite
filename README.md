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

### ✅ カテゴリ機能
- カテゴリの追加／編集／削除
- カテゴリによる絞り込み表示

### ✅ UX向上
- htmx による非同期更新（部分的な再描画）
- Tailwind CSS で簡易UIスタイル

## 起動方法

1. 依存関係のインストール
```bash
go mod tidy
```

2. アプリケーションのビルド
```bash
go build -o todo-app main.go
```

3. アプリケーションの起動
```bash
./todo-app
```

4. ブラウザで http://localhost:8080 にアクセス

## ディレクトリ構成

```
/
├── main.go                 # アプリケーションエントリーポイント
├── go.mod                  # Go modules設定
├── src/
│   ├── templates/          # HTMLテンプレート
│   │   ├── layout.html
│   │   └── index.html
│   ├── static/             # 静的ファイル
│   │   └── style.css      # Tailwind CSS
│   ├── models/             # データモデル
│   │   └── models.go
│   ├── handlers/           # HTTPハンドラー
│   │   └── todo.go
│   └── db/                 # データベーススキーマ
│       └── schema.sql
├── template-admin/         # デザインテンプレート
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

### CSS開発（Tailwind）
```bash
cd template-admin/html
npm run dev    # 監視モード
npm run build  # ビルド
```