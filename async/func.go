package async

import "runtime"

func Await() {
	runtime.Gosched()
}
