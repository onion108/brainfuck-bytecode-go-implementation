package main

const (
	// VMFileHeaderLength The header's length
	VMFileHeaderLength = 16
	// -----------Instructions------------
	VMBFInstructionOutput     = 0x18
	VMBFInstructionInput      = 0x19
	VMBFInstructionIncrease   = 0x1A
	VMBFInstructionDecrease   = 0x1B
	VMBFInstructionPointerR   = 0x1C
	VMBFInstructionPointerL   = 0x1D
	VMBFInstructionFlag       = 0x1E
	VMBFInstructionJumpToFlag = 0x1F

	VMEBFInstructionRandomNum = 0x20
	VMEBFInstructionExitProg  = 0x21
	VMEBFInstructionCall      = 0x22
	VMEBFInstructionReturn    = 0x23
	VMEBFInstructionPush      = 0x24
	VMEBFInstructionPop       = 0x25
	VMNopInstruction          = 0x00

	VMDebugInstructionBreakpoint = 0x30

	VMMultibytePrefix = 0xE0
	VMMultibyteAdd    = 0x1A
	VMMultibyteSub    = 0x1B
	VMMultibyteAssign = 0x1C
	VMMultibyteJmpB   = 0x1D
	VMMultibyteJmpI   = 0x1E
)
