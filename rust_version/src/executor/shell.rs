use tokio::process::Command;
use anyhow::{Result, Context};
use std::process::Stdio;

pub async fn execute_command(cmd: &str) -> Result<String> {
    let (shell, arg) = if cfg!(target_os = "windows") {
        ("cmd", "/C")
    } else {
        ("sh", "-c")
    };

    let output = Command::new(shell)
        .arg(arg)
        .arg(cmd)
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .output()
        .await
        .context("Failed to execute process")?;

    let stdout = String::from_utf8_lossy(&output.stdout);
    let stderr = String::from_utf8_lossy(&output.stderr);
    
    let combined = format!("{}{}", stdout, stderr);
    Ok(clean_output(&combined))
}

fn clean_output(output: &str) -> String {
    output.trim().to_string()
}
