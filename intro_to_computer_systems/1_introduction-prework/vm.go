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

	// Keep looping, like a physical computer's clock
	for {
		instruction := memory[registers[0]]
		switch instruction {
		case Load:
			load(&registers, &memory)
		case Store:
			store(&registers, &memory)
		case Add:
			add(&registers, &memory)
		case Sub:
			subtract(&registers, &memory)
		case Addi:
			addi(&registers, &memory)
		case Subi:
			subi(&registers, &memory)
		case Jump:
			pc := registers[0]
			registers[0] = memory[pc+1]
		case Beqz:
			beqz(&registers, &memory)
		case Halt:
			return
		default:
			fmt.Printf("did not find instruction\n")
			return
		}
	}
}

func load(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	setRegister(reg, *memory, pc+2)
	registers[0] += 3
}

func store(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	storeFromRegister(reg, *memory, pc+2)
	registers[0] += 3
}

func add(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg1 := getRegisterFromPC(registers, (*memory)[pc+1])
	reg2 := getRegisterFromPC(registers, (*memory)[pc+2])
	*reg1 = *reg1 + *reg2
	registers[0] += 3
}

func subtract(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg1 := getRegisterFromPC(registers, (*memory)[pc+1])
	reg2 := getRegisterFromPC(registers, (*memory)[pc+2])
	*reg1 = *reg1 - *reg2
	registers[0] += 3
}

func addi(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	addVal := (*memory)[pc+2]
	*reg = *reg + addVal
	registers[0] += 3
}

func subi(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	addVal := (*memory)[pc+2]
	*reg = *reg - addVal
	registers[0] += 3
}

func beqz(registers *[3]byte, memory *[]byte) {
	pc := registers[0]
	reg := getRegisterFromPC(registers, (*memory)[pc+1])
	if *reg == 0 {
		offset := (*memory)[pc+2]
		registers[0] += offset
	}
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
