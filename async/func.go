package async

import (
	"context"
	"runtime"
)

func Await(_ context.Context) {
	runtime.Gosched()
}
