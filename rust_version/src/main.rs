use std::env;
use std::fs;
use anyhow::{Result, Context};
use serde::Deserialize;
use std::sync::Arc;
use tokio::sync::Mutex;

mod brain;
mod executor;
mod core;
mod web;
mod capabilities;

use crate::brain::{OllamaProvider, ModelProvider};
use crate::core::PrometheusLoop;
use crate::capabilities::manager::CapabilityManager;
use crate::web::cdp::WebDriver;

#[derive(Deserialize)]
struct Config {
    model_provider: String,
    model_name: String,
    parallel_requests: u32,
}

#[tokio::main]
async fn main() -> Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Usage: prometheus \"<goal>\"");
        return Ok(());
    }
    let goal = &args[1];

    // Load Config
    let config_str = fs::read_to_string("config.yaml").context("Failed to read config.yaml")?;
    let config: Config = serde_yaml::from_str(&config_str).context("Failed to parse config.yaml")?;

    // Load System Prompt
    let prompt_content = fs::read_to_string("prompt.md").context("Failed to read prompt.md")?;

    // Initialize Brain
    let url = "http://localhost:11434".to_string();
    let model: Box<dyn ModelProvider + Send + Sync> = Box::new(OllamaProvider::new(url, config.model_name.clone()));

    // Initialize Capability Manager
    let capability_manager = CapabilityManager::new("capabilities")
        .context("Failed to initialize capability manager")?;

    // Initialize WebDriver (Auto-launches browser if needed)
    let web_driver = WebDriver::connect(9222).await.ok();
    let web_driver_arc = web_driver.map(|wd| Arc::new(Mutex::new(wd)));

    // Start Prometheus Loop
    let mut prometheus = PrometheusLoop::new(model, capability_manager, web_driver_arc);

    println!("Starting Prometheus loop with goal: {}", goal);
    let result = prometheus.run(goal, &prompt_content).await?;
    println!("\nFinal Result: {}", result);

    Ok(())
}