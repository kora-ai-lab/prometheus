use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct CapabilityMeta {
    pub name: String,
    pub description: String,
    pub command: String,
    pub version: String,
}
