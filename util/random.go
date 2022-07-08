package util

import (
	"math"
	"math/rand"
	"time"
)

const (
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits  = "0123456789"

	codes   = Letters + Digits
	codeLen = len(codes)
)

var (
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenRandomNumber 生成随机数字
func GenRandomNumber(length int64) int {
	rand.Seed(time.Now().Unix())
	r := math.Pow(10, float64(length))
	return rand.Intn(int(r))
}

// GenRandomString 生成随机字符串
func GenRandomString(length int) string {
	data := make([]byte, length)

	for i := 0; i < length; i++ {
		idx := random.Intn(codeLen)
		data[i] = codes[idx]
	}
	return string(data)
}
