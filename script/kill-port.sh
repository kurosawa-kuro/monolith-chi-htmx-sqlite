#!/bin/bash

# ポート8080で実行中のプロセスを終了するスクリプト
PORT=8080

echo "ポート $PORT で実行中のプロセスを検索中..."

# ポート8080で実行中のプロセスのPIDを取得
PID=$(lsof -ti:$PORT)

if [ -z "$PID" ]; then
    echo "ポート $PORT で実行中のプロセスが見つかりませんでした。"
    exit 0
else
    echo "ポート $PORT で実行中のプロセス (PID: $PID) を終了します..."
    
    # プロセスを強制終了
    kill -9 $PID
    
    if [ $? -eq 0 ]; then
        echo "プロセス (PID: $PID) を正常に終了しました。"
    else
        echo "プロセスの終了に失敗しました。"
        exit 1
    fi
fi
