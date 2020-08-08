package vm

import "fmt"

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	var pc byte
	var reg *byte
	var reg1 *byte
	var reg2 *byte

	// Keep looping, like a physical computer's clock
	for {
		pc = registers[0]
		instruction := memory[registers[0]]
		switch instruction {
		case Load:
			load(&registers, &memory, pc)
		case Store:
			reg = getRegisterFromPC(&registers, memory[pc+1])
			storeFromRegister(reg, memory, pc+2)
			registers[0] += 3
			store(&registers, &memory, pc)
		case Add:
			reg1 = getRegisterFromPC(&registers, memory[pc+1])
			reg2 = getRegisterFromPC(&registers, memory[pc+2])
			*reg1 = *reg1 + *reg2
			registers[0] += 3
		case Sub:
			reg1 = getRegisterFromPC(&registers, memory[pc+1])
			reg2 = getRegisterFromPC(&registers, memory[pc+2])
			*reg1 = *reg1 - *reg2
			registers[0] += 3
		case Halt:
			fmt.Println("halting")
			return
		}
	}
}

func load(registers *[3]byte, memory *[]byte, pc byte) {
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	setRegister(reg, *memory, pc+2)
	registers[0] += 3
}

func store(registers *[3]byte, memory *[]byte, pc byte) {
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	storeFromRegister(reg, *memory, pc+2)
	registers[0] += 3
}

func setRegister(reg *byte, memory []byte, pc byte) {
	*reg = memory[memory[pc]]
}

func storeFromRegister(reg *byte, memory []byte, pc byte) {
	memory[memory[pc]] = *reg
}

func getRegisterFromPC(registers *[3]byte, val byte) *byte {
	regMap := map[byte]*byte{
		0x01: &registers[1],
		0x02: &registers[2],
	}
	if reg, ok := regMap[val]; ok {
		return reg
	}
	panic("invalid register value given: " + string(val))
}