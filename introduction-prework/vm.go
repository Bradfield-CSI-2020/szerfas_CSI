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
	//var r1 byte = 0
	//var r2 byte = 0

	var reg *byte
	var reg1 *byte
	var reg2 *byte

	//instrMap := make(map[byte]string)
	instructionSet := map[byte]string{
		1: "load",
		2: "store",
		3: "add",
		4: "subtract",
		255: "halt",
	}

	fmt.Printf("memory is %v\n", memory)
	// Keep looping, like a physical computer's clock
	for {
	//for i := 0; i < 20; i++ {
		pc = registers[0]
		instruction := memory[registers[0]]
		fmt.Printf("pc is %v and instruction is %v\n", pc, instructionSet[instruction])
		switch instruction {
		case Load:
			reg = getRegisterFromPC(&registers, memory[pc+1])
			setRegister(reg, memory, pc+2)
			registers[0] += 3
		case Store:
			//registers[0]++
			//switch memory[registers[0]] {
			//case 0x01: reg = &registers[1]
			//case 0x02: reg = &registers[2]
			//default:
			//	panic("invalid register value given: " + string(memory[registers[0]]))
			//}
			//fmt.Printf("pc is %v and register addr is %v\n", registers[0], reg)
			//registers[0]++
			//memory[memory[registers[0]]] = *reg
			//fmt.Printf("stored value: %v at 'mem addr': %v\n", memory[registers[0]], registers)
			//registers[0]++
			reg = getRegisterFromPC(&registers, memory[pc+1])
			storeFromRegister(reg, memory, pc+2)
			registers[0] += 3
		case Add:
			//registers[0]++
			//switch memory[registers[0]] {
			//case 0x01: reg1 = &registers[1]
			//case 0x02: reg1 = &registers[2]
			//default:
			//	panic("invalid register value given: " + string(memory[registers[0]]))
			//}
//			fmt.Printf("registers addr: %v, reg1 addr %v, reg2 addr %v\n", &registers, &registers[1], &registers[2])
			//registers[0]++
			//switch memory[registers[0]] {
			//case 0x01: reg2 = &registers[1]
			//case 0x02: reg2 = &registers[2]
			//default:
			//	panic("invalid register value given: " + string(memory[registers[0]]))
			//}
//			fmt.Printf("reg1: %v, reg2: %v\n", *reg1, *reg2)
			//registers[0]++
//			fmt.Printf("Adding: value at first provided register is now: %v\n", *reg1)
			reg1 = getRegisterFromPC(&registers, memory[pc+1])
			reg2 = getRegisterFromPC(&registers, memory[pc+2])
			*reg1 = *reg1 + *reg2
			registers[0] += 3
		case Sub:
			fmt.Printf("subtracting\n")
			registers[0]++
			switch memory[registers[0]] {
			case 0x01: reg1 = &registers[1]
			case 0x02: reg1 = &registers[2]
			default:
				panic("invalid register value given: " + string(memory[registers[0]]))
			}
			registers[0]++
			switch memory[registers[0]] {
			case 0x01: reg2 = &registers[1]
			case 0x02: reg2 = &registers[2]
			default:
				panic("invalid register value given: " + string(memory[registers[0]]))
			}
			registers[0]++
			*reg1 = *reg1 - *reg2
		case Halt:
			fmt.Println("halting")
			return
		}
		// op := TODO // fetch the opcode

		// // decode and execute
		// switch op {
		// case Load:
		//   TODO
		// ...
	}
}

func getRegister(memory []byte, registers [3]byte) byte {
	var register byte
	switch memory[registers[0]] {
	case 0x01: register = registers[1]
	case 0x02: register = registers[2]
	}
	return register
}

func incrementPC(pc *byte) byte {
	*pc++
	return *pc
}

func setRegister(reg *byte, memory []byte, pc byte) {
	*reg = memory[memory[pc]]
}

func storeFromRegister(reg *byte, memory []byte, pc byte) {
	memory[memory[pc]] = *reg
}

func getRegisterFromPC(registers *[3]byte, val byte) *byte {
	fmt.Printf("inside getRegisterFromPC with array: %v and val: %v\n", *registers, val)
	fmt.Printf("registers addr: %v, reg1 addr %v, reg2 addr %v\n", &registers, &registers[1], &registers[2])
	regMap := map[byte]*byte{
		0x01: &registers[1],
		0x02: &registers[2],
	}
	if reg, ok := regMap[val]; ok {
		fmt.Printf("reg addr: %v, reg value: %v\n", reg, *reg)
		return reg
	}
	panic("invalid register value given: " + string(val))
}