package security

type Pattern struct {
	Needle string
	Score  int
}

var dangerousPatterns = []Pattern{
	{Needle: "rm -rf /", Score: 100},
	{Needle: "curl", Score: 40},
	{Needle: "| sh", Score: 80},
	{Needle: "chmod 777 /", Score: 100},
	{Needle: "/etc/passwd", Score: 95},
	{Needle: ":(){ :|:& };:", Score: 100},
}
