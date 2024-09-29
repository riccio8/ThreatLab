rule LogStop
{
    meta:
        description = "Detects attempts to stop or tamper with Windows Event Logging"
        author = "blackmagic"
        date = "2024-09-24"
        threat_level = 5
        in_the_wild = true
        reference = "https://attack.mitre.org/techniques/T1562/002/"

    strings:
        $etw_disable = {FF 15 ?? ?? ?? ??}          // Generic call to external functions (suspicious)
        $etw_patch = {B8 00 00 00 00 C3}            // Common patch disabling ETW via function return
        $ntdll_string = "ntdll.dll"                 
        $etw_control_string = "EtwEventWrite"       // Function related to ETW event logging
        $etw_event_suppress = {BA 00 00 00 00 48 89 D8 C3}

    condition:
        // Match either the generic external call or the specific ETW tampering signatures
        $etw_disable or $etw_patch or $etw_control_string or $etw_event_suppress
}
