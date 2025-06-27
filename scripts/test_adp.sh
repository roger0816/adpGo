#!/bin/sh



killall adpGoTest

cd ../bin

cp adpGo adpGoTest


./adpGo 6103 172.104.112.34 adp > /dev/null 2>&1 &
