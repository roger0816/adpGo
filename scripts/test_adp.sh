#!/bin/sh



killall adpTest

cd ../bin

cp adpGo adpTest


./adpTest 6103 172.104.112.34 adp > /dev/null 2>&1 &
