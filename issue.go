package main

func issue() {
	if len(PreIssueBuff) > 0 {
		var k = PreIssueBuff[0]
		if line[k].op == "LDUR" || line[k].op == "STUR" {
			// send to premembuff
			if len(PreMemBuff) < 2 {
				PreMemBuff = enqueue(PreMemBuff, k)
				PreIssueBuff = dequeue(PreIssueBuff)
			}
		} else {
			if len(PreALUBuff) < 2 {
				PreALUBuff = enqueue(PreALUBuff, k)
				PreIssueBuff = dequeue(PreIssueBuff)
			}

			// send to prealubuff
		}
	}
	if len(PreIssueBuff) > 0 {
		var k = PreIssueBuff[0]
		if line[k].op == "LDUR" || line[k].op == "STUR" {
			// send to premembuff
			if len(PreMemBuff) < 2 {
				PreMemBuff = enqueue(PreMemBuff, k)
				PreIssueBuff = dequeue(PreIssueBuff)
			}
		} else {
			if len(PreALUBuff) < 2 {
				PreALUBuff = enqueue(PreALUBuff, k)
				PreIssueBuff = dequeue(PreIssueBuff)
			}

			// send to prealubuff
		}
	}
}
//func HazardCheck() bool{
//	//
//	for i := 0; i < 2; i++ {
//		if i < len(PreALUBuff) {
//			buffers += OpController(int(PreALUBuff[i]))
//		}
//	}
//	if postALUBuff[1] != -1 {
//
//	}
//
//	for i := 0; i < 2; i++ {
//		if i < len(PreMemBuff) {
//			buffers += OpController(int(PreMemBuff[i]))
//		}
//		buffers += "\n"
//	}
//	buffers += "Post_MEM Queue:\n"
//	buffers += "\tEntry 0:"
//	if postMemBuff[1] != -1 {
//		buffers += OpController(int(postMemBuff[2]))
//	}
//	return false
//}

//func HazardCheckHelper(g int64){
//	if line[g].typeofInstruction == "R1"{
//
//	}else if line[g].typeofInstruction == "R2" {
//
//	}else if line[g].typeofInstruction == "I" {
//
//	}else if line[g].typeofInstruction == "IM" {
//
//	}else if line[g].typeofInstruction == "D" {
//
//	}
//
//}
