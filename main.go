package main

import (
	"fmt"
	"net/http"
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


var url = "https://id.en-courage.com/admin/users/search?univ_name_or_graduate_name_cont%5D=%E5%90%8C%E5%BF%97%E7%A4%BE&q%5Buser_educational_last_educational_year_eq%5D=2020"

func main() {
	page := getLastPage()
	for i:=1 ;i<=page ;i++{
		resp := Get(url)
	defer resp.Body.Close()
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	users := []User{}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			user := User{}
			user.getUserStatus(s)
			user.getUserDetail()
			users = append(users, user)
		}
	})
	fmt.Printf("%+v \n", users)
}

func connectKey(m map[string][string])string{
	key := ""
	for k,v := range m {
		str= "&" + k + "=" + v
		key += str
	}
	return key
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
			if strings.Index(s.Text(), "理") != -1 {
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
	resp := Get(url)
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
