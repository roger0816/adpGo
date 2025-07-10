#!/bin/sh


killall adpGo

cd ../bin

./adpGo 6101 db.adp.idv.tw adp > /dev/null 2>&1 &
./adpGo 6201 db.adp.idv.tw adp > /dev/null 2>&1 &
 
