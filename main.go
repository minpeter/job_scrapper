package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	location string
	title string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=linux&limit=50"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i:=0; i<totalPages; i++ {
		extractedJob := getPage(i)
		jobs = append(jobs, extractedJob...)
	}
	
	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs), "jobs")
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(utf8bom)
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"LINK", "Title", "Location", "Salary", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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
		title := CleanString(card.Find(".jobTitle>span").Text())
		location := CleanString(card.Find(".companyLocation").Text())
		salary := CleanString(card.Find(".salary-snippet").Text())
		summary := CleanString(card.Find(".job-snippet").Text())

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
	
	searchCountStr := doc.Find("#searchCountPages").Text()
	searchCountStr = strings.Replace(searchCountStr, ",", "", -1)
	searchCountStr = strings.Replace(searchCountStr, "1페이지 결과 ", "", -1)
	searchCountStr = strings.Replace(searchCountStr, "건", "", -1)
	pages, err = strconv.Atoi(CleanString(searchCountStr))
	checkErr(err)
	fmt.Println("find ", pages / 50, "pages")
	
	return pages / 50
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