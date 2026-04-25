package security

import (
	"net"
	"strconv"
	"time"
)

type PortResult struct {
	Port    int
	Service string
	Risk    string
}

var portServices = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:  "http",
	110: "pop3",
	143: "imap",
	443: "https",
	445: "smb",
	993: "imaps",
	995: "pop3s",
	1433: "mssql",
	1521: "oracle",
	3306: "mysql",
	3389: "rdp",
	5432: "postgres",
	6379: "redis",
	8080: "http-proxy",
	8443: "https-alt",
}

var highRiskPorts = []int{23, 135, 136, 137, 138, 139, 445, 3389, 5900}
var mediumRiskPorts = []int{21, 110, 143, 1433, 1521, 3306, 5432}

func ScanPort(host string, port int) (*PortResult, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return nil, err
	}
	conn.Close()

	service := portServices[port]
	if service == "" {
		service = "unknown"
	}

	risk := "low"
	for _, p := range highRiskPorts {
		if port == p {
			risk = "high"
			break
		}
	}
	if risk == "low" {
		for _, p := range mediumRiskPorts {
			if port == p {
				risk = "medium"
				break
			}
		}
	}

	return &PortResult{
		Port:    port,
		Service: service,
		Risk:    risk,
	}, nil
}

func ScanPorts(host string, ports []int) []*PortResult {
	var results []*PortResult
	for _, port := range ports {
		result, err := ScanPort(host, port)
		if err == nil {
			results = append(results, result)
		}
	}
	return results
}

func ScanCommonPorts(host string) []*PortResult {
	ports := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 993, 995, 1433, 1521, 3306, 3389, 5432, 6379, 8080, 8443}
	return ScanPorts(host, ports)
}