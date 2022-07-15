package floor

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/demoManito/bilibiliscript/utils"
)

type floorInfo struct {
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

type Floor struct {
	Conf *Conf

	max     int
	pageNum int
}

func New(fileName string) *Floor {
	f := &Floor{
		Conf: Init(fileName),
	}
	f.pageNum = int(f.Conf.FloorNum / 10)
	return f
}

func (f *Floor) Report() {
	var info *floorInfo
	for {
		f.max++
		info = f.findFloorInfo()
		if info != nil {
			break
		}
		if f.max >= 2 {
			log.Println("未查询到相关数据 ⚠️")
			return
		}
	}

	log.Printf("\n楼层信息：\n 楼层号：%d \n 时间：%s \n 昵称：%s \n 姓名：%s",
		info.FloorNum, info.CommentTime, info.SourceNickname, info.SourceUserAccount)
}

func (f *Floor) findFloorInfo() *floorInfo {
	req, err := http.NewRequest(http.MethodGet, f.parseURL(f.pageNum), nil)
	if err != nil {
		log.Fatalf("[request err] err: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF", f.Conf.XCSRF)
	req.Header.Set("Cookie", f.Conf.Cookie)

	resp := new(utils.Resp)
	response, err := new(http.Client).Do(req)
	if err != nil {
		log.Fatalf("[http err] http client do err: %s", err)
	}
	ioBody, _ := io.ReadAll(response.Body)
	err = json.Unmarshal(ioBody, &resp)
	if err != nil {
		log.Fatalf("[unmarshal err] resp json unmarshal err: %s", err)
	}
	if resp.Code != 0 {
		log.Fatalf("[resp err] code: %d; message: %s; err: %s \n", resp.Code, resp.Message, err)
	}

	list := resp.Data["commentReplyList"]
	jsonList, err := json.Marshal(list)
	if err != nil {
		log.Fatalf("[marshal err] comment_reply_list json marshal err: %s", err)
	}
	floors := make([]*floorInfo, 0)
	err = json.Unmarshal(jsonList, &floors)
	if err != nil {
		log.Fatalf("[unmarshal err] floors json unmarshal err: %s", err)
	}

	if floors[len(floors)-1].FloorNum < f.Conf.FloorNum {
		f.pageNum = f.pageNum + 1
		return nil
	}
	if floors[0].FloorNum > f.Conf.FloorNum {
		f.pageNum = f.pageNum - 1
		return nil
	}
	var floor *floorInfo
	for _, fls := range floors {
		if fls.FloorNum == f.Conf.FloorNum {
			floor = fls
			break
		}
	}
	return floor
}

func (f *Floor) parseURL(pageNum int) string {
	up, err := url.Parse(f.Conf.URL)
	if err != nil {
		log.Fatal(err)
	}
	q := up.Query()
	q.Set("articleBusinessId", f.Conf.ArticleBusinessID)
	q.Set("pageSize", "10")
	q.Set("pageNum", strconv.FormatInt(int64(pageNum), 10))
	q.Set("order", "2")
	q.Set("scrollId", "null")
	up.RawQuery = q.Encode()
	return up.String()
}
