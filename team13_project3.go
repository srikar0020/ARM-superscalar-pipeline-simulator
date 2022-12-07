package main

//go run team13_project3.go -i test1_bin.txt -o team13_out
import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	typeofInstruction string
	rawInstruction    string
	linevalue         int64
	programCnt        int64
	opcode            int64
	op                string
	rd                uint8
	rn                uint8
	rm                uint8
	shamt             uint8
	address           uint16
	offset            int64
	immN              int16
	immS              string
	field             int64
}
type Simulation struct {
	r           [32]uint64
	cycleN      int32
	OffsetValue int64
	DataMem     []uint64
}

type lineN struct {
	valid int8
	dirty int8
	tag   int64
	data1 uint64
	data2 uint64
}
type setN struct {
	lru  int8
	line [2]lineN
}
type cacheN struct {
	sets [4]setN
}

var CachelineAddress [8]uint64

var postMemBuff = [3]int64{-1, -1, -1} //first number is value, second is instr index
var postALUBuff = [3]int64{-1, -1, -1} //first number is value, second is instr index
var PreMemBuff []int64                 // two instruction indexes
var PreALUBuff []int64
var PreIssueBuff []int64
var breakflag bool
var cache cacheN
var tempSlice []uint64
var line []Instruction
var sim Simulation
var StrInst []string
var k int
var programCountConst int64
var programCount int64
var startData, endData int64
var startPrgCnt, endPrgCnt int64
var startPrgCntFlag bool

func enqueue(queue []int64, element int64) []int64 {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}
func dequeue(queue []int64) []int64 {
	//element := queue[0] // The first element is the one to be dequeued.
	if len(queue) == 1 {
		var tmp = []int64{}
		return tmp
	}
	return queue[1:] // Slice off the element once it is dequeued.
}

func registers() string {
	var registersLine string = "registers:\n"
	for j := 0; j < 4; j++ {
		switch j {
		case 0:
			registersLine += "r00:\t"
		case 1:
			registersLine += "r08:\t"
		case 2:
			registersLine += "r16:\t"
		case 3:
			registersLine += "r24:\t"
		}
		for i := 0; i < 8; i++ {
			registersLine += strconv.FormatInt(int64(sim.r[i+(j*8)]), 10) + "\t"
		}
		registersLine += "\n"
	}
	registersLine += "\n"
	return registersLine
}
func firstLine() string {
	var FLine string = "====================" + "\n\n"
	FLine += "cycle:" + strconv.FormatInt(int64(sim.cycleN), 10) + "\t"
	FLine += strconv.FormatInt(int64(line[k].programCnt), 10) + "\t"
	FLine += line[k].op + "\t"
	return FLine
}

func r1Sim() string {
	sim.cycleN++
	var SecondLine string = "	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rn), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rm), 10) + "\n"
	return firstLine() + SecondLine + registers() + "\n" + DataOut()

}
func r2Sim() string {
	sim.cycleN++
	var SecondLine string = "	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rn), 10)
	SecondLine += ", #" + strconv.FormatInt(int64(line[k].shamt), 10) + "\n"
	return firstLine() + SecondLine + registers() + "\n" + DataOut()

}

func immediate() string {
	sim.cycleN++
	var SecondLine string = " 	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rn), 10)
	SecondLine += ", #" + strconv.FormatInt(int64(line[k].immN), 10) + "\n"

	return firstLine() + SecondLine + registers() + "\n" + DataOut()
}

func BOffset() string {
	sim.cycleN++
	var SecondLine string
	if line[k].op == "B" {
		SecondLine = "#" + strconv.FormatInt(int64(line[k].offset), 10) + "\n"
	} else {
		SecondLine = "R" + strconv.FormatInt(int64(line[k].rd), 10) + "\t"
		SecondLine += "#" + strconv.FormatInt(int64(line[k].offset), 10) + "\n"
	}
	return firstLine() + SecondLine + registers() + "\n" + DataOut()
}
func Move() string {
	sim.cycleN++
	var SecondLine string
	SecondLine += " 	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", " + strconv.FormatInt(line[k].field, 10)
	SecondLine += ", LSL " + strconv.FormatInt(int64(line[k].shamt), 10) + "\n"
	return firstLine() + SecondLine + registers() + "\n" + DataOut()
}

func DataManipulation() string {
	sim.cycleN++
	var SecondLine string
	SecondLine += "	R" + strconv.FormatInt(int64(line[k].linevalue&0x1F), 10)
	SecondLine += ", [R" + strconv.FormatInt(int64((line[k].linevalue&0x3E0)>>5), 10)
	SecondLine += ", #" + strconv.FormatInt(int64((line[k].linevalue&0x1FF000)>>12), 10) + "]\n"
	return firstLine() + SecondLine + registers() + "\n" + DataOut()
}
func NOPprint() string {
	sim.cycleN++

	return firstLine() + "\n" + registers() + "\n" + DataOut()
}
func BREAKprint() string {
	sim.cycleN++
	return firstLine() + "\n" + registers() + "\n" + DataOut()
}

func DataOut() string {
	var dataout string
	dataout = "Data:\n"
	for j := startData; j <= endData; j = j + 32 {
		dataout += strconv.FormatInt(j, 10) + ":\t"
		for i := 0; i < 32 && (int64(i)+j) <= endData; i = i + 4 {
			dataout += strconv.FormatInt(int64(sim.DataMem[int64(i)+j]), 10) + "\t"
		}
		dataout += "\n"
	}
	return dataout
}

func main() {
	var InputFileName *string
	var OutputFileName *string
	InputFileName = flag.String("i", "", "Gets the input file name")
	OutputFileName = flag.String("o", "", "Gets the input file name")
	//*InputFileName = "test1_bin.txt"
	//*OutputFileName = "team13_out"

	flag.Parse()
	fmt.Println(flag.NArg())
	if flag.NArg() != 0 {
		os.Exit(200)
	}
	f, err := os.Create(*OutputFileName + "_dis.txt")
	if err != nil {
		log.Fatal(err)
	}
	s, err := os.Create(*OutputFileName + "_pipeline.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	file, err := os.Open(*InputFileName)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()
	var breakFlag bool = false
	//StrInst = make([]string, 10)
	//xyz := 0
	//sim.r[1] = 3
	//sim.r[2] = 200
	programCount = 96
	programCountConst = programCount
	line = make([]Instruction, 1000)
	sim.DataMem = make([]uint64, 30)
	k = 0
	startPrgCntFlag = true
	for _, eachline := range txtlines {
		line[k].linevalue, err = strconv.ParseInt(eachline, 2, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		//	Finding opcode by masking and shifting
		line[k].opcode = (line[k].linevalue & 0xFFE00000) >> 21

		line[k].programCnt = programCount
		if startPrgCntFlag == true {
			startPrgCnt = programCount
			startPrgCntFlag = false
		}

		if !breakFlag {
			// code for R instruction sets
			if line[k].opcode == 1104 || line[k].opcode == 1112 || line[k].opcode == 1360 || line[k].opcode == 1624 || line[k].opcode == 1872 {
				line[k].typeofInstruction = "R1"
				line[k].rd = uint8(line[k].linevalue & 0x1F)
				line[k].rn = uint8((line[k].linevalue & 0x3E0) >> 5)
				line[k].rm = uint8((line[k].linevalue & 0x1F0000) >> 16)
				//var line.op[] string
				if line[k].opcode == 1104 { // print Instruction
					line[k].op = "AND"
				} else if line[k].opcode == 1112 {
					line[k].op = "ADD"
				} else if line[k].opcode == 1360 {
					line[k].op = "ORR"
				} else if line[k].opcode == 1624 {
					line[k].op = "SUB"
				} else if line[k].opcode == 1872 {
					line[k].op = "EOR"
				}
				// AND Rd, Rn, Rm
				strArr := []string{eachline[0:11], eachline[11:16], eachline[16:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "	" + strconv.FormatInt(int64(line[k].programCnt), 10)
				strOut = strOut + "	" + line[k].op
				strOut = strOut + "	R" + strconv.FormatInt(int64(line[k].rd), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line[k].rn), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line[k].rm), 10) + "\n"
				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			// code for R instruction set - LSR, LSL, ASR
			if line[k].opcode == 1690 || line[k].opcode == 1691 || line[k].opcode == 1692 {
				//var line.op[] string
				line[k].typeofInstruction = "R2"
				line[k].rd = uint8(line[k].linevalue & 0x1F)
				line[k].rn = uint8((line[k].linevalue & 0x3E0) >> 5)
				line[k].shamt = uint8((line[k].linevalue & 0xFC00) >> 10)
				if line[k].opcode == 1690 {
					line[k].op = "LSR"
				} else if line[k].opcode == 1691 {
					line[k].op = "LSL"
				} else if line[k].opcode == 1692 {
					line[k].op = "ASR"
				}
				// AND Rd, Rn, Shamt
				strArr := []string{eachline[0:11], eachline[11:16], eachline[16:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "	" + strconv.FormatInt(int64(line[k].programCnt), 10)
				strOut = strOut + "	" + line[k].op
				strOut = strOut + "	R" + strconv.FormatInt(int64(line[k].rd), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line[k].rn), 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64(line[k].shamt), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			// Code for D Instruction set
			if line[k].opcode == 1984 || line[k].opcode == 1986 {
				//var line.op[] string
				line[k].typeofInstruction = "D"
				if line[k].opcode == 1984 {
					line[k].op = "STUR"
				} else if line[k].opcode == 1986 {
					line[k].op = "LDUR"
				}
				//LDUR Rt, [Rn, #line.programCnt]
				line[k].rd = uint8(line[k].linevalue & 0x1F)
				line[k].rn = uint8((line[k].linevalue & 0x3E0) >> 5)
				line[k].address = uint16((line[k].linevalue & 0x1FF000) >> 12)
				strArr := []string{eachline[0:11], eachline[11:20], eachline[20:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "	" + strconv.FormatInt(int64(line[k].programCnt), 10)
				strOut = strOut + "	" + line[k].op
				strOut = strOut + "	R" + strconv.FormatInt(int64(line[k].linevalue&0x1F), 10)
				strOut = strOut + ", [R" + strconv.FormatInt(int64((line[k].linevalue&0x3E0)>>5), 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64((line[k].linevalue&0x1FF000)>>12), 10) + "]\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}
			// I format instructions
			if line[k].opcode == 1160 || line[k].opcode == 1161 || line[k].opcode == 1672 || line[k].opcode == 1673 {
				// ADDI
				// code for 2's compliment
				line[k].typeofInstruction = "I"
				var x uint16 = uint16((line[k].linevalue & 0x3FFC00) >> 10)
				var xs int16 = int16((line[k].linevalue & 0x3FFC00) >> 10)
				if x > 0x7FF {
					x = ^x + 1
					x = x << 4
					x = x >> 4
					xs = int16(x) * -1
				}
				line[k].immN = xs
				line[k].rd = uint8(line[k].linevalue & 0x1F)
				line[k].rn = uint8((line[k].linevalue & 0x3E0) >> 5)
				//var line.op[] string
				if line[k].opcode == 1160 || line[k].opcode == 1161 {
					line[k].op = "ADDI"
				} else if line[k].opcode == 1672 || line[k].opcode == 1673 {
					line[k].op = "SUBI"
				}
				// ADDI/SUBI Rd, Rn, #immediate
				strArr := []string{eachline[0:10], eachline[10:22], eachline[22:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "  	" + strconv.FormatInt(int64(line[k].programCnt), 10) + " 	" + line[k].op
				strOut = strOut + " 	R" + strconv.FormatInt(int64(line[k].rd), 10)
				strOut = strOut + ", R" + strconv.FormatInt(int64(line[k].rn), 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64(line[k].immN), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			// for B instruction
			if line[k].opcode >= 160 && line[k].opcode <= 191 {
				line[k].typeofInstruction = "B"
				var bx uint32 = uint32(line[k].linevalue & 0x3FFFFFF)
				var bxs int64 = int64(line[k].linevalue & 0x3FFFFFF)
				if bx > 0x1FFFFFF {
					bx = ^bx + 1
					bx = bx << 6
					bx = bx >> 6
					bxs = int64(bx) * -1
				}
				// B #offset
				line[k].op = "B"
				line[k].offset = bxs
				strArr := []string{eachline[0:6], eachline[6:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "   	" + strconv.FormatInt(int64(line[k].programCnt), 10) + " 	" + line[k].op
				strOut = strOut + " 	#" + strconv.FormatInt(int64(bxs), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			// for CB instructions
			if line[k].opcode >= 1440 && line[k].opcode <= 1455 {
				//var line.op[] string
				line[k].typeofInstruction = "B"
				if line[k].opcode >= 1440 && line[k].opcode <= 1447 {
					line[k].op = "CBZ"
				} else if line[k].opcode >= 1448 && line[k].opcode <= 1455 {
					line[k].op = "CBNZ"
				}
				var cbx uint32 = uint32((line[k].linevalue & 0xFFFFE0) >> 5)
				var cbxs int32 = int32((line[k].linevalue & 0xFFFFE0) >> 5)
				if cbx > 0x3FFFFF {
					cbx = ^cbx + 1
					cbx = cbx << 8
					cbx = cbx >> 8
					cbxs = int32(cbx) * -1
				}

				// B #offset
				line[k].rd = uint8(line[k].linevalue & 0x1F)
				line[k].offset = int64(cbxs)
				strArr := []string{eachline[0:8], eachline[8:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "   	" + strconv.FormatInt(int64(line[k].programCnt), 10) + " 	" + line[k].op
				strOut = strOut + " 	R" + strconv.FormatInt(line[k].linevalue&0x1F, 10)
				strOut = strOut + ", #" + strconv.FormatInt(int64(cbxs), 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			// for IM instructions
			if (line[k].opcode >= 1684 && line[k].opcode <= 1687) || (line[k].opcode >= 1940 && line[k].opcode <= 1943) {
				//var line.op[] string
				line[k].typeofInstruction = "IM"
				if line[k].opcode >= 1684 && line[k].opcode <= 1687 {
					line[k].op = "MOVZ"
				} else if line[k].opcode >= 1940 || line[k].opcode <= 1943 {
					line[k].op = "MOVK"
				}
				// B #offset
				line[k].shamt = uint8(((line[k].linevalue & 0x600000) >> 21) * 16)
				line[k].field = (line[k].linevalue & 0x1FFFE0) >> 5
				line[k].rd = uint8(line[k].linevalue & 0x1F)
				strArr := []string{eachline[0:9], eachline[9:11], eachline[11:27], eachline[27:32]}
				strOut := strings.Join(strArr, " ")
				strOut = strOut + "  	" + strconv.FormatInt(int64(line[k].programCnt), 10) + " 	" + line[k].op
				strOut = strOut + " 	R" + strconv.FormatInt(line[k].linevalue&0x1F, 10)
				strOut = strOut + ", " + strconv.FormatInt((line[k].linevalue&0x1FFFE0)>>5, 10)
				strOut = strOut + ", LSL " + strconv.FormatInt(((line[k].linevalue&0x600000)>>21)*16, 10) + "\n"

				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			if line[k].opcode == 0x0 {
				line[k].op = "NOP"
				strOut := "00000000000000000000000000000000\t"
				strOut += strconv.FormatInt(int64(line[k].programCnt), 10) + "\tNOP" + "\n"
				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
			}

			// for Break Instruction
			if line[k].opcode == 2038 {
				line[k].op = "BREAK"
				endPrgCnt = line[k].programCnt
				startData = line[k].programCnt + 4
				strOut := eachline + "     " + strconv.FormatInt(int64(line[k].programCnt), 10) + "	 BREAK\n"
				_, err := f.WriteString(strOut)
				if err != nil {
					log.Fatal(err)
				}
				breakFlag = true
			}
		} else {
			var x uint32 = uint32(line[k].linevalue)
			var xs int32 = int32(line[k].linevalue)
			if x > 0x7FFFFFFF {
				x = ^x + 1
				xs = int32(x) * -1
			}
			endData = line[k].programCnt
			if len(sim.DataMem) <= int(line[k].programCnt) {
				tempSlice = make([]uint64, int(line[k].programCnt+8), int(line[k].programCnt+8))
				for i := range sim.DataMem {
					tempSlice[i] = sim.DataMem[i]
				}
				sim.DataMem = tempSlice
				sim.DataMem[line[k].programCnt] = uint64(xs)
			} else {
				sim.DataMem[line[k].programCnt] = uint64(xs)
			}

			strOut := eachline + "     " + strconv.FormatInt(int64(line[k].programCnt), 10) + "	 " + strconv.FormatInt(int64(xs), 10) + "\n"
			_, err := f.WriteString(strOut)
			if err != nil {
				log.Fatal(err)
			}

		}
		programCount = programCount + 4
		k++
	}

	if endData == 0 {
		tempSlice = make([]uint64, int(programCount+40+1), int(programCount+40+1))
		for l := 0; l < 40; l = l + 4 {
			tempSlice[l+int(programCount)] = uint64(l)
		}
		sim.DataMem = tempSlice
		startData = programCount
		endData = programCount + 40
	}

	//sim.DataMem = []uint64{0}
	//fmt.Println(line.op)
	//(programCount/4)-95
	for k = 0; k < int((programCount-96)/4); k++ {
		if len(sim.DataMem) <= int((k*4)+96) {
			tempSlice = make([]uint64, int((k*4)+96+1), int((k*4)+96+1))
			for i := range sim.DataMem {
				tempSlice[i] = sim.DataMem[i]
			}
			sim.DataMem = tempSlice
			//endData = int64((k * 4) + 96 + 1)

		}
		sim.DataMem[(k*4)+96] = uint64(line[k].linevalue)
	}
	//processor()
	//for k = 0; k < len(line); k++
	for k = 0; true; {
		_, err = s.WriteString(printresult())
		if err != nil {
			log.Fatal(err)
		}
		cycleNum++
		writeBack()
		mem()
		ALU()
		issue()
		if breakflag == false {
			k = int(fetching(int64(k)))
		}

		if breakflag == true && len(PreIssueBuff) == 0 && len(PreALUBuff) == 0 && len(PreMemBuff) == 0 && postMemBuff[0] == -1 && postALUBuff[0] == -1 {
			break
		}
	}
	_, err = s.WriteString(printresult())
	if err != nil {
		log.Fatal(err)
	}
	for k = 0; false; k++ {

		if line[k].typeofInstruction == "R1" {
			if line[k].op == "AND" {
				sim.r[line[k].rd] = sim.r[line[k].rm] & sim.r[line[k].rn]
			} else if line[k].op == "ADD" {
				sim.r[line[k].rd] = sim.r[line[k].rm] + sim.r[line[k].rn]
			} else if line[k].op == "ORR" {
				sim.r[line[k].rd] = sim.r[line[k].rm] | sim.r[line[k].rn]
			} else if line[k].op == "SUB" {
				sim.r[line[k].rd] = sim.r[line[k].rn] - sim.r[line[k].rm]
			} else if line[k].op == "EOR" {
				sim.r[line[k].rd] = sim.r[line[k].rm] ^ sim.r[line[k].rn]
			}
			_, err = s.WriteString(r1Sim())
			if err != nil {
				log.Fatal(err)
			}
		} else if line[k].typeofInstruction == "R2" {
			if line[k].op == "LSR" {
				sim.r[line[k].rd] = sim.r[line[k].rn] >> line[k].shamt
			} else if line[k].op == "LSL" {
				sim.r[line[k].rd] = sim.r[line[k].rn] << line[k].shamt
			} else if line[k].op == "ASR" {
				//sim.r[line[k].rd] = int64(sim.r[line[k].rn]) / uint64(math.Pow(2, float64(line[k].shamt)))
				sim.r[line[k].rd] = uint64(int64(sim.r[line[k].rn]) >> line[k].shamt)
			}
			_, err = s.WriteString(r2Sim())
			if err != nil {
				log.Fatal(err)
			}
		} else if line[k].typeofInstruction == "I" {
			if line[k].op == "ADDI" {
				sim.r[line[k].rd] = sim.r[line[k].rn] + uint64(line[k].immN)
			} else if line[k].op == "SUBI" {
				sim.r[line[k].rd] = sim.r[line[k].rn] - uint64(line[k].immN)
			}
			_, err = s.WriteString(immediate())
			if err != nil {
				log.Fatal(err)
			}
		} else if line[k].typeofInstruction == "B" {
			if k+int(line[k].offset) < 0 || k+int(line[k].offset) > int((endPrgCnt-startPrgCnt)/4) {
				fmt.Fprintf(os.Stderr, "error: branch out of bounds%v\n", err)
				os.Exit(1)
			}
			_, err = s.WriteString(BOffset())
			if err != nil {
				log.Fatal(err)
			}
			if line[k].op == "B" {
				k = k - 1 + int(line[k].offset)
			} else if line[k].op == "CBZ" && sim.r[line[k].rd] == 0 {
				k = k - 1 + int(line[k].offset)
			} else if line[k].op == "CBNZ" && sim.r[line[k].rd] != 0 {
				k = k - 1 + int(line[k].offset)
			}

		} else if line[k].typeofInstruction == "IM" {
			if line[k].op == "MOVZ" {
				sim.r[line[k].rd] = 0
				sim.r[line[k].rd] = sim.r[line[k].rd] + uint64(line[k].field<<line[k].shamt)

			} else if line[k].op == "MOVK" {

				if line[k].shamt == 0 {
					sim.r[line[k].rd] = sim.r[line[k].rd] & 0xFFFFFFFFFFFF0000
					sim.r[line[k].rd] = sim.r[line[k].rd] + uint64(line[k].field<<line[k].shamt)
				} else if line[k].shamt == 16 {
					sim.r[line[k].rd] = sim.r[line[k].rd] & 0xFFFFFFF0000FFFF
					sim.r[line[k].rd] = sim.r[line[k].rd] + uint64(line[k].field<<line[k].shamt)
				} else if line[k].shamt == 32 {
					sim.r[line[k].rd] = sim.r[line[k].rd] & 0xFFFF0000FFFFFFFF
					sim.r[line[k].rd] = sim.r[line[k].rd] + uint64(line[k].field<<line[k].shamt)
				} else if line[k].shamt == 48 {
					sim.r[line[k].rd] = sim.r[line[k].rd] & 0x0000FFFFFFFFFFFF
					sim.r[line[k].rd] = sim.r[line[k].rd] + uint64(line[k].field<<line[k].shamt)
				}
			}
			_, err = s.WriteString(Move())
			if err != nil {
				log.Fatal(err)
			}
		} else if line[k].typeofInstruction == "D" {
			if line[k].op == "LDUR" {
				if len(sim.DataMem) <= int(sim.r[line[k].rn]+uint64(4*line[k].address)) {
					tempSlice = make([]uint64, int(sim.r[line[k].rn]+uint64(4*line[k].address)+1), int(sim.r[line[k].rn]+uint64(4*line[k].address)+1))
					for i := range sim.DataMem {
						tempSlice[i] = sim.DataMem[i]
					}
					sim.DataMem = tempSlice
					endData = int64(sim.r[line[k].rn] + uint64(4*line[k].address) + 1)

				}
				sim.r[line[k].rd] = sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)]

			} else if line[k].op == "STUR" {
				if len(sim.DataMem) <= int(sim.r[line[k].rn]+uint64(4*line[k].address)) {
					tempSlice = make([]uint64, int(sim.r[line[k].rn]+uint64(4*line[k].address)+1), int(sim.r[line[k].rn]+uint64(4*line[k].address)+1))
					for i := range sim.DataMem {
						tempSlice[i] = sim.DataMem[i]
					}
					sim.DataMem = tempSlice
					endData = int64(sim.r[line[k].rn] + uint64(4*line[k].address) + 1)

				}
				sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
			}
			_, err = s.WriteString(DataManipulation())
			if err != nil {
				log.Fatal(err)
			}
		} else if line[k].op == "NOP" {

			_, err = s.WriteString(NOPprint())
			if err != nil {
				log.Fatal(err)
			}
		} else if line[k].op == "BREAK" {

			_, err = s.WriteString(BREAKprint())
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}
}
