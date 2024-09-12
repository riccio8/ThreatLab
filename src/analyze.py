import os
import platform
import pefile as pe
import peutils as pes

if platform.system == 'Linux':
    os.system("pip3 install psutil")
    import psutil
else: 
  os.system("pip install psutil")
  import psutil


proc_name = input("[X]name of the process...: \n")


def get_process_path(pid):
    try:
        process = psutil.Process(pid)
        return process.exe()  
    except psutil.NoSuchProcess:
        return f"Process with PID {pid} not found."
    except psutil.AccessDenied:
        return "AccessDenied to the proces."
    except Exception as e:
        return str(e)


for processes in psutil.process_iter():
    if processes.name() == proc_name:
        print("[INFO]: ", processes, "with pid: \n", processes.pid)
        process_pid = processes.pid
        kill = input("[*] Do u want to kill the process? (Y/N): \n")
        if kill == lower("Y")
            end_process = os.system("taskkill /in" + str(processes.pid))
        else:
            print("[X] Analysis of", proc_name, ".")
            print("[X] Analysis could take some time, please wait...")
            time.sleep(2)
            print("[X] Analysis in progress...")
            path = get_process_path(process_pid)
            pes.peutils.is_suspicious(path)
            
else: 
    print("[x]No processes found...")

