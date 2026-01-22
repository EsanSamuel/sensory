package workers

import (
	"fmt"
	"log"
	"os"

	"github.com/EsanSamuel/sensory/jobs"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Redis connection
func NewRedisPool(redisURL string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisURL)
		},
	}
}

var RedisURL = os.Getenv("REDIS_URL")
var redisPool *redis.Pool = NewRedisPool(RedisURL)

func SendEmailQueue(email string, userId string, logId string) {
	var enqueuer = work.NewEnqueuer("emailQueue", redisPool)

	_, err := enqueuer.Enqueue("send_email", work.Q{"email_addr": email, "user_id": userId, "log_id": logId})

	if err != nil {
		fmt.Println("Error queuing email", err.Error())
		log.Println(err)
	}
}

func EmailWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "emailQueue", redisPool)

	worker.Middleware((*jobs.Context).Log)
	worker.Middleware((*jobs.Context).FindUser)

	worker.Job("send_email", (*jobs.Context).SendEmail)

	worker.Start()
}

func StopEmailWorker() {
	worker := work.NewWorkerPool(jobs.Context{}, 10, "emailQueue", redisPool)
	worker.Stop()
}
