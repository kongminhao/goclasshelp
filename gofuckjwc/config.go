package gofuckjwc

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type courseConfig struct {
	courseId     string
	courseSeqNum string
	IsSelect     bool
}

type ConfigInfo struct {
	Password string
	CourseIdList []string
	CourseSeqNumList []string
}


func NewPasswdAndCourseConfig() (passwd string, config []*courseConfig) {
	confData, err := ioutil.ReadFile("conf.toml")

	fmt.Println(confData)
	var conf ConfigInfo
	if err != nil{
		fmt.Println("配置文件解析错误")
		panic(err)
	}
	if _, err := toml.Decode(string(confData), &conf); err != nil {
		fmt.Println("配置文件解析错误")
		panic(err)
	}
	for k,v := range conf.CourseIdList{
		var course = new(courseConfig)
		course.IsSelect = false
		course.courseId = v
		course.courseSeqNum = conf.CourseSeqNumList[k]
		config = append(config, course)
	}
	return conf.Password, config
}
