# Goユニットテスト 導入対応

このディレクトリには、Goアプリケーションのユニットテストと統合テストが含まれています。

## テスト構造

```
src/tests/
├── helpers/           # テストヘルパーとモック
├── models/           # モデルレイヤーのテスト
├── services/         # サービスレイヤーのテスト
├── handlers/         # ハンドラーレイヤーのテスト
├── integration/      # 統合テスト
└── README.md         # このファイル
```

## 依存関係

テストには以下の依存関係が必要です：

- `github.com/stretchr/testify` - テストアサーション
- `github.com/DATA-DOG/go-sqlmock` - SQLモック
- `github.com/mattn/go-sqlite3` - SQLiteドライバー（統合テスト用）

## テストの実行

### 全テストの実行
```bash
make test-all-go
```

### ユニットテストのみ
```bash
make test-unit
```

### 統合テストのみ
```bash
make test-integration
```

### カバレッジ付きテスト
```bash
make test-coverage
```

### 特定のテストパッケージ
```bash
go test ./src/tests/models/...
go test ./src/tests/services/...
go test ./src/tests/handlers/...
go test ./src/tests/integration/...
```

## テストの種類

### 1. ユニットテスト

#### Models テスト
- `TodoRepository` の各メソッドのテスト
- `CategoryRepository` の各メソッドのテスト
- SQLモックを使用したデータベース操作のテスト

#### Services テスト
- `TodoService` のビジネスロジックのテスト
- バリデーションとエラーハンドリングのテスト
- モックリポジトリを使用したテスト

#### Handlers テスト
- HTTPハンドラーのテスト
- リクエスト/レスポンスのテスト
- モックサービスを使用したテスト

### 2. 統合テスト

- 実際のSQLiteデータベースを使用
- エンドツーエンドのワークフローテスト
- 一時的なデータベースファイルを使用

## テストヘルパー

### TestDB
SQLモックを使用したテストデータベースを提供します。

```go
testDB := helpers.NewTestDB(t)
defer testDB.Close()
```

### CreateTempDB
統合テスト用の一時的なSQLiteデータベースを作成します。

```go
db, dbPath := helpers.CreateTempDB(t)
defer db.Close()
defer helpers.CleanupTempDB(dbPath)
```

### MockTime
テスト用の固定時刻を提供します。

```go
mockTime := helpers.MockTime()
```

## テストの書き方

### 基本的なテスト構造

```go
func TestFunctionName(t *testing.T) {
    // セットアップ
    testDB := helpers.NewTestDB(t)
    defer testDB.Close()
    
    // テスト実行
    result, err := functionToTest()
    
    // アサーション
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### モックの使用

```go
// モックリポジトリの作成
mockRepo := &MockTodoRepository{
    todos: []models.Todo{...},
    err:   nil,
}

// モックサービスの作成
mockService := &MockTodoService{
    todos: []models.Todo{...},
    err:   nil,
}
```

### SQLモックの使用

```go
// クエリの期待値を設定
testDB.Mock.ExpectQuery("SELECT.*FROM todos").
    WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
        AddRow(1, "Test Todo"))

// テスト実行
todos, err := repo.GetAllTodos()

// 期待値の検証
require.NoError(t, err)
require.NoError(t, testDB.Mock.ExpectationsWereMet())
```

## ベストプラクティス

1. **テストの独立性**: 各テストは独立して実行できるようにする
2. **クリーンアップ**: テスト後にリソースを適切にクリーンアップする
3. **モックの使用**: 外部依存関係はモックを使用する
4. **エラーケース**: 正常系と異常系の両方をテストする
5. **アサーション**: 適切なアサーションを使用して結果を検証する

## トラブルシューティング

### SQLiteドライバーのエラー
統合テストで `sql: unknown driver "sqlite3"` エラーが発生する場合：

```go
import _ "github.com/mattn/go-sqlite3"
```

### SQLモックのエラー
SQLモックでクエリが一致しない場合：
- 実際のSQLクエリを確認
- 正規表現パターンを調整
- 引数の数を確認

### テストの失敗
テストが失敗する場合：
- テストデータの準備を確認
- モックの設定を確認
- アサーションの期待値を確認 