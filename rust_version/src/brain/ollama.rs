use crate::brain::ModelProvider;
use anyhow::{Result, anyhow};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::time::Duration;

#[derive(Serialize)]
struct OllamaRequest {
    model: String,
    prompt: String,
    stream: bool,
}

#[derive(Deserialize, Debug)]
#[serde(untagged)]
enum OllamaResponse {
    Success { response: String },
    Error { error: String },
}

pub struct OllamaProvider {
    url: String,
    model: String,
}

impl OllamaProvider {
    pub fn new(url: String, model: String) -> Self {
        Self {
            url,
            model,
        }
    }
}

#[async_trait::async_trait]
impl ModelProvider for OllamaProvider {
    async fn generate_response(&self, prompt: &str) -> Result<String> {
        let client = Client::builder()
            .use_http1_only()
            .timeout(Duration::from_secs(90))
            .build()
            .map_err(|e| anyhow!("Failed to create HTTP client: {}", e))?;

        let request_body = OllamaRequest {
            model: self.model.clone(),
            prompt: prompt.to_string(),
            stream: false,
        };

        let full_url = format!("{}/api/generate", self.url);
        
        let res = client
            .post(&full_url)
            .json(&request_body)
            .send()
            .await
            .map_err(|e| anyhow!("HTTP Error: {}", e))?;

        let response_text = res.text().await.map_err(|e| anyhow!("Response read error: {}", e))?;

        println!("Ollama raw response: {}", response_text);

        match serde_json::from_str::<OllamaResponse>(&response_text) {
            Ok(OllamaResponse::Success { response }) => Ok(response),
            Ok(OllamaResponse::Error { error }) => anyhow::bail!("Ollama API Error: {}", error),
            Err(e) => anyhow::bail!("Failed to decode Ollama response: {} - Content: {}", e, response_text),
        }
    }
}