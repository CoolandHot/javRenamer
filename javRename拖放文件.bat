@echo off
pushd %~dp0
if [%1]==[] goto :eof
set count=1
:loop
REM echo ����·��: %~1
REM echo �ļ�·��: %~dp1
REM echo �ļ�ȫ��: %~nx1
REM echo �ļ����� %~n1
REM echo ��չ���� %~x1
echo -----------------------------
python javRename.py "%~1"
shift
set /a count+=1
if not [%1]==[] goto loop
pause