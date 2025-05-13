package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// 创建上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second) // 将超时时间增加到 120 秒
	defer cancel()

	// 定义变量来存储渲染后的HTML和图片URL
	var htmlContent string
	var imgURLs []string

	// 执行任务
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://pixabay.com/zh/images/search/%E5%A5%B3%E7%94%9F%E5%8A%A8%E6%BC%AB%E5%A4%B4%E5%83%8F/"),
		chromedp.WaitVisible("body"),             // 等待页面加载完成
		chromedp.OuterHTML("html", &htmlContent), // 获取渲染后的HTML
		chromedp.Evaluate(`Array.from(document.querySelectorAll('img')).map(img => img.src)`, &imgURLs), // 获取所有图片URL
	)
	if err != nil {
		log.Fatal(err)
	}

	// 输出渲染后的HTML
	fmt.Println("Rendered HTML:", htmlContent)

	// 下载图片
	for _, url := range imgURLs {
		if url != "" {
			fmt.Println("Downloading image:", url)
			downloadImage(url)
		}
	}
}

func downloadImage(url string) {
	// 创建上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second) // 将超时时间增加到 60 秒
	defer cancel()

	// 定义变量来存储图片数据
	var imgData []byte

	// 执行任务
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("img"),          // 等待图片加载完成
		chromedp.CaptureScreenshot(&imgData), // 捕获图片数据
	)
	if err != nil {
		log.Println("Error downloading image:", err)
		return
	}

	// 创建资源文件夹
	resourceDir := "resources"
	if _, err := os.Stat(resourceDir); os.IsNotExist(err) {
		err := os.Mkdir(resourceDir, 0755)
		if err != nil {
			log.Println("Error creating resource directory:", err)
			return
		}
	}

	// 保存图片
	fileName := filepath.Join(resourceDir, filepath.Base(url))
	err = os.WriteFile(fileName, imgData, 0644)
	if err != nil {
		log.Println("Error saving image:", err)
	} else {
		fmt.Println("Image saved:", fileName)
	}
}
