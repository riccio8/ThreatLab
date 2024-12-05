
# Cutter and Rizin Cheat Sheet

## Cutter Cheat Sheet

Cutter is a GUI tool for reverse engineering powered by Rizin.

### Navigation and Views:
- **Ctrl + 1**: Switch to "Disassembly" view.
- **Ctrl + 2**: Switch to "Graph" view.
- **Ctrl + 3**: Switch to "Hexdump" view.
- **Ctrl + 4**: Switch to "Strings" view.
- **Ctrl + 5**: Switch to "Symbols" view.
- **Ctrl + 6**: Switch to "Functions" view.

### Control the Execution Flow:
- **F7**: Execute until the next breakpoint.
- **F8**: Execute the program until it completes (single step).
- **Shift + F7**: Step into the code.
- **Shift + F8**: Step out of the code.

### Exploring the Disassembly:
- **Ctrl + G**: Go to a specific address.
- **Ctrl + P**: Go to a function (if identified).
- **Ctrl + B**: Show the "Bookmarks" window.
- **Shift + F**: Search a function in the disassembly.

### Analysis and Searching:
- **Ctrl + F**: Search for text, strings, or patterns.
- **Ctrl + R**: Search for functions or symbols.
- **Shift + R**: Search for references to a symbol or address.
- **Ctrl + I**: View information about a selected address or symbol.
- **Shift + I**: View information for the current address.

### Advanced Analysis Commands:
- **Ctrl + A**: Start a new program analysis.
- **Ctrl + N**: Create a new project.
- **Ctrl + E**: Export the analysis results.
- **Ctrl + M**: Open the memory search module.
- **Ctrl + T**: Modify data type analysis.

### Scripting and Advanced Functions:
- **Shift + T**: Open the scripting terminal.
- **Ctrl + Shift + S**: Run a script command for analysis.

### Memory and Section Management:
- **Ctrl + Shift + M**: Show memory map.
- **Ctrl + Shift + F**: Show control flow graph.
- **Ctrl + D**: Open the "Data" window to explore binary data.
- **Ctrl + L**: Show a list of loaded files.

### Additional Tools:
- **Shift + V**: Show disassembly in graphical form.
- **Ctrl + W**: Open disassembly options menu.
- **Shift + S**: Enable symbol view for function analysis.

### Other Useful Commands:
- **Ctrl + Shift + C**: Copy the selected text from disassembly or output.
- **Ctrl + Shift + X**: Export the selected output.
- **Ctrl + Q**: Exit Cutter.

---

## Rizin Cheat Sheet

Rizin is a reverse engineering framework forked from Radare2. It also has a GUI tool, Cutter, powered by Rizin.

### Using Cutter

Cutter is a GUI tool for reverse engineering powered by Rizin.
It can also have a decompiler, so it’s recommended to use it first.

```bash
cutter <file>
```

To use the Ghidra decompiler, install the package:

```bash
sudo apt install rizin-plugin-ghidra
# or
sudo apt install rz-ghidra
```

### Start Debugging

Start Rizin with the program:

```bash
rizin ./example
```

- **Debug mode**: `rizin -d ./example`
- **Write mode**: `rizin -w ./example`

### Analyze

Analyze the program after starting the debugger:

- **Analyze all calls**: 
  ```bash
  aaa
  ```
- **Analyze function**: 
  ```bash
  af
  ```
- **List all functions**: 
  ```bash
  afl
  afl | grep main
  ```
- **Show address of current function**: 
  ```bash
  afo
  ```

### Print Usage

- **Print usage**: 
  ```bash
  ?
  ```
- Add `?` suffix to print the usage of a specific command:
  ```bash
  i?
  p?
  ```

### Visual Mode

You can enter visual mode for more intuitive operation:

- **Enter visual mode**: 
  ```bash
  v
  ```
- **Visual Debugger Mode**: 
  ```bash
  Vpp
  ```

### Basic Debugging Commands

- **Toggle print mode**: `p` or `P`
- **Step**: `s`
- **Toggle cursor mode**: `c`
- **Exit**: `q`
- **Enable regular Rizin commands**: `:`

### Debugging Flow

- **Step**: 
  ```bash
  ds
  ```
- **Step 3 times**: 
  ```bash
  ds 3
  ```
- **Step back**: 
  ```bash
  dsb
  ```

### Breakpoints

- **Setup a breakpoint**: 
  ```bash
  db @ 0x8048920
  ```
- **Remove a breakpoint**: 
  ```bash
  db @ -0x8048920
  ```
- **Remove all breakpoints**: 
  ```bash
  db-*
  ```
- **List all breakpoints**: 
  ```bash
  dbl
  ```

- **Continue to execute until breakpoint**: 
  ```bash
  dc
  ```
- **Continue until syscall**: 
  ```bash
  dcs
  ```

- **Read all registers values**: 
  ```bash
  dr
  dr=
  ```
- **Read given register value**: 
  ```bash
  dr eip
  dr rip
  ```
- **Set a register value**: 
  ```bash
  dr eax=24
  ```
- **Show register references**: 
  ```bash
  drr
  ```

### Seeking

- **Print current address**: 
  ```bash
  s
  ```
- **Seek to given function**: 
  ```bash
  s main
  s sym.main
  ```
- **Seek to given address**: 
  ```bash
  s 0x1360
  s 0x0x00001360
  ```
- **Seek to register address**: 
  ```bash
  s esp
  s esp+0x40
  s rsp
  s rsp+0x40
  ```

- **Seek 8 positions**: 
  ```bash
  sd 8
  ```

- **Show the seek history**: 
  ```bash
  sh
  ```
- **Undoing**: 
  ```bash
  shu
  ```
- **Redoing**: 
  ```bash
  shr
  ```

### Printing

- **Disassemble at current address**: 
  ```bash
  pd
  ```
- **Disassemble 10 instructions at current address**: 
  ```bash
  pd 10
  ```
- **Disassemble all possible opcodes at current address**: 
  ```bash
  pda
  pda 10
  ```
- **Disassemble at the given function**: 
  ```bash
  pd @ main
  pd 20 @ main
  ```
- **Disassemble a function at current address**: 
  ```bash
  pdf
  ```
- **Disassemble at given address**: 
  ```bash
  pdf @ 0x401005
  ```
- **Disassemble the main function**: 
  ```bash
  pdf @ main
  ```

- **Print string at address**: 
  ```bash
  ps @ 0x2100
  ```
- **Print zero-terminated string**: 
  ```bash
  psz @0x2100
  ```

- **Show 200 hex bytes**: 
  ```bash
  px 200
  ```
- **Show hex bytes at given register**: 
  ```bash
  px @ eip
  px @ esp
  ```

### Decompiling

To decompile functions, install the Ghidra plugin:

```bash
sudo apt install rizin-plugin-ghidra
```

- **Decompile the "main" function**: 
  ```bash
  pdg @ main
  ```

### Writing

To write, use the `-w` option when starting the debugger:

- **Write string at address**: 
  ```bash
  w Hello World\n @ 0x2100
  ```
- **Write opcodes at given address**: 
  ```bash
  wa 'mov eax, 1' @ 0x2100
  wa 'mov byte [rbp-0x1], 0x61' @ 0x2100
  ```

### Expressions

- **Evaluate expression**: 
  ```bash
  ?vi 0x000011a4
  ```
- **Evaluate simple expression**: 
  ```bash
  ?vi 1+2
  ```

### Information about Binary File

- **Information about the binary file**: 
  ```bash
  i
  ```
- **Show all summary**: 
  ```bash
  ia
  ```
- **Show main address**: 
  ```bash
  iM
  ```
- **Symbols**: 
  ```bash
  is
 

 ```
- **Show metadata**: 
  ```bash
  iM
  ```
- **Show dependencies**: 
  ```bash
  iD
  ```

---

## Finding ROP Gadgets in Rizin and Cutter

### Rizin Commands for Finding ROP Gadgets:

1. **Search for ROP Gadgets**:
   To search for ROP gadgets, use the `rop` plugin in Rizin, which can find sequences of instructions that can be used for a ROP chain.

   ```bash
   rz-rop
   ```

   This command will search for available ROP gadgets in the binary.

2. **List Available ROP Gadgets**:
   After finding the gadgets, list them using the following command:

   ```bash
   rz-rop -l
   ```

   This will display all available ROP gadgets.

3. **Search for Specific Gadgets**:
   To search for specific gadgets, filter by instruction type (e.g., `pop`, `ret`):

   ```bash
   rz-rop -s pop
   ```

   This command will search for gadgets containing the instruction `pop`.

4. **Using the Rizin Console**:
   You can also use the `rop` command directly in the Rizin console to search for ROP gadgets in the binary:

   ```bash
   rop
   ```

   This will search for ROP gadgets in the program's memory.

5. **Find ROP Gadgets in a Specific Section**:
   To search for ROP gadgets in a specific section (e.g., `.text`):

   ```bash
   rop -s .text
   ```

6. **Find ROP Gadgets in a Range of Addresses**:
   If you know a range of addresses where you expect ROP gadgets, specify them like this:

   ```bash
   rop 0x08048000 0x08049000
   ```

   This will search for ROP gadgets between the specified address range.

---

### Cutter (GUI) for Finding ROP Gadgets:

Cutter provides a graphical interface to search for ROP gadgets.

1. **Open Cutter and Load the Binary**:
   After loading the binary in Cutter, navigate to the **Rizin** console by clicking on the **Console** tab at the bottom panel.

2. **Use the `rop` Plugin**:
   In the Cutter console, type the following command to find ROP gadgets:

   ```bash
   rz-rop
   ```

3. **List Available Gadgets**:
   To list the available gadgets:

   ```bash
   rz-rop -l
   ```

4. **Filter Gadgets**:
   You can search for specific types of gadgets, such as `pop` or `ret`:

   ```bash
   rz-rop -s pop
   ```

5. **Explore Gadgets in Cutter**:
   Cutter's **"Gadgets"** tab also allows you to visually inspect and analyze ROP gadgets, making it easier to select the useful ones for constructing ROP chains.

---


# References

- Cutter: https://cutter.re
- Rizin: https://rizin.re
```
