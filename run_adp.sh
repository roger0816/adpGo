#!/bin/sh


pkill -f "./adpGo"

./adpGo 6101 172.104.117.7 adp > /dev/null 2>&1 &
 