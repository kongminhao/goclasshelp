package gofuckjwc

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"strings"
	"time"
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
	user.session.SetDebug(false)
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
	for {
		_, respbody, _ := this.session.Get("http://202.115.47.141/student/courseSelect/courseSelect/index").
			AddCookies(this.cookies).End()
		if strings.Contains(respbody, "非选课阶段") {
			fmt.Println("当前为非选课阶段，10s后重试")
			time.Sleep(time.Duration(10)* time.Second)
		} else {
			break
		}
	}
	return nil
}

func (this *User) LoopSelectCourse(config []*courseConfig) {
	for _, course := range config {
		if course.IsSelect == true {
			continue
		}
		queryData := fmt.Sprintf("searchtj=%s&xq=0&jc=0&kyl=0&kclbdm=", course.courseId)
		_, respBody, _ := this.session.Post("http://202.115.47.141/student/courseSelect/freeCourse/courseList").
			Type("form").
			Send(queryData).
			AddCookies(this.cookies).End()
		js, err := simplejson.NewFromReader(strings.NewReader(respBody))
		if err != nil {
			fmt.Println(err) // error
		}
		courseInfo, err := js.Get("rwRxkZlList").String()
		if err != nil {
			fmt.Println(err) // error
		}
		courseJs, err := simplejson.NewFromReader(strings.NewReader(courseInfo))
		courseArray, err := courseJs.Array()
		if err != nil {
			fmt.Println(err) // error
		}
		if len(courseArray) == 0 {
			fmt.Printf("选课成功课程号%s\n", course.courseId)
			course.IsSelect = true
		}
		for _, v := range courseArray {
			m := v.(map[string]interface{})
			if kxh, ok := m["kxh"]; ok {
				if kxh.(string) != course.courseSeqNum {
					continue
				}
			}
			if kyl, ok := m["bkskyl"]; ok {
				kckyl, _ := kyl.(json.Number).Int64()
				if kckyl > 0 {
					kcms := m["kcm"].(string) + "_" + course.courseId
					kcids := course.courseId + "_" + course.courseSeqNum + "_" + m["zxjxjhh"].(string)
					this.doSelect(kcms, course.courseId, kcids)
				}else {
					fmt.Println("尝试选课失败，原因:课余量不足")
				}
			}
		}
	}
}

func (this *User) getTokenAndFajhh() (token, fajhh string) {
	_, respBody, _ := this.session.Get("http://202.115.47.141/student/courseSelect/freeCourse/index").
		AddCookies(this.cookies).End()
	document, err := goquery.NewDocumentFromReader(strings.NewReader(respBody))
	if err != nil {
		panic(err)
	}
	token, _ = document.Find("input[name=tokenValue]").First().Attr("value")
	fajhh, _ = document.Find("input[name=fajhh]").First().Attr("value")
	return
}

func (this *User) doSelect(kcms string, searchtj string, kcids string) {
	token, fajhh := this.getTokenAndFajhh()
	queryData := fmt.Sprintf("tokenValue=%s&kcIds=%s&kcms=%s&fajhh=%s&sj=0_0&searchtj=%s&kclbdm=",
		token, kcids, kcms, fajhh, searchtj)
	this.session.Post("http://202.115.47.141/student/courseSelect/freeCourse/waitingfor?dealType=5").
		Type("form").
		Send(queryData).
		AddCookies(this.cookies).End()
	fmt.Println("尝试选课中......等待30s以便服务器处理结果")
	time.Sleep(time.Second * time.Duration(30)) // sleep 30s, wait server finish.
}
