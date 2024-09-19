package main

import (
	academic "github.com/west2-online/fzuhelper-server/kitex_gen/academic/academicservice"
	"log"
)

func main() {
	svr := academic.NewServer(new(AcademicServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
