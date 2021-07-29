package comment

import (
	"container/heap"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func ReplaceStr(s string) string {
	caches.Emote.Range(func(key, value interface{}) bool {
		s = strings.Replace(s, key.(string), value.(string), -1)
		return true
	})
	return s
}

func ReplaceRune(s string) string {
	caches.Emote.Range(func(key, value interface{}) bool {
		s = strings.Replace(s, value.(string), key.(string), -1)
		return true
	})
	return s
}

func Compare(s string) CompareResults {
	commResults := make(CompareResults, 0, conf.HeapLength)

	counts := make(map[int64]float64)
	for _, v := range Hash(s) {
		if uids, ok := caches.Comm.Load(v); ok {
			for uid := range *uids.(*Set) {
				counts[uid] += 1.0
			}
		}
	}

	for uid, count := range counts {
		charNum := utf8.RuneCountInString(s)
		if comm, ok := caches.Reply.Load(uid); ok && utf8.RuneCountInString(comm.(*db.Comment).Comment) >= charNum {
			heap.Push(&commResults, CompareResult{
				Comment:    comm.(*db.Comment),
				Similarity: count / float64(utf8.RuneCountInString(comm.(*db.Comment).Comment)-conf.DefaultK+1),
			})
		}
	}

	for i, v := range commResults {
		if v.Similarity < 0.2 {
			return commResults[:i]
		}

		v.Comment.Comment = ReplaceRune(v.Comment.Comment)
	}

	return commResults
}
