@echo off
echo Starting TPT Titan Desktop App Experience...
echo.
echo This creates a desktop-app-like experience in your browser.
echo.
echo Press any key to continue or Ctrl+C to cancel...
pause >nul

echo.
echo Starting backend server...
start /B "TPT Titan Backend" cmd /C "cd backend && copy simple_go.mod go.mod >nul 2>&1 && go mod tidy >nul 2>&1 && go run simple_main.go"

echo Waiting for backend to start...
timeout /t 3 /nobreak >nul

echo.
echo Starting frontend development server...
start /B "TPT Titan Frontend" cmd /C "cd frontend && npm run dev"

echo Waiting for frontend to start...
timeout /t 5 /nobreak >nul

echo.
echo Opening TPT Titan in desktop-like browser window...

REM Try to open in Chrome app mode first (most desktop-like)
start "" "chrome.exe" --app="http://localhost:5173" --window-size=1400,900 --window-position=50,50

REM If Chrome not found, try Edge
if errorlevel 1 (
    start "" "msedge.exe" --app="http://localhost:5173" --window-size=1400,900 --window-position=50,50
)

REM If Edge not found, try Firefox
if errorlevel 1 (
    start "" "firefox.exe" -P "TPT Titan" "http://localhost:5173"
)

REM If no special mode works, just open in default browser
if errorlevel 1 (
    start "" "http://localhost:5173"
)

echo.
echo ===========================================
echo TPT Titan Desktop Experience Active!
echo ===========================================
echo.
echo - Backend: Running on port 8080
echo - Frontend: Running on port 5173
echo - Browser: Opened in app-like window
echo.
echo The browser should now be open in a desktop-app style window.
echo It will look and feel like a native desktop application!
echo.
echo To close everything:
echo 1. Close the browser window
echo 2. Close this command prompt
echo.
echo Press any key to exit this launcher...
pause >nul
