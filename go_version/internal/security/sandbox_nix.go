//go:build !linux
// +build !linux

package security

func canUseNamespaces() bool {
	return false
}

func newNamespaceSandbox(cfg SandboxConfig) Sandbox {
	return nil
}