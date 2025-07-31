package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "./src/db/todo.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Create sample categories (INSERT OR IGNORE to avoid duplicates)
	categories := []string{"仕事", "プライベート", "学習", "家事", "趣味", "健康", "買い物", "旅行"}
	categoryIDs := make([]int, 0)

	for _, categoryName := range categories {
		result, err := db.Exec("INSERT OR IGNORE INTO categories (title) VALUES (?)", categoryName)
		if err != nil {
			log.Printf("Failed to insert category %s: %v", categoryName, err)
			continue
		}
		id, _ := result.LastInsertId()
		if id > 0 {
			categoryIDs = append(categoryIDs, int(id))
		}
	}

	// If no new categories were inserted, get existing ones
	if len(categoryIDs) == 0 {
		rows, err := db.Query("SELECT id FROM categories")
		if err != nil {
			log.Fatal("Failed to get existing categories:", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				log.Printf("Failed to scan category ID: %v", err)
				continue
			}
			categoryIDs = append(categoryIDs, id)
		}
	}

	if len(categoryIDs) == 0 {
		log.Fatal("No categories available")
	}

	// Create sample todos (現実的な量に削減)
	todoTitles := []string{
		"プロジェクトの要件定義を完了する",
		"週次レポートを提出する",
		"チームミーティングに参加する",
		"コードレビューを実施する",
		"ドキュメントを更新する",
		"テストケースを作成する",
		"バグ修正を行う",
		"パフォーマンス最適化を実施する",
		"セキュリティ監査を実行する",
		"デプロイメントを準備する",
		"ユーザーフィードバックを収集する",
		"新しい機能の設計を行う",
		"データベースのバックアップを実行する",
		"システム監視を設定する",
		"ログ分析を実施する",
	}

	// Insert todos with random categories and status
	for i, title := range todoTitles {
		// Random status
		status := "incomplete"
		if i%3 == 0 {
			status = "complete"
		}

		// Insert todo
		result, err := db.Exec("INSERT INTO todos (title, status) VALUES (?, ?)", title, status)
		if err != nil {
			log.Printf("Failed to insert todo %s: %v", title, err)
			continue
		}

		todoID, _ := result.LastInsertId()

		// Assign 1-3 random categories
		numCategories := (i % 3) + 1
		for j := 0; j < numCategories; j++ {
			categoryIndex := (i + j) % len(categoryIDs)
			_, err := db.Exec("INSERT INTO todo_category (todo_id, category_id) VALUES (?, ?)", todoID, categoryIDs[categoryIndex])
			if err != nil {
				log.Printf("Failed to assign category to todo: %v", err)
			}
		}
	}

	fmt.Printf("Successfully created %d categories and %d todos\n", len(categories), len(todoTitles))
}
