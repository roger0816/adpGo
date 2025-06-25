#!/bin/sh


killall adpGo

./adpGo 6101 172.104.117.7 adp > /dev/null 2>&1 &
./adpGo 6201 172.104.117.7 adp > /dev/null 2>&1 &
 
