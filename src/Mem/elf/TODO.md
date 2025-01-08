# TODO List for ELF Utils Project

This file outlines future enhancements and ideas to improve the functionality of the ELF Utils project.

---

## **Features to Add**

1. **Multiple Command Execution**
   - Allow users to execute multiple commands in a single invocation.
   - Example: 
     ```bash
     ./elfutils myfile.elf sections sym fileHeader
     ```

2. **Output Formats**
   - Add support for exporting output in different formats (besides JSON):
     - **XML**
     - **YAML**
     - Allow users to specify the format with a flag:
       ```bash
       ./elfutils myfile.elf sections --output-format=yaml
       ```

3. **Better Logging Options**
   - Log output files in specific directories based on OS:
     - **Linux:** `/var/log/`
     - **Windows:** Current working directory.
   - Add a customizable log file path option:
     ```bash
     ./elfutils myfile.elf sections --log-path=/custom/log/directory/
     ```

4. **Verbose Mode**
   - Implement a verbose mode (`--verbose` or `-v`) for detailed output of each step:
     ```bash
     ./elfutils myfile.elf sym --verbose
     ```

5. **Suggestion from user**

6. **Enhanced Help Command**
   - Include more examples in the `help` command.
   - Categorize commands based on functionality (e.g., debugging, symbol management, headers).

7. **Improved Output Parsing**
   - Provide better human-readable descriptions for commands like `machine` and `entryPoint`.
   - Link to relevant ELF documentation for advanced properties:
     - Example: `https://pkg.go.dev/debug/elf#Machine`

8. **Support for Compressed ELF Files**
   - Implement support for `.gz` and `.xz` compressed ELF binaries.

9. **Cross-Platform Compatibility**
   - Ensure seamless operation on **Windows**, **Linux**, and **macOS**.
   - Add more OS-specific functionalities (e.g., logging, directory paths).

10. **Interactive Mode**
    - Create an interactive shell to explore ELF files without needing to specify the file each time:
      ```bash
      ./elfutils myfile.elf
      ELFUtils> sections
      ELFUtils> headers
      ELFUtils> exit
      ```

---

## **Code Improvements**

1. **Error Handling**
   - Improve error messages and provide actionable feedback (Do u think it's needed?).

2. **Logging Enhancements**
   - Add timestamps to log files.
   - Example:
     ```json
     {
       "timestamp": "2025-01-08T12:34:56Z",
       "result": {...}
     }
     ```


---

## **Testing and Quality**

1. **Benchmarking**
   - Profile performance for large ELF files and optimize processing, this is missing.

---

## **Future Ideas**

1. **GUI Application**
   - Develop a graphical interface for easier ELF exploration.

2. **Web-Based Tool**
   - Build a web-based ELF file analyzer with drag-and-drop support.

3. **Integration with Debuggers**
   - Provide options to interact with popular debuggers like `gdb`.

4. **Support for Additional Formats**
   - Extend support to other binary formats like PE (Windows) or Mach-O (macOS).

---

## **References**

- [Go ELF Documentation](https://pkg.go.dev/debug/elf)
- [DWARF Debugging Standard](https://dwarfstd.org/)
- [ELF File Format Specification](https://refspecs.linuxfoundation.org/elf/)

---

### **Contributors**
If you have suggestions or would like to contribute, feel free to open a pull request or issue on the [GitHub repository](https://github.com/riccio8/ThreatLab).
