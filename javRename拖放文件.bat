@echo off
pushd %~dp0
if [%1]==[] goto :eof
set count=1
set "AllinOne= "
:loop
@REM echo ����·��: %~1
@REM echo Ŀ¼·��: %~dp1
@REM echo �ļ���+��׺: %~nx1
@REM echo �ļ���: %~n1
@REM echo ��׺: %~x1
@REM echo ȫ���ļ��ľ���·��: %*
echo -----------------------------
@REM python javRename.py "%~1"
set AllinOne=%AllinOne%***%~1
@REM add a delimiter between files, "***"
shift
set /a count+=1
if not [%1]==[] goto loop
go run javRename.go "%AllinOne%"
pause