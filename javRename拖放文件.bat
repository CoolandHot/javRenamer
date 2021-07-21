@echo off
pushd %~dp0
if [%1]==[] goto :eof
set count=1
set "AllinOne= "
:loop
@REM echo 绝对路径: %~1
@REM echo 目录路径: %~dp1
@REM echo 文件名+后缀: %~nx1
@REM echo 文件名: %~n1
@REM echo 后缀: %~x1
@REM echo 全部文件的绝对路径: %*
echo -----------------------------
@REM python javRename.py "%~1"
set AllinOne=%AllinOne%***%~1
@REM add a delimiter between files, "***"
shift
set /a count+=1
if not [%1]==[] goto loop
go run javRename.go "%AllinOne%"
pause