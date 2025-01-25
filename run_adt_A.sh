#!/bin/bash

start() {
    echo "Starting adpGo processes..."

    # 启动 adpGo 进程，不保存日志
    ./adpGo 6101 159.223.54.194 adt1 > /dev/null 2>&1 &
    echo "adpGo 6101 started with PID $!"

    ./adpGo 6201 159.223.54.194 adt2 > /dev/null 2>&1 &
    echo "adpGo 6102 started with PID $!"

    ./adpGo 6301 159.223.54.194 adt3 > /dev/null 2>&1 &
    echo "adpGo 6103 started with PID $!"

    echo "All processes are running in the background."
}

stop() {
    echo "Stopping all adpGo processes..."
    # 停止所有 adpGo 进程
    pkill -f "./adpGo"
    echo "All adpGo processes have been stopped."
}

status() {
    echo "Checking running adpGo processes:"
    # 查看 adpGo 进程状态
    ps aux | grep "./adpGo" | grep -v "grep"
}

# 检查用户输入的命令
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    *)
        echo "Usage: $0 {start|stop|status}"
        exit 1
        ;;
esac