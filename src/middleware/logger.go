package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// FilteredLogger Chrome DevToolsのリクエストをフィルタリングするログミドルウェア
func FilteredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// レスポンスをラップしてステータスコードを取得
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Chrome DevToolsのリクエストをフィルタリング
		if shouldLogRequest(r) {
			log.Printf("%s %s %d %v", r.Method, r.URL.Path, ww.statusCode, duration)
		}
	})
}

// shouldLogRequest リクエストをログに出力するかどうかを判定
func shouldLogRequest(r *http.Request) bool {
	// Chrome DevToolsのリクエストを除外
	if strings.Contains(r.URL.Path, "/.well-known/appspecific/com.chrome.devtools.json") {
		return false
	}

	// その他の不要なリクエストも除外
	if strings.Contains(r.URL.Path, "/favicon.ico") {
		return false
	}

	return true
}

// responseWriter レスポンスライターのラッパー
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader ステータスコードを記録
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write レスポンスを書き込み
func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
