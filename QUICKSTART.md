# TPT Titan - Quick Start Guide

TPT Titan is a comprehensive open-source productivity suite. Here's how to get started:

## 🚀 **Recommended: Desktop App Experience**
**File:** `run-desktop-app.bat`
**Requirements:** Go 1.19+ and Node.js 18+

**What it does:**
- Starts the backend API server
- Starts the frontend development server
- Opens Chrome/Edge in "app mode" (feels like a native desktop app)

**To use:** Double-click `run-desktop-app.bat` and wait for your browser to open in desktop-app mode!

## 🐳 **Alternative: Docker (Production Ready)**
**File:** `run-tpt-titan.bat`
**Requirements:** Docker Desktop

**What it does:**
- Runs everything in containers
- Most stable and production-ready

## 🛠️ **Development Mode (For Coders)**
**File:** `run-dev.bat`
**Requirements:** Go 1.19+ and Node.js 18+

**What it does:**
- Runs backend and frontend separately
- Best for development and debugging

## 📋 **What You'll Get:**

✅ **Text Editor** - AI-powered writing with rich formatting
✅ **Spreadsheet** - 50+ formulas with charts and Excel import
✅ **Forms Builder** - Drag-and-drop form creation
✅ **Email Client** - PGP encryption and IMAP/SMTP
✅ **Calendar** - Event management and sharing
✅ **Contacts** - Address book with vCard support
✅ **Chat** - Real-time messaging rooms
✅ **Tasks** - Kanban board with AI suggestions

## 🎯 **Quick Start:**

1. **Double-click `run-desktop-app.bat`**
2. Wait 10-15 seconds for servers to start
3. Your browser opens as a desktop app
4. Enjoy your complete Office 365 alternative!

**Your TPT Titan productivity suite is ready to use!** 🎉

## 📋 Features Included

Once running, you'll have access to:
- **Text Editor** with AI assistance and rich formatting
- **Spreadsheet** with 50+ formulas and charting
- **Forms & Templates** with drag-and-drop builder
- **Email Client** with PGP encryption
- **Calendar** with event management
- **Contacts** with vCard support
- **Chat** with real-time messaging
- **Tasks** with Kanban board
- **File Sync** with P2P capabilities

## 🔧 Requirements Summary

| Option | Docker | Node.js | Go | Rust | Complexity |
|--------|--------|---------|----|------|------------|
| Docker | ✅ Required | ❌ | ❌ | ❌ | Lowest |
| Desktop | ❌ | ✅ 18+ | ❌ | ✅ | Medium |
| Dev Mode | ❌ | ✅ 18+ | ✅ 1.19+ | ❌ | Medium |

## 🛑 To Stop

- **Docker:** Run `stop-tpt-titan.bat`
- **Desktop/Dev:** Close the application windows or press Ctrl+C in terminals

## 📁 Data Storage

All options use SQLite database stored locally. Your data persists between sessions.

## ❓ Having Issues?

Check the main README.md for detailed installation instructions or visit the documentation.
