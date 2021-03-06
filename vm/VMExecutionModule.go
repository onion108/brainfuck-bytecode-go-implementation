package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type VMExecutionModule struct {
	callStack                  *VMStack
	bfjmpStack                 *VMStack
	userStack                  *VMStack
	inputFile                  *os.File
	pc                         int
	buffer                     []byte
	pointer                    int
	stateCode                  int // Prepare for multibyte instructions support.
	multibyteParseInfo_length  int
	multibyteParseInfo_content []byte
}

func VMExecutionModuleInit() *VMExecutionModule {
	return &VMExecutionModule{
		VMStackInit(),
		VMStackInit(),
		VMStackInit(),
		nil,
		0,
		[]byte{0},
		0,
		0,
		0,
		nil,
	}
}

// BindToFile Bind a file to the virtual machine.
func (vm *VMExecutionModule) BindToFile(reader *os.File) {
	if vm != nil && reader != nil {
		vm.inputFile = reader
	}
}

// Reset reset the machine
func (vm *VMExecutionModule) Reset() {
	if vm != nil {
		vm.callStack = VMStackInit()
		vm.bfjmpStack = VMStackInit()
		vm.userStack = VMStackInit()
		vm.inputFile = nil
		vm.pc = 0
	}
}

// Execute Start the execution loop.
func (vm *VMExecutionModule) Execute() {
	// Prepare for resources
	stdin := bufio.NewReader(os.Stdin)
	randSeed := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(randSeed)
	// Check required variables
	if vm == nil || vm.inputFile == nil {
		return
	}
	// Initialize user variables
	loader := VMCodeLoaderInit(vm.inputFile)
	// Check the file head (16bytes).
	var header [16]byte
	var ok bool
	var expectedHeader [16]byte = [16]byte{
		0x27,
		0x26,
		0x4A,
		0x00, // Would be the version
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	for i := 0; i < 16; i++ {
		header[i], ok = loader.At(i)
		if !ok {
			// File ended when reading the file header.
			fmt.Fprintf(os.Stderr, "Unexpected eof at 0x%X. Aborted.", i)
			vm.exitWithException()
		}
		if header[i] != expectedHeader[i] {
			fmt.Fprintf(os.Stderr, "Unexpected header number %X (expected %X) at 0x%X. Abort.", header[i], expectedHeader[i], i)
			vm.exitWithException()
		}
	}
	// Start validate the file header
	//
	// Start execute
	for {
		b, ok := loader.At(vm.pc + 16)
		if !ok {
			break
		}
		// Parse b
		if vm.stateCode == 0 {
			// State 0, parsing the first code (Usually native brainfuck instructions).
			switch b {
			// Brainfuck Native Instructions
			case VMBFInstructionOutput:
				fmt.Printf("%c", vm.buffer[vm.pointer])
				break
			case VMBFInstructionInput:
				ch, _ := stdin.ReadByte()
				vm.buffer[vm.pointer] = byte(ch)
				break
			case VMBFInstructionIncrease:
				vm.buffer[vm.pointer]++
				break
			case VMBFInstructionDecrease:
				vm.buffer[vm.pointer]--
				break
			case VMBFInstructionFlag:
				vm.bfjmpStack.Push(vm.pc)
				break
			case VMBFInstructionJumpToFlag:
				if vm.bfjmpStack.Empty() {
					fmt.Fprintf(os.Stderr, "Stack empty.")
					vm.exitWithException()
				}
				if vm.buffer[vm.pointer] != 0 {
					newPc := vm.bfjmpStack.Pop()
					vm.pc = newPc - 1 // Because PC will automatically increase, so assign newPc - 1 to it.
				}
			case VMBFInstructionPointerL:
				if vm.pointer > 0 {
					vm.pointer--
				}
			case VMBFInstructionPointerR:
				if vm.pointer < len(vm.buffer)-1 {
					vm.pointer++
				} else {
					// Out-of-range, increase the length
					vm.buffer = append(vm.buffer, 0)
					vm.pointer++
				}
			// End BNI
			case VMNopInstruction:
				break
			case VMEBFInstructionRandomNum:
				// Generate a random and put it into the current memory space
				vm.buffer[vm.pointer] = byte(randomizer.Int() % 0x100)
				break
			case VMEBFInstructionExitProg:
				vm.inputFile.Close()
				os.Exit(0)
			case VMEBFInstructionCall:
				vm.callStack.Push(vm.pc)
				vm.pc = int(vm.buffer[vm.pointer] - 1)
				break
			case VMEBFInstructionReturn:
				if vm.callStack.Empty() {
					fmt.Fprintf(os.Stderr, "Stack empty.")
					vm.exitWithException()
				}
				vm.pc = int(vm.callStack.Pop() - 1)
				break
			case VMEBFInstructionPush:
				vm.userStack.Push(int(vm.buffer[vm.pointer]))
				break
			case VMEBFInstructionPop:
				if vm.userStack.Empty() {
					fmt.Fprintf(os.Stderr, "Stack empty.")
					vm.exitWithException()
				}
				vm.buffer[vm.pointer] = byte(vm.userStack.Pop())
				break
			case VMDebugInstructionBreakpoint:
				// TODO: Add more debugging features here
				vm.PrintMachineState()
				stdin.ReadByte()
				break
			case VMMultibytePrefix:
				// Prepare for multibyte instructions
				vm.stateCode = 1 // Prepare to parse the next length
				break
			default:
				// Unknown Command
				fmt.Fprintf(os.Stderr, "Unknown instruction at 0x%X: 0x%X. Abort. \n", vm.pc, b)
				vm.exitWithException()
			}
			vm.pc++
			continue
		} else if vm.stateCode == 1 {
			// Parse the length
			vm.multibyteParseInfo_length = int(b)
			vm.multibyteParseInfo_content = []byte{}
			// Move to the next step
			vm.stateCode = 2
			// Continue
			vm.pc++
			continue
		} else if vm.stateCode == 2 {
			// Parse it
			if vm.multibyteParseInfo_length > 0 {
				vm.multibyteParseInfo_length--
				vm.multibyteParseInfo_content = append(vm.multibyteParseInfo_content, b)
				vm.pc++
				continue
			} else {
				vm.stateCode = 3 // Call the multibyte parser
				vm.pc++
				continue
			}
		} else if vm.stateCode == 3 {
			// Parse the multibyte command
			command := vm.multibyteParseInfo_content
			if len(command) < 1 {
				fmt.Fprintf(os.Stderr, "Invalid multibyte command at 0x%X.\n", vm.pc)
				vm.exitWithException()
			}
			// Use command
			switch command[0] {
			// TODO: Complete the command parser
			case VMMultibyteAdd:
				if len(command) != 2 {
					fmt.Errorf("Unexpected Instruction Length: %d\n", len(command))
					vm.exitWithException()
				}
				vm.buffer[vm.pointer] += command[1]
				break
			case VMMultibyteSub:
				vm.assertInstructionLength(command, 2)
				vm.buffer[vm.pointer] -= command[1]
				break
			case VMMultibyteAssign:
				vm.assertInstructionLength(command, 2)
				vm.buffer[vm.pointer] = command[1]
				break
			default:
				// Undefined Command
				fmt.Errorf("Undefined Command %X \n", command[0])
				vm.exitWithException()
			}
			// End use
			// Reset the variables
			vm.multibyteParseInfo_content = nil
			vm.multibyteParseInfo_length = 0
			// No need to increase vm.pc
			continue
		}
	}
}

func (vm *VMExecutionModule) assertInstructionLength(command []byte, length int) {
	if len(command) != length {
		fmt.Errorf("Unexpected Instruction Length: %d\n", len(command))
		vm.exitWithException()
	}
}

func (vm *VMExecutionModule) exitWithException() {
	vm.PrintMachineState()
	vm.inputFile.Close()
	os.Exit(-1)
}

func (vm *VMExecutionModule) PrintMachineState() {
	fmt.Printf("========MACHINE STATE========\n")
	fmt.Printf("Buffer[size:%d]: %v\n^(0x%X at 0x%X)\n", len(vm.buffer), vm.buffer, vm.pointer, vm.buffer[vm.pointer])
	fmt.Printf("PSC: 0x%X\n", vm.stateCode)
	fmt.Printf("Brainfuck Jumping Stack: %v\n", vm.bfjmpStack.content)
	fmt.Printf("User Stack: %v\n", vm.userStack.content)
	fmt.Printf("Call Stack: %v\n", vm.callStack.content)
	fmt.Printf("Program Counter: 0x%X\n", vm.pc)
	fmt.Printf("=============================\n")
}
