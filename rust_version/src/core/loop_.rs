use crate::brain::ModelProvider;
use crate::capabilities::manager::CapabilityManager;
use crate::executor::shell::execute_command;
use crate::web::cdp::WebDriver;
use anyhow::{Result, anyhow};
use std::collections::VecDeque;
use std::sync::Arc;
use tokio::sync::Mutex;
use tokio::time::{timeout, Duration};

pub struct PrometheusLoop {
    model: Box<dyn ModelProvider + Send + Sync>,
    history: VecDeque<String>,
    capability_manager: CapabilityManager,
    web_driver: Option<Arc<Mutex<WebDriver>>>,
}

impl PrometheusLoop {
    pub fn new(model: Box<dyn ModelProvider + Send + Sync>, capability_manager: CapabilityManager, web_driver: Option<Arc<Mutex<WebDriver>>>) -> Self {
        Self {
            model,
            history: VecDeque::new(),
            capability_manager,
            web_driver,
        }
    }

    pub async fn run(&mut self, goal: &str, system_prompt: &str) -> Result<String> {
        let current_goal = goal.to_string();
        let mut iteration = 0;
        let max_iterations = 15;

        loop {
            iteration += 1;
            if iteration > max_iterations {
                return Err(anyhow!("Max iterations ({}) reached without completing goal. Progress: {:?}", max_iterations, self.history));
            }

            println!("--- Iteration {} ---", iteration);

            let prompt = self.build_prompt(system_prompt, &current_goal);
            let response = match timeout(Duration::from_secs(90), self.model.generate_response(&prompt)).await {
                Ok(Ok(resp)) => resp,
                Ok(Err(e)) => return Err(anyhow!("LLM generation failed: {}", e)),
                Err(_) => return Err(anyhow!("LLM generation timed out after 30s. Check if Ollama is responding.")),
            };
            
            println!("LLM Response:\n{}\n", response);

            if response.contains("FINISH:") {
                let result = response.split("FINISH:").nth(1).unwrap_or("").trim().to_string();
                return Ok(result);
            }

            if let Some(cmd) = self.extract_command(&response) {
                println!("Executing: {}\n", cmd);
                
                let execution_result = if cmd.starts_with("USE_CAPABILITY:") {
                    let rest = cmd.strip_prefix("USE_CAPABILITY:").unwrap().trim();
                    let cap_name = rest.split_whitespace().next().unwrap_or("").trim();
                    
                    if cap_name == "web_browser" {
                        match &self.web_driver {
                            Some(driver) => {
                                let params = self.parse_web_params(rest);
                                let mut driver_locked = driver.lock().await;
                                
                                let action = params.get("action").map(|s| s.as_str()).unwrap_or("get_content");
                                
                                match action {
                                    "navigate" => {
                                        if let Some(url) = params.get("url") {
                                            let url_owned = url.clone();
                                            driver_locked.navigate_to(&url_owned).await?;
                                            Ok("Navigated successfully".to_string())
                                        } else {
                                            Err(anyhow!("navigate action requires url parameter"))
                                        }
                                    },
                                    "click" => {
                                        if let Some(sel) = params.get("selector") {
                                            let sel_owned = sel.clone();
                                            driver_locked.click(&sel_owned).await?;
                                            Ok("Clicked successfully".to_string())
                                        } else {
                                            Err(anyhow!("click action requires selector parameter"))
                                        }
                                    },
                                    "type" => {
                                        let sel = params.get("selector").cloned().unwrap_or_else(|| "".to_string());
                                        let txt = params.get("text").cloned().unwrap_or_else(|| "".to_string());
                                        driver_locked.type_text(&sel, &txt).await?;
                                        Ok("Typed successfully".to_string())
                                    },
                                    "get_content" => {
                                        let mut d = driver_locked;
                                        d.get_page_content().await
                                    },
                                    "get_visual" | "accessibility" => {
                                        let mut d = driver_locked;
                                        d.get_accessibility_tree().await
                                    },
                                    _ => {
                                        let mut d = driver_locked;
                                        d.get_page_content().await
                                    }
                                }
                            },
                            None => Err(anyhow!("WebDriver not initialized")),
                        }
                    } else {
                        match self.capability_manager.get_capability_path(cap_name) {
                            Ok(path) => {
                                let path_str = path.to_string_lossy().to_string();
                                execute_command(&format!("sh {}", path_str)).await
                            },
                            Err(e) => Err(e),
                        }
                    }
                } else {
                    execute_command(&cmd).await
                };

                match execution_result {
                    Ok(output) => {
                        let entry = format!("COMMAND: {}\nOUTPUT: {}\n", cmd, output);
                        self.history.push_back(entry);
                    }
                    Err(e) => {
                        let entry = format!("COMMAND: {}\nERROR: {}\n", cmd, e);
                        self.history.push_back(entry);
                    }
                }
            } else {
                let entry = format!("LLM_NO_COMMAND: {}\n", response);
                self.history.push_back(entry);
            }

            if self.history.len() > 10 {
                self.history.pop_front();
            }
        }
    }

    fn build_prompt(&self, system_prompt: &str, goal: &str) -> String {
        let caps_desc = self.capability_manager.get_capabilities_description();
        let caps_section = if !caps_desc.is_empty() {
            format!("\nYour available capabilities are: {}\n", caps_desc)
        } else {
            "".to_string()
        };

        let mut full_prompt = format!("{}\n{}\nGoal: {}\n\nHistory:\n", system_prompt, caps_section, goal);
        for entry in &self.history {
            full_prompt.push_str(entry);
        }
        full_prompt.push_str("\nNext Action (COMMAND: <cmd>, USE_CAPABILITY: <name param=value>, or FINISH: <result>):\n");
        full_prompt
    }

    fn extract_command(&self, response: &str) -> Option<String> {
        if let Some(idx) = response.find("COMMAND:") {
            let start = idx + "COMMAND:".len();
            let end = response[start..].find('\n').unwrap_or(response[start..].len());
            return Some(response[start..start + end].trim().to_string());
        }
        if let Some(idx) = response.find("USE_CAPABILITY:") {
            let start = idx + "USE_CAPABILITY:".len();
            let end = response[start..].find('\n').unwrap_or(response[start..].len());
            return Some(format!("USE_CAPABILITY: {}", response[start..start + end].trim()));
        }
        None
    }

    fn parse_web_params(&self, rest: &str) -> std::collections::HashMap<String, String> {
        let mut params = std::collections::HashMap::new();
        let cmd_part = rest.split_whitespace().skip(1).collect::<Vec<_>>().join(" ");
        
        for pair in cmd_part.split(',') {
            if let Some((key, val)) = pair.split_once('=') {
                params.insert(key.trim().to_string(), val.trim().to_string());
            }
        }
        params
    }
}