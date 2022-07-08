package util

// Size 返回页面大小
func Size(size uint64) uint64 {
	if size == 0 {
		return 15
	}
	return size
}

// Offset 计算分页便宜量
func Offset(page uint64, size uint64) int {
	return int((page - 1) * size)
}

// GenPaginationParams 处理分页参数 s: size 页面大小 o: offset 偏移
func GenPaginationParams(size uint64, page uint64) (s int, o int) {
	if size == 0 {
		size = 15
	}
	o = int((page - 1) * size)
	return int(size), o
}
