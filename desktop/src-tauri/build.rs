fn main() {
    // This tells cargo to rerun this build script if the tauri.conf.json changes
    println!("cargo:rerun-if-changed=tauri.conf.json");

    // Set OUT_DIR for tauri::generate_context!() macro
    // OUT_DIR is automatically set by cargo when a build script exists
    let out_dir = std::env::var("OUT_DIR").unwrap();
    println!("cargo:rustc-env=TAURI_OUT_DIR={}", out_dir);

    // Run tauri-build to generate context and handle resources
    // This is required for tauri::generate_context!() to work
    tauri_build::build()
}
