import os
import platform
import pefile as pe
import peutils as pes
import time
import psutil
import argparse

parser = argparse.ArgumentParser(description='Process analysis tool.')
parser.add_argument('--logfile', help='Optional. Save output to a log file. Usage: python proc_analysis.py --logfile=log.txt', default=None)
args = parser.parse_args()

def log(message):
    print(message)
    if args.logfile:
        with open(args.logfile, 'a') as logfile:
            logfile.write(message + '\n')

if platform.system() == 'Linux':
    os.system("clear")
else:
    os.system("cls")

proc_name = input("[X] Name of the process...: \n")

def get_process_path(pid):
    try:
        process = psutil.Process(pid)
        return process.exe()  
    except psutil.NoSuchProcess:
        return f"[ERROR] Process with PID {pid} not found."
    except psutil.AccessDenied:
        return "[ERROR] AccessDenied to the process."
    except Exception as e:
        return str(e)

def main():  
    found_process = False
    try: 
        for processes in psutil.process_iter():
            if processes.name() == proc_name:
                found_process = True
                log(f"[INFO]: {processes}")
                process_pid = processes.pid
                kill = input("[*] Do you want to kill the process? (Y/N): \n")
                if kill.lower() == "y":
                    os.system("taskkill /pid " + str(processes.pid))
                else:
                    log(f"[X] Analysis of {proc_name}.")
                    log("[X] Analysis could take some time, please wait...")
                    time.sleep(2)
                    log("[X] Analysis in progress...")

                    path1 = get_process_path(process_pid)

                    try:
                        path = pe.PE(path1, True)
                        suspicious = pes.is_suspicious(path)
                        if suspicious:
                            log(f"[INFO] {proc_name} is suspicious...")
                            log("[INFO] It could be that: import tables are in unusual locations, or section names are unrecognized, or there is a presence of long ASCII strings....")
                        else:
                            log("[INFO] File is not suspicious...\n")
                            continue

                        time.sleep(3)
                        log("[X] Analyzing signature... \n")
                        time.sleep(1)
                        valid = pes.is_valid(path)
                        if not valid:
                            log("[INFO] File's signature is not valid...\n")
                        else:
                            log("[INFO] File's signature is valid...\n")

                        log("[X] Analyzing sections...\n")

                        for section in path.sections:
                            time.sleep(1)
                            log(f"[INFO] Section: {section.Name}, Virtual Address: {hex(section.VirtualAddress)}, Virtual Size: {hex(section.Misc_VirtualSize)}, Size of Raw Data: {section.SizeOfRawData}")

                        log("[INFO] Analysis of the libraries... \n")
                        time.sleep(1)

                        path.parse_data_directories()
                        for entry in path.DIRECTORY_ENTRY_IMPORT:
                            log(f"[INFO] Libraries loaded: {entry.dll}")
                            for imp in entry.imports:
                                log(f"\t {hex(imp.address)} {imp.name}")

                        log("[X] Analyzing the file format... \n")

                        time.sleep(2)    

                        strip = pes.is_probably_packed(path)
                        if strip:   
                            log("[INFO] File is probably packed or contains compressed data... \n")
                        else:
                            log("[INFO] File isn't compressed or packed... \n")

                        log("[X] Saving the PE dump file... \n")

                        with open(proc_name + "_dump.txt", 'w') as dump:
                            dump.write(path.dump_info())
                            dump.close()
                    except Exception as e:
                        log("[ERROR] " + str(e))

        if not found_process:
            log(f"[INFO] No processes found with the name {proc_name}.")
    except Exception as e:
        log("[ERROR] " + str(e))
   
if __name__ == "__main__":
    main()

def net_analysis():
    pass
