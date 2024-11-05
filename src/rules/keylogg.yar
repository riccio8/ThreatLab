rule Detect_Keylogger_Patterns
{
    meta:
        description = "Detects patterns commonly used in keylogger software"
        author = "0x90"
        date = "2024-11-05"
    
    strings:
        $log1 = "GetAsyncKeyState"
        $log2 = "SaveLogs"
        $log3 = "KeyboardHook"
    
    condition:
        any of ($log*)
}
