package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sfv/sfv"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-xorm/xorm"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/gomail.v1"
)

var Db *sql.DB

type Res struct {
	Count  int64       `json:"count"`
	Sfvers []sfv.Sfver `json:"sfvers"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, err := os.Open("public/html/index.html")
	defer file.Close()
	if err != nil {
		fmt.Fprint(w, "Error: %v\n", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprint(w, "Error: %v\n", err)
	}
	fmt.Println("Get Request for Index ....")
	fmt.Fprint(w, string(data))
}

func UserGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func UserPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func SfversGetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	count := sfverService.Count()
	log.Println("count:", count)
	sfvers := sfverService.List()
	json.NewEncoder(w).Encode(Res{count, sfvers})
}

func SfversPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sfver := sfv.Sfver{}
	json.NewDecoder(r.Body).Decode(&sfver)
	fmt.Println("email:", sfver)
	if ok := len(sfver.Email) > 0; !ok {
		return
	}
	defer r.Body.Close()
	// sfver := sfv.Sfver{Email: sfver.Email}
	ok := sfverService.Insert(sfver)
	if ok {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		count := sfverService.Count()
		log.Println("count2:", count)
		sfvers := sfverService.List()
		json.NewEncoder(w).Encode(Res{count, sfvers})
	}
}

var (
	Engine       *xorm.Engine
	sfverService *sfv.SfvService
)

var NOT_AVAILABLE = "No places available"

func sfvTask(url string) {
	res, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Println("got error:", err)
		return
	}
	defer res.Body.Close()

	// use utfBody using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		// handler error
		fmt.Println("got error:", err)
		return
	}
	doc.Find("div.textblock").Each(func(i int, s *goquery.Selection) {
		alt, ok := s.Find("IMG").Attr("alt")
		if !ok {
			fmt.Println(alt)
			return
		}
		fmt.Println("get from html as result:", alt)
		if strings.Contains(alt, NOT_AVAILABLE) {
			// notify("Hi Sfver:</br> No available space currently. please try one via the link below: </br><a href=\"http://www.immigration.govt.nz/migrant/stream/work/silverfern/jobsearch.htm\">http://www.immigration.govt.nz/migrant/stream/work/silverfern/jobsearch.htm</a>")
			fmt.Println(time.Now().Format("2006-01-02 15:04:06"), alt, ", please waiting for available")
		} else {
			notify("Hi Sfver:</br> Some of SFVs are available! please try one via the link below: </br><a href=\"http://www.immigration.govt.nz/migrant/stream/work/silverfern/jobsearch.htm\">http://www.immigration.govt.nz/migrant/stream/work/silverfern/jobsearch.htm</a>")
			fmt.Println(time.Now().Format("2006-01-02 15:04:06"), "Some of SFVs are available! :-)")
		}
	})
}

func sfvTaskSchedule() {
	// backupUrl := "https://www.immigration.govt.nz/secure/Login+Silver+Fern.htm"
	url := "http://www.immigration.govt.nz/migrant/stream/work/silverfern/jobsearch.htm"
	for {
		sfvTask(url)
		time.Sleep(120 * 1e9)
		fmt.Println("Next task will be started 120s later...")
	}
}

func notify(body string) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "sfver@qq.com")
	msg.SetHeader("To", "ckjacket@163.com")
	msg.SetAddressHeader("Bcc", "879939101@qq.com", "nehe")
	sfvers := sfverService.List()
	for _, sfver := range sfvers {
		log.Println("Ready to send notification to ", sfver.Email)
		msg.SetAddressHeader("Bcc", sfver.Email, sfver.Name)
	}
	msg.SetHeader("Subject", "Silver Fern Visa Available Notification")
	msg.SetBody("text/html", body)

	mailer := gomail.NewMailer("smtp.qq.com", "sfver", "sfver2015", 25)
	if err := mailer.Send(msg); err != nil {
		fmt.Println("sendng notification failed:", err)
	}
}

func main() {
	router := httprouter.New()
	router.ServeFiles("/js/*filepath", http.Dir("public/js"))
	router.ServeFiles("/css/*filepath", http.Dir("public/css"))
	router.ServeFiles("/bower_components/*filepath", http.Dir("public/bower_components"))
	router.GET("/", Index)
	router.GET("/user", UserGetHandler)
	router.POST("/user", UserPostHandler)

	// REST API for sfver
	router.GET("/sfvers", SfversGetHandler)
	router.POST("/sfvers", SfversPostHandler)

	go sfvTaskSchedule()

	var err error
	Engine, err = xorm.NewEngine("mysql", "test:1234@/sfver?charset=utf8")
	if err != nil {
		log.Fatal("database err:", err)
	}
	Engine.ShowDebug = true
	Engine.ShowErr = true
	Engine.ShowWarn = true
	Engine.ShowSQL = true
	f, err := os.Create("sql.log")
	if err != nil {
		log.Fatal("sql log create failed", err)
	}
	defer f.Close()
	Engine.Logger = xorm.NewSimpleLogger(f)

	err = Engine.Sync2(new(sfv.Sfver))
	if err != nil {
		log.Fatal("database sync failed", err)
	}
	// Engine.DropTables("sfver")
	sfverService = sfv.DefaultSfvService(Engine)
	log.Println("Listening at 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
