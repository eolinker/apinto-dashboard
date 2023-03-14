package page_entry

// PageIndex 用于gorm  MYSQL 分页查询
// 示例如下
// db.Limit(pageSize).Offset(entry.PageIndex(pageNum,pageSize))
func PageIndex(pageNum, pageSize int) int {
	if pageNum == 0 {
		pageNum = 1
	}
	return (pageNum - 1) * pageSize
}
