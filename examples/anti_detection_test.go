package examples

import (
	"context"
	"github.com/Ted-Young/rabida/lib"
	"github.com/Ted-Young/rabida/useragent"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"testing"
)

func TestAntiDetection(t *testing.T) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.UserAgent(useragent.RandomMacChromeUA()))
	ctx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	var buf []byte
	var tasks chromedp.Tasks
	tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		_, err = page.AddScriptToEvaluateOnNewDocument(lib.Script).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}))
	tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		_, err = page.AddScriptToEvaluateOnNewDocument(lib.AntiDetectionJS).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}))

	tasks = append(tasks,
		chromedp.Navigate("https://bot.sannysoft.com/"),
		chromedp.FullScreenshot(&buf, 100))

	if err := chromedp.Run(ctx, tasks); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
