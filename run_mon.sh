#!/bin/bash
log=/var/log/concurrency_masops.log
cd /root/go/src/github.com/pprisn/test_masops/
export PATH=/usr/local/bin:/usr/local/sbin:/bin:/sbin:/usr/bin:/usr/sbin 
./concurrency_work -ufps C00 
./concurrency_work -ufps R00 
./concurrency_work -ufps R01 
./concurrency_work -ufps R02 
./concurrency_work -ufps R03 
./concurrency_work -ufps R04 
./concurrency_work -ufps R05 
./concurrency_work -ufps R06 
./concurrency_work -ufps R07 
./concurrency_work -ufps R08 
./concurrency_work -ufps R09 
./concurrency_work -ufps R10 
./concurrency_work -ufps R11 
./concurrency_work -ufps R12 
./concurrency_work -ufps R13 
./concurrency_work -ufps R14 
./concurrency_work -ufps R15 
./concurrency_work -ufps R16 
./concurrency_work -ufps R17 
./concurrency_work -ufps R18 
./concurrency_work -ufps R19 
./concurrency_work -ufps R21 
./concurrency_work -ufps R22 
./concurrency_work -ufps R23 
./concurrency_work -ufps R24 
./concurrency_work -ufps R25 
./concurrency_work -ufps R26 
./concurrency_work -ufps R27 
./concurrency_work -ufps R28 
./concurrency_work -ufps R29 
./concurrency_work -ufps R30 
./concurrency_work -ufps R31 
./concurrency_work -ufps R32 
./concurrency_work -ufps R33 
./concurrency_work -ufps R34 
./concurrency_work -ufps R35
./concurrency_work -ufps R36
./concurrency_work -ufps R37
./concurrency_work -ufps R38
./concurrency_work -ufps R39
