@echo off
echo 正在查找并终止 main.exe 进程...

:: 查找 main.exe 的 PID
for /f "tokens=2 delims==," %%a in ('wmic process where "name='main.exe'" get processid^,status /format:csv ^| findstr [0-9]') do (
    echo 正在终止 PID: %%a
    taskkill /PID %%a /F
)

echo 所有 main.exe 进程已终止。
pause
