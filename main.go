package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-juno/kit/pool"
	"golang.org/x/xerrors"
)

const (
	FinishApi = "https://cat-match.easygame2021.com/sheep/v1/game/game_over?rank_score=1&rank_state=1&rank_time=%d&rank_role=1&skin=1"
)

const (
	HeaderT         = "eyJXX.XXX"   // header t
	HeaderUserAgent = "MozillaXXXX" // header UserAgent
	CostTime        = -1            // 花费的时间，-1 随机生成
	CycleCount      = 100           // 需要通关的次数，默认1
	Concurrent      = 10            // 同时闯关多少个
)

type Result struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Data    int    `json:"data"`
}

func FinishGame(costTime int) (err error) {
	url := fmt.Sprintf(FinishApi, costTime)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Host", "cat-match.easygame2021.com")
	req.Header.Set("User-Agent", HeaderUserAgent)
	req.Header.Set("t", HeaderT)
	client := &http.Client{}
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		return
	}
	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		err = xerrors.Errorf("%w", err)
		return
	}
	return
}

func Start(i int) {
	log.Printf("...第{%d}次开始闯关...", i+1)
	costTime := CostTime
	if CostTime == -1 {
		costTime = rand.Intn(3600)
	}
	err := FinishGame(costTime)
	if err != nil {
		log.Printf("...第{%d}次闯关失败... err:%+v", i+1, err)
	} else {
		log.Printf("...第{%d}次完成闯关...", i+1)

	}

}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("【羊了个羊一键闯关启动】")
	p := pool.New(Concurrent)
	for i := 0; i < CycleCount; i++ {
		p.Add(1)
		go func(i int) {
			Start(i)
			p.Done()
		}(i)
	}
	p.Wait()
	log.Println("【羊了个羊一键闯关结束】")

}
