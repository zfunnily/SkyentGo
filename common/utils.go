package common

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Timex() int64 {
	return time.Now().Unix()
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func RandomName(name [][]string) string {
	idx1 := rand.Intn(len(name[0]))
	idx2 := rand.Intn(len(name[1]))

	return name[0][idx1] + name[1][idx2]
}

//英文字符串，第一个字符大写
func FirstCharToUpper(key string) string {
	first := strings.ToUpper(key[0:1])
	return first + key[1:]
}

func RandomCode() string {
	var str bytes.Buffer
	for i := 0; i < 6; i++ {
		str.WriteString(strconv.Itoa(rand.Intn(9)))
	}
	return str.String()
}
