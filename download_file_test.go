package service

import (
	"context"
	"fmt"
	"github.com/Ted-Young/rabida/config"
	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"log"
	"strings"
	"testing"
	"time"
)

func TestRabidaImpl_DownloadFile(t *testing.T) {
	conf := config.LoadFromEnv()
	fmt.Printf("%+v\n", conf)
	rabi := NewRabida(conf)
	job := Job{
		Link: "http://dnr.sc.gov.cn/scdnr/sctdlwj/2022/1/7/d555de82b08a42fbb7006f4cc567d296/files/8ab0e1fad55c459c9742448b1f975795.pdf",
	}
	var err error
	defer func() {
		if val := recover(); val != nil {
			var ok bool
			err, ok = val.(error)
			if !ok {
				err = errors.New(fmt.Sprint(val))
			} else {
				err = errors.Wrap(err, "recover from panic")
			}
		}
	}()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = rabi.DownloadFile(timeoutCtx, job, func(file string) {
		fmt.Println(file)
	}, conf)
	if err != nil {
		t.Log(err)
	}
}

func TestRabidaImpl_Download(t *testing.T) {
	var (
		allocCancel   context.CancelFunc
		contextCancel context.CancelFunc
		ctx           context.Context
	)
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
	}
	opts = append(opts, chromedp.Headless)
	ctx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	//ctx, contextCancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	ctx, contextCancel = chromedp.NewContext(ctx)
	defer contextCancel()

	// set up a channel so we can block later while we monitor the download progress
	downloadComplete := make(chan bool)

	// this will be used to capture the file name later
	var fileName string

	// set up a listener to watch the download events and close the channel when complete
	// this could be expanded to handle multiple downloads through creating a guid map,
	// monitor download urls via EventDownloadWillBegin, etc
	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*browser.EventDownloadWillBegin); ok {
			fileName = ev.SuggestedFilename
			return
		}
		if ev, ok := v.(*browser.EventDownloadProgress); ok {
			fmt.Printf("current download state: %s\n", ev.State.String())
			if ev.State == browser.DownloadProgressStateCompleted {
				close(downloadComplete)
			}
		}
	})

	var err error
	defer func() {
		if val := recover(); val != nil {
			var ok bool
			err, ok = val.(error)
			if !ok {
				err = errors.New(fmt.Sprint(val))
			} else {
				err = errors.Wrap(err, "recover from panic")
			}
		}
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err = chromedp.Run(timeoutCtx,
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath("out").
			WithEventsEnabled(true),
		chromedp.Navigate("http://dnr.sc.gov.cn/scdnr/sctdlwj/2022/1/7/d555de82b08a42fbb7006f4cc567d296/files/8ab0e1fad55c459c9742448b1f975795.pdf"),
	); err != nil && !strings.Contains(err.Error(), "net::ERR_ABORTED") {
		// Note: Ignoring the net::ERR_ABORTED page error is essential here since downloads
		// will cause this error to be emitted, although the download will still succeed.
		t.Log(err)
	}

	// This will block until the chromedp listener closes the channel
	<-downloadComplete

	log.Printf("Download Complete: %v/%v", "out", fileName)
}
