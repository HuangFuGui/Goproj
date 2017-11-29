/*
 github.com/astaxie/beego/session模块设计思路来源于database/sql/driver，先定义好接口，
 然后具体的存储session的结构实现相应的接口并注册后，相应功能这样就可以使用了。
*/
package main

import (
	"net/http"
	"log"
	"fmt"
	"html/template"
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

func init(){
	config := &session.ManagerConfig{CookieName:"gosessionid", EnableSetCookie:true, Maxlifetime:60, CookieLifeTime:60, Gclifetime:60}
	globalSessions, _ = session.NewManager("memory", config)

	//GC充分利用了time包中的定时器功能，当超时maxLifeTime之后调用GC函数，这样就可以保证Gclifetime时间内的session都是可用的，类似的方案也可以用于统计在线用户数之类的。
	go globalSessions.GC()
}

func ServeHTTP(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
		return
	}

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		panic(err)
	}
	defer sess.SessionRelease(w)

	//先读取该用户session信息，若不存在/已过期为<nil>
	username := sess.Get("username")
	fmt.Printf("从session中获取username：%v\n", username)

	//根据表单信息重新设置session
	r.ParseForm()
	username = r.FormValue("username")
	fmt.Printf("从表单数据中获取username：%v\n", username)
	sess.Set("username", username)

	fmt.Fprintf(w, "hello %v~", username)
}

func SessionBackend(w http.ResponseWriter, r *http.Request){
	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		panic(err)
	}
	defer sess.SessionRelease(w)

	username := sess.Get("username")
	fmt.Fprintf(w, "session情况：%v", username)
}

func main(){
	http.HandleFunc("/index", ServeHTTP)
	http.HandleFunc("/sessionbackend", SessionBackend)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("main main() err:=http.ListenAndServe error")
	}
}