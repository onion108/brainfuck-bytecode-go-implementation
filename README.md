# Brainfuck VM

It's a Golang version of my old project [Brainfuck bytecode](https://github.com/onion108/Brainfuck-bytecode).

## Bytecode Format
### Header(16-Bytes)
```
27 26 4A [Version Number(1-byte)] 
00 00 00 00
00 00 00 00
00 00 00 00
```

### Instructions
#### Single-byte Instructions
Commands(Standard Brainfuck Commands - start with `0x1`):<br>
```
18: bf_output(.)
19: bf_input(,)
1A: bf_inc(+)
1B: bf_dec(-)
1C: bf_rshift(>)
1D: bf_lshift(<)
1E: bf_flag([)
1F: bf_jnz(])
```
Extended Brainfuck Commands: 
```
20 - give a random number to the current pointed memory.
21 - Exit the program.
24 - Push the value in the current pointed memory to the stack.
25 - Pop the value from stack to the current pointed memory.
22 - Call a function that the position is the content in the current pointed memory.
23 - Return the from the function.
```
Debug Commands:
```
30 - Set a breakpoint to show some messages about the vm.
```
#### Multi-byte commands (Coming soon)

Multi-byte commands are all in the follow format:
```
E0 [length(Excluding E0 and the byte comes after E0):byte] [...instruction content&arguments]
```
For example:
```
E0 02 00 00
```
