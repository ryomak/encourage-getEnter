package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type User struct {
	ID         string
	Name       string
	Yomi       string
	Mentor     string
	Phone      string
	Univ       string
	Department string
	Intern     string
	Gender     string
	Url        string
	Eval       string
	Science    bool
}

var fetchUrl = ""
var typeFunc func(int) interface{}
var config Config

func init() {
	config = GetConfig()
	//query
	values := url.Values{}
	for _, v := range config.Qs {
		values.Add(v.Key, v.Val)
	}
	fetchUrl = config.Url + values.Encode()
	//functype
	switch config.WriteType {
	case "light":
		typeFunc = light
	case "detail":
		typeFunc = detail
	}
}

func main() {
	config := GetConfig()
	//page数
	page := getLastPage()
	fmt.Printf("url:%v\npage:%v\n", fetchUrl, page)
	typeFunc(page)
}

//メソッド一覧
func detail(page int) []User {
	var users []User
	//スクレイピング
	for i := 1; i <= page; i++ {
		users = append(users, fetchDetailBody(i)...)
	}
	WriteCsv(config.WriteFile, users)
	return users
}

func light(page int) []User {
	var users []User
	//スクレイピング
	for i := 1; i <= page; i++ {
		users = append(users, fetchBody(i)...)
	}
	WriteCsv(config.WriteFile, users)
	return users
}

func printData(page int) err {
	var users []User
	//スクレイピング
	for i := 1; i <= page; i++ {
		users = append(users, fetchDetailBody(i)...)
	}
	scienceNum := 0
	dojo := 0
	for _, v = range users {
		if v.Sciense == "理系" {
			scienceNum++
		}
		if strings.Index(v.Univ, "女") != -1 {
			dojo++
		}
	}
	fmt.Printf("登録者数:%d人\n理系:%d人\n同女:%d人", len(users), scienceNum, dojo)
}

//便利関数
func fetchBody(p int) []User {
	p := strconv.Itoa(i)
	resp := Get(fetchUrl + "&page=" + p)
	users := []User{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			user := User{}
			user.getUserStatus(s)
			users = append(users, user)
		}
	})
	return users
}

func fetchDetailBody(p int) []User {
	p := strconv.Itoa(i)
	resp := Get(fetchUrl + "&page=" + p)
	users := []User{}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			user := User{}
			user.getUserStatus(s)
			user.getUserDetail()
			users = append(users, user)
		}
	})
	return users
}

func (u *User) getUserDetail() {
	resp := Get(u.Url)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:

		case 1:
			u.Yomi = s.Text()
		case 2:
			u.Gender = s.Text()
		case 3:
			u.Department = s.Text()
			if strings.Index(s.Text(), "理工") != -1 {
				u.Science = true
			} else if strings.Index(s.Text(), "医") != -1 {
				u.Science = true
			} else {
				u.Science = false
			}
		default:
		}
	})
}

func getLastPage() int {
	resp := Get(fetchUrl + "&page=1")
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	u, _ := doc.Find("li").Last().Find("a").Attr("href")
	in := strings.Index(u, "page=")
	page := u[(in + 5):(in + 7)]
	iPage, _ := strconv.Atoi(page)
	return iPage
}

func (u *User) getUserStatus(s *goquery.Selection) {
	s.Find("td").Each(func(k int, se *goquery.Selection) {
		switch k {
		case 1:
			u.Name = se.Text()
			path, _ := se.Find("a").Attr("href")
			u.Url = "https://id.en-courage.com" + path
		case 2:
			u.ID = se.Text()
		case 3:
			u.Mentor = se.Text()
		case 4:
			u.Phone = se.Text()
		case 5:
			u.Eval = se.Text()
		case 6:
			u.Univ = se.Text()
		case 8:
			u.Intern = se.Text()
		default:

		}
	})
}

func Get(url string) *http.Response {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "_ga=GA1.2.1457449303.1520512918; _uma_session=32d98c4dac82bd2041e12f6f3b152324; _hjIncludedInSample=1; _gid=GA1.2.1305894491.1525069773")
	client := new(http.Client)
	resp, _ := client.Do(req)
	return resp
}
