package utils

// Resp response
type Resp struct {
	Code    int64                  `json:"code"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

// FloorInfo 楼层信息
type FloorInfo struct {
	ID                   int           `json:"id"`
	ArticleBusinessID    string        `json:"articleBusinessId"`
	AuditStatus          int           `json:"auditStatus"`
	AvatarURL            string        `json:"avatarUrl"`
	BusinessID           string        `json:"businessId"`
	CommentStatus        int           `json:"commentStatus"`
	CommentTime          string        `json:"commentTime"`
	CommentType          int           `json:"commentType"`
	Content              string        `json:"content"`
	FloorNum             int64         `json:"floorNum"` // 楼层数
	IsDeleted            int           `json:"isDeleted"`
	IsLike               int           `json:"isLike"`
	IsOwner              int           `json:"isOwner"`
	IsPublisher          int           `json:"isPublisher"`
	IsRead               int           `json:"isRead"`
	IsReplyToPublisher   int           `json:"isReplyToPublisher"`
	IsTop                int           `json:"isTop"`
	LikeCount            int           `json:"likeCount"`
	ParentBusinessId     string        `json:"parentBusinessId"`
	RelationBusinessId   string        `json:"relationBusinessId"`
	RelationBusinessType int           `json:"relationBusinessType"`
	ReplyCount           int           `json:"replyCount"`
	ReplyList            []interface{} `json:"replyList"`
	SourceNickname       string        `json:"sourceNickname"`    // 昵称
	SourceUserAccount    string        `json:"sourceUserAccount"` // 姓名拼音
	SourceWorkCode       string        `json:"sourceWorkCode"`
}
