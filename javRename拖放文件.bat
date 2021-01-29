@echo off
pushd %~dp0
if [%1]==[] goto :eof
set count=1
:loop
REM echo 绝对路径: %~1
REM echo 文件路径: %~dp1
REM echo 文件名+扩展: %~nx1
REM echo 文件名： %~n1
REM echo 扩展名： %~x1
echo -----------------------------
python javRename_v7.py "%~1"
shift
set /a count+=1
if not [%1]==[] goto loop
pause