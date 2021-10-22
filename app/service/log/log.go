package log

import (
	"fmt"
	"github.com/wtifs/room-booking/app/utils"
	"reflect"
	"strings"
	"time"
)

var (
	isDebugMode      bool //true的情况下会打印sql
)

func DebugSql(ctx string, t time.Time, sql string, v ...interface{}) {
	if isDebugMode {
		sql = strings.ReplaceAll(sql, "%", "%%")
		if ctx != "" {
			sql = ctx + " SQL: " + sql
		} else {
			sql = "SQL: " + sql
		}
		if len(v) >= 0 && len(v) < 50 {
			for _, val := range v {
				if val == nil || reflect.TypeOf(val).Kind() == reflect.String {
					sql = strings.Replace(sql, "?", "'%v'", 1)
				} else {
					sql = strings.Replace(sql, "?", "%v", 1)
				}
			}
			sql = fmt.Sprintf(sql, v...)
		} else {
			firstQuestionMark := strings.Index(sql, "?")
			sql = sql[:utils.MinInt(firstQuestionMark, len(sql))]
		}

		duration := time.Now().Sub(t).Seconds()
		Debug("%s\nSQL execution time: %.4fs", sql, duration)
	}
}

func Debug(format string, v ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", v...)
}

func Info(format string, v ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", v...)
}

func Warning(format string, v ...interface{}) {
	fmt.Printf("[WARNING] "+format+"\n", v...)
}

func Err(format string, v ...interface{}) {
	fmt.Printf("[ERR] "+format+"\n", v...)
}


// 重大错误
func Fatal(format string, v ...interface{}) {
	fmt.Printf("[FATAL] "+format+"\n", v...)
	panic(fmt.Sprintf(format, v...))
}
