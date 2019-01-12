package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func Test_Glob(t *testing.T) {
	filename := "sakura.log"
	file, _ := os.Create(filename)
	file.Close()
	for i := 0; i < 10; i++ {
		file, _ = os.Create(fmt.Sprintf("%s.%d", filename, i))
		file.Close()
	}
	files, _ := filepath.Glob(filename + "*")
	t.Logf("%v", files)
	f := make([]string, len(files))
	for i, item := range files {
		f[i] = item[len(filename):]
	}

	t.Logf("%v", f)

	for _, item := range files {
		if err := os.Remove(item); err != nil {
			t.Error(err)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(f)))
	t.Logf("%v", f)
}
