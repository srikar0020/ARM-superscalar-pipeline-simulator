package main

func mem() bool {
	if postMemBuff[0] == -1 && len(PreMemBuff) >= 1 {
		var k = PreMemBuff[0]
		postMemBuff[2] = k
		PreMemBuff = dequeue(PreMemBuff)
		if line[k].op == "LDUR" {
			if !loadCache(int64(sim.r[uint16(line[k].rn)] + uint64(4*line[k].address))) {
				fetchMem(int64(sim.r[uint16(line[k].rn)] + uint64(4*line[k].address)))
				return false
			}
			postMemBuff[1] = int64(line[k].rd)
			return true
			//sim.r[line[k].rd] = sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)]

		} else if line[k].op == "STUR" {
			if !sturCache(int64(sim.r[uint16(line[k].rn)] + uint64(4*line[k].address))) {
				fetchMem(int64(sim.r[uint16(line[k].rn)] + uint64(4*line[k].address)))
				return false
			}
			return true
			//if len(sim.DataMem) <= int(sim.r[line[k].rn]+uint64(4*line[k].address)) {
			//	tempSlice = make([]uint64, int(sim.r[line[k].rn]+uint64(4*line[k].address)+1), int(sim.r[line[k].rn]+uint64(4*line[k].address)+1))
			//	for i := range sim.DataMem {
			//		tempSlice[i] = sim.DataMem[i]
			//	}
			//	sim.DataMem = tempSlice
			//	endData = int64(sim.r[line[k].rn] + uint64(4*line[k].address) + 1)
			//
			//}
			//sim.DataMem[sim.r[uint16(line[k].rn)]+uint64(4*line[k].address)] = sim.r[line[k].rd]
		}
		return true
		//_, err = s.WriteString(DataManipulation())
		//if err != nil {
		//	log.Fatal(err)
		//}
	}
	return false
}
