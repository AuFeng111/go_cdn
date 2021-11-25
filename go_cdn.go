package main

import (
	"bufio"
	"flag"
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

var file string
var t int

func init() {
	flag.StringVar(&file, "file", "D:\\vscode\\python\\杂技\\webalive\\WebAliveScan-master\\target.txt", "")
	flag.IntVar(&t, "t", 100, "thread")
	flag.Usage = func() {
		fmt.Printf("\nUsage: \n-file 1.txt")
		flag.PrintDefaults() //输出flag
	}
	flag.Parse() //解析flag
}

func main() {
	var wg sync.WaitGroup
	var domains []string

	//file, err := os.Open("D:\\vscode\\python\\杂技\\webalive\\WebAliveScan-master\\target.txt")
	file, err := os.Open(file)
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
	fmt.Println("                      --create by aufeng")
	fmt.Println("开始探测: ")
	a := make(chan string, len(domains))
	for i := 0; i < t; i++ {
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
func get_hash(i string, wg *sync.WaitGroup) { //获取服务器的hash
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

func check_ping(i string, a []string) { //利用各地的服务器ip节点去进行ping，进行比较
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
		//c, err := regexp.MatchString("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}", string(body))
		if strings.Index(string(body), "ipAddress") > 0 { //如果服务器报错，就不会存在ipAddress的字段
			c := regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}").FindAllStringSubmatch(string(body), -1)
			ip = append(ip, c[0][0])
		}
	}
	b := 1
	if len(ip) > 0 { //ip长度大于0再进行比较
		//fmt.Print(len(ip), "   ", i)
		for a := 1; a < len(ip); a++ {
			if ip[0] != ip[a] {
				b = 0
				break
			}
		}
	}
	if b == 1 {
		fmt.Printf("[+]no cdn: %-30s %10s\n", i, ip[0])
	}
	// } else {
	// 	fmt.Printf("[-]have cdn: %-20s\n", i)
	// }
}
