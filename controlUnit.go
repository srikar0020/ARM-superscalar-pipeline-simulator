package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var cycleNum int64 = 0

func processor() {
	var OutputFileName *string
	OutputFileName = flag.String("o", "", "Gets the input file name")
	s, err := os.Create(*OutputFileName + "_Cycle.txt")
	if err != nil {
		log.Fatal(err)
	}
	var n int64

	for n = 0; true; {
		cycleNum++
		writeBack()
		mem()
		ALU()
		issue()
		if breakflag == false {
			n = fetching(n)
		}
		_, err = s.WriteString(printresult())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(PreIssueBuff))
		fmt.Println(len(PreALUBuff))
		fmt.Println(len(PreMemBuff))
		fmt.Println(postMemBuff[1])
		fmt.Println(postALUBuff[1])

		if breakflag == true && len(PreIssueBuff) == 0 && len(PreALUBuff) == 0 && len(PreMemBuff) == 0 && postMemBuff[1] == -1 && postALUBuff[1] == -1 {
			break
		}
	}

}

func printresult() string {

	return firstLineNew() + buffersOut() + registersNew() + cacheOut() + DataOutNew()

}

func registersNew() string {
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
func firstLineNew() string {
	var FLine string = "--------------------" + "\n"
	FLine += "cycle:" + strconv.FormatInt(cycleNum, 10) + "\n\n"
	return FLine
}
func DataOutNew() string {
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
func cacheOut() string {
	var cacheoutput string
	cacheoutput += "Cache" + "\n"
	for i := 0; i < 4; i++ {
		cacheoutput += "Set " + strconv.FormatInt(int64(i), 10) + ": LRU= " + strconv.FormatInt(int64(cache.sets[i].lru), 10) + "\n"
		for j := 0; j < 2; j++ {
			cacheoutput += "\t" + "Entry " + strconv.FormatInt(int64(j), 10) + ": "
			cacheoutput += "[(" + strconv.FormatInt(int64(cache.sets[i].line[j].valid), 10) + "," + strconv.FormatInt(int64(cache.sets[i].line[j].dirty), 10) + ","
			cacheoutput += strconv.FormatInt(int64(cache.sets[i].line[j].tag), 10) + ")<" + strconv.FormatInt(int64(cache.sets[i].line[j].data1), 2) + ","
			cacheoutput += strconv.FormatInt(int64(cache.sets[i].line[j].data2), 2) + ">\n"
		}
	}
	return cacheoutput
}

func r1SimNew(k int) string {
	var SecondLine string = line[k].op + "\t" + "	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rn), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rm), 10)
	return SecondLine

}
func r2SimNew(k int) string {
	var SecondLine string = line[k].op + "\t" + "	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rn), 10)
	SecondLine += ", #" + strconv.FormatInt(int64(line[k].shamt), 10)
	return SecondLine

}

func immediateNew(k int) string {
	var SecondLine string = line[k].op + "\t" + " 	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", R" + strconv.FormatInt(int64(line[k].rn), 10)
	SecondLine += ", #" + strconv.FormatInt(int64(line[k].immN), 10)

	return SecondLine
}

func BOffsetNew(k int) string {
	var SecondLine string
	if line[k].op == "B" {
		SecondLine = line[k].op + "\t" + "#" + strconv.FormatInt(int64(line[k].offset), 10)
	} else {
		SecondLine = line[k].op + "\t" + "R" + strconv.FormatInt(int64(line[k].rd), 10) + "\t"
		SecondLine += "#" + strconv.FormatInt(int64(line[k].offset), 10)
	}
	return SecondLine
}
func MoveNew(k int) string {

	var SecondLine string
	SecondLine += line[k].op + "\t" + " 	R" + strconv.FormatInt(int64(line[k].rd), 10)
	SecondLine += ", " + strconv.FormatInt(line[k].field, 10)
	SecondLine += ", LSL " + strconv.FormatInt(int64(line[k].shamt), 10) + "\n"
	return SecondLine
}

func DataManipulationNew(k int) string {
	var SecondLine string
	SecondLine += line[k].op + "\t" + "	R" + strconv.FormatInt(int64(line[k].linevalue&0x1F), 10)
	SecondLine += ", [R" + strconv.FormatInt(int64((line[k].linevalue&0x3E0)>>5), 10)
	SecondLine += ", #" + strconv.FormatInt(int64((line[k].linevalue&0x1FF000)>>12), 10) + "]"
	return SecondLine
}
func NOPprintNew(k int) string {

	return line[k].op
}
func BREAKprintNew(k int) string {
	return line[k].op + "\t"
}
func OpController(k int) string {
	if line[k].typeofInstruction == "R1" {
		return r1SimNew(k)
	} else if line[k].typeofInstruction == "R2" {
		return r2SimNew(k)
	} else if line[k].typeofInstruction == "D" {
		return DataManipulationNew(k)
	} else if line[k].typeofInstruction == "I" {
		return immediateNew(k)
	} else if line[k].typeofInstruction == "B" {
		return BOffsetNew(k)
	} else if line[k].typeofInstruction == "IM" {
		return MoveNew(k)
	} else if line[k].op == "NOP" {
		return NOPprintNew(k)
	} else if line[k].op == "BREAK" {
		return BREAKprintNew(k)
	} else {
		return ""
	}

}
func buffersOut() string {
	var buffers string
	buffers = "Pre-Issue Buffer:\n"
	for i := 0; i < 4; i++ {
		buffers += "\tEntry " + strconv.FormatInt(int64(i), 10) + ":\t"
		if i < len(PreIssueBuff) {
			buffers += "[" + OpController(int(PreIssueBuff[i])) + "]"
		}
		buffers += "\n"
	}
	buffers += "Pre_ALU Queue:\n"
	for i := 0; i < 2; i++ {
		buffers += "\tEntry " + strconv.FormatInt(int64(i), 10) + ":\t"
		if i < len(PreALUBuff) {
			buffers += "[" + OpController(int(PreALUBuff[i])) + "]"
		}
		buffers += "\n"
	}
	buffers += "Post_ALU Queue:\n"
	buffers += "\tEntry 0:\t"
	if postALUBuff[1] != -1 {
		buffers += "[" + OpController(int(postALUBuff[2])) + "]"
	}
	buffers += "\n"
	buffers += "Pre_MEM Queue:\n"
	for i := 0; i < 2; i++ {
		buffers += "\tEntry " + strconv.FormatInt(int64(i), 10) + ":\t"
		if i < len(PreMemBuff) {
			buffers += "[" + OpController(int(PreMemBuff[i])) + "]"
		}
		buffers += "\n"
	}
	buffers += "Post_MEM Queue:\n"
	buffers += "\tEntry 0:\t"
	if postMemBuff[1] != -1 {
		buffers += "[" + OpController(int(postMemBuff[2])) + "]"
	}
	buffers += "\n\n"

	return buffers
}
