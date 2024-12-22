/*
 * Copyright 2023-2024 Riccardo Adami. All rights reserved.
 * License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
 */


rule Encrypt_Attempt
{
    meta:
        description = "Detects attempts to encrypt"
        author = "0x90"
        date = "2024-09-24"
        threat_level = 5
        in_the_wild = true
        reference = "https://attack.mitre.org/techniques/T1486/"
    
    strings:
        $api_call_1 = "CryptEncrypt"
        $api_call_2 = "Rijndael"
        $api_call_3 = "AES"
        $api_call_4 = "EncryptFileA"
        $suspicious_string = "encrypt"

    condition:
        any of them
}
