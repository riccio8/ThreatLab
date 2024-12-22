/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


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
        $exec4 = "WinExec"
    
    condition:
        any of ($exec*)
}
