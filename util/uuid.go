package util

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
)

var snowFlake *SnowFlake

// GenUUID 生成唯一ID
func GenUUID() (uuid string, err error) {
	var i uint64
	snowFlake, err = NewSnowFlake(1)
	if err != nil {
		return "", err
	}
	if i, err = snowFlake.Generate(); err != nil {
		return "", err
	}
	m := md5.New()
	m.Write([]byte(strconv.Itoa(int(i))))
	uuid = hex.EncodeToString(m.Sum(nil))
	return
}
