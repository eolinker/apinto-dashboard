package common

// ListPage 返回截取下标  [i:y]
//
//		if i == 0 && y == 0 {
//		  "分页错误"
//		}
//	 list[i:y]
func ListPage(pageNo, pageSize int, length int) (int, int) {
	if pageNo == 0 {
		pageNo = 1
	}
	if (pageNo-1)*pageSize > length {
		return 0, 0
	}
	if pageNo == 1 {
		if pageNo*pageSize > length {
			return 0, (pageSize * pageNo) - ((pageSize * pageNo) - length)
		}
		return 0, pageSize * pageNo
	} else {
		if pageNo*pageSize > length {
			return pageSize * (pageNo - 1), (pageSize * pageNo) - ((pageSize * pageNo) - length)
		} else {
			return pageSize * (pageNo - 1), pageSize * pageNo
		}
	}
}
