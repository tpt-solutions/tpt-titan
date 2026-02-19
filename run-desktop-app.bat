@echo off
setlocal EnableDelayedExpansion

echo ============================================
echo  TPT Titan Desktop App Experience Launcher
echo ============================================
echo.
echo This creates a desktop-app-like experience in your browser.
echo.
echo Press any key to continue or Ctrl+C to cancel...
pause >nul

REM -----------------------------------------------
REM  Setup .env if missing
REM -----------------------------------------------
if not exist "%~dp0.env" (
    if exist "%~dp0.env.example" (
        copy "%~dp0.env.example" "%~dp0.env" >nul
        echo Created .env file from template.
    )
)

REM -----------------------------------------------
REM  Start Backend
REM  Use start /d to set working directory directly,
REM  which avoids the nested-quote bug with cmd /k "cd /d "..."".
REM -----------------------------------------------
echo.
echo Starting backend server...
start "TPT Titan Backend" /d "%~dp0backend" cmd /k "echo [Backend] Starting Go server... && go run simple_main.go"

echo Waiting for backend to initialise (5 seconds)...
timeout /t 5 /nobreak >nul

REM -----------------------------------------------
REM  Start Frontend
REM  Use start /d to set working directory directly
REM -----------------------------------------------
echo.
echo Starting frontend development server...

if not exist "%~dp0frontend\node_modules" (
    echo [Frontend] node_modules not found - running npm install first...
    start "TPT Titan Frontend" /d "%~dp0frontend" cmd /k "npm install && npm run dev"
) else (
    start "TPT Titan Frontend" /d "%~dp0frontend" cmd /k "npm run dev"
)

echo Waiting for frontend to initialise (8 seconds)...
timeout /t 8 /nobreak >nul

REM -----------------------------------------------
REM  Open browser in app mode
REM  Use 'where' to detect which browser is available
REM -----------------------------------------------
echo.
echo Opening TPT Titan in desktop-like browser window...

set BROWSER_LAUNCHED=0
set APP_URL=http://localhost:5173

REM --- Try Google Chrome ---
where chrome.exe >nul 2>&1
if !errorlevel! == 0 (
    echo Found Chrome - launching in app mode...
    start "" "chrome.exe" --app="%APP_URL%" --window-size=1400,900 --window-position=50,50
    set BROWSER_LAUNCHED=1
    goto :browser_done
)

REM Chrome may be installed but not on PATH - check common locations
set CHROME_PATH=
if exist "%ProgramFiles%\Google\Chrome\Application\chrome.exe" (
    set CHROME_PATH=%ProgramFiles%\Google\Chrome\Application\chrome.exe
)
if exist "%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe" (
    set CHROME_PATH=%ProgramFiles(x86)%\Google\Chrome\Application\chrome.exe
)
if exist "%LocalAppData%\Google\Chrome\Application\chrome.exe" (
    set CHROME_PATH=%LocalAppData%\Google\Chrome\Application\chrome.exe
)

if defined CHROME_PATH (
    echo Found Chrome at: !CHROME_PATH!
    start "" "!CHROME_PATH!" --app="%APP_URL%" --window-size=1400,900 --window-position=50,50
    set BROWSER_LAUNCHED=1
    goto :browser_done
)

REM --- Try Microsoft Edge ---
where msedge.exe >nul 2>&1
if !errorlevel! == 0 (
    echo Chrome not found - using Microsoft Edge in app mode...
    start "" "msedge.exe" --app="%APP_URL%" --window-size=1400,900 --window-position=50,50
    set BROWSER_LAUNCHED=1
    goto :browser_done
)

set EDGE_PATH=
if exist "%ProgramFiles%\Microsoft\Edge\Application\msedge.exe" (
    set EDGE_PATH=%ProgramFiles%\Microsoft\Edge\Application\msedge.exe
)
if exist "%ProgramFiles(x86)%\Microsoft\Edge\Application\msedge.exe" (
    set EDGE_PATH=%ProgramFiles(x86)%\Microsoft\Edge\Application\msedge.exe
)
if exist "%LocalAppData%\Microsoft\Edge\Application\msedge.exe" (
    set EDGE_PATH=%LocalAppData%\Microsoft\Edge\Application\msedge.exe
)

if defined EDGE_PATH (
    echo Found Edge at: !EDGE_PATH!
    start "" "!EDGE_PATH!" --app="%APP_URL%" --window-size=1400,900 --window-position=50,50
    set BROWSER_LAUNCHED=1
    goto :browser_done
)

REM --- Try Firefox ---
where firefox.exe >nul 2>&1
if !errorlevel! == 0 (
    echo Chrome/Edge not found - using Firefox...
    start "" "firefox.exe" "%APP_URL%"
    set BROWSER_LAUNCHED=1
    goto :browser_done
)

set FIREFOX_PATH=
if exist "%ProgramFiles%\Mozilla Firefox\firefox.exe" (
    set FIREFOX_PATH=%ProgramFiles%\Mozilla Firefox\firefox.exe
)
if exist "%ProgramFiles(x86)%\Mozilla Firefox\firefox.exe" (
    set FIREFOX_PATH=%ProgramFiles(x86)%\Mozilla Firefox\firefox.exe
)

if defined FIREFOX_PATH (
    echo Found Firefox at: !FIREFOX_PATH!
    start "" "!FIREFOX_PATH!" "%APP_URL%"
    set BROWSER_LAUNCHED=1
    goto :browser_done
)

REM --- Fallback: open with default browser via shell ---
echo No specific browser found - opening with default browser...
start "" "%APP_URL%"
set BROWSER_LAUNCHED=1

:browser_done

echo.
echo ===========================================
echo  TPT Titan Desktop Experience Active!
echo ===========================================
echo.
echo  Backend API : http://localhost:8080
echo  Frontend    : http://localhost:5173
echo.
echo  Two terminal windows are open:
echo    - "TPT Titan Backend"  (Go server)
echo    - "TPT Titan Frontend" (Vite dev server)
echo.
echo  To stop everything:
echo    1. Close the browser window
echo    2. Close the "TPT Titan Backend" terminal
echo    3. Close the "TPT Titan Frontend" terminal
echo.
echo Press any key to exit this launcher...
pause >nul

endlocal
