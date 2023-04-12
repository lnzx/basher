package main

import (
	"crypto/tls"
	"flag"
	"github.com/lnzx/faker/ip"
	"github.com/lnzx/faker/useragent"
	"github.com/pterm/pterm"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var (
	concurrency int
	url         string
	fake        bool
	debug       bool
)

var ua = useragent.New()
var client *http.Client

func init() {
	flag.IntVar(&concurrency, "c", 8, "concurrent")
	flag.StringVar(&url, "s", "https://proof.ovh.net/files/1Mb.dat", "target url")
	flag.BoolVar(&fake, "fake", true, "Random X-Forwarded-For and X-Real-IP")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()

	if concurrency < 1 {
		concurrency = runtime.NumCPU()
	}

	transport := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   1000,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	client = &http.Client{Transport: transport}
}

func main() {
	tree := pterm.TreeNode{Text: "Basher"}
	var children []pterm.TreeNode

	for i := 0; i < concurrency; i++ {
		children = append(children, pterm.TreeNode{
			Text: "worker-" + strconv.Itoa(i) + " running",
		})
		go work(i)
	}
	tree.Children = children
	pterm.DefaultTree.WithRoot(tree).Render()
	select {}
}

func work(i int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("work %d has error: %v\n", i, r)
			work(i)
		}
	}()

	task := 0
	for {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Printf("work %d NewRequest error: %v\n", i, err)
			continue
		}
		req.Header.Set("User-Agent", ua.Random())
		if fake {
			fakeIp := ip.IPv4()
			req.Header.Add("X-Forwarded-For", fakeIp)
			req.Header.Add("X-Real-IP", fakeIp)
		}

		rsp, err := client.Do(req)
		if err != nil {
			log.Printf("work %d client.Do error: %v\n", i, err)
			if rsp != nil { // 重定向错误时, response和err变量都会变为非nil值
				rsp.Body.Close()
			}
			continue
		}
		_, err = io.Copy(io.Discard, rsp.Body)
		if err != nil && err != io.EOF {
			log.Printf("work %d io.Copy error: %v\n", i, err)
		}
		rsp.Body.Close()
		if debug {
			task++
			log.Printf("task %d-%d ok\n", i, task)
		}
	}
}
