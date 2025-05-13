package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
	// 移除异步设置
	// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
		RandomDelay: 5 * time.Second,
	})

	// 设置代理
	c.SetProxy("http://proxy1:port")
	c.SetProxy("http://proxy2:port")

	// 添加请求头
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
	})

	// 回调函数，当爬虫找到一个 <title> 元素时调用
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Page Title:", e.Text)
	})
	// 提取所有链接
	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("Link:", link)
	})

	c.OnHTML("img", func(e *colly.HTMLElement) {
	    imgURL := e.Attr("src")
	    if imgURL == "" {
	        // 处理懒加载图片
	        imgURL = e.Attr("data-src")
	    }
	    if imgURL == "" {
	        // 处理其他可能的动态加载图片
	        imgURL = e.Attr("data-lazy")
	    }
	    if imgURL == "" {
	        // 处理JavaScript动态加载的图片
	        imgURL = e.Attr("data-original")
	    }
	    if imgURL != "" {
	        fmt.Println("Found image:", imgURL)
	        c.Visit(imgURL)
	    }
	})

	c.OnHTML("script", func(e *colly.HTMLElement) {
	    // 解析script内容，提取图片URL
	    // 这里可以使用正则表达式或其他方法提取图片URL
	    // 例如：regexp.MustCompile(`https?://[^\s]+\.(jpg|png|gif)`)
	    // 然后使用c.Visit访问提取到的图片URL
	})

	// 分页爬取
	c.OnHTML(".next a", func(e *colly.HTMLElement) {
		nextUrl := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextUrl) // 自动爬取下一页[[6]]
	})

	c.OnResponse(func(r *colly.Response) {
	    if strings.Contains(r.Headers.Get("Content-Type"), "application/json") {
	        // 解析JSON响应，提取图片URL
	        // 然后使用c.Visit访问提取到的图片URL
	    }
	})

	// 处理错误
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	// 访问目标页面
	c.Visit("https://fe-1.test.funnymamu.com/yari-fe-app-feature-2024-test/cp/index.html")

	// 如果必须使用异步模式，可以添加等待
	// c.Wait()
}
