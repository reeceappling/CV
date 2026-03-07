package this

import (
	"os"
	"path"
	"runtime"
)

func Dir() string {
	_, f, _, _ := runtime.Caller(1)
	return path.Dir(f) + string(os.PathSeparator)
}
