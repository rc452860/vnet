package cachex

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"reflect"
	"time"

	"github.com/rc452860/vnet/common/cache"
)

var instance = cache.New(time.Second * 60)

func GetCache() *cache.Cache {
	return instance
}

// 构建缓存key
func BuildKey(element ...interface{}) (string, error) {
	buf := new(bytes.Buffer)
	for _, item := range element {
		itemType := reflect.TypeOf(item)
		var data interface{}
		data = item
		switch itemType.Kind() {
		case reflect.Int:
			data = int32(item.(int))
		case reflect.String:
			data = []byte(item.(string))
		}
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			return "", err
		}
	}
	bufDig := md5.Sum(buf.Bytes())
	return base64.StdEncoding.EncodeToString(bufDig[:]), nil
}
