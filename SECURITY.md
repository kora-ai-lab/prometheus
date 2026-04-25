# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

1. **Do NOT** open a public GitHub issue for security vulnerabilities
2. Email security concerns to the maintainers through GitHub's private vulnerability reporting
3. Include as much detail as possible:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes (optional)

### What to Expect

- Acknowledgment within 48 hours
- Regular updates on progress
- Credit in the security advisory (unless you prefer anonymity)
- Public disclosure after a fix is available

## Security Features

Prometheus includes multiple layers of security:

- **Command auditing**: Every executed command is analyzed before execution
- **Rate limiting**: Maximum 10 command executions per second
- **SAST scanning**: Generated code is scanned for vulnerabilities
- **Encrypted vault**: Credentials stored with AES-256-GCM encryption
- **Namespace isolation**: Linux systems use syscall isolation for command execution

## Best Practices

When using Prometheus:

1. Keep your API keys in environment variables, never hardcode them
2. Review confirmation prompts for high-risk commands
3. Run in an isolated environment when testing untrusted code
4. Keep Prometheus updated to the latest version