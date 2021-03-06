package check

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/vo/response"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
)

const (
	_video   = `https://www.bilibili.com/video/av`
	_column  = `https://www.bilibili.com/read/cv`
	_dynamic = `https://t.bilibili.com/`
)

// Compare 查重
func Compare(s string) response.Relateds {
	s = check.ReplaceStr(s)

	commResults := make(check.CompareResults, 0, config.HeapLength)
	counts := make(map[string]float64)

	for _, v := range check.Hash(s) {
		val, err := cache.GetInstance().Check.Get(cache.CheckKey, fmt.Sprint(v))
		if err != nil {
			continue
		}
		for _, id := range strings.Split(val, ",") {
			if len(id) == 0 {
				continue
			}
			counts[id] += 1.0
		}
	}

	for id, count := range counts {
		charNum := utf8.RuneCountInString(s)

		content, err := cache.GetInstance().Content.Get(cache.ContentKey, id)
		if err != nil {
			continue
		}

		n := utf8.RuneCountInString(check.ReplaceStr(string(content)))
		if n >= charNum {
			heap.Push(&commResults, check.CompareResult{
				ID:         id,
				Similarity: count / float64(n-config.DefaultK+1),
			})
		}
	}

	related := response.Relateds{
		Related: make([]response.Related, 0, len(commResults)),
	}

	for len(commResults) > 0 {
		result := heap.Pop(&commResults).(check.CompareResult)
		if result.Similarity < 0.2 {
			break
		}

		reply, err := dao.Check(result.ID)
		if err != nil {
			continue
		}

		related.Related = append(related.Related, response.Related{
			Rate:  result.Similarity,
			Reply: reply,
			URL:   buildURL(info.Type(reply.Type), reply.Rid, reply.Rpid),
		})
	}

	sort.Slice(related.Related, func(i, j int) bool {
		if related.Related[i].Rate == related.Related[j].Rate {
			return related.Related[i].Reply.Rid < related.Related[j].Reply.Rid
		}

		return related.Related[i].Rate > related.Related[j].Rate
	})

	if len(related.Related) > 0 {
		related.Rate = related.Related[0].Rate
	}

	startTime, endTime := dao.GetTimeInfo()
	related.StartTime = startTime
	related.EndTme = endTime

	return related
}

func buildURL(typ info.Type, rid, rpid int64) string {
	switch typ {
	case info.CommentViedo:
		return fmt.Sprintf("%s%d#reply%d", _video, rid, rpid)
	case info.CommentColumn:
		return fmt.Sprintf("%s%d#reply%d", _column, rid, rpid)
	case info.CommentDynamic:
		return fmt.Sprintf("%s%d#reply%d", _dynamic, rid, rpid)
	default:
		return ""
	}
}
