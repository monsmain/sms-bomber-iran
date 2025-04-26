package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	// یک آدرس پروکسی نمونه (این آدرس واقعی نیست و فقط برای تست ساختار کد است)
	proxyAddr := "socks5://example.com:1080"

	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		fmt.Println("Error parsing proxy URL:", err)
		return
	}

	// تلاش برای ساخت Dialer از طریق پروکسی
	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		fmt.Println("Error creating proxy dialer:", err)
		return
	}

	// تلاش برای استفاده از DialContext (اینجا خطا رخ می دهد اگر DialContext نباشد)
    // ما در واقع قصد اتصال نداریم، فقط میخواهیم ببینیم متد وجود دارد یا نه
    fmt.Println("Attempting to use DialContext method...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // این خط همان جایی است که کامپایلر شما خطا می گیرد
	conn, err := dialer.DialContext(ctx, "tcp", "google.com:80")
	if err != nil {
        // اگر اینجا خطا گرفتیم، به احتمال زیاد پروکسی کار نمیکند یا اتصال موفق نبوده
        // اما هدف اصلی ما دیدن خطای کامپایل نیست، دیدن خطای زمان اجرا است
        fmt.Println("Error during DialContext (expected if proxy is not real):", err)
    } else {
        fmt.Println("DialContext successful (proxy might be working or test inconclusive).")
        conn.Close() // اتصال را ببندید
    }


	fmt.Println("Test finished.")
}
