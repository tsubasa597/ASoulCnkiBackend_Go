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

// select dynamic.type, dynamic.rid, user.uid, comment.rpid as origin_rpi, comment.content, comment.total_like, comment.num, commentator.uname, commentator.dynamic_id, commentator.rpid, commentator.like, commentator.uid, commentator.uname, commentator.time, comment.like from commentator commentator left join comment comment on comment.rpid = commentator.rpid, dynamic, user where dynamic.id = commentator.dynamic_id and dynamic.user_id = user.id order by comment.num desc limit 10 offset 20;

// select dynamic.type, dynamic.rid, user.uid, comment.rpid as origin_rpi, comment.content, comment.total_like, comment.num, commentator.uname, commentator.dynamic_id, commentator.rpid, commentator.like, commentator.uid, commentator.uname, commentator.time, comment.like from commentator commentator left join comment comment on comment.rpid = commentator.rpid left join dynamic on dynamic.id = commentator.dynamic_id left join user on user.id = dynamic.user_id order by comment.num desc limit 10 offset 20;

// select dynamic.type, dynamic.rid, user.uid, comment.rpid as origin_rpi, comment.content, comment.total_like, comment.num, commentator.uname, commentator.dynamic_id, commentator.rpid, commentator.like, commentator.uid, commentator.uname, commentator.time, comment.like from comment comment left join commentator commentator on comment.rpid = commentator.rpid left join dynamic dynamic on dynamic.rid = commentator.dynamic_id left join user user on user.id = dynamic.user_id
