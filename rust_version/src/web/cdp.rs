use anyhow::{anyhow, Result};
use futures_util::{SinkExt, StreamExt};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use tokio_tungstenite::MaybeTlsStream;
use tokio_tungstenite::{connect_async, tungstenite::protocol::Message};

#[derive(Debug)]
pub struct WebDriver {
    ws_stream: tokio_tungstenite::WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>,
    id_counter: u64,
}

#[derive(Serialize)]
struct CDPRequest {
    id: u64,
    method: String,
    params: Value,
}

impl WebDriver {
    pub async fn connect(port: u16) -> Result<Self> {
        crate::web::launcher::ensure_browser_running().await?;
        let url = format!("ws://localhost:{}/devtools/browser/{}", port, "placeholder"); 
        // Note: In a real CDP flow, we'd first hit http://localhost:port/json/version to get the websocketDebuggerUrl
        // For this implementation, we'll simplify or assume a known path if provided.
        // Actually, we should fetch the debugger URL first.
        
        let client = reqwest::Client::new();
        let resp: Value = client.get(format!("http://localhost:{}/json/version", port))
            .send().await?
            .json().await?;
        
        let ws_url = resp["webSocketDebuggerUrl"].as_str()
            .ok_or_else(|| anyhow!("Failed to get webSocketDebuggerUrl from browser"))?;
            
        let (ws_stream, _) = connect_async(ws_url).await?;
        Ok(Self {
            ws_stream,
            id_counter: 0,
        })
    }

    async fn send_command(&mut self, method: &str, params: Value) -> Result<Value> {
        self.id_counter += 1;
        let id = self.id_counter;
        
        let req = CDPRequest {
            id,
            method: method.to_string(),
            params,
        };
        
        let json_req = serde_json::to_string(&req)?;
        self.ws_stream.send(Message::Text(json_req.into())).await?;
        
        while let Some(msg) = self.ws_stream.next().await {
            let msg = msg?;
            if let Message::Text(text) = msg {
                let resp: Value = serde_json::from_str(&text.to_string())?;
                if resp["id"] == id {
                    return Ok(resp["result"].clone());
                }
            }
        }
        Err(anyhow!("Timeout or connection lost while waiting for CDP response"))
    }

    pub async fn navigate_to(&mut self, url: &str) -> Result<()> {
        self.send_command("Page.navigate", json!({ "url": url })).await?;
        // Wait for page to load
        self.send_command("Page.waitForLoadState", json!({ "state": "load" })).await?;
        Ok(())
    }

    pub async fn click(&mut self, selector: &str) -> Result<()> {
        // 1. Find element
        let root = self.send_command("DOM.getDocument", json!({})).await?;
        let _doc_id = root["root"].as_i64().ok_or_else(|| anyhow!("Failed to get doc id"))?;
        
        let _query_res = self.send_command("DOM.querySelector", json!({
            "nodeId": _doc_id,
            "selector": selector
        })).await?;
        
        // Note: node_id is retrieved but we use JS evaluate for the actual click to simplify.
        
        // 2. Click using Input.dispatchEvent or Mouse
        // Simplest via Runtime.evaluate for a quick click
        let js = format!("document.querySelector('{}').click()", selector);
        self.send_command("Runtime.evaluate", json!({ "expression": js })).await?;
        
        Ok(())
    }

    pub async fn type_text(&mut self, selector: &str, text: &str) -> Result<()> {
        let js = format!(
            "var el = document.querySelector('{}'); el.value = '{}'; el.dispatchEvent(new Event('input', {{ bubbles: true }})); el.dispatchEvent(new Event('change', {{ bubbles: true }}));",
            selector, text
        );
        self.send_command("Runtime.evaluate", json!({ "expression": js })).await?;
        Ok(())
    }

    pub async fn get_page_content(&mut self) -> Result<String> {
        let js = "document.documentElement.outerHTML";
        let res = self.send_command("Runtime.evaluate", json!({ 
            "expression": js, 
            "returnByValue": true 
        })).await?;
        
        let html = res["returnValue"]["value"].as_str()
            .ok_or_else(|| anyhow!("Failed to get page content"))?;
            
        Ok(html.to_string())
    }

    pub async fn get_accessibility_tree(&mut self) -> Result<String> {
        // Simplified visual fallback: Extract interactive elements and their labels
        let js = r#"
            (function() {
                const interactive = ['BUTTON', 'A', 'INPUT', 'SELECT', 'TEXTAREA'];
                const elements = Array.from(document.querySelectorAll('*')).filter(el => 
                    interactive.includes(el.tagName) || el.getAttribute('role')
                );
                return elements.map(el => {
                    const rect = el.getBoundingClientRect();
                    const label = el.innerText || el.getAttribute('aria-label') || el.placeholder || el.value || 'No label';
                    return `Element <${el.tagName}> '${label.trim()}' at (${Math.round(rect.left)}, ${Math.round(rect.top)})`;
                }).join('\n');
            })()
        "#;
        
        let res = self.send_command("Runtime.evaluate", json!({ 
            "expression": js, 
            "returnByValue": true 
        })).await?;
        
        let text = res["returnValue"]["value"].as_str()
            .ok_or_else(|| anyhow!("Failed to extract accessibility tree"))?;
            
        Ok(text.to_string())
    }
}
