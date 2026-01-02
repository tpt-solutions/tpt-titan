@echo off
echo Stopping TPT Titan...
echo.

echo Stopping services...
docker-compose down

echo.
echo ===========================================
echo TPT Titan has been stopped.
echo ===========================================
echo.
echo Your data is preserved in the ./data directory.
echo To start again, run: run-tpt-titan.bat
echo.
echo Press any key to exit...
pause >nul
