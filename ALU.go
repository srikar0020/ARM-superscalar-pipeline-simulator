package main

var IMvariable uint64

func ALU() {
	if postALUBuff[0] == -1 && len(PreALUBuff) > 0 {
		var k = PreALUBuff[0]
		PreALUBuff = dequeue(PreALUBuff)
		postALUBuff[1] = int64(line[k].rd)
		postALUBuff[2] = k

		if line[k].typeofInstruction == "R1" {
			if line[k].op == "AND" {
				postALUBuff[0] = int64(sim.r[line[k].rm] & sim.r[line[k].rn])
			} else if line[k].op == "ADD" {
				//sim.r[line[k].rd] = sim.r[line[k].rm] + sim.r[line[k].rn]
				postALUBuff[0] = int64(sim.r[line[k].rm] + sim.r[line[k].rn])
			} else if line[k].op == "ORR" {
				//sim.r[line[k].rd] = sim.r[line[k].rm] | sim.r[line[k].rn]
				postALUBuff[0] = int64(sim.r[line[k].rm] | sim.r[line[k].rn])
			} else if line[k].op == "SUB" {
				//sim.r[line[k].rd] = sim.r[line[k].rn] - sim.r[line[k].rm]
				postALUBuff[0] = int64(sim.r[line[k].rn] - sim.r[line[k].rm])
			} else if line[k].op == "EOR" {
				//sim.r[line[k].rd] = sim.r[line[k].rm] ^ sim.r[line[k].rn]
				postALUBuff[0] = int64(sim.r[line[k].rm] ^ sim.r[line[k].rn])
			}

		} else if line[k].typeofInstruction == "R2" {
			if line[k].op == "LSR" {
				//sim.r[line[k].rd] = sim.r[line[k].rn] >> line[k].shamt
				postALUBuff[0] = int64(sim.r[line[k].rn] >> line[k].shamt)
			} else if line[k].op == "LSL" {
				//sim.r[line[k].rd] = sim.r[line[k].rn] << line[k].shamt
				postALUBuff[0] = int64(sim.r[line[k].rn] << line[k].shamt)
			} else if line[k].op == "ASR" {
				//sim.r[line[k].rd] = int64(sim.r[line[k].rn]) / uint64(math.Pow(2, float64(line[k].shamt)))
				//sim.r[line[k].rd] = uint64(int64(sim.r[line[k].rn]) >> line[k].shamt)
				postALUBuff[0] = int64(uint64(int64(sim.r[line[k].rn]) >> line[k].shamt))
			}

		} else if line[k].typeofInstruction == "I" {
			if line[k].op == "ADDI" {
				//sim.r[line[k].rd] = sim.r[line[k].rn] + uint64(line[k].immN)
				postALUBuff[0] = int64(sim.r[line[k].rn] + uint64(line[k].immN))
			} else if line[k].op == "SUBI" {
				//sim.r[line[k].rd] = sim.r[line[k].rn] - uint64(line[k].immN)
				postALUBuff[0] = int64(sim.r[line[k].rn] - uint64(line[k].immN))
			}

		} else if line[k].typeofInstruction == "IM" {
			if line[k].op == "MOVZ" {
				IMvariable = 0
				IMvariable = IMvariable + uint64(line[k].field<<line[k].shamt)
				//sim.r[line[k].rd] = 0
				//sim.r[line[k].rd] = sim.r[line[k].rd] + uint64(line[k].field<<line[k].shamt)
				postALUBuff[0] = int64(IMvariable)

			} else if line[k].op == "MOVK" {
				IMvariable = sim.r[line[k].rd]

				if line[k].shamt == 0 {
					IMvariable = IMvariable & 0xFFFFFFFFFFFF0000
					IMvariable = IMvariable + uint64(line[k].field<<line[k].shamt)
				} else if line[k].shamt == 16 {
					IMvariable = IMvariable & 0xFFFFFFF0000FFFF
					IMvariable = IMvariable + uint64(line[k].field<<line[k].shamt)
				} else if line[k].shamt == 32 {
					IMvariable = IMvariable & 0xFFFF0000FFFFFFFF
					IMvariable = IMvariable + uint64(line[k].field<<line[k].shamt)
				} else if line[k].shamt == 48 {
					IMvariable = IMvariable & 0x0000FFFFFFFFFFFF
					IMvariable = IMvariable + uint64(line[k].field<<line[k].shamt)
				}
				postALUBuff[0] = int64(IMvariable)
			}

		}
	}
}
