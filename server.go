package main

import (
	"Blog/views"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//整个project的执行入口，用于url匹配正确的执行函数，执行函数通过models获取数据库中的数据
//再通过template渲染给浏览器，参照MVC架构

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/index", views.Index)
	router.HandleFunc("/detail/{id:[0-9]+}", views.Detail)
	router.HandleFunc("/cate/{id:[0-9]+}", views.Cate)
	router.HandleFunc("/tag/{id:[0-9]+}", views.Tag)
	router.HandleFunc("/archives/{year:[0-9]+}/{month:[0-9]+}", views.Archive)

	// 启动静态服务
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 标准库启用静态服务
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("The server has been started...\n Please access the url http://127.0.0.1/index to browse the pages\n ")
	log.Fatal(http.ListenAndServe("", router))
}
