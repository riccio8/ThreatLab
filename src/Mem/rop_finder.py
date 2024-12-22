"""
Copyright 2023-2024 Riccardo Adami. All rights reserved.
License: https://github.com/riccio8/ThreatLab/blob/main/LICENSE
"""


import psutil  # type: ignore
import os
import time
import platform


def read_bin_and_find_ret(file_path):
    try:
        with open(file_path, 'rb') as f:
            binary_data = f.read()

        hex_data = binary_data.hex()
        print(f"Hexadecimal content of the file: \n{hex_data}")
        print("------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
        time.sleep(2)

        ret_offsets = []
        for i, byte in enumerate(binary_data):
            if byte == 0xC3:
                ret_offsets.append(i)

        if ret_offsets:
            print(f"'ret' instructions found at offsets: {ret_offsets}")
        else:
            print("No 'ret' instructions found.")
    except FileNotFoundError:
        print(f"[ERROR] File {file_path} not found.")
    except Exception as e:
        print(f"[ERROR] An error occurred: {e}")

def get_process_path(proc_name):
    try:
        for proc in psutil.process_iter(['pid', 'name', 'exe']):
            if proc.info['name'].lower() == proc_name.lower():
                return proc.info['exe']
        print(f"[ERROR] Process {proc_name} not found.")
        return None
    except psutil.NoSuchProcess:
        print(f"[ERROR] No such process: {proc_name}.")
    except psutil.AccessDenied:
        print(f"[ERROR] Access denied to the process: {proc_name}.")
    except Exception as e:
        print(f"[ERROR] An unexpected error occurred: {e}")
        return None

def main():
    proc_name = input("[X] Name of the process: ").strip()
    file_path = get_process_path(proc_name)
    
    if file_path:
        if platform.system() == 'Linux':
            os.system("clear")
        else:
            os.system("cls")
        print(f"Analyzing binary at: {file_path}")
        time.sleep(2)
        read_bin_and_find_ret(file_path)

if __name__ == "__main__":
    main()
