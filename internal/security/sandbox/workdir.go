//go:build !linux
// +build !linux

package sandbox

func canUseNamespaces() bool {
	return false
}

func newNamespaceSandbox(cfg SandboxConfig) Sandbox {
	return nil
}