package main

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_main(t *testing.T) {

	// 启动一个流服务器，用于统计熔断降级结果，实时发送到服务器上。
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8074", hystrixStreamHandler)

	// 熔断器配置
	hystrix.ConfigureCommand("aaa", hystrix.CommandConfig{
		Timeout:                1000, // 单次请求 超时时间
		MaxConcurrentRequests:  1,    // 最大并发量
		SleepWindow:            5000, // 熔断后多久去尝试服务是否可用
		RequestVolumeThreshold: 1,    // 验证熔断的 请求数量, 10秒内采样
		ErrorPercentThreshold:  1,    // 验证熔断的 错误百分比
	})

	for i := 0; i < 10000; i++ {
		//异步调用使用 hystrix.Go
		err := hystrix.Do("aaa", func() error {
			//test case 1 并发测试
			if i == 0 {
				return errors.New("service error")
			}
			//test case 2 超时测试
			//time.Sleep(2 * time.Second)
			log.Println("do services")
			return nil
		}, func(err error) error {
			// 当第一个函数出错的时候，会触发这个函数，做处理，如果这个函数也出错了。或者直接return err
			// 如果这个错误回调函数也出错了。或者直接return err 那就就是直接给到外面的 err:=hystrix.Do
			log.Println("hystrix err:" + err.Error())
			time.Sleep(1 * time.Second)
			log.Println("sleep 1 second")
			//return err
			return nil
		})
		if err != nil {
			log.Println("hystrix err:" + err.Error())
			time.Sleep(1 * time.Second)
			log.Println("sleep 1 second")
		}
	}
	time.Sleep(100 * time.Second)
}
