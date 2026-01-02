@echo off
echo Starting TPT Titan in Development Mode...
echo.
echo This will start both the backend (Go) and frontend (Svelte) in development mode.
echo Make sure you have Go 1.19+ and Node.js 18+ installed.
echo.
echo Press any key to continue or Ctrl+C to cancel...
pause >nul

echo.
echo Setting up environment...
if not exist ".env" (
    copy .env.example .env
    echo Created .env file from template
)

echo.
echo Starting backend server in new window...
start "TPT Titan Backend" cmd /k "cd backend && copy simple_go.mod go.mod && go mod tidy && go run simple_main.go"

echo Waiting 5 seconds for backend to start...
timeout /t 5 /nobreak >nul

echo.
echo Starting frontend development server in new window...
start "TPT Titan Frontend" cmd /k "cd frontend && npm install && npm run dev"

echo.
echo ===========================================
echo TPT Titan Development Servers Started!
echo ===========================================
echo.
echo Backend API: http://localhost:8080
echo Frontend:    http://localhost:5173 (or check the frontend terminal for the exact port)
echo.
echo Close the terminal windows to stop the servers.
echo Your data is stored locally in SQLite database.
echo.
echo Press any key to exit this launcher...
pause >nul
