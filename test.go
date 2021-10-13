package main

import (
	"Blog/models"
	"Blog/mylog"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
)

func main() {

	//test1()
	//test2()
	//test3()
	test4()
}

func test1() {
	r := mux.NewRouter()
	r.HandleFunc("/test/{id:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		fmt.Println(vars["id"])
		writer.Write([]byte("github.com/gorilla/mux"))
	})

	log.Fatal(http.ListenAndServe("", r))
	//srv := http.Server{
	//	Addr:              "",
	//	Handler:           r,
	//}
	//srv.ListenAndServe()
}

func test2() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "wo ta ma de lie kai %s \n", "le")
	})
	log.Fatal(http.ListenAndServe("", nil))
}

func test3() {
	p := models.Post{
		Id:          5,
		Title:       "test",
		Author:      "carrot",
		Body:        "玩玩的",
		Category:    2,
		Tag:         1,
		Img:         "green.png",
		Excerpt:     "玩l525252",
		CreatedTime: models.Time{2020,8,1},
		ModifyTime:  models.Time{2020,9,1},
	}

	//c := models.Category{
	//	Id:   8,
	//	Name: "ssssss",
	//}
	//
	//t := models.Tag{
	//	Id:   14,
	//	Name: "sdsdad",
	//}
	//
	//models.Insert(p)  //通过
	//models.Insert(c)   //通过
	//models.Insert(t)   //通过
	//
	//models.Delete("Post", 5)   //通过，删除不存在的记录不会报错
	//models.Delete("Tag", 13)  //通过
	//models.Delete("Category", 8)   //通过
	//
	//models.UpdateForPost(p, 4)  //通过
	//models.UpdateForCategory(c, 8)  //通过
	//models.UpdateForTag(t,13)   //通过

	//t := reflect.TypeOf(p)
	//fmt.Println("Type: ", t.Name())
	//v := reflect.ValueOf(p)
	//fmt.Println("Fields: ")
	//for i := 0; i < t.NumField(); i++ {
	//	f := t.Field(i)
	//	val := v.Field(i).Interface()
	//	fmt.Printf("%6s: %v = %v\n",f.Name,f.Type,val)
	//}

	t := reflect.TypeOf(models.Tag{
		Id:   1,
		Name: "test",
	})
	if t.Kind() == reflect.Struct {
		fmt.Printf("name is: %v, kind is: %v , string is: %v \n", t.Name(), t.Kind(), t.String())
		num := t.NumField()
		for i := 0; i < num; i++ {
			fmt.Printf("name is : %v, index is : %v, tag is :%v, type is : %v\n",
				t.Field(i).Name, t.Field(i).Index, t.Field(i).Tag, t.Field(i).Type)
		}
	}

	v := reflect.ValueOf(models.Tag{
		Id:   2,
		Name: "demo",
	})
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fmt.Printf("kind of field is:%v, type is:%v, value is:%v \n",
				v.Field(i).Kind(), v.Field(i).Type(), v.Field(i).Interface())
		}
	}

	v2 := reflect.ValueOf(p)
	if v2.Field(8).Kind() == reflect.Struct {
		fmt.Println(v2.Field(8).Interface() )
	}
}

func test4() {
	l := mylog.NewLogger(mylog.INFO, "./", "server.log")
	defer l.Close()
	for i:=0; i<5000; i++{
		l.DEBUG("%d号正在执行中", 5)
		l.ERROR("有点小错误%d", 12)
		l.WARN("管他呢")
		l.FATAL("完了，我傻了")
	}
}
