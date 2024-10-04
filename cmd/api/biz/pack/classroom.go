package pack

import (
	classroomModel "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildClassroom(res *model.Classroom) *classroomModel.Classroom {
	return &classroomModel.Classroom{
		Build:    res.Build,
		Location: res.Location,
		Capacity: res.Capacity,
		Type:     res.Type,
	}
}
func BuildClassroomList(res []*model.Classroom) []*classroomModel.Classroom {
	list := make([]*classroomModel.Classroom, 0, len(res))
	for _, v := range res {
		list = append(list, BuildClassroom(v))
	}
	return list
}
