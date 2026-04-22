use async_trait::async_trait;
use anyhow::Result;

pub mod ollama;
pub use ollama::OllamaProvider;

#[async_trait]
pub trait ModelProvider {
    async fn generate_response(&self, prompt: &str) -> Result<String>;
}

