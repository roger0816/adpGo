#!/bin/bash

start() {
    echo "Starting adpGo processes..."

    cd ../bin

    # 启动 adpGo 进程，不保存日志
    ./adpGo 6101 159.223.54.194 adt1 > /dev/null 2>&1 &

    ./adpGo 6201 159.223.54.194 adt2 > /dev/null 2>&1 &

    ./adpGo 6301 159.223.54.194 adt3 > /dev/null 2>&1 &

    ./adpGo 6401 159.223.54.194 adt4 > /dev/null 2>&1 &

    ./adpGo 6501 159.223.54.194 adt5 > /dev/null 2>&1 &

    ./adpGo 6601 159.223.54.194 adt6 > /dev/null 2>&1 &
 
 

    # 启动 adpGo 进程，不保存日志
    ./adpGo 6102 159.223.54.194 adt1 > /dev/null 2>&1 &

    ./adpGo 6202 159.223.54.194 adt2 > /dev/null 2>&1 &

    ./adpGo 6302 159.223.54.194 adt3 > /dev/null 2>&1 &

    ./adpGo 6402 159.223.54.194 adt4 > /dev/null 2>&1 &

    ./adpGo 6502 159.223.54.194 adt5 > /dev/null 2>&1 &

    ./adpGo 6602 159.223.54.194 adt6 > /dev/null 2>&1 &


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

stop
start