import os
import platform

if platform.system == 'Linux':
    os.system("pip3 install psutil")
    import psutil
else: 
  os.system("pip install psutil")
  import psutil

proc_name = input("type the name of the process u want to analyze, like notepad.exe: \n")

for processes in psutil.process_iter():
    if processes.name() == proc_name:
        print(processes, "with pid: \n", processes.pid)
        process_pid = processes.pid
        
