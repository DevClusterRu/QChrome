package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	_ "github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type Instance struct {
	Ctx               context.Context
	Screenshot        []byte
	ScreenshotCounter int
	Mode              string
	Nodes             []*cdp.Node
	Data              []map[string]string //For iterable only!
	CustomData        map[string]string   //Any other data
	ConditionResult   bool

	Cancel1 context.CancelFunc
	Cancel2 context.CancelFunc
	Cancel3 context.CancelFunc
}

func (d *Instance) debug(command string, str ...interface{}) {
	if d.Mode == "debug" {
		fmt.Println("==> Command ", command, str)
	}
}

func MakeBrowser(mode string) (dp *Instance, err error) {
	dp = &Instance{
		Mode:       mode,
		CustomData: make(map[string]string),
	}

	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"),
		chromedp.WindowSize(1200, 711), // init with a mobile view
		chromedp.UserDataDir("userdata"),
	)

	var browserCtx, ctx context.Context

	browserCtx, dp.Cancel3 = chromedp.NewExecAllocator(context.Background(), options...)
	ctx, dp.Cancel2 = chromedp.NewContext(browserCtx)
	// create a timeout
	dp.Ctx, dp.Cancel1 = context.WithTimeout(ctx, 30*time.Second)
	return dp, nil
}

func (dp *Instance) RunPipeline(r Request) (err error) {

Loop:
	for k, chainStep := range r.Chain {

		switch chainStep.Command {
		case "URL":
			dp.debug("Navigate to", chainStep.Params[0])
			_, err := dp.DoChromium("simple", chromedp.Navigate(chainStep.Params[0]))
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "WAIT":
			dp.debug("Wait for", chainStep.Params[0])
			_, err := dp.DoChromium("simple", chromedp.WaitVisible(chainStep.Params[0], chromedp.ByQuery))
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "ITERABLE":
			dp.debug("Iter for", chainStep.Params[0])
			_, err := dp.DoChromium("simple", chromedp.Nodes(chainStep.Params[0], &dp.Nodes, chromedp.ByQueryAll))
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "GET":
			allSets := strings.Join(chainStep.Params[1:], " ")
			dp.debug("GET for", allSets)
			err := dp.FindNodes(chainStep.Params[0], allSets)
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "KEY":
			err := dp.Keys(chainStep.Params)
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "MAKE":
			switch chainStep.Params[0] {
			case "SCREENSHOT":
				dp.debug("FullScreenshot", strconv.Itoa(dp.ScreenshotCounter))
				_, err := dp.DoChromium("simple", chromedp.FullScreenshot(&dp.Screenshot, 90))
				if err != nil {
					log.Println(err)
					break Loop
				}
				fname := fmt.Sprintf("%d-%d.png", time.Now().Unix(), dp.ScreenshotCounter)
				dp.ScreenshotCounter++
				err = ioutil.WriteFile(fname, dp.Screenshot, 0777)
				if err != nil {
					log.Println(err)
					break Loop
				}
				dp.CustomData["ss_"+fname] = ""
				break
			}
			break
		case "SLEEP":
			msec, err := strconv.Atoi(chainStep.Params[0])
			if err != nil {
				log.Println(err)
				break Loop
			}
			dp.debug("Sleep", chainStep.Params[0])
			_, err = dp.DoChromium("simple", chromedp.Sleep(time.Duration(msec)*time.Millisecond))
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "CLICK":
			dp.debug("Click", chainStep.Params[0])
			_, err := dp.DoChromium("simple", chromedp.Click(chainStep.Params[0], chromedp.ByQuery))
			if err != nil {
				log.Println(err)
				break Loop
			}
			break
		case "IF":
			dp.ConditionResult, err = dp.DoChromium("condition", chromedp.Nodes(chainStep.Params[0], &dp.Nodes, chromedp.AtLeast(0)))
			if err != nil {
				log.Println(err)
				break Loop
			}
			dp.debug("Condition", chainStep.Params[0], " : ", dp.ConditionResult)

			//???????????? ?? ???????????????? ?????????? ELSE
			var i int
			for i = k + 1; i < len(r.Chain); i++ {
				if r.Chain[i].Command == "ELSE" {
					break
				}
			}

			if dp.ConditionResult {
				r.Chain = r.Chain[k+1 : i]
			} else {
				r.Chain = r.Chain[i:]
			}

			dp.RunPipeline(r)

			break Loop //???? ?????????????????? condition, ?? ?????????? ?????? ??????????????????. ???????????? ???????????? ?????????????????? ??????, ???? ???????????????? ??????????, ???????? ??????????????
		}

	}

	return

}
