package main

import (
	"bufio"
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var replayData string
var replyString string
var stap = 1
var waitTime = 0

func main() {
	defer func() {
		// recover() 可以将捕获到的panic信息打印
		if err := recover(); err != nil {
			log.Println(err)
			// 这里不能用 continue  所以直接 返回，这样程序就不会因为有错误而退出  recover 可以捕获异常数据
			log.Println("我捕获到了错误,程序重新运行.....")
			main()
		}

	}()
	if !isLogin() {
		log.Println("你的cookie无效，请从新输入！！~~")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	for {

		list := getList()
		//log.Println(list)
		if list == nil {
			log.Println("没有获取到连接...... 即将等待40秒..")
			time.Sleep(40 * time.Second)
			continue
		}
		//创建一个新文件，写入内容 5 句 “http://c.biancheng.net/golang/”

		file, err := os.OpenFile("./replyed.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("文件打开失败", err)
		}

		//写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)
		//及时关闭file句柄
		file.Close()
		file1, err := os.Open("./replyed.txt")
		if err != nil {
			return
		}
		allReply, err := ioutil.ReadAll(file1)
		if err != nil {
			return
		}
		s := string(allReply)
		file1.Close()
		for i := 0; i < 25; i++ {

			if strings.Contains(s, list[i]) {
				continue
			}
			if i == 24 {
				log.Println("已重试25次了。。。跳过请等待...")
			}
			write.WriteString(list[i] + " \n")
			//Flush将缓存的文件真正写入到文件中
			write.Flush()
			log.Println("正在浏览与回复" + list[i])
			go newTask("https://t66y.com/" + list[i])
			time.Sleep(50 * time.Second)
			// 退出所有携程
			// 超时抛出 context deadline exceeded
			//  这会调用 os.exit 1  我人麻了 不会给我捕获错误的机会
			// ackCtx, ackCancel := context.WithTimeout(ctx, 50*time.Second)
			// defer ackCancel()
			// 所以改用这个方法
			go runtime.Goexit()
			stap++
			break
		}
		if stap == 15 {
			log.Println("已经回复了10次了。坐等明天继续")
			// 等待 21个小时 1024/60*10/60=2.84444444444 大概是3小时 加上随机等待时间已经足够3小时了
			alltime := 1024*15 + waitTime // 总共等待时间
			// 一天是 多少秒 86400 ？
			alltime = 86400 - alltime
			// 等待 剩下的时间 通过管道  光写 sleep 感觉不够高级
			<-time.After(time.Duration(alltime) * time.Second)
			waitTime = rand.Intn(7200)
			// 不要每天同一个时间回复帖子。 // 在两个小时内随机一下 比较可靠
			<-time.After(time.Duration(waitTime) * time.Second)
			stap = 1
		}
		log.Println("回复帖子返回的值： " + replayData)
		log.Println("正在等待1024秒..")
		time.Sleep(1024 * time.Second)

		rand.Seed(time.Now().Unix())
		t1 := rand.Intn(600)
		waitTime = waitTime + t1
		log.Println("等待1024秒后正在随机等待" + strconv.Itoa(t1) + "秒")
		time.Sleep(time.Duration(t1) * time.Second)

	}

}

func getList() []string {

	url := "https://t66y.com/thread0806.php?fid=7&search=today"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	//log.Println(string(body))
	// https://t66y.com/htm_data/2206/7/5121232.html
	//解析正则表达式，如果成功返回解释器
	reg1 := regexp.MustCompile(`htm_data/\d*/\d/\d*\.html`)
	if reg1 == nil {
		log.Println("regexp err")
		return nil
	}
	//根据规则提取关键信息
	result1 := reg1.FindAllString(string(body), -1)
	//log.Println("result1 = ", result1)
	return result1
}

func newTask(newurl string) {

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),                        // 控制是否显示浏览界面
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //是否加载图片
		chromedp.Flag("disable-gpu", true),                     //禁用显卡 （服务器端）
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create context
	// 设置50秒超时

	ackCtx, allCancel := chromedp.NewContext(allocCtx)
	defer allCancel()
	// 超时抛出 context deadline exceeded
	//  这会调用 os.exit 1  我人麻了 不会给我捕获错误的机会
	//ackCtx, ackCancel := context.WithTimeout(ctx, 50*time.Second)
	//defer ackCancel()

	//cookieFile, err2 := os.Open("G:\\GitHub\\GitHub\\t66y_auto_Submit\\t66ycookie")
	cookieFile, err2 := os.OpenFile("./t66ycookie", os.O_RDWR|os.O_CREATE, 0666)
	if err2 != nil {
		log.Fatal("获取cookie出错")
	}
	defer cookieFile.Close()
	allcookie, err := ioutil.ReadAll(cookieFile)
	if err != nil {
		log.Fatal("读取cookie出错")
	}
	if string(allcookie) == "" {
		log.Println("在t66ycookie中填入cookie  例如")
		log.Println("227c9_winduser=xxxxxxxx")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	cookies := strings.Split(string(allcookie), "=")

	rand.Seed(time.Now().Unix())
	reply := [...]string{"感谢分享~", "1024 感谢", "感谢分享", "感谢分享!", "楼主辛苦", "很好的分享", "多谢分享", "1024感谢分享"}
	// 保证每次回复不重复
	for {
		_s1 := reply[rand.Intn(len(reply))]
		if _s1 != replyString {
			replyString = _s1
			break
		}
	}
	log.Println("回复内容是 " + replyString)

	err = chromedp.Run(ackCtx,
		chromedp.Tasks{
			network.SetCookie(cookies[0], cookies[1]).WithDomain("t66y.com"),

			chromedp.Navigate(newurl),

			chromedp.WaitVisible("body", chromedp.ByQuery), // 等待 body 加载完毕

			chromedp.Sleep(5 * time.Second),

			chromedp.Click("t_like"),

			chromedp.SendKeys(`document.querySelector("#main > form > div > table > tbody > tr:nth-child(2) > td:nth-child(2) > textarea")`, replyString, chromedp.ByJSPath),

			chromedp.Sleep(5 * time.Second),

			chromedp.Submit("Submit"),

			chromedp.WaitVisible("body", chromedp.ByQuery), // 等待 body 加载完毕
			chromedp.Sleep(5 * time.Second),
			chromedp.Text("main", &replayData, chromedp.ByID),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(newurl + "    已回复结束.. 即将等待....")
}

func isLogin() bool {
	cookieFile, err2 := os.OpenFile("./t66ycookie", os.O_RDWR|os.O_CREATE, 0666)
	if err2 != nil {
		log.Fatal("获取cookie出错")
	}
	defer cookieFile.Close()
	allcookie, err := ioutil.ReadAll(cookieFile)
	if err != nil {
		log.Fatal("读取cookie出错")
	}
	if string(allcookie) == "" {
		log.Println("在t66ycookie中填入cookie  例如")
		log.Println("227c9_winduser=xxxxxxxx")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	url := "https://t66y.com/index.php"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Add("authority", "t66y.com")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("cookie", string(allcookie))
	req.Header.Add("dnt", "1")
	req.Header.Add("referer", "https://t66y.com/")
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	//log.Println(string(body))
	if strings.Contains(string(body), "上次登錄時間") {
		log.Println("你已登陆.... cookie 正确")
		return true
	}
	return false

}
