#!/bin/sh



killall adpGoTest

cp ../bin/adpGo ../bin/adpGoTest


../bin/adpGo 6103 172.104.112.34 adp > /dev/null 2>&1 &
