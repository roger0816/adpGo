#!/bin/sh



killall adpTest

cd ../bin

cp adpGo adpTest


./adpTest 6103 db.adp.idv.tw adpTest > /dev/null 2>&1 &
