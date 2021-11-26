# go_cdn
#### 批量域名筛查CDN(最准确的CDN判断脚本)
- go build request1.go 编译对应版本就好
- -h查看使用方法,主要是-file和-t参数
- 主要方便与批量域名中筛出无cdn的网站
- 爬取的https://wepcc.com/ 上服务器的节点，进行多个ip的比较，如果ip相等证明无cdn
- 代码可能存在bug，文本中的domain尽量不要有多余的换行。后期会继续修改
- 只输出无cdn的域名ip

![Image text](https://github.com/AuFeng111/go_cdn/blob/main/cdn2.png)
