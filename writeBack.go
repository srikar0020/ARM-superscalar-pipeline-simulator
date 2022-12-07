package main

func writeBack() {
	if postMemBuff[1] != -1 {
		sim.r[postMemBuff[1]] = uint64(postMemBuff[0])
		postMemBuff[0] = -1
		postMemBuff[1] = -1
		postMemBuff[2] = -1
	}
	if postALUBuff[1] != -1 {
		sim.r[postALUBuff[1]] = uint64(postALUBuff[0])
		postALUBuff[0] = -1
		postALUBuff[1] = -1
		postALUBuff[2] = -1
	}
}
