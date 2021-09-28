package check

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
)

// Hash 提取评论 hash
func Hash(s string) []int64 {
	s = ReplaceStr(s)
	n := utf8.RuneCountInString(s)
	hashs := make([]int64, n)
	for i := 0; i < utf8.RuneCountInString(s)-setting.DefaultK+1; i++ {
		var ans int64
		for j, v := range ([]rune(s))[i : i+setting.DefaultK] {
			ans += int64(v) * int64(math.Pow(setting.DefaultB, float64(setting.DefaultK-1-j)))
		}
		hashs[i] = ans
	}
	return hashs
}

// HashSet 提取评论 hash(去重结果)
func HashSet(s string) map[int64]struct{} {
	s = ReplaceStr(s)
	hashs := make(map[int64]struct{})
	for i := 0; i < utf8.RuneCountInString(s)-setting.DefaultK+1; i++ {
		var ans int64
		for j, v := range ([]rune(s))[i : i+setting.DefaultK] {
			ans += int64(v) * int64(math.Pow(setting.DefaultB, float64(setting.DefaultK-1-j)))
		}
		hashs[ans] = struct{}{}
	}
	return hashs
}

// CompareStr 比较两个字符串相似度
func CompareStr(s1, s2 string) float64 {
	h1, h2 := Hash(s1), HashSet(s2)
	set := make(map[int64]struct{})
	count := 0.0
	charNum := utf8.RuneCountInString(s1)
	for i := 0; i < charNum-setting.DefaultK+1; i++ {
		if _, ok := h2[h1[i]]; ok {
			for j := 0; j < setting.DefaultK; j++ {
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

// ReplaceStr 去除多余字符
func ReplaceStr(s string) string {
	return replacer.Replace(s)
}

var (
	replacer = strings.NewReplacer("\n", "", " ", "")
)
