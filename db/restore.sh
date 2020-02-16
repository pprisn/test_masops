#!/bin/bash
log=/var/log/dump_masops.log 
export PATH=/usr/local/bin:/usr/local/sbin:/bin:/sbin:/usr/bin:/usr/sbin 
mysql -u ppri -p masops <masops-dump.sql
