package main

func fetching(k int64) int64 {
	if instructionCheck(96 + (k * 4)) {

		if len(PreIssueBuff) <= 2 { //for non branch, break, nop instructions
			if (96+(k*4))%8 == 0 {
				//send 2 addresses to issue buffer
				if BreakCheck(k) {
					return k
				}
				if BranchCheck(k) {
					return k - 1 + int64(line[k].offset)
				}

				PreIssueBuff = enqueue(PreIssueBuff, k)

				if BreakCheck(k + 1) {
					return k + 1
				}
				if BranchCheck(k + 1) {
					return k - 1 + int64(line[k].offset)
				}
				PreIssueBuff = enqueue(PreIssueBuff, k+1)

				return k + 2
			} else if (96+(k*4))%4 == 0 {
				if BreakCheck(k) {
					return k
				}
				if BranchCheck(k) {
					return k - 1 + int64(line[k].offset)
				}
				PreIssueBuff = enqueue(PreIssueBuff, k)
				return k + 1
			}
		} else if len(PreIssueBuff) == 3 {
			if BreakCheck(k) {
				return k
			}
			if BranchCheck(k) {
				return k - 1 + int64(line[k].offset)
			}
			PreIssueBuff = enqueue(PreIssueBuff, k)
			return k + 1
		}
	} else {
		fetchMem(96 + (k * 4))
		return k
	}
	return k

}

func BreakCheck(b int64) bool {
	if line[b].op == "BREAK" {
		breakflag = true
		return true
	}
	return false
}

func BranchCheck(b int64) bool {
	if line[b].op == "B" {
		return true
	} else if line[b].op == "CBZ" && sim.r[line[k].rd] == 0 {
		return true
	} else if line[b].op == "CBNZ" && sim.r[line[k].rd] != 0 {
		return true
	} else {
		return false
	}

}
