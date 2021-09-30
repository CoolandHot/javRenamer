@echo off
>NUL 2>&1 REG.exe query "HKU\S-1-5-19" || (
    ECHO SET UAC = CreateObject^("Shell.Application"^) > "%TEMP%\Getadmin.vbs"
    ECHO UAC.ShellExecute "%~f0", "%1", "", "runas", 1 >> "%TEMP%\Getadmin.vbs"
    "%TEMP%\Getadmin.vbs"
    DEL /f /q "%TEMP%\Getadmin.vbs" 2>NUL
    Exit /b
)

title Add right-click menu --- JavRename it

reg add "HKEY_CLASSES_ROOT\SystemFileAssociations\.mp4\shell\subtitle" /ve /d "jav¸ÄÃû" /f
reg add "HKEY_CLASSES_ROOT\SystemFileAssociations\.mp4\shell\subtitle" /v "Icon" /d "D:/ProgramFiles/go/favicon.ico" /f
reg add "HKEY_CLASSES_ROOT\SystemFileAssociations\.mp4\shell\subtitle\command" /ve /d "cmd.exe /s /k pushd D:\\Coding_Programs\\javRenamer&& go run javRename.go \"%%1\""  /f

pause