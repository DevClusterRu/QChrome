package main

import (
	"QChromium/internal"
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

var actions []chromedp.Action

func TestDP(t *testing.T) {

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
		chromedp.WindowSize(1200, 711), // init with a mobile view
	)

	browserCtx, browserCancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer browserCancel()

	ctx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	var nodes []*cdp.Node
	var res []byte

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	actions = append(actions, chromedp.Navigate(`https://emex.ru/products/A0024203483/Daimler%20AG/27199`))
	actions = append(actions, chromedp.WaitVisible(`.e-inputField`, chromedp.ByQuery))
	actions = append(actions, chromedp.SendKeys(`.e-inputField`, "A0024203483", chromedp.ByQuery))
	actions = append(actions, chromedp.SendKeys(`.e-inputField`, kb.Enter, chromedp.ByQuery))
	actions = append(actions, chromedp.Nodes(`[data-testid="Offers:block:tableoriginals"]`, &nodes, chromedp.ByQueryAll))

	actions = append(actions, chromedp.ActionFunc(func(c context.Context) error {
		for _, node := range nodes {
			dom.RequestChildNodes(node.NodeID).WithDepth(100).Do(c)
		}
		for _, node := range nodes {
			fmt.Println(node.Children[0].Children[0].Children[1].Children[0].Children[0].Children[0].NodeValue)
		}
		return nil
	}))

	//actions = append(actions, chromedp.ActionFunc(func(c context.Context) error {
	//	return dom.RequestChildNodes(nodes[0].NodeID).WithDepth(-1).Do(c)
	//}))

	//actions = append(actions, chromedp.Sleep(1*time.Second))

	//actions = append(actions, chromedp.ActionFunc(func(c context.Context) error {
	//	for _, node := range nodes {
	//		fmt.Println(node.Children[0].NodeValue)
	//	}
	//	return nil
	//}))
	//
	////actions = append(actions, chromedp.WaitVisible(`.smibiyl`, chromedp.ByQuery))
	////actions = append(actions, chromedp.Nodes(`h3`, &nodes, chromedp.ByQueryAll))
	//actions = append(actions, chromedp.FullScreenshot(&res, 90))

	//actions = append(actions, chromedp.Text(`input[class="e-inputField"]`, &example))

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("fullScreenshot.png", res, 0777); err != nil {
		log.Fatal(err)
	}

	fmt.Print(len(nodes))

	//for _, v := range nodes {
	//	fmt.Println(v.Children[0].NodeValue)
	//}

}

func TestDPStages(t *testing.T) {

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
		chromedp.WindowSize(1200, 711), // init with a mobile view
	)

	browserCtx, browserCancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer browserCancel()

	ctx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	var res []byte

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	//actions = append(actions, )
	//actions = append(actions, chromedp.WaitVisible(`.e-inputField`, chromedp.ByQuery))
	//actions = append(actions, chromedp.SendKeys(`.e-inputField`, "A0024203483", chromedp.ByQuery))
	//actions = append(actions, chromedp.SendKeys(`.e-inputField`, kb.Enter, chromedp.ByQuery))
	//actions = append(actions, chromedp.Nodes(`[data-testid="Offers:block:tableoriginals"]`, &nodes, chromedp.ByQueryAll))

	err := chromedp.Run(ctx, chromedp.Navigate(`https://emex.ru`))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx, chromedp.Sleep(1*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx, chromedp.FullScreenshot(&res, 90))
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("ss.png", res, 0777); err != nil {
		log.Fatal(err)
	}

}

func TestPipeline(t *testing.T) {

	chain := internal.Request{}
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "URL", Params: []string{"https://freelance.ru/"}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "WAIT", Params: []string{`a[href="/login/"]`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "CLICK", Params: []string{`a[href="/login/"]`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "WAIT", Params: []string{`[name="login"]`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "KEY", Params: []string{`[name="login"]`, `lifeon`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "KEY", Params: []string{`[name="passwd"]`, `Vivasan123`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "CLICK", Params: []string{`[name="submit"]`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "WAIT", Params: []string{`.navbar-brand`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "SLEEP", Params: []string{`1000`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "MAKE", Params: []string{"SCREENSHOT"}})

	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "IF", Params: []string{`li.button-notification-flag`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "CLICK", Params: []string{`li.button-notification-flag`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "SLEEP", Params: []string{`5000`}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "MAKE", Params: []string{"SCREENSHOT"}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "ELSE", Params: []string{}})
	chain.Chain = append(chain.Chain, internal.ChainElem{Command: "MAKE", Params: []string{"SCREENSHOT"}})

	br, err := internal.MakeBrowser("debug")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		br.Cancel1()
		br.Cancel2()
		br.Cancel3()
	}()

	err = br.RunPipeline(chain)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(br.CustomData)
}
