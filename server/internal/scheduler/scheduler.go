package scheduler

import (
	"log"

	"business-report-system/internal/service"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron   *cron.Cron
	alertSvc *service.AlertService
}

func New(alertSvc *service.AlertService) *Scheduler {
	return &Scheduler{
		cron:   cron.New(),
		alertSvc: alertSvc,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 9 * * *", func() {
		log.Println("[Scheduler] 执行预警扫描...")
		if err := s.alertSvc.ScanAlerts(); err != nil {
			log.Printf("[Scheduler] 预警扫描失败: %v", err)
		} else {
			log.Println("[Scheduler] 预警扫描完成")
		}
	})

	s.cron.AddFunc("0 0 1 * *", func() {
		log.Println("[Scheduler] 执行月度汇总...")
	})

	s.cron.Start()
	log.Println("[Scheduler] 定时任务已启动")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
