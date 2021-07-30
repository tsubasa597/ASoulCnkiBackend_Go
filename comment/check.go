package comment

import (
	"container/heap"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
)

func ReplaceStr(s string) string {
	for _, emote := range *db.Get(&db.Emote{}).(*[]db.Emote) {
		s = strings.Replace(s, emote.EmoteText, string(emote.EmoteID), -1)
	}
	return s
}

func ReplaceRune(s string) string {
	for _, emote := range *db.Get(&db.Emote{}).(*[]db.Emote) {
		s = strings.Replace(s, string(emote.EmoteID), emote.EmoteText, -1)
	}

	return s
}

func Compare(s string) CompareResults {
	commResults := make(CompareResults, 0, conf.HeapLength)
	counts := make(map[uint64]float64)
	for _, v := range Hash(s) {
		if uids, ok := caches.Comm.Load(v); ok {
			for uid := range *uids.(*Set) {
				counts[uint64(uid)] += 1.0
			}
		}
	}

	for id, count := range counts {
		charNum := utf8.RuneCountInString(s)

		comm := &db.Comment{
			Model: db.Model{
				ID: id,
			},
		}
		if err := db.Find(comm); err != nil {
			continue
		}

		if utf8.RuneCountInString(comm.Comment) >= charNum {
			heap.Push(&commResults, CompareResult{
				Comment:    comm,
				Similarity: count / float64(utf8.RuneCountInString(comm.Comment)-conf.DefaultK+1),
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
