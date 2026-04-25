$pinfo = New-Object System.Diagnostics.ProcessStartInfo
$pinfo.FileName = "C:\Users\junio\OneDrive\AI AGENT HACKATON\Prometheus\go_version\bin\prometheus.exe"
$pinfo.Arguments = "echo test"
$pinfo.RedirectStandardOutput = $true
$pinfo.RedirectStandardError = $true
$pinfo.UseShellExecute = $false
$pinfo.CreateNoWindow = $true

$p = New-Object System.Diagnostics.Process
$p.StartInfo = $pinfo
$p.Start() | Out-Null

$stdout = $p.StandardOutput.ReadToEnd()
$stderr = $p.StandardError.ReadToEnd()
$timeout = 60
$p.WaitForExit($timeout) | Out-Null

"STDOUT: $stdout" | Out-File test_output.txt
"STDERR: $stderr" | Out-File -Append test_output.txt
"ExitCode: $($p.ExitCode)" | Out-File -Append test_output.txt