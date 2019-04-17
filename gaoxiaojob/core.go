package gaoxiaojob

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/earlzo/colly-bolt-storage/colly/bolt"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"

	"github.com/earlzo/skr/dingtalk"
	"github.com/earlzo/skr/utils"
)

var URL = url.URL{
	Scheme: "http",
	Host:   "www.gaoxiaojob.com",
}

var LOCATION *time.Location

func init() {
	LOCATION, _ = time.LoadLocation("Asia/Shanghai")
}

func FetchJobs(storageFilename string) []*Job {
	var mut sync.Mutex
	var jobs []*Job
	var options = []func(*colly.Collector){
		colly.DetectCharset(),
		colly.AllowedDomains(URL.Host),
		colly.Debugger(&utils.LogrusCollectorDebugger{}),
	}

	jobListCollector := colly.NewCollector(options...)
	options = append(options, colly.Async(true))
	jobCollector := colly.NewCollector(options...)
	storage := &bolt.Storage{
		Path: storageFilename,
	}
	if err := jobCollector.SetStorage(storage); err != nil {
		panic(err)
	}
	jobListCollector.OnHTML("ul.last_updated > li > span > a:nth-child(2)", func(e *colly.HTMLElement) {
		jobURL := e.Attr("href")
		if err := jobCollector.Visit(jobURL); err != nil {
			logrus.WithError(err).WithField("JobURL", jobURL).Warningln("error when visit job url")
		}
	})
	jobCollector.OnHTML("body.articleview", func(e *colly.HTMLElement) {
		job := Job{Meta: make(map[string]string, 5)}
		job.URL = e.Request.URL.String()
		e.DOM.Find("div.position a:nth-child(n+2)").Each(func(_ int, selection *goquery.Selection) {
			job.Categories = append(job.Categories, selection.Text())
		})
		job.Title = e.DOM.Find("div.article_left.border > h1.title-a").Text()
		job.Body = e.DOM.Find("div.article_body").Text()
		e.DOM.Find("ul.article_fenlei > li").Each(func(_ int, selection *goquery.Selection) {
			meta := strings.SplitN(selection.Text(), "：", 2)
			job.Meta[strings.TrimSpace(meta[0])] = strings.TrimSpace(meta[1])
		})
		if v, ok := job.Meta["所属省份"]; ok {
			job.Provinces = strings.Split(v, " ")
		}
		if v, ok := job.Meta["工作地点"]; ok {
			job.Locations = strings.Split(v, " ")
		}
		if v, ok := job.Meta["需求学科"]; ok {
			job.Subjects = strings.Split(v, " ")
		}
		if v, ok := job.Meta["发布时间"]; ok {
			t, err := time.ParseInLocation("2006-01-02", v, LOCATION)
			if err != nil {
				panic(err)
			}
			job.PublishedAt = &t
		}
		if v, ok := job.Meta["截止日期"]; ok {
			t, err := time.ParseInLocation("2006年1月2日", v, LOCATION)
			if err == nil {
				job.ExpireAt = &t
			}
		}
		logrus.WithFields(logrus.Fields{
			"Title":       job.Title,
			"Categories":  job.Categories,
			"Provinces":   job.Provinces,
			"Locations":   job.Locations,
			"Subjects":    job.Subjects,
			"PublishedAt": job.PublishedAt,
			"ExpireAt":    job.ExpireAt,
		}).Infoln("job fetched")
		mut.Lock()
		defer mut.Unlock()
		jobs = append(jobs, &job)
	})
	err := jobListCollector.Visit(URL.String())
	if err != nil {
		logrus.WithError(err).WithField("JobListURL", URL.String()).Panicln("error when visit job list url")
	}
	jobCollector.Wait()
	return jobs
}

func FilterJobs(jobs []*Job, keywords []string) []*Job {
	if len(keywords) == 0 {
		return jobs
	}
	var filteredJobs []*Job
	for _, job := range jobs {
		text := strings.Join(append(job.Provinces, append(job.Locations, append(job.Categories, append(job.Subjects)...)...)...), " ")
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				filteredJobs = append(filteredJobs, job)
			}
		}
	}
	return filteredJobs
}

func Notify(webhookURL string, jobs []*Job) error {
	feedCardLinks := make([]dingtalk.FeedCardLink, len(jobs))
	for i, job := range jobs {
		pictureURL := utils.DummyImage(400, 400, "1ff22a", "0011ff", job.Categories[len(job.Categories)-1], "png")
		feedCardLinks[i] = dingtalk.FeedCardLink{
			Title:      job.Title,
			MessageURL: job.URL,
			PictureURL: pictureURL,
		}
	}

	for _, message := range []dingtalk.Message{
		{
			Type: "feedCard",
			FeedCard: &dingtalk.FeedCard{
				Links: feedCardLinks,
			},
			At: &dingtalk.At{IsAtAll: true},
		},
	} {
		resp, err := dingtalk.Send(http.DefaultClient, webhookURL, &message)
		if err != nil {
			return err
		}
		if resp.ErrorCode != 0 {
			return fmt.Errorf(
				"%d %s", resp.ErrorCode, resp.ErrorMessage)
		}
	}
	return nil
}
