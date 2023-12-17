package main

import (
	"fmt"
	"log"
	"time"
	"context"
	"sync"
	"strings"
	"regexp"
	"net/http"
	"io/ioutil"
	"encoding/json"
	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/go-redis/redis/v8"
)

type SafeMap struct {
	mu    sync.Mutex
	items map[string]string
}

type Subscription struct {				   
	subscription string `json:"subscription"`
}

const (
	vapidPublicKey  = "발급한 vapidPublicKey  값으로 넣어주세요."
  vapidPrivateKey = "발급한 vapidPrivateKey 값으로 넣어주세요."
)



var ctx = context.Background()


func (m *SafeMap) Set(key string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

func (m *SafeMap) Get(key string) (string, bool) {
        m.mu.Lock()
        defer m.mu.Unlock()
        value, exists := m.items[key]
        return value, exists
}



func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379",})

	for {
		var cursor uint64
		safeMap := SafeMap{items: make(map[string]string)}
		for {
			var keys []string
			var err error
			keys, cursor, err = rdb.Scan(ctx, cursor, "", 10).Result()
			if err != nil {
				panic(err)
			}
	
			for _, key := range keys {
				parts := strings.Split(key, "_$$_")

				if err != nil {
					log.Fatal(err)
				}

				val, err := rdb.Get(ctx, key).Result()
				if err == nil {
					go webPush(key, rdb, &safeMap, parts[0], val, parts[1]);
				}

			}	

			if cursor == 0 {
				break
			}
		}
		time.Sleep(30 * time.Second)
	}	

	/*
	// Decode subscription
	s := &webpush.Subscription{}
	json.Unmarshal([]byte(subscription), s)

	// Send Notification
	resp, err := webpush.SendNotification([]byte("Test"), s, &webpush.Options{
		Subscriber:      "example@example.com", // Do not include "mailto:"
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	*/
}

func webPush(key string,rdb *redis.Client, m *SafeMap, url string,  keyword string, subscription string) {
	
	fmt.Println(url)
	fmt.Println(keyword)
	fmt.Println(subscription)

	var document string
	var exists bool
	if document, exists =  m.Get(url); exists {
		fmt.Println("exist!!!");
	} else {
		fmt.Println("no!!")		

		c := http.Client{}		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

		resp, err := c.Do(req)
		if err != nil {
	        	fmt.Println(err)
			return
	 	}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		document = string(body)		
		m.Set(url, document)
	}

	pattern  := keyword
	re, err  := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("정규식 컴파일 오류", err)
		return
	}

	matches := re.FindAllString(document, -1)
	if matches != nil {
		fmt.Println("Find!");
	        s := &webpush.Subscription{}
	        json.Unmarshal([]byte(subscription), s)

		_, err = webpush.SendNotification([]byte(keyword), s, &webpush.Options{
        	        Subscriber:      "example@example.com", // Do not include "mailto:"
        	        VAPIDPublicKey:  vapidPublicKey,
        	        VAPIDPrivateKey: vapidPrivateKey,
        	        TTL:             30,
        	})
	 	rdb.Del(ctx,key) 
        	if err != nil {
		        fmt.Println(err)
		}

	}

}
