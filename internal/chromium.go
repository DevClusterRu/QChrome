package internal

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"log"
	"strconv"
	"strings"
	"time"
)

func (i *Instance) DoChromium(method string, action chromedp.Action) (cond bool, err error) {

	cond = false
	err = chromedp.Run(i.Ctx, action)
	if err != nil {
		return
	}

	fmt.Println("DONE")

	switch method {
	case "simple":
		break
	case "condition":
		if len(i.Nodes) > 0 {
			cond = true
		}
		break
	}
	return
}

func (dp *Instance) IKey(params []string) error {

	if len(params) < 3 {
		return fmt.Errorf("Wrong key format")
	}

	nn, err := strconv.Atoi(params[0])
	if err != nil {
		return err
	}

	if len(dp.Nodes) < nn {
		return fmt.Errorf("Finded %d needed %d", len(dp.Nodes), nn)
	}

	dp.DoChromium("simple", chromedp.ActionFunc(func(c context.Context) error {

		node := dp.Nodes[nn]

		set := strings.Split(strings.TrimSpace(params[1]), "-")

		dom.RequestChildNodes(node.NodeID).WithDepth(-1).Do(c) //Пробежимся по выбранной ноде и соберем всех дочек
		time.Sleep(100 * time.Millisecond)

		good := true
		child := node

		for _, p := range set {
			i, _ := strconv.Atoi(p)
			if len(node.Children) > i && len(child.Children) > 0 {
				child = child.Children[i]
				log.Println(child.Attributes)
			} else {
				log.Println("Can't find node ", i)
				good = false
				break
			}
		}

		if good {
			keys := params[2]
			switch keys {
			case "ENTER":
				keys = kb.Enter
				break
			case "DEL":
				keys = kb.Delete
				break
			case "BS":
				keys = kb.Backspace
				break
			}
			dp.debug("Enter key", keys)
			dp.DoChromium("simple", chromedp.SendKeys([]cdp.NodeID{child.NodeID}, keys, chromedp.ByNodeID))
			if err != nil {
				return err
			}
		}

		return nil
	}))

	return nil
}

func (dp *Instance) IClick(node_number, set_s string) error {

	nn, err := strconv.Atoi(node_number)
	if err != nil {
		return err
	}

	if len(dp.Nodes) < nn {
		return fmt.Errorf("Finded %d needed %d", len(dp.Nodes), nn)
	}

	dp.DoChromium("simple", chromedp.ActionFunc(func(c context.Context) error {

		node := dp.Nodes[nn]

		set := strings.Split(strings.TrimSpace(set_s), "-")

		dom.RequestChildNodes(node.NodeID).WithDepth(-1).Do(c) //Пробежимся по выбранной ноде и соберем всех дочек
		time.Sleep(100 * time.Millisecond)

		good := true
		child := node

		for _, p := range set {
			i, _ := strconv.Atoi(p)
			if len(node.Children) > i {
				child = child.Children[i]
				log.Println(child.Attributes)
			} else {
				log.Println("Can't find node ", i)
				good = false
				break
			}
		}

		if good {
			_, err := dp.DoChromium("simple", chromedp.Click([]cdp.NodeID{child.NodeID}, chromedp.ByNodeID))
			if err != nil {
				return err
			}
		}

		return nil
	}))

	return nil
}

func (dp *Instance) FindNodes(sets []string) error {

	dp.DoChromium("simple", chromedp.ActionFunc(func(c context.Context) error {

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

				child := node
				good := true
				for _, p := range path {
					i, _ := strconv.Atoi(p)
					if len(child.Children) > i {
						child = child.Children[i]
					} else {
						good = false
						break
					}
				}

				if good && len(child.Children) > 0 {
					if child.Children[0] != nil {
						dp.Data[mapKey][key] = strings.Join(child.Children[0].Attributes, ",")
						dp.Data[mapKey][key] = dp.Data[mapKey][key] + "[" + child.Children[0].NodeValue + "]"
					}
				}
			}

		}
		return nil
	}))

	return nil
}

func (dp *Instance) Keys(params []string) error {
	if len(params) < 2 {
		return fmt.Errorf("Wrong key format")
	}
	keys := params[1]
	switch keys {
	case "ENTER":
		keys = kb.Enter
		break
	default:
		keys = strings.Join(params[1:], " ")
		break
	}
	dp.debug("Enter key", "")
	dp.DoChromium("simple", chromedp.SendKeys(params[0], keys, chromedp.ByQuery))
	return nil
}
