package service

import (
	"context"
	"fmt"
	"github.com/Ted-Young/rabida/config"
	"github.com/Ted-Young/rabida/lib"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"log"
	"testing"
	"time"
)

func TestRabidaImpl_CrawlWithListeners(t *testing.T) {
	conf := config.LoadFromEnv()
	fmt.Printf("%+v\n", conf)
	rabi := NewRabida(conf)
	job := Job{
		Link: "http://mfb.sh.gov.cn/zwgk/jcgk/zcfg/gfxwj/index.html",
		CssSelector: CssSelector{
			Scope: `#Datatable-1>tbody>tr`,
			Attrs: map[string]CssSelector{
				"title": {
					Css: "td:first-child",
				},
				"link": {
					Css:  "td:first-child",
					Attr: "node",
				},
				"date": {
					Css: "td:last-child",
				},
			},
		},
		Paginator: CssSelector{
			Css: "div[name='whj_nextPage']:not(.whj_hoverDisable)",
		},
		Limit: 1,
	}

	linkCh := make(chan string, 1)
	err := rabi.CrawlWithListeners(context.Background(), job, func(ctx context.Context, ret []interface{}, nextPageUrl string, currentPageNo int) bool {
		for _, item := range ret {
			value, ok := item.(map[string]interface{})
			if !ok {
				panic(errors.New("cast failed"))
			}
			node := value["link"].(*cdp.Node)
			timeoutCtx, jsClickCancel := context.WithTimeout(ctx, 10*time.Second)
			err := chromedp.Run(timeoutCtx, lib.JsClickNode(node))
			if err != nil {
				jsClickCancel()
				panic(err)
			}
			jsClickCancel()
			link := <-linkCh
			log.Println(link)
		}
		if currentPageNo >= job.Limit {
			return true
		}
		return false
	}, nil, []chromedp.Action{
		chromedp.EmulateViewport(1777, 903, chromedp.EmulateLandscape),
	}, nil,
		[]chromedp.ExecAllocatorOption{chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")},
		func(v interface{}) {
			if ev, ok := v.(*page.EventWindowOpen); ok {
				linkCh <- ev.URL
			}
		},
	)
	if err != nil {
		t.Error(fmt.Sprintf("%+v", err))
	}
}

func TestRabidaImpl_CrawlWithListeners2(t *testing.T) {
	conf := config.LoadFromEnv()
	fmt.Printf("%+v\n", conf)
	rabi := NewRabida(conf)
	job := Job{
		Link: "http://minzheng.hebei.gov.cn/policyDatabase",
		CssSelector: CssSelector{
			Scope: `div.el-table__body-wrapper.is-scrolling-none > table > tbody > tr`,
			Attrs: map[string]CssSelector{
				"title": {
					Css: "td:first-child a",
				},
				"link": {
					Css:  "td:first-child a",
					Attr: "node",
				},
				"date": {
					Css: "td:last-child>div>div",
				},
			},
		},
		Paginator: CssSelector{
			Css: "button.btn-next>span",
		},
		Limit: 3,
	}

	linkCh := make(chan string, 1)
	err := rabi.CrawlWithListeners(context.Background(), job, func(ctx context.Context, ret []interface{}, nextPageUrl string, currentPageNo int) bool {
		for _, item := range ret {
			value, ok := item.(map[string]interface{})
			if !ok {
				panic(errors.New("cast failed"))
			}
			node := value["link"].(*cdp.Node)
			timeoutCtx, jsClickCancel := context.WithTimeout(ctx, 10*time.Second)
			err := chromedp.Run(timeoutCtx, lib.JsClickNode(node))
			if err != nil {
				jsClickCancel()
				panic(err)
			}
			jsClickCancel()
			link := <-linkCh
			log.Println(link)
		}
		if currentPageNo >= job.Limit {
			return true
		}
		return false
	}, nil, []chromedp.Action{
		chromedp.EmulateViewport(1777, 903, chromedp.EmulateLandscape),
	}, nil,
		[]chromedp.ExecAllocatorOption{chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")},
		func(v interface{}) {
			if ev, ok := v.(*page.EventWindowOpen); ok {
				linkCh <- ev.URL
			}
		},
	)
	if err != nil {
		t.Error(fmt.Sprintf("%+v", err))
	}
}

func TestRabidaImpl_WaitNewTarget(t *testing.T) {
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

	ctx, contextCancel = chromedp.NewContext(ctx)
	defer contextCancel()

	downloadBegin := make(chan bool)
	var downloadUrl string

	chromedp.ListenTarget(ctx, func(v interface{}) {
		if ev, ok := v.(*page.EventWindowOpen); ok {
			downloadUrl = ev.URL
			close(downloadBegin)
		}
	})

	var res []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate("http://minzheng.hebei.gov.cn/policyDatabase"),
		chromedp.EvaluateAsDevTools("document.querySelector('.cell>a').click()", &res),
	); err != nil {
		t.Error(fmt.Sprintf("%+v", err))
	}

	// This will block until the chromedp listener closes the channel
	<-downloadBegin

	// We can predict the exact file location and name here because of how we configured
	// SetDownloadBehavior and WithDownloadPath
	log.Printf("Download url: %s", downloadUrl)
}
