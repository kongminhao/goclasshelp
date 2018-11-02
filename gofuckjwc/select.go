package gofuckjwc

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"strings"
)

type User struct {
	stu_no  string
	passwd  string
	session *gorequest.SuperAgent
	cookies []*http.Cookie
}

func NewUser(stu_no, passwd string) User {
	user := User{stu_no: stu_no, passwd: passwd, session: gorequest.New(), cookies: make([]*http.Cookie, 0)}
	user.session.Transport.DisableKeepAlives = false
	user.session.SetDebug(true)
	return user
}

func (this *User) Login() (err error) {
	queryString := fmt.Sprintf("j_username=%s&j_password=%s&j_captcha1=error", this.stu_no, this.passwd)
	fmt.Println(queryString)
	resp, _, _ := this.session.Get("http://202.115.47.141/login").End()
	cookies := (* http.Response)(resp).Cookies()
	// Get cookies from first request
	this.cookies = cookies
	// login
	this.session.Post("http://202.115.47.141/j_spring_security_check").
		Type("form").
		Send(queryString).
		AddCookies(this.cookies).End()

	// enter course select
	this.session.Get("http://202.115.47.141/student/courseSelect/courseSelect/index").
		AddCookies(this.cookies).End()
	return nil
}

func (this *User) LoopSelectCourse(config []CourserConfig) {
	for _, v := range config {
		queryData := fmt.Sprintf("searchtj=%s&xq=0&jc=0&kyl=0&kclbdm=", v.courserId)
		_, respBody, _ := this.session.Post("http://202.115.47.141/student/courseSelect/freeCourse/courseList").
			Type("form").
			Send(queryData).
			AddCookies(this.cookies).End()
		var jsonInterface interface{}
		json.Unmarshal([]byte(respBody), &jsonInterface)
		// 剩下的等选课开放再写吧。此处判断是否有课余量
	}
}

func (this *User) getTokenAndFajhh(queryData string) (token, fajhh string) {
	_, respBody, _ :=this.session.Get("http://202.115.47.141/student/courseSelect/freeCourse/index").
		AddCookies(this.cookies).End()
	document, err := goquery.NewDocumentFromReader(strings.NewReader(respBody))
	if err == nil{
		panic(err)
	}
	token = document.Find("input[name=tokenValue]").Get(0).Data
	fajhh = document.Find("input[name=fajhh]").Get(0).Data
	return
}
