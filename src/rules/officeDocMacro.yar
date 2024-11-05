rule Detect_Malicious_Macro_Code
{
    meta:
        description = "Detects suspicious macro code in Office documents"
        author = "0x90"
        date = "2024-11-05"
    
    strings:
        $macro1 = "AutoOpen"                   
        $macro2 = "Shell.Application"        
        $macro3 = "CreateObject(\"WScript.Shell\")" 
        $macro4 = "Regsvr32"                 
        $macro5 = "powershell -exec bypass"   
        
    condition:
        any of ($macro*)
}
