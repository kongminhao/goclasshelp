package main

import (
	"fmt"
	"gofuckjwc/gofuckjwc"
)

func main()  {
	var passwd,config = gofuckjwc.NewPasswdAndCourseConfig()
	user := gofuckjwc.NewUser("", passwd)
	user.Login()
	for {
		count := 0
		user.LoopSelectCourse(config)
		for _, course := range config{
			if course.IsSelect == true {
				count += 1
			}
		}
		if count == len(config){
			fmt.Println("选课完成")
			break
		}
	}
}