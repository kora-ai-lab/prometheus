package security

import (
	"os"
)

type PermissionFinding struct {
	Path     string
	Mode     os.FileMode
	Severity string
	Issue    string
}

func CheckPermissions(paths []string) []*PermissionFinding {
	var findings []*PermissionFinding
	
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			continue
		}
		
		mode := info.Mode()
		
		if mode&0002 != 0 {
			findings = append(findings, &PermissionFinding{
				Path:     path,
				Mode:     mode,
				Severity: "high",
				Issue:    "world-writable",
			})
		}
		
		if mode&0044 != 0 {
			if isSensitive(path) {
				findings = append(findings, &PermissionFinding{
					Path:     path,
					Mode:     mode,
					Severity: "high",
					Issue:    "world-readable sensitive file",
				})
			}
		}
		
		if mode&os.ModeSetuid != 0 || mode&os.ModeSetgid != 0 {
			findings = append(findings, &PermissionFinding{
				Path:     path,
				Mode:     mode,
				Severity: "high",
				Issue:    "SUID/SGID set",
			})
		}
	}
	
	return findings
}

func CheckSensitiveDirs(dirs []string) []*PermissionFinding {
	var findings []*PermissionFinding
	sensitive := []string{"/etc/passwd", "/etc/shadow"}
	
	for _, dir := range sensitive {
		findings = append(findings, CheckPermissions([]string{dir})...)
	}
	
	return findings
}

func isSensitive(path string) bool {
	sensitive := []string{".ssh", ".aws", "credentials", "secret", "key", "vault"}
	for _, s := range sensitive {
		if contains(path, s) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		 contains(s[1:], substr))))
}