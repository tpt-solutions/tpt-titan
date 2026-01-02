@echo off
echo Starting TPT Titan Desktop Experience...
echo.
echo This will start the backend and open the frontend in a browser window.
echo It will feel like a desktop app!
echo.
echo Press any key to continue or Ctrl+C to cancel...
pause >nul

echo.
echo Starting backend server...
start /B "TPT Titan Backend" cmd /C "cd backend && copy simple_go.mod go.mod >nul 2>&1 && go mod tidy >nul 2>&1 && go run simple_main.go"

echo Waiting for backend to start...
timeout /t 3 /nobreak >nul

echo.
echo Opening TPT Titan in browser window...
start "" "C:\Windows\System32\mshta.exe" "javascript:window.resizeTo(1400,900);window.moveTo(100,100);window.location='http://localhost:5173';"

echo.
echo ===========================================
echo TPT Titan Desktop Experience Started!
echo ===========================================
echo.
echo - Backend running on port 8080
echo - Frontend opened in browser window
echo.
echo The browser window should appear maximized and feel like a desktop app.
echo.
echo To close: Close the browser window and this terminal.
echo.
echo Press any key to exit this launcher...
pause >nul
