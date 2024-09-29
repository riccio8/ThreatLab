rule LogStop
{
    meta:
        description = "Stop the windows log system"
        threat_level = 3
        in_the_wild = true
        author = "riccioadami@gmail.com"

    strings:
        $d = {FF 15 ?? ?? ?? ??} 

    condition:
        $d
}
