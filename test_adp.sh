#!/bin/sh



cp adpGo adpGoTest

pkill -f "./adpGoTest"

./adpGoTest 6103 172.104.112.34 adp > /dev/null 2>&1 &
