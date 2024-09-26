package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/config"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"net/http"
	"testing"
)

func getIDAndCookies() (string, []*http.Cookie) {
	return jwch.NewStudent().WithUser("", "").GetIdentifierAndCookies()
}

func Init() {
	// config init
	config.Init()

	utils.LoggerInit()
	dal.Init()
	klog.SetLevel(klog.LevelDebug)
}

func TestGetQiShanEmptyRooms(t *testing.T) {
	Init()
	ctx := context.Background()
	id, cookies := getIDAndCookies()
	s := NewClassroomService(ctx, id, cookies)
	res, err := s.GetEmptyRooms(&classroom.EmptyRoomRequest{
		Date:      "2024-09-19",
		Campus:    "旗山校区",
		StartTime: "1",
		EndTime:   "2",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pack.BuildClassRooms(res, "旗山校区"))
}

func TestGuLangYuEmptyRooms(t *testing.T) {
	Init()
	ctx := context.Background()
	id, cookies := getIDAndCookies()
	s := NewClassroomService(ctx, id, cookies)
	res, err := s.GetEmptyRooms(&classroom.EmptyRoomRequest{
		Date:      "2024-09-18",
		Campus:    "鼓浪屿校区",
		StartTime: "1",
		EndTime:   "2",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pack.BuildClassRooms(res, "鼓浪屿校区"))
}

func TestJiMeiEmptyRooms(t *testing.T) {
	Init()
	ctx := context.Background()
	id, cookies := getIDAndCookies()
	s := NewClassroomService(ctx, id, cookies)
	res, err := s.GetEmptyRooms(&classroom.EmptyRoomRequest{
		Date:      "2024-09-18",
		Campus:    "集美校区",
		StartTime: "1",
		EndTime:   "2",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pack.BuildClassRooms(res, "集美校区"))
}

func TestTongPanEmptyRooms(t *testing.T) {
	Init()
	ctx := context.Background()
	id, cookies := getIDAndCookies()
	s := NewClassroomService(ctx, id, cookies)
	res, err := s.GetEmptyRooms(&classroom.EmptyRoomRequest{
		Date:      "2024-09-18",
		Campus:    "铜盘校区",
		StartTime: "1",
		EndTime:   "2",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pack.BuildClassRooms(res, "铜盘校区"))
}

func TestYiShanEmptyRooms(t *testing.T) {
	Init()
	ctx := context.Background()
	id, cookies := getIDAndCookies()
	s := NewClassroomService(ctx, id, cookies)
	res, err := s.GetEmptyRooms(&classroom.EmptyRoomRequest{
		Date:      "2024-09-18",
		Campus:    "怡山校区",
		StartTime: "1",
		EndTime:   "2",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pack.BuildClassRooms(res, "怡山校区"))
}

func TestQuanGangEmptyRooms(t *testing.T) {
	Init()
	ctx := context.Background()
	id, cookies := getIDAndCookies()
	s := NewClassroomService(ctx, id, cookies)
	res, err := s.GetEmptyRooms(&classroom.EmptyRoomRequest{
		Date:      "2024-09-18",
		Campus:    "泉港校区",
		StartTime: "1",
		EndTime:   "2",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pack.BuildClassRooms(res, "泉港校区"))
}
