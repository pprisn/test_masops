#!/bin/bash
log=/var/log/test_masops.log 
export LOGDB=ppri:bsd123ufps
export PATH=/usr/local/bin:/usr/local/sbin:/bin:/sbin:/usr/bin:/usr/sbin 
#exec &>>$log echo $(date +"%D %T")+'Starting...'
exec ./server  &>>$log  & exit 0
