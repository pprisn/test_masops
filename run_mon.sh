#!/bin/bash
log=/var/log/test_masops.log 
cd /root/go/src/github.com/pprisn/test_masops/
export PATH=/usr/local/bin:/usr/local/sbin:/bin:/sbin:/usr/bin:/usr/sbin 
exec ./main -ufps C00,R00,R01 &>>$log  & 
exec ./main -ufps R02,R03,R04 &>>$log  & 
exec ./main -ufps R05,R06,R07 &>>$log  & 
exec ./main -ufps R08,R09,R10 &>>$log  & 
exec ./main -ufps R11,R12,R13 &>>$log  & 
exec ./main -ufps R14,R15,R16 &>>$log  & 
exec ./main -ufps R17,R18,R19 &>>$log  & 
exec ./main -ufps R21,R22,R23 &>>$log  & 
exec ./main -ufps R24,R25,R26 &>>$log  & 
exec ./main -ufps R27,R28,R29 &>>$log  & 
exec ./main -ufps R30,R31,R32 &>>$log  & 
exec ./main -ufps R33,R34,R35 &>>$log  & 
exec ./main -ufps R36,R37,R38 &>>$log  & 
exec ./main -ufps R39,R40,R41 &>>$log  & 
exec ./main -ufps R42,R43,R44 &>>$log  & 
exec ./main -ufps R45,R46,R48 &>>$log  & 
exec ./main -ufps R49,R50,R51 &>>$log  & 
exec ./main -ufps R52,R53,R54 &>>$log  & 
exec ./main -ufps R55,R56,R57 &>>$log  & 
exec ./main -ufps R58,R59,R60 &>>$log  & 
exec ./main -ufps R61,R62,R63 &>>$log  & 
exec ./main -ufps R64,R65,R67 &>>$log  & 
exec ./main -ufps R68,R69,R70 &>>$log  & 
exec ./main -ufps R71,R72,R73 &>>$log  & 
exec ./main -ufps R74,R75,R76 &>>$log  & 
exec ./main -ufps R77,R78,R79 &>>$log  & 
exec ./main -ufps R83,R86,R87 &>>$log  & 
exec ./main -ufps R89,R95,R96 &>>$log  & 
