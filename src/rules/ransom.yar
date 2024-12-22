/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


rule Detect_Ransomware_Strings
{
    meta:
        description = "Detects common ransomware strings"
        author = "0x90"
        date = "2024-11-05"
    
    strings:
        $ransom1 = "All your files have been encrypted"
        $ransom2 = "Send Bitcoin to the following address"
        $ransom3 = "Decryptor key"
        $extension = ".locked"
    
    condition:
        any of ($ransom*) or $extension
}
