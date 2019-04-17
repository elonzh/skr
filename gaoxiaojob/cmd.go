package gaoxiaojob

import "github.com/sirupsen/logrus"

func Run(storageFilename, webhookURL string, keywords []string) error {
	var err error
	jobs := FetchJobs(storageFilename)
	logrus.WithField("JobsNum", len(jobs)).Infoln("抓取最新招聘信息完成")
	jobs = FilterJobs(jobs, keywords)
	logrus.WithFields(logrus.Fields{
		"FilteredJobsNum": len(jobs),
		"Keywords":        keywords,
	}).Infoln("过滤最新招聘信息完成")
	if len(jobs) > 0 {
		err = Notify(webhookURL, jobs)
		logrus.WithError(err).WithField("WebhookURL", webhookURL).Warningln("推送招聘信息")
	}
	return err
}
