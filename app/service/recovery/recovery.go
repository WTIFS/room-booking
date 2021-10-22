package recovery

import (
	"fmt"
	"github.com/wtifs/room-booking/app/service/log"
	"runtime"
)

func Recovery(format string, values ...interface{}) {
	if r := recover(); r != nil {
		buf := make([]byte, 1<<18)
		n := runtime.Stack(buf, false)
		ctx := fmt.Sprintf(format, values...)
		log.Fatal(ctx+". error: %+v. stack: %s", r, buf[0:n])
	}
}
