use crate::capabilities::meta::CapabilityMeta;
use anyhow::{Result, Context};
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

pub struct CapabilityManager {
    capabilities: HashMap<String, CapabilityMeta>,
    vault_path: PathBuf,
}

impl CapabilityManager {
    pub fn new(vault_path: impl Into<PathBuf>) -> Result<Self> {
        let vault_path = vault_path.into();
        let mut manager = Self {
            capabilities: HashMap::new(),
            vault_path,
        };
        manager.reload()?;
        Ok(manager)
    }

    pub fn reload(&mut self) -> Result<()> {
        if !self.vault_path.exists() {
            return Ok(());
        }

        self.capabilities.clear();
        let entries = fs::read_dir(&self.vault_path)
            .context("Failed to read capabilities directory")?;

        for entry in entries {
            let entry = entry?;
            let path = entry.path();
            if path.is_dir() {
                let meta_file = path.join("meta.json");
                if meta_file.exists() {
                    let content = fs::read_to_string(&meta_file)
                        .context(format!("Failed to read meta file: {:?}", meta_file))?;
                    let meta: CapabilityMeta = serde_json::from_str(&content)
                        .context(format!("Failed to parse meta file: {:?}", meta_file))?;
                    self.capabilities.insert(meta.name.clone(), meta);
                }
            }
        }
        Ok(())
    }

    pub fn get_capabilities_description(&self) -> String {
        let mut desc = String::new();
        let mut sorted_meta: Vec<_> = self.capabilities.values().collect();
        sorted_meta.sort_by_key(|m| &m.name);

        for meta in sorted_meta {
            desc.push_str(&format!("{} ({}), ", meta.name, meta.description));
        }
        if !desc.is_empty() {
            desc.pop(); // Remove trailing comma
            desc.pop(); // Remove trailing space
        }
        desc
    }

    pub fn get_capability_path(&self, name: &str) -> Result<PathBuf> {
        let _meta = self.capabilities.get(name)
            .ok_or_else(|| anyhow::anyhow!("Capability {} not found", name))?;
        
        Ok(self.vault_path.join(name).join("script.sh"))
    }

    pub fn list_capabilities(&self) -> Vec<CapabilityMeta> {
        self.capabilities.values().cloned().collect()
    }
}
