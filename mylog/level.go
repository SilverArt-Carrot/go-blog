package mylog

type Level int8

//debug 级别最低，详细的了解系统程序的运行情况，在任何地方都可以用；
//info  重要，输出信息：用来反馈系统的当前状态给用户，以便定位问题；
//后三个，警告、错误、严重错误，这三者应该都在系统运行时检测到了一个不正常的状态。
//warn, 可修复，系统可继续运行下去；
//Error, 可修复性，但无法确定系统会正常的工作下去;
//Fatal, 相当严重，可以肯定这种错误已经无法修复，并且如果系统继续运行下去的话后果严重。

const (
	// 日志等级
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

//返回日志级别对应的字符串表达
func getLogLevel(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}
