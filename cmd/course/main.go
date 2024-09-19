package main

import (
	course "github.com/west2-online/fzuhelper-server/kitex_gen/course/courseservice"
	"log"
)

func main() {
	svr := course.NewServer(new(CourseServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
