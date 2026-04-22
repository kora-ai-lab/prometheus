use anyhow::{anyhow, Result};
use std::net::{TcpStream, SocketAddr};
use std::path::{PathBuf};
use tokio::process::Command;
use std::time::Duration;
use tokio::time::sleep;

pub async fn ensure_browser_running() -> Result<()> {
    let port = 9222;
    if is_port_open(port) {
        return Ok(());
    }
    
    check_ram_pressure().await;

    let browser_path = find_browser_executable()
        .await
        .ok_or_else(|| anyhow!("Could not find Chrome or Edge executable on this system"))?;

    println!("Launching browser from: {:?}", browser_path);
    
    let mut cmd = Command::new(browser_path);
    cmd.arg(format!("--remote-debugging-port={}", port))
       .arg("--disable-gpu")
       .arg("--disable-dev-shm-usage")
       .arg("--no-first-run")
       .arg("--disable-extensions")
       .arg("--disable-software-rasterizer");

    cmd.spawn()
        .map_err(|e| anyhow!("Failed to spawn browser process: {}", e))?;

    for _ in 0..10 {
        sleep(Duration::from_millis(500)).await;
        if is_port_open(port) {
            return Ok(());
        }
    }

    Err(anyhow!("Browser failed to open debugging port {} within timeout", port))
}

async fn check_ram_pressure() {
    #[cfg(target_os = "windows")]
    {
        let output = Command::new("cmd")
            .args(["/C", "wmic OS get FreePhysicalMemory"])
            .output()
            .await;
        
        if let Ok(out) = output {
            let stdout = String::from_utf8_lossy(&out.stdout);
            if let Some(line) = stdout.lines().find(|l| !l.trim().is_empty() && *l != "FreePhysicalMemory") {
                if let Ok(free_kb) = line.trim().parse::<u64>() {
                    if free_kb < 1024 * 1024 { // Less than 1GB
                        println!("WARNING: Critically low RAM detected ({} KB free). Browser may crash.", free_kb);
                    }
                }
            }
        }
    }
}

fn is_port_open(port: u16) -> bool {
    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    TcpStream::connect_timeout(&addr, Duration::from_millis(200)).is_ok()
}

async fn find_browser_executable() -> Option<PathBuf> {
    let common_paths = [
        r"C:\Program Files\Google\Chrome\Application\chrome.exe",
        r"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe",
        r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe",
        r"C:\Program Files\Microsoft\Edge\Application\msedge.exe",
    ];

    for path in common_paths {
        let p = PathBuf::from(path);
        if p.exists() {
            return Some(p);
        }
    }

    install_browser().await
}

async fn install_browser() -> Option<PathBuf> {
    println!("No browser found. Attempting sovereign recovery (auto-install)...");
    
    if cfg!(target_os = "windows") {
        let status = Command::new("choco")
            .args(["install", "chromium", "--yes", "-y"])
            .status()
            .await;
        
        if status.map(|s| s.success()).unwrap_or(false) {
            let path = PathBuf::from(r"C:\ProgramData\chocolatey\lib\chromium\tools\chrome.exe");
            if path.exists() {
                return Some(path);
            }
        }
    } else if cfg!(target_os = "linux") {
        let status = Command::new("sudo")
            .args(["apt-get", "install", "-y", "chromium-browser"])
            .status()
            .await;
            
        if status.map(|s| s.success()).unwrap_or(false) {
            let path = PathBuf::from("/usr/bin/chromium-browser");
            if path.exists() {
                return Some(path);
            }
        }
    }
    
    None
}
