package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
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
	Update     string
}
type us []User

//global変数
var fetchUrl = ""
var typeFunc func(int) []User
var config Config
var EnterList []User

func init() {
	config = GetConfig()
	//query
	values := url.Values{}
	for _, v := range config.Qs {
		values.Add(v.Key, v.Val)
	}
	EnterList = []User{}
	fetchUrl = config.Url + values.Encode()
	//functype
	switch config.WriteType {
	case "light":
		typeFunc = light
	case "detail":
		typeFunc = detail
	case "printData":
		typeFunc = printData
	default:
		fmt.Println("関数が存在しません")
	}
}

func main() {
	//page数
	page := getLastPage()
	fmt.Printf("url:%v\npage:%v\n", fetchUrl, page)
	typeFunc(page)
}

//concurrency
func fetchDetailWorker(num int) []User {
	workerNum := 4
	var wg sync.WaitGroup
	q := make(chan int, 6)
	users := []User{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go users.fetchDetail(&wg, q)
	}
	for i := 1; i <= num; i++ {
		q <- i
	}
	close(q)
	wg.Wait()
	return users
}

func (users []User) fetchDetail(wg *sync.WaitGroup, page int) {
	defer wg.Done()
	for {
		s, ok := <-q
		if !ok {
			return
		}
		users = append(users, fetchDetailBody(s)...)
	}
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

func printData(page int) []User {
	//スクレイピング
	users := fetchDetailWorker(page)
	scienceNum := 0
	dojo := 0
	ge := 0
	for _, v := range users {
		if v.Science {
			scienceNum++
		}
		if strings.Index(v.Department, "女") != -1 {
			fmt.Println(v)
			dojo++
		}
		if strings.Index(v.Eval, "GE") != -1 {
			ge++
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"全体登録者数", "理系登録者数", "同女数", "GE数"})
	table.Append([]string{strconv.Itoa(len(users)), strconv.Itoa(scienceNum), strconv.Itoa(dojo), strconv.Itoa(ge)})
	table.Render()
	return users
}

//便利関数
func fetchBody(i int) []User {
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

func fetchDetailBody(i int) []User {
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
		case 9:
			u.Update = se.Text()
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
