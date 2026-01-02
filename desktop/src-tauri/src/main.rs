// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::process::Command;
use std::sync::Mutex;
// use tauri::{AppHandle, State};

// App state
#[derive(Default)]
struct AppState {
    backend_process: Mutex<Option<std::process::Child>>,
}

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

// Start the backend server
#[tauri::command]
async fn start_backend(state: State<'_, AppState>) -> Result<String, String> {
    let mut backend_process = state.backend_process.lock().unwrap();

    if backend_process.is_some() {
        return Ok("Backend is already running".to_string());
    }

    // For now, assume backend is running externally
    // In production, you would bundle the backend binary
    Ok("Backend should be started externally (run backend separately)".to_string())
}

// Stop the backend server
#[tauri::command]
async fn stop_backend(state: State<'_, AppState>) -> Result<String, String> {
    let mut backend_process = state.backend_process.lock().unwrap();

    if let Some(mut child) = backend_process.take() {
        child.kill().map_err(|e| format!("Failed to stop backend: {}", e))?;
        Ok("Backend stopped successfully".to_string())
    } else {
        Ok("Backend is not running".to_string())
    }
}

// Get backend status
#[tauri::command]
async fn get_backend_status(state: State<'_, AppState>) -> Result<String, String> {
    let mut backend_process = state.backend_process.lock().unwrap();

    if let Some(child) = backend_process.as_mut() {
        match child.try_wait() {
            Ok(Some(status)) => {
                if status.success() {
                    Ok("Backend exited successfully".to_string())
                } else {
                    Ok(format!("Backend exited with error: {:?}", status))
                }
            }
            Ok(None) => Ok("Backend is running".to_string()),
            Err(e) => Err(format!("Failed to check backend status: {}", e)),
        }
    } else {
        Ok("Backend is not running".to_string())
    }
}

// Open file dialog
#[tauri::command]
async fn open_file_dialog() -> Result<String, String> {
    // Use Tauri's dialog API
    // This would be implemented with Tauri's dialog plugin
    Ok("File dialog not implemented yet".to_string())
}

// Save file dialog
#[tauri::command]
async fn save_file_dialog() -> Result<String, String> {
    // Use Tauri's dialog API
    // This would be implemented with Tauri's dialog plugin
    Ok("Save dialog not implemented yet".to_string())
}

// Get system information
#[tauri::command]
async fn get_system_info() -> Result<serde_json::Value, String> {
    let info = serde_json::json!({
        "platform": std::env::consts::OS,
        "arch": std::env::consts::ARCH,
        "version": env!("CARGO_PKG_VERSION"),
        "name": "TPT Titan Desktop"
    });

    Ok(info)
}

// Handle file operations
#[tauri::command]
async fn read_file(path: String) -> Result<String, String> {
    std::fs::read_to_string(&path)
        .map_err(|e| format!("Failed to read file {}: {}", path, e))
}

#[tauri::command]
async fn write_file(path: String, content: String) -> Result<(), String> {
    std::fs::write(&path, content)
        .map_err(|e| format!("Failed to write file {}: {}", path, e))
}

// Handle notifications
#[tauri::command]
async fn show_notification(title: String, body: String) -> Result<(), String> {
    // Use Tauri's notification API
    // This would be implemented with Tauri's notification plugin
    println!("Notification: {} - {}", title, body);
    Ok(())
}

fn main() {
    tauri::Builder::default()
        .manage(AppState::default())
        .invoke_handler(tauri::generate_handler![
            greet,
            start_backend,
            stop_backend,
            get_backend_status,
            open_file_dialog,
            save_file_dialog,
            get_system_info,
            read_file,
            write_file,
            show_notification
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
