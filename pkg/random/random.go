package random

import "math/rand"

var chars []byte = []byte("abcdefghijklmnopqrstuvwxyz1234567890")

func GenerateRandomString(n int) string {
	var res []byte
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(chars))
		res = append(res, chars[idx])
	}
	return string(res)
}
