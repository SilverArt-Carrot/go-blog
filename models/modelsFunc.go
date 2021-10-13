package models

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)
	return db
}

func GetAllPosts() []Post {
	db := openDatabase()

	posts := make([]Post, 0)

	rows, err := db.Query("select * from post;")
	checkErr(err)

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Id, &p.Title, &p.Author,
			&p.Body, &p.Category, &p.Tag, &p.Img, &p.Excerpt,
			&p.CreatedTime.Year, &p.CreatedTime.Month, &p.CreatedTime.Day,
			&p.ModifyTime.Year, &p.ModifyTime.Month, &p.ModifyTime.Day)
		checkErr(err)
		posts = append(posts, p)
	}

	defer db.Close()
	return posts
}

func GetAllTags() []Tag {
	db := openDatabase()

	tags := make([]Tag, 0)

	rows, err := db.Query("select * from Tag;")
	checkErr(err)

	for rows.Next() {
		var t Tag
		err := rows.Scan(&t.Id, &t.Name)
		checkErr(err)
		tags = append(tags, t)
	}

	defer db.Close()
	return tags
}

func GetAllCategories() []Category {
	db := openDatabase()

	cates := make([]Category, 0)

	rows, err := db.Query("select * from Category;")
	checkErr(err)

	for rows.Next() {
		var c Category
		err := rows.Scan(&c.Id, &c.Name)
		checkErr(err)
		cates = append(cates, c)
	}

	defer db.Close()
	return cates
}

func GetArchives() []Archives {
	db := openDatabase()

	archives := make([]Archives, 0)

	rows, err := db.Query("select DISTINCT CreatedTimeYear, CreatedTimeMonth from post;")
	checkErr(err)

	for rows.Next() {
		var a Archives
		err := rows.Scan(&a.Year, &a.Month)
		checkErr(err)
		archives = append(archives, a)
	}

	defer db.Close()
	return archives
}

func GetInfo()(numOfPost, numOfCate, numOfTag int) {
	db := openDatabase()

	err := db.QueryRow("select COUNT(*) from Post;").Scan(&numOfPost)
	checkErr(err)
	err = db.QueryRow("select COUNT(*) from Category;").Scan(&numOfCate)
	checkErr(err)
	err = db.QueryRow("select COUNT(*) from Tag;").Scan(&numOfTag)
	checkErr(err)

	defer db.Close()
	return numOfPost,numOfCate,numOfTag
}

func GetPostsByTagId(id int) []Post {
	db := openDatabase()

	posts := make([]Post, 0)

	stmt, err := db.Prepare("select * from Post where Tag=?;")
	checkErr(err)

	rows, err := stmt.Query(id)
	checkErr(err)

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Id, &p.Title, &p.Author,
			&p.Body, &p.Category, &p.Tag, &p.Img, &p.Excerpt,
			&p.CreatedTime.Year, &p.CreatedTime.Month, &p.CreatedTime.Day,
			&p.ModifyTime.Year, &p.ModifyTime.Month, &p.ModifyTime.Day)
		checkErr(err)
		posts = append(posts, p)
	}

	defer db.Close()
	return posts
}

func GetPostsByCateId(id int) []Post {
	db := openDatabase()

	posts := make([]Post, 0)

	stmt, err := db.Prepare("select * from Post where Category=?;")
	checkErr(err)

	rows, err := stmt.Query(id)
	checkErr(err)

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Id, &p.Title, &p.Author,
			&p.Body, &p.Category, &p.Tag, &p.Img, &p.Excerpt,
			&p.CreatedTime.Year, &p.CreatedTime.Month, &p.CreatedTime.Day,
			&p.ModifyTime.Year, &p.ModifyTime.Month, &p.ModifyTime.Day)
		checkErr(err)
		posts = append(posts, p)
	}

	defer db.Close()
	return posts
}

func GetPostsByArchive(year, month int) []Post {
	db := openDatabase()

	posts := make([]Post, 0)

	stmt, err := db.Prepare("select * from Post where CreatedTimeYear=? and CreatedTimeMonth=?;")
	checkErr(err)

	rows, err := stmt.Query(year, month)
	checkErr(err)

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Id, &p.Title, &p.Author,
			&p.Body, &p.Category, &p.Tag, &p.Img, &p.Excerpt,
			&p.CreatedTime.Year, &p.CreatedTime.Month, &p.CreatedTime.Day,
			&p.ModifyTime.Year, &p.ModifyTime.Month, &p.ModifyTime.Day)
		checkErr(err)
		posts = append(posts, p)
	}

	defer db.Close()
	return posts
}

func GetDetailByPostID(id int) Post {
	db := openDatabase()

	stmt, err := db.Prepare("select * from Post where id=?")
	checkErr(err)

	var p Post
	err = stmt.QueryRow(id).Scan(&p.Id, &p.Title, &p.Author,
		&p.Body, &p.Category, &p.Tag, &p.Img, &p.Excerpt,
		&p.CreatedTime.Year, &p.CreatedTime.Month, &p.CreatedTime.Day,
		&p.ModifyTime.Year, &p.ModifyTime.Month, &p.ModifyTime.Day)
	checkErr(err)

	defer db.Close()
	return p
}

func createQueryForInsert(q interface{}) (string, error) {   // 用于创建插入SQL语句
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		query := fmt.Sprintf("insert into %s values(",t)
		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i==0 {
					query = fmt.Sprintf("%s%d",query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %d",query, v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s'%s'", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, '%s'", query, v.Field(i).String())
				}
			case reflect.Struct:
				v2 := reflect.ValueOf(v.Field(i).Interface())
				for i := 0; i < v2.NumField(); i++ {
					query = fmt.Sprintf("%s, %v", query, v2.Field(i).Int())
				}
			default:
				return "", errors.New("unsupported type from reflect")
			}
		}
		query = fmt.Sprintf("%s);", query)
		return query, nil
	}
	return "", errors.New("unsupported type from reflect")
}

func Insert(q interface{})  {
	db := openDatabase()

	stmt, err := createQueryForInsert(q)
	checkErr(err)

	_, err = db.Exec(stmt)
	if err == nil { log.Println("insert successfully...") }
	checkErr(err)
	defer db.Close()
}

func Delete(tableName string, id int) {
	db := openDatabase()

	q := fmt.Sprintf("delete from %s where id=%d;", tableName, id)
	_, err := db.Exec(q)
	if err == nil { log.Println("delete successfully...") }
	checkErr(err)
	defer db.Close()
}

func createQueryForUpdate(q interface{}) (string, error) {
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		tn := reflect.TypeOf(q).Name()
		query := fmt.Sprintf("update %s set", tn)
		t := reflect.TypeOf(q)
		for i := 1; i < t.NumField(); i++ {
			if i == 1 {
				query = fmt.Sprintf("%s %s=?", query, t.Field(i).Name)
			} else {
				query = fmt.Sprintf("%s, %s=?", query, t.Field(i).Name)
			}
		}
		query = fmt.Sprintf("%s where id=?;", query)
		return query, nil
	}
	return "", errors.New("unsupported type from reflect")
}

func UpdateForPost(p Post, id int) {
	db := openDatabase()

	stmt, err := createQueryForUpdate(p)
	checkErr(err)
	stmt = strings.Replace(stmt, "CreatedTime=?, ModifyTime=?",
		"CreatedTimeYear=?, CreatedTimeMonth=?, CreatedTimeDay=?, ModifyTimeYear=?, ModifyTimeMonth=?, ModifyTimeDay=?", 1)

	_, err = db.Exec(stmt, p.Title, p.Author, p.Body, p.Category, p.Tag, p.Img, p.Excerpt,
		p.CreatedTime.Year, p.CreatedTime.Month, p.CreatedTime.Day,
		p.ModifyTime.Year, p.ModifyTime.Month, p.ModifyTime.Day, id)
	if err == nil { log.Println("update successfully...") }
	checkErr(err)
	defer db.Close()
}

func UpdateForCategory(c Category, id int) {
	db := openDatabase()

	stmt, err := createQueryForUpdate(c)
	checkErr(err)

	_, err = db.Exec(stmt, c.Name, id)
	if err == nil { log.Println("update successfully...") }
	checkErr(err)
	defer db.Close()
}

func UpdateForTag(t Tag, id int) {
	db := openDatabase()

	stmt, err := createQueryForUpdate(t)
	checkErr(err)

	_, err = db.Exec(stmt, t.Name, id)
	if err == nil { log.Println("update successfully...") }
	checkErr(err)
	defer db.Close()
}
