package main

// find set number set no. = address>>3%4
// tag[set no.] == address>>3
// yes then fetch data
// no then fetch data from main mem and send it in next cycle
// place the tag
func instructionCheck(address int64) bool {

	var setNum int64
	setNum = address >> 3 % 4
	if cache.sets[setNum].line[0].tag == address>>5 {
		cache.sets[setNum].lru = 1
		//cache.sets[setNum].line[0].tag = address >> 5
		//if address>>2%2 == 0 {
		//	postMemBuff[0] = int64(cache.sets[setNum].line[0].data1)
		//} else if address>>2%2 == 1 {
		//	postMemBuff[0] = int64(cache.sets[setNum].line[0].data2)
		//}
		return true
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]

	} else if cache.sets[setNum].line[1].tag == address>>5 {
		cache.sets[setNum].lru = 0
		//cache.sets[setNum].line[1].tag = address >> 5
		//if address>>2%2 == 0 {
		//	cache.sets[setNum].line[1].data1 = sim.r[line[k].rd]
		//} else if address>>2%2 == 1 {
		//	cache.sets[setNum].line[1].data2 = sim.r[line[k].rd]
		//}
		return true
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
	}
	return false

}
func writeback(address int64) {
	var setNum int64
	setNum = address >> 3 % 4
	sim.DataMem[address] = cache.sets[setNum].line[cache.sets[setNum].lru].data1
	sim.DataMem[address+4] = cache.sets[setNum].line[cache.sets[setNum].lru].data2
}
func fetchMem(address int64) {

	if len(sim.DataMem) <= int(address+4) {
		tempSlice = make([]uint64, int(address+4+1), int(address+4+1))
		for i := range sim.DataMem {
			tempSlice[i] = sim.DataMem[i]
		}
		sim.DataMem = tempSlice
		endData = int64(address + 4 + 1)
	}

	var setNum int64
	setNum = address >> 3 % 4
	if cache.sets[setNum].line[0].valid == 0 {
		cache.sets[setNum].lru = 1
		cache.sets[setNum].line[0].tag = address >> 5
		cache.sets[setNum].line[0].valid = 1

		cache.sets[setNum].line[0].data1 = sim.DataMem[address]
		cache.sets[setNum].line[0].data2 = sim.DataMem[address+4]

		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]

	} else if cache.sets[setNum].line[1].valid == 0 {
		cache.sets[setNum].lru = 0
		cache.sets[setNum].line[1].tag = address >> 5
		cache.sets[setNum].line[1].valid = 1
		cache.sets[setNum].line[1].data1 = sim.DataMem[address]
		cache.sets[setNum].line[1].data2 = sim.DataMem[address+4]
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
	} else {
		writeback(int64(CachelineAddress[setNum*2+int64(cache.sets[setNum].lru)]))
		cache.sets[setNum].line[cache.sets[setNum].lru].tag = address >> 5
		cache.sets[setNum].line[cache.sets[setNum].lru].valid = 1
		cache.sets[setNum].line[cache.sets[setNum].lru].dirty = 0

		cache.sets[setNum].line[cache.sets[setNum].lru].data1 = sim.DataMem[address]
		cache.sets[setNum].line[cache.sets[setNum].lru].data2 = sim.DataMem[address+4]
		cache.sets[setNum].lru = (cache.sets[setNum].lru * -1) + 1
	}
}
func loadCache(address int64) bool {
	var setNum int64
	setNum = address >> 3 % 4
	if cache.sets[setNum].line[0].tag == address>>5 {
		cache.sets[setNum].lru = 1
		cache.sets[setNum].line[0].tag = address >> 5
		if address>>2%2 == 0 {
			postMemBuff[0] = int64(cache.sets[setNum].line[0].data1)
		} else if address>>2%2 == 1 {
			postMemBuff[0] = int64(cache.sets[setNum].line[0].data2)
		}
		return true
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]

	} else if cache.sets[setNum].line[1].tag == address>>5 {
		cache.sets[setNum].lru = 0
		cache.sets[setNum].line[1].tag = address >> 5
		if address>>2%2 == 0 {
			postMemBuff[0] = int64(cache.sets[setNum].line[0].data1)
		} else if address>>2%2 == 1 {
			postMemBuff[0] = int64(cache.sets[setNum].line[0].data1)
		}
		return true
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
	}
	return false

}

func sturCache(address int64) bool {
	var setNum int64
	setNum = address >> 3 % 4
	if cache.sets[setNum].line[0].tag == address>>5 {
		cache.sets[setNum].lru = 1
		cache.sets[setNum].line[0].dirty = 1
		cache.sets[setNum].line[0].tag = address >> 5
		if address>>2%2 == 0 {
			cache.sets[setNum].line[0].data1 = sim.r[line[k].rd]
		} else if address>>2%2 == 1 {
			cache.sets[setNum].line[0].data2 = sim.r[line[k].rd]
		}
		return true
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]

	} else if cache.sets[setNum].line[1].tag == address>>5 {
		cache.sets[setNum].lru = 0
		cache.sets[setNum].line[1].dirty = 1
		cache.sets[setNum].line[1].tag = address >> 5
		if address>>2%2 == 0 {
			cache.sets[setNum].line[1].data1 = sim.r[line[k].rd]
		} else if address>>2%2 == 1 {
			cache.sets[setNum].line[1].data2 = sim.r[line[k].rd]
		}
		return true
		//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
	}
	return false
	// else {
	//	sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
	//	cache.sets[setNum].line[cache.sets[setNum].lru].dirty = 1
	//	cache.sets[setNum].line[cache.sets[setNum].lru].tag = address >> 5
	//	if address>>2%2 == 0 {
	//		cache.sets[setNum].line[cache.sets[setNum].lru].data1 = sim.r[line[k].rd]
	//	} else if address>>2%2 == 1 {
	//		cache.sets[setNum].line[cache.sets[setNum].lru].data2 = sim.r[line[k].rd]
	//	}
	//
	//	cache.sets[setNum].lru = (cache.sets[setNum].lru * -1) + 1 // changing lru
	//
	//}

	//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]

}

//func instructionFetch(address int32) {
//	var setNum int32
//	setNum = address >> 3 % 4
//	if cache.sets[setNum].line[0].tag == address>>3 {
//
//	} else if cache.sets[setNum].line[1].tag == address>>3 {
//
//	}
//
//}
