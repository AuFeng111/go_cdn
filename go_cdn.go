package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var domains []string
	file, err := os.Open("D:\\vscode\\python\\杂技\\webalive\\WebAliveScan-master\\target.txt")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		domains = append(domains, strings.TrimSpace(string(line)))
	}

	fmt.Println("已经获取以下域名: ")
	fmt.Println(domains)
	fmt.Println("开始探测: ")
	fmt.Println("                      --create by aufeng")
	a := make(chan string, len(domains))
	for i := 0; i < 50; i++ {
		go func() {
			for i := range a {
				get_hash(i, &wg)
				//wg.Done()
			}
		}()
	}
	for _, domain := range domains {
		wg.Add(1)
		a <- domain
	}
	wg.Wait()
	close(a)
}
func get_hash(i string, wg *sync.WaitGroup) {
	defer wg.Done()
	//i := "ffbebbs.kingsoft.com"
	resp, err := http.PostForm("https://www.wepcc.com/", url.Values{"host": {i}, "node": {"2,3,6"}})
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	//fmt.Print(string(body))
	reg := regexp.MustCompile(`<tr class="item" data-id="(?s:(.*?))">`)
	if reg == nil {
		fmt.Println("MustCompile err")
		return
	}
	//提取关键信息
	result := reg.FindAllStringSubmatch(string(body), -1)
	//fmt.Println(result)
	//过滤<></>
	var a []string
	for _, text := range result {
		a = append(a, text[1])
		//fmt.Println(text[1])
	}
	//fmt.Println(a)
	check_ping(i, a)
}

func check_ping(i string, a []string) {
	var ip []string
	for _, b := range a {
		resp, err := http.PostForm("https://www.wepcc.com/check-ping.html", url.Values{"host": {i}, "node": {b}})
		if err != nil {
			fmt.Printf("get failed, err:%v\n", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("read from resp.Body failed, err:%v\n", err)
			return
		}
		//var data[] string
		//fmt.Println(string(body))
		//matched, err := regexp.MatchString("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}", string(body))
		c := regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}").FindAllStringSubmatch(string(body), -1)
		ip = append(ip, c[0][0])
		//fmt.Println("判断ip中: ", i, " ", c[0][0])
	}
	if ip[0] == ip[1] && ip[0] == ip[2] && ip[0] == ip[3] && ip[0] == ip[4] {
		fmt.Printf("[+]no cdn: %-30s %10s\n", i, ip[0])
	}
	// } else {
	// 	fmt.Printf("[-]have cdn: %-20s\n", i)
	// }
}
