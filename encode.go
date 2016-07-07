// Copyright 2016 Michael Van Brummen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencode

import (
	"bytes"
	"fmt"
)

func Marshal(data interface{}) (bencoding []byte, err error) {
	defer func() {
		if ex := recover(); ex != nil {
			err = fmt.Errorf("%v", ex)
		}
	}()
	bencoding = encodeEntity(data)
	return bencoding, nil
}

func encodeEntity(data interface{}) []byte {
	var buf bytes.Buffer
	switch v := data.(type) {
	case string:
		buf.WriteString(fmt.Sprintf("%d:%s", len(v), v))
	case []byte:
		buf.WriteString(fmt.Sprintf("%d:%s", len(v), v))
	case int, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		buf.WriteString(fmt.Sprintf("i%de", v))
	case []interface{}:
		buf.WriteRune('l')
		for _, item := range v {
			buf.Write(encodeEntity(item))
		}
		buf.WriteRune('e')
	case map[string]interface{}:
		buf.WriteRune('d')
		for k, v := range v {
			buf.Write(encodeEntity(k))
			buf.Write(encodeEntity(v))
		}
		buf.WriteRune('e')
	default:
		panic(fmt.Sprintf("Failed to BEncode type: %T", data))
	}
	return buf.Bytes()
}
