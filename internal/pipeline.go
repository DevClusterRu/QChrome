package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	_ "github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Instance struct {
	Actions    []chromedp.Action
	Screenshot []byte
	Mode       string
	Nodes      []*cdp.Node
	Data       []map[string]string
}

func (d *Instance) debug(command, str string) {
	if d.Mode == "debug" {
		fmt.Println("==> Command ", command, str)
	}
}

func parseScript(name string) []string {
	body, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(body), "\n")
}

func RunPipeline(find, script, mode string) (response string, err error) {

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
		chromedp.WindowSize(1200, 711), // init with a mobile view
	)

	browserCtx, browserCancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer browserCancel()

	ctx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	dp := Instance{Mode: mode}

	instructions := parseScript(script)

	for _, step := range instructions {
		step = strings.TrimSpace(step)
		if step == "" || (len(step) > 1 && step[0:2] == "--") {
			continue
		}
		directives := strings.Split(step, " ")
		if len(directives) < 2 {
			log.Println("Directive missing: ", step)
			continue
		}

		switch directives[0] {
		case "URL":
			dp.debug("Navigate to", directives[1])
			dp.Actions = append(dp.Actions, chromedp.Navigate(directives[1]))
			break
		case "WAIT":
			dp.debug("Wait for", directives[1])
			dp.Actions = append(dp.Actions, chromedp.WaitVisible(directives[1], chromedp.ByQuery))
			break
		case "ITERABLE":
			dp.debug("Iter for", directives[1])
			dp.Actions = append(dp.Actions, chromedp.Nodes(directives[1], &dp.Nodes, chromedp.ByQueryAll))
			break
		case "GET":
			//GET orig_name 0-0-1-0-0
			allSets := strings.Join(directives[1:], " ")
			dp.debug("GET for", allSets)

			sets := strings.Split(allSets, "|")

			dp.Actions = append(dp.Actions, chromedp.ActionFunc(func(c context.Context) error {
				for _, node := range dp.Nodes {
					dom.RequestChildNodes(node.NodeID).WithDepth(100).Do(c)
				}
				time.Sleep(100 * time.Millisecond)

				//Т.к. мы знаем количество нод, мы создаем это количество пустых мап для заполнения
				dp.Data = make([]map[string]string, len(dp.Nodes))

				for mapKey, node := range dp.Nodes {
					dp.Data[mapKey] = make(map[string]string)

					for _, set := range sets {

						el := strings.Split(strings.TrimSpace(set), " ")
						key := strings.TrimSpace(el[0])
						path := strings.Split(strings.TrimSpace(el[1]), "-")

						fmt.Println("Set: ", key, path)

						child := node
						good := true
						for _, p := range path {
							i, _ := strconv.Atoi(p)
							if len(child.Children) > 0 {
								child = child.Children[i]
							} else {
								good = false
								break
							}
						}

						if good {
							dp.Data[mapKey][key] = child.Children[0].NodeValue
						}
					}

				}
				return nil
			}))
			break

		case "KEY":
			if len(directives) < 3 {
				break
			}
			keys := directives[2]
			switch keys {
			case "ENTER":
				keys = kb.Enter
				break
			case "$1":
				keys = find
				break
			}
			dp.debug("Enter key", directives[1]+" : "+directives[2])
			dp.Actions = append(dp.Actions, chromedp.SendKeys(directives[1], keys, chromedp.ByQuery))
			break
		case "MAKE":
			switch directives[1] {
			case "SCREENSHOT":
				dp.debug("FullScreenshot", "")
				dp.Actions = append(dp.Actions, chromedp.FullScreenshot(&dp.Screenshot, 90))
			}
			break
		case "SLEEP":
			msec, err := strconv.Atoi(directives[1])
			if err != nil {
				log.Fatal(err)
			}
			dp.debug("Sleep", directives[1])
			dp.Actions = append(dp.Actions, chromedp.Sleep(time.Duration(msec)*time.Millisecond))
			break
		case "CLICK":
			dp.debug("Click", directives[1])
			dp.Actions = append(dp.Actions, chromedp.Click(directives[1], chromedp.ByQuery))
			break
		}
	}

	dp.debug("Start", "")
	err = chromedp.Run(ctx, dp.Actions...)
	if err != nil {
		return
	}
	dp.debug("End", "")

	if dp.Screenshot != nil {
		err = ioutil.WriteFile("fullScreenshot.png", dp.Screenshot, 0777)
		if err != nil {
			return
		}
	}

	b, err := json.Marshal(dp.Data)
	if err != nil {
		return
	}
	return string(b), nil

}
