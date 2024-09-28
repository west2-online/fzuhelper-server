package pack

import (
	model2 "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildClassroom(res *model.Classroom) *model2.Classroom {
	return &model2.Classroom{
		Build:    res.Build,
		Location: res.Location,
		Capacity: res.Capacity,
		Type:     res.Type,
	}
}
func BuildClassroomList(res []*model.Classroom) []*model2.Classroom {
	var list []*model2.Classroom
	for _, v := range res {
		list = append(list, BuildClassroom(v))
	}
	return list
}
