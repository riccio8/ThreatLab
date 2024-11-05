rule Detect_Packed_Executable
{
    meta:
        description = "Detects common packer signatures in executables"
        author = "0x90"
        date = "2024-11-05"
    
    strings:
        $packer1 = "UPX0"           // UPX packer
        $packer2 = "FSG!"           // FSG packer
        $packer3 = "PECompact2"     // PECompact packer
        $packer4 = { 50 45 00 00 4C 01 03 00 00 00 00 00 } // Suspicious PE header
        
    condition:
        any of ($packer*)
}
