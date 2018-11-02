package main

import (
	"flag"
	"net/http"
	"regexp"
	"time"

	"github.com/jelinden/rssfetcher/app/domain"
	"github.com/jelinden/rssfetcher/app/handler"
	"github.com/jelinden/rssfetcher/app/mongo"
)

var (
	validPath    = regexp.MustCompile("^/(edit|save|view|delete)/([a-zA-Z0-9]*)$")
	mongoAddress = flag.String("address", "localhost", "mongo address")
	env          = flag.String("env", "dev", "environment")
)

func main() {
	flag.Parse()
	mongo.InitMongo(*mongoAddress)
	runFeedFetcher(*env)
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(handler.ViewHandler))
	http.HandleFunc("/delete/", makeHandler(handler.RemoveHandler))
	http.HandleFunc("/edit/", makeHandler(handler.EditHandler))
	http.HandleFunc("/save/", makeHandler(handler.SaveHandler))
	http.HandleFunc("/save/category", makeHandler(handler.SaveCategoryHandler))
	http.HandleFunc("/save/subcategory", makeHandler(handler.SaveSubCategoryHandler))

	http.ListenAndServe(":9200", nil)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func runFeedFetcher(env string) {
	if env == "dev" {
		go doEvery(60*time.Second, mongo.GetFeeds)
	} else {
		go doEvery(80*time.Second, mongo.GetFeeds)
	}
}

func doEvery(d time.Duration, feeds func(args ...bool) []domain.Feed) {
	for _ = range time.Tick(d) {
		feedList := feeds(true)
		mongo.GetNews(feedList)
	}
}