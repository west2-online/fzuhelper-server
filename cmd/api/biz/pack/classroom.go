package pack

import (
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
)

func BuildClassroom(res *classroom.Classroom) *api.Classroom {
	return &api.Classroom{
		Build:    res.Build,
		Location: res.Location,
		Capacity: res.Capacity,
		Type:     res.Type,
	}
}
func BuildClassroomList(res []*classroom.Classroom) []*api.Classroom {
	var list []*api.Classroom
	for _, v := range res {
		list = append(list, BuildClassroom(v))
	}
	return list
}
