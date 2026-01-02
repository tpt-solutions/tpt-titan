@echo off
echo Starting TPT Titan...
echo.
echo This will start the TPT Titan productivity suite using Docker Compose.
echo The application will be available at:
echo - Frontend: http://localhost:3000
echo - Backend API: http://localhost:8080
echo.
echo Press any key to continue or Ctrl+C to cancel...
pause >nul

echo.
echo Pulling latest images and starting services...
docker-compose up -d

echo.
echo Waiting for services to start...
timeout /t 10 /nobreak >nul

echo.
echo Checking service status...
docker-compose ps

echo.
echo ===========================================
echo TPT Titan is now running!
echo ===========================================
echo.
echo Access the application at: http://localhost:3000
echo.
echo To stop the application, run: stop-tpt-titan.bat
echo.
echo Press any key to exit...
pause >nul
