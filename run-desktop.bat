@echo off
echo Starting TPT Titan Desktop Application...
echo.
echo This will start the TPT Titan desktop app using Tauri.
echo Make sure you have Node.js and Rust installed.
echo.
echo Press any key to continue or Ctrl+C to cancel...
pause >nul

echo.
echo Installing dependencies and starting desktop app...
cd desktop
npm install
npm run dev

echo.
echo Desktop app should now be running!
echo If it doesn't start automatically, check for any error messages above.
echo.
pause
