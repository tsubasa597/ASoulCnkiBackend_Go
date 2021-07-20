package check

import (
	"math"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
)

type Set map[int64]struct{}

func Hash(s string) []int64 {
	n := utf8.RuneCountInString(s)
	hashs := make([]int64, n)
	for i := 0; i < utf8.RuneCountInString(s)-conf.DefaultK+1; i++ {
		var ans int64
		for j, v := range ([]rune(s))[i : i+conf.DefaultK] {
			// if utf8.RuneLen(v) == 4 {
			// 	r1, r2 := utf16.EncodeRune(v)
			// 	ans += int64(r1)*int64(math.Pow(conf.DefaultB, float64(conf.DefaultK-1-j))) +
			// 		int64(r2)*int64(math.Pow(conf.DefaultB, float64(conf.DefaultK-1-j)))
			// } else {
			ans += int64(v) * int64(math.Pow(conf.DefaultB, float64(conf.DefaultK-1-j)))
			// }
		}
		hashs[i] = ans
	}
	return hashs
}

func HashSet(s string) Set {
	hashs := make(Set)
	for i := 0; i < utf8.RuneCountInString(s)-conf.DefaultK+1; i++ {
		var ans int64
		for j, v := range ([]rune(s))[i : i+conf.DefaultK] {
			// if utf8.RuneLen(v) == 4 {
			// 	r1, r2 := utf16.EncodeRune(v)
			// 	ans += int64(r1)*int64(math.Pow(conf.DefaultB, float64(conf.DefaultK-1-j))) +
			// 		int64(r2)*int64(math.Pow(conf.DefaultB, float64(conf.DefaultK-1-j)))
			// } else {
			ans += int64(v) * int64(math.Pow(conf.DefaultB, float64(conf.DefaultK-1-j)))
			// }
		}
		hashs[ans] = struct{}{}
	}
	return hashs
}

func CompareStr(s1, s2 string) float64 {
	h1, h2 := Hash(s1), HashSet(s2)
	set := make(Set)
	count := 0.0
	charNum := utf8.RuneCountInString(s1)
	for i := 0; i < charNum-conf.DefaultK+1; i++ {
		if _, ok := h2[h1[i]]; ok {
			for j := 0; j < conf.DefaultK; j++ {
				set[int64(i+j)] = struct{}{}
			}
		}
	}
	for i := 0; i < charNum; i++ {
		if _, ok := set[int64(i)]; ok {
			count++
		}
	}
	return count / float64(charNum)
}
