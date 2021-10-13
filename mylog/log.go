package mylog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

//日志信息通道logChan，用于传入日志信息
type logData struct {
	timeStr 	string		    //触发时间
	level 		Level			//日志等级
	callerFile 	string			//触发文件
	callerFunc 	string			//触发函数
	line 		int				//触发行号
	format 		string			//详情信息
}

type Logger struct {
	level 		Level      		//日志等级
	filePath 	string	   		//日志文件路径
	fileName 	string	   		//日志文件名
	file 		*os.File   		//日志文件句柄
	errFile		*os.File   		//ERROR、FATAL级别日志文件句柄
	maxSize     int64	   		//单个日志文件最大大小(单位：字节)
	logChan		chan *logData	//并发写日志channel
}

//封装Logger结构体内部字段，初始化一个Logger结构体
func NewLogger(level Level, filePath, fileName string) *Logger {
	l := &Logger{
		level:    level,    //这里的level相当于能写入日志的最低日志级别，低于这个级别的日志将无法被写入
		filePath: filePath,
		fileName: fileName,
		maxSize:  20 * 1024 * 1024,   //20MB
		logChan:  make(chan *logData, 10000),   //通道容量必须足够大，为了减少被丢弃的日志数
	}
	l.file, l.errFile = initFile(filePath, fileName)
	for i := 0; i < 5; i++ {   //可以设置goroutine数量
		go l.writeLog()   	//开启goroutine
	}
	return l
}

//获得文件大小
func getFileSize(f *os.File) int64 {
	fileInfo, err := f.Stat()
	if err != nil { panic(err) }
	return fileInfo.Size()
}

//初始化日志文件，ERROR及以上级别日志文件
func initFile(filePath, fileName string) (file, errFile *os.File) {
	f := fmt.Sprintf("%s/%s", filePath, fileName)
	file, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil { panic(err) }
	ef := fmt.Sprintf("%s/%s.err", filePath, fileName)
	errFile, err = os.OpenFile(ef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil { panic(err) }
	return
}

//获得调用写日志方法的文件名、函数、行号
func getCallerInfo() (fileName, funcName string, line int) {
	pc, fileName, line, ok := runtime.Caller(3)   //获取调用DEBUG、INFO、WARN等方法的位置信息
	if !ok { return }

	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	return fileName, funcName, line
}

//当日志文件的大小超过最大限度maxSize时切分文件
func (log *Logger)splitLogFile() {
	logFileSize := getFileSize(log.file)
	errLogFileSize := getFileSize(log.errFile)

	// 判断日志文件大小是否需要切分
	if logFileSize > log.maxSize {
		fileName := log.file.Name()
		backupName := fmt.Sprintf("%s_%v.back", fileName, time.Now().Unix())
		log.file.Close()

		os.Rename(fileName, backupName)  //保存备份
		newLogFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil { panic(err) }
		log.file = newLogFile
	}
	if errLogFileSize > log.maxSize {
		fileName := log.errFile.Name()
		backupName := fmt.Sprintf("%s_%v.back", fileName, time.Now().Unix())
		log.errFile.Close()

		os.Rename(fileName, backupName)
		newErrLogFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
		if err != nil { panic(err) }
		log.errFile = newErrLogFile
	}
	return
}

//以goroutine的方式写入日志
func (log *Logger)writeLog() {
	for data := range log.logChan {   //从通道中取数据

		//写入的日志格式为：
		//[触发时间] [日志等级] [触发文件:触发函数] [触发行号] [详情信息]
		format := fmt.Sprintf("%s [%s] [%s : %s] [line:%d] %s\n",
			data.timeStr, getLogLevel(data.level), data.callerFile, data.callerFunc, data.line, data.format)
		log.splitLogFile()
		fmt.Fprintf(log.file, format)
		if data.level >= ERROR {    //如果日志级别为ERROR或者FATAL级别的话，会再次写入错误日志文件
			fmt.Fprintf(log.errFile, format)
		}
	}
}

//写入日志
func (log *Logger)log(level Level, format string, args ...interface{}) {
	if log.level > level { return }   //如果初始化的Level等级大于写入日志的等级，则退出，也就是说日志等级过低无法写入

	timeStr := time.Now().Format("[2006-01-02 15:04:05]")   //获取格式化的当前时间
	callerFile, callerFunc, line := getCallerInfo()   //获取触发文件，触发函数，行号
	format = fmt.Sprintf(format, args...)	//将日志信息格式化

	ld := &logData{
		timeStr:    timeStr,
		level:      level,
		callerFile: callerFile,
		callerFunc: callerFunc,
		line:       line,
		format:     format,
	}

	select {
	case log.logChan <- ld:   	//将一条日志信息加入channel
	default:
		<-log.logChan             //如果channel满了，则丢弃最前面的一条日志
		log.logChan <- ld
	}
}

//关闭打开的的文件句柄
func (log *Logger)Close() {
	time.Sleep(time.Second * 2)   //显式等待goroutine的关闭，完成日志的写入
	log.file.Close()
	log.errFile.Close()
}

//5种日志级别的写入方法
func (log *Logger)DEBUG(format string, args ...interface{}) {
	log.log(DEBUG, format, args...)
}

func (log *Logger)INFO(format string, args ...interface{}) {
	log.log(INFO, format, args...)
}

func (log *Logger)WARN(format string, args ...interface{}) {
	log.log(WARN, format, args...)
}

func (log *Logger)ERROR(format string, args ...interface{}) {
	log.log(ERROR, format, args...)
}

func (log *Logger)FATAL(format string, args ...interface{}) {
	log.log(FATAL, format, args...)
}
