rule Detect_Suspicious_PowerShell
{
    meta:
        description = "Detects PowerShell scripts with suspicious behavior"
        author = "0x90"
        date = "2024-11-05"
    
    strings:
        $exec1 = "Invoke-Expression"
        $exec2 = "DownloadString"
        $exec3 = "System.Net.WebClient"
    
    condition:
        any of ($exec*)
}
