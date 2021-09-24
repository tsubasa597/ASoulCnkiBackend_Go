package vo

type Reply struct {
	Type       uint8  `json:"type_id"`
	Rid        int64  `json:"oid"`
	UUID       int64  `json:"uid"`
	Rpid       int64  `json:"rpid"`
	UID        int64  `json:"mid"`
	Time       int32  `json:"ctime"`
	Name       string `json:"m_name"`
	Content    string `json:"content"`
	Like       uint32 `json:"like_num"`
	OriginRpid int64  `json:"origin_rpid"`
	Num        uint32 `json:"similar_count"`
	TotalLike  uint32 `json:"similar_like_sum"`
}
