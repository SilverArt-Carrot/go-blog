package views

import (
	"Blog/models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

const domain string = "http://localhost/"

// 错误检查
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 反转义
func safe(s string) template.HTML {
	return template.HTML(s)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.New("detail").Funcs(template.FuncMap{
		"getUrlPrefix":getUrlPrefix,
	}).ParseFiles("./templates/base.html", "./templates/index.html")
	if err != nil {
		fmt.Fprintf(writer, "parse :%v", err)
		return
	}

	posts := models.GetAllPosts()
	cates := models.GetAllCategories()
	tags := models.GetAllTags()
	archives := models.GetArchives()
	numOfPost, numOfCate, numOfTag := models.GetInfo()

	//当执行的模板不取名字时可以用Execute，此时go为模板指定名字为模板文件扩展名
	//err = tmpl.Execute(writer, "hello world")
	//ExecuteTemplate与Execute类似，但可以执行指定名字的模板，适用于在多模板嵌套的情况下使用，在模板文件中还需define模板的名字
	err = tmpl.ExecuteTemplate(writer, "base", map[string]interface{}{
		"Posts": 	posts,
		"Cates": 	cates,
		"Tags":		tags,
		"Archives":	archives,
		"numOfPost":numOfPost,
		"numOfCate":numOfCate,
		"numOfTag":	numOfTag,
	})
	if err != nil {
		fmt.Fprintf(writer, "execute :%v", err)
		return
	}
}

func Detail(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.New("detail").Funcs(template.FuncMap{
		"getUrlPrefix":	getUrlPrefix,
		"safe":			safe,
	}).ParseFiles("./templates/base.html", "./templates/detail.html")

	if err != nil {
		fmt.Fprintf(writer, "parse :%v", err)
		return
	}

	id, _ := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
	post := models.GetDetailByPostID(int(id))
	post.Body = string(blackfriday.MarkdownCommon([]byte(post.Body)))  //html转markdown
	tags := models.GetAllTags()
	cates := models.GetAllCategories()
	archives := models.GetArchives()
	numOfPost, numOfCate, numOfTag := models.GetInfo()

	err = tmpl.ExecuteTemplate(writer, "base", map[string]interface{}{
		"Img":		post.Img,
		"Title":	post.Title,
		"Year":		post.CreatedTime.Year,
		"Month":	post.CreatedTime.Month,
		"Day":		post.CreatedTime.Day,
		"Body":		post.Body,
		"Tags":		tags,
		"Cates":	cates,
		"Archives":	archives,
		"numOfPost":numOfPost,
		"numOfCate":numOfCate,
		"numOfTag":	numOfTag,
	})
	if err != nil {
		fmt.Fprintf(writer, "execute :%v", err)
		return
	}
}

func Cate(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.New("cate").Funcs(template.FuncMap{
		"getUrlPrefix":getUrlPrefix,
	}).ParseFiles("./templates/base.html", "./templates/index.html")
	if err != nil {
		fmt.Fprintf(writer, "parse :%v", err)
		return
	}

	id, _ := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
	posts := models.GetPostsByCateId(int(id))
	tags := models.GetAllTags()
	cates := models.GetAllCategories()
	archives := models.GetArchives()
	numOfPost, numOfCate, numOfTag := models.GetInfo()

	err = tmpl.ExecuteTemplate(writer, "base", map[string]interface{}{
		"Posts":	posts,
		"Tags":		tags,
		"Cates":	cates,
		"Archives":	archives,
		"numOfPost":numOfPost,
		"numOfCate":numOfCate,
		"numOfTag":	numOfTag,
	})
	if err != nil {
		fmt.Fprintf(writer, "execute :%v", err)
		return
	}
}

func Tag(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.New("tag").Funcs(template.FuncMap{
		"getUrlPrefix":getUrlPrefix,
	}).ParseFiles("./templates/base.html", "./templates/index.html")
	if err != nil {
		fmt.Fprintf(writer, "parse :%v", err)
		return
	}

	id, _ := strconv.ParseInt(mux.Vars(request)["id"], 10, 64)
	posts := models.GetPostsByTagId(int(id))
	tags := models.GetAllTags()
	cates := models.GetAllCategories()
	archives := models.GetArchives()
	numOfPost, numOfCate, numOfTag := models.GetInfo()

	err = tmpl.ExecuteTemplate(writer, "base", map[string]interface{}{
		"Posts":	posts,
		"Tags":		tags,
		"Cates":	cates,
		"Archives":	archives,
		"numOfPost":numOfPost,
		"numOfCate":numOfCate,
		"numOfTag":	numOfTag,
	})
	if err != nil {
		fmt.Fprintf(writer, "execute :%v", err)
		return
	}
}

func Archive(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.New("arch").Funcs(template.FuncMap{
		"getUrlPrefix":getUrlPrefix,
	}).ParseFiles("./templates/base.html", "./templates/index.html")
	if err != nil {
		fmt.Fprintf(writer, "parse :%v", err)
		return
	}

	year, _ := strconv.ParseInt(mux.Vars(request)["year"], 10, 64)
	month, _ := strconv.ParseInt(mux.Vars(request)["month"], 10, 64)
	posts := models.GetPostsByArchive(int(year), int(month))
	tags := models.GetAllTags()
	cates := models.GetAllCategories()
	archives := models.GetArchives()
	numOfPost, numOfCate, numOfTag := models.GetInfo()

	err = tmpl.ExecuteTemplate(writer, "base", map[string]interface{}{
		"Posts":	posts,
		"Tags":		tags,
		"Cates":	cates,
		"Archives":	archives,
		"numOfPost":numOfPost,
		"numOfCate":numOfCate,
		"numOfTag":	numOfTag,
	})
	if err != nil {
		fmt.Fprintf(writer, "execute :%v", err)
		return
	}
}

//测试静态文件服务
func TestStaticFileServer(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("./templates/test.html")
	if err != nil {
		fmt.Fprintf(writer, "parse :%v", err)
		return
	}

	err = tmpl.Execute(writer, nil)
	if err != nil {
		fmt.Fprintf(writer, "execute :%v", err)
		return
	}
}

//ParseFiles与Funcs联合执行出现问题，所以特地写了一个读取模板文件的函数，用于配合Funcs
func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	return string(content)
}

//模板函数，在html页面中得到url前缀
func getUrlPrefix() string {
	return domain
}

//创建的模板函数，用于在模板中生成index页面的url
func redirectToIndex() string {
	return domain + "index"
}
