package security

type SandboxLevel string

const (
	SandboxNone  SandboxLevel = "none"
	SandboxAudit SandboxLevel = "audit"
)
