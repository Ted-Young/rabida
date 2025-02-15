<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>

### Rabida

Rabida 是一个基于 [chromedp](https://github.com/chromedp/chromedp) 简单易用的爬虫框架。

### 目前支持的特性

- `分页`:  使用css 选择器获取下页数据。
- `预分页`: 在获取分页数据前做一些操作，比如点击日期按钮获取最新数据。
- `浏览器cookie`: 启用浏览器cookie，针对于某些网站需要登录状态。
- `延迟跟超时`:  能够自定义延迟跟超时时间。
- `反爬虫检测`:
  每个任务默认加载了反爬虫检测脚本，脚本来源于[puppeteer-extra-stealth](https://github.com/berstend/puppeteer-extra/tree/master/packages/extract-stealth-evasions#readme)。
- `严格模式`: useragent、浏览器、浏览器的平台必须匹配，如果设置成true，将设置为chrome-mac相关的useragent、chrome浏览器、浏览器平台为Mac。针对于某些网站的反爬机制。
- `Xpath表达式`: 使用xpath表达式获取元素
- `Iframe`: 指定iframe选择器，获取页面某个iframe作为父级元素
- `Scroll`: 滚动当前页面，ScrollType类型为scrollBy和scrollTo，默认是scrollBy，跟window.scrollBy, window.scrollTo表现行为一样.

### 安装

```go
go get -u github.com/Ted-Young/rabida
```

### 配置

添加.env 文件到你的项目

```shell
RABI_DELAY=1s,2s
RABI_CONCURRENCY=1
RABI_THROTTLE_NUM=2
RABI_THROTTLE_DURATION=1s
RABI_TIMEOUT=3s
RABI_MODE=headless
RABI_DEBUG=false
RABI_OUT=out
RABI_STRICT=false
RABI_PROXY=
```

### 用法

这里看更多的例子 [examples](https://github.com/Ted-Young/rabida/blob/master/examples)

Css选择器使用：

```go
func TestRabidaImplCrawl(t *testing.T) {
    conf := config.LoadFromEnv()
    fmt.Printf("%+v\n", conf)
    rabi := NewRabida(conf)
    job := Job{
        Link: "https://tieba.baidu.com/f?kw=nba",
        CssSelector: CssSelector{
            Scope: `#thread_list > li.j_thread_list`,
            Attrs: map[string]CssSelector{
                "title": {
                    Css: "div.threadlist_title > a",
                },
                "date": {
                    Css: "span.threadlist_reply_date",
                },
            },
        },
        Paginator: CssSelector{
            Css: "#frs_list_pager > a.next.pagination-item",
        },
        Limit: 3,
    }
    err := rabi.Crawl(context.Background(), job, func(ret []interface{}, nextPageUrl string, currentPageNo int) bool {
        for _, item := range ret {
            fmt.Println(gabs.Wrap(item).StringIndent("", "  "))
        }
        if currentPageNo >= job.Limit {
            return true
        }
        return false
    }, nil, []chromedp.Action{
        chromedp.EmulateViewport(1777, 903, chromedp.EmulateLandscape),
    })
    if err != nil {
        panic(fmt.Sprintf("%+v", err))
    }
}
```

Xpath表达式：

```go
func TestRabidaXpathImpl_Crawl(t *testing.T) {
    conf := config.LoadFromEnv()
    fmt.Printf("%+v\n", conf)

    rabi := NewRabida(conf)
    job := Job{
        Link: "https://you.ctrip.com/sight/shenzhen26/2778.html",
        CssSelector: CssSelector{
            XpathScope: `//*[@id="commentModule"]/div[@class='commentList']/div`,
            Attrs: map[string]CssSelector{
                "content": {
                    Xpath: "//div[@class='commentDetail']",
                },
                "date": {
                    Xpath: `//div[@class='commentTime']`,
                },
            },
        },
        Paginator: CssSelector{
            Xpath: "//*[@id='commentModule']//li[@class=' ant-pagination-next' and not(@aria-disabled='true')]",
        },
        Limit: 3,
    }
    err := rabi.Crawl(context.Background(), job, func(ret []interface{}, nextPageUrl string, currentPageNo int) bool {
        for _, item := range ret {
            fmt.Println(gabs.Wrap(item).StringIndent("", "  "))
        }
        logrus.Printf("currentPageNo: %d\n", currentPageNo)
        if currentPageNo >= job.Limit {
            return true
        }
        return false
    }, nil, []chromedp.Action{
        chromedp.EmulateViewport(1777, 903, chromedp.EmulateLandscape),
    })
    if err != nil {
        t.Error(fmt.Sprintf("%+v", err))
    }
}
```

Scorll用法:

```go
func TestRabidaImplCrawlScrollSmooth(t *testing.T) {
    t.Run("CrawlScrollSmooth", func(t *testing.T) {
        conf := config.LoadFromEnv()
        fmt.Printf("%+v\n", conf)
        rabi := NewRabida(conf)
        job := Job{
            Link: "https://twitter.com/NASA",
            CssSelector: CssSelector{
                Scope: `div[data-testid='cellInnerDiv'] article[data-testid='tweet']`,
                Attrs: map[string]CssSelector{
                    "title": {
                        Css: `div[data-testid="tweetText"]`,
                    },
                    "date": {
                        Css:  `a > time`,
                        Attr: `datetime`,
                    },
                    "link": {
                        Css:  `a[role="link"][href*=status]`,
                        Attr: `href`,
                    },
                    "reply": {
                        Css:  `div[data-testid="reply"]`,
                        Attr: `aria-label`,
                    },
                    "retweet": {
                        Css:  `div[data-testid="retweet"]`,
                        Attr: `aria-label`,
                    },
                    "like": {
                        Css:  `div[data-testid="like"]`,
                        Attr: `aria-label`,
                    },
                },
            },
            Limit: 5,
        }
        err := rabi.CrawlScrollSmooth(context.Background(), job, func(ret []interface{}, currentPageNo int) bool {
            for _, item := range ret {
                fmt.Println(gabs.Wrap(item).StringIndent("", "  "))
            }
            if currentPageNo >= job.Limit {
                return true
            }
            return false
        }, nil, nil)
        if err != nil {
            t.Errorf("%+v", err)
        }

    })
}
```