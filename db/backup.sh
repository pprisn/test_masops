#!/bin/bash
log=/var/log/dump_masops.log 
export PATH=/usr/local/bin:/usr/local/sbin:/bin:/sbin:/usr/bin:/usr/sbin 
mysqldump -u root -p masops >masops-dump.sql