package crontab

// Run 启动计划任务
func Run() {
	WriteSQLLog() // 写入SQL日志
}

// Stop 停止计划任务
func Stop() {
	StopWriteSQLLog <- true
	close(StopWriteSQLLog)
}
