package http

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/utils/datasize"
)

type FakeFile struct {
	Size   int64
	Offset int64
}

func NewFakeFile(size int64) *FakeFile {
	return &FakeFile{
		Size:   size,
		Offset: 0,
	}
}
func (f *FakeFile) Read(p []byte) (n int, err error) {
	if f.Offset >= f.Size {
		return 0, io.EOF
	}
	remain := f.Size - f.Offset
	if int64(cap(p)) < remain {
		n, err = rand.Read(p)
	} else {
		n, err = rand.Read(p[:int(remain)])
	}

	f.Offset = f.Offset + int64(n)
	return n, err
}

func (f *FakeFile) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekEnd {
		f.Offset = f.Size + offset

	}
	if whence == io.SeekStart {
		f.Offset = 0 + offset
	}
	if whence == io.SeekCurrent {
		f.Offset = f.Offset + offset
	}
	if f.Offset > f.Size || f.Offset < 0 {
		return 0, fmt.Errorf("offset is out of bounds")
	}
	return f.Offset, nil
}

func StartFakeFileServer() *http.Server {
	r := gin.Default()
	r.GET("download", func(c *gin.Context) {
		sizeStr := c.Query("size")
		size, err := datasize.Parse(sizeStr)
		if err != nil {
			log.Err(err)
		}
		file := NewFakeFile(int64(size))
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.test"`, sizeStr))
		http.ServeContent(c.Writer, c.Request, fmt.Sprintf("%s.test", sizeStr), time.Now(), file)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("FakeFileServer closed")
			} else {
				log.Error("Server closed unexpect")
			}
		}
	}()
	return server
}
