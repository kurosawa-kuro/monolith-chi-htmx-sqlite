#!/bin/bash

# ポート番号を環境変数から取得、デフォルトは8080
PORT=${PORT:-8080}

# 色付きの出力関数
print_info() {
    echo -e "\033[36m[INFO]\033[0m $1"
}

print_success() {
    echo -e "\033[32m[SUCCESS]\033[0m $1"
}

print_warning() {
    echo -e "\033[33m[WARNING]\033[0m $1"
}

print_error() {
    echo -e "\033[31m[ERROR]\033[0m $1"
}

# ヘルプメッセージ
show_help() {
    echo "Usage: $0 [PORT]"
    echo ""
    echo "Kill process running on specified port"
    echo ""
    echo "Arguments:"
    echo "  PORT    Port number (default: 8080)"
    echo ""
    echo "Environment Variables:"
    echo "  PORT    Port number to kill process on"
    echo ""
    echo "Examples:"
    echo "  $0           # Kill process on port 8080"
    echo "  $0 3000      # Kill process on port 3000"
    echo "  PORT=5000 $0 # Kill process on port 5000"
}

# ヘルプオプションの処理
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_help
    exit 0
fi

# コマンドライン引数でポート番号が指定された場合
if [[ -n "$1" ]]; then
    PORT="$1"
fi

# ポート番号の検証
if ! [[ "$PORT" =~ ^[0-9]+$ ]] || [[ "$PORT" -lt 1 ]] || [[ "$PORT" -gt 65535 ]]; then
    print_error "Invalid port number: $PORT (must be 1-65535)"
    exit 1
fi

print_info "Searching for processes on port $PORT..."

# ポートで実行中のプロセスのPIDを取得
PIDS=$(lsof -ti:$PORT 2>/dev/null)

if [ -z "$PIDS" ]; then
    print_warning "No processes found running on port $PORT"
    exit 0
fi

# 複数のプロセスが見つかった場合の処理
PID_COUNT=$(echo "$PIDS" | wc -w)
print_info "Found $PID_COUNT process(es) running on port $PORT"

# 各プロセスを終了
for PID in $PIDS; do
    print_info "Terminating process (PID: $PID)..."
    
    # プロセス情報を表示
    PROCESS_INFO=$(ps -p "$PID" -o pid,ppid,cmd --no-headers 2>/dev/null)
    if [ -n "$PROCESS_INFO" ]; then
        print_info "Process details: $PROCESS_INFO"
    fi
    
    # プロセスを終了
    if kill -TERM "$PID" 2>/dev/null; then
        print_success "Process (PID: $PID) terminated successfully"
    else
        print_warning "Failed to terminate process (PID: $PID) with SIGTERM, trying SIGKILL..."
        if kill -KILL "$PID" 2>/dev/null; then
            print_success "Process (PID: $PID) killed with SIGKILL"
        else
            print_error "Failed to kill process (PID: $PID)"
        fi
    fi
done

# 最終確認
sleep 1
REMAINING_PIDS=$(lsof -ti:$PORT 2>/dev/null)
if [ -z "$REMAINING_PIDS" ]; then
    print_success "All processes on port $PORT have been terminated"
else
    print_warning "Some processes may still be running on port $PORT"
    print_info "Remaining PIDs: $REMAINING_PIDS"
fi
