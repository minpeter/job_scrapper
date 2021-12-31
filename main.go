package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	location string
	title string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()    //페이지 전체를 불러오도록 수정해야됨

	for i:=0; i<totalPages; i++ {
		extractedJob := getPage(i)
		jobs = append(jobs, extractedJob...)
	}
	
	fmt.Println(jobs)
}

func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".mosaic-zone").Find("a").Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		if job.id != "" {
			jobs = append(jobs, job)
		}
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {

	id, tf := card.Attr("data-jk")
	
	if tf {
		title := card.Find(".jobTitle>span").Text()
		location := card.Find(".companyLocation").Text()
		salary := card.Find(".salary-snippet").Text()
		summary := card.Find(".job-snippet").Text()

		return extractedJob{id:id, title:title, location:location, salary:salary, summary:summary}
	} else {
		return extractedJob{}
	}
	
}


func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Status code error: ", res.StatusCode , res.Status)
	}
}