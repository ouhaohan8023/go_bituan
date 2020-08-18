package main

import (
	"fmt"
	// "net/http"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	// "encoding/base64"
	// "encoding/hex"
	"encoding/json"
	"os"
	"io"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"time"

)

func main() {
	CreateOrder(12355.00, "SELL", 0.001)
	// t := GetLastTransaction()
 //    fmt.Println(t.Price, t.Types, t.Status)

 //    if t.Status != 3 {
 //    	fmt.Println("等待上一笔订单执行完毕，终止本次任务")
 //    } else {
	//     c := GetPrice()
	//     fmt.Println("买入价格：" + strconv.FormatFloat(c.BidPrice, 'f', 2, 64))
	//     fmt.Println("可买入数量：" + strconv.FormatFloat(c.BidQty, 'f', 4, 64))
	//     fmt.Println("卖出数量：" + strconv.FormatFloat(c.AskPrice, 'f', 2, 64))
	//     fmt.Println("可卖出数量：" + strconv.FormatFloat(c.AskQty, 'f', 4, 64))

	//     money := GetMoney()
	//     fmt.Println(money.UsdtQty, money.BtcQty)

	//     if t.Types == 1 {
	//     	fmt.Println("上一笔为 USDT => BTC，当前需要将 BTC => USDT ， 上一次买入价格：" + strconv.FormatFloat(t.Price, 'f', 4, 64) + " 本次最快卖出价格 " + strconv.FormatFloat(c.BidPrice, 'f', 2, 64))
	//     	sub := c.BidPrice - t.Price
	//     	if sub > 0 {
	//     		fmt.Println("触发订单api")
	//     		CreateOrder(c.BidPrice, "SELL", money.BtcQty)
	//     	} else {
	//     		fmt.Println("无法触发")
	//     	}
	//     } else {
	//     	fmt.Println("上一笔为 BTC => USDT，当前需要将 USDT => BTC ， 上一次卖出价格：" + strconv.FormatFloat(t.Price, 'f', 4, 64) + " 本次最快买入价格 " + strconv.FormatFloat(c.AskPrice, 'f', 2, 64))
	//     	sub := c.AskPrice - t.Price
	//     	if sub > 0 {
	//     		fmt.Println("触发订单api")
	//     	} else {
	//     		fmt.Println("无法触发")
	//     	}
	//     }
 //    }


}

type Content struct {
    BidPrice    float64
    BidQty      float64
    AskPrice    float64
    AskQty      float64
}

type Money struct {
	UsdtQty     float64
    BtcQty      float64
    LastPrice   float64
}

type Transaction struct {
    Types       int
    Price       float64
    Status      int
}

func CreateOrder(price float64, side string, volume float64) {
	callback := `http://124.156.119.73:3777`
	timestamp := strconv.FormatInt(time.Now().UnixNano() / 1e6,10)
	sign := md5V("api_keyeyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJkYXRhSWQiOiIyMDIwMDgxNjA5MjYzNUEzRDlFMDM4LUI4MDctNDhGMi1BNURBLTE2M0MzMUY1NjRCMCIsImxhYmVsIjoidGVzdCIsImV4cCI6MTYwMDE4NTYwMCwianRpIjoiMjAyNDg1OCJ9._x-GwxVfsJhunTMPrSR4wQS2mbHD1X4Jlamg5XjKt24callbackUrl"+callback+"noticeId"+timestamp+"price"+strconv.FormatFloat(price, 'f', 2, 64)+"side"+side+"symbolbtcusdttime"+timestamp+"type1volume"+strconv.FormatFloat(volume, 'f', 4, 64)+"7f56b48d295b4bc2a8aec7b55a75fcea")
	
	url := `https://open.bituan.io/v1/create_order?`
	url += `api_key=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJkYXRhSWQiOiIyMDIwMDgxNjA5MjYzNUEzRDlFMDM4LUI4MDctNDhGMi1BNURBLTE2M0MzMUY1NjRCMCIsImxhYmVsIjoidGVzdCIsImV4cCI6MTYwMDE4NTYwMCwianRpIjoiMjAyNDg1OCJ9._x-GwxVfsJhunTMPrSR4wQS2mbHD1X4Jlamg5XjKt24&`
	url += `sign=` + sign + `&`
	url += `time=` + timestamp + `&`
	url += `symbol=btcusdt&`
	url += `price=` + strconv.FormatFloat(price, 'f', 2, 64) + `&`
	url += `side=` + side + `&`
	url += `type=1&`
	url += `volume=` + strconv.FormatFloat(volume, 'f', 4, 64) + `&`
	url += `callbackUrl=` + callback + `&`
	url += `noticeId=` + timestamp + `&`

	ctx, _ := Get(url, nil)

	// fmt.Println(ctx, ctx["code"], ctx["code"] == 0, ctx["code"] == "0")
    code := ctx["code"]

	fmt.Println(time.Now())

    if code != "0" {
    	fmt.Println("下单失败，失败原因", ctx["msg"])
    } else {
    	fmt.Println("下单成功，等待异步", timestamp)
    }
}

func md5V(str string) string  {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}


func GetPrice() (c *Content)  {
	ctx, _ := Get(`https://open.bituan.io/v1/market_dept?symbol=btcusdt&type=0`, nil)
    bid := ctx["data"].(map[string]interface{})["tick"].(map[string]interface{})["bids"].([]interface{})[0]
    bidPrice := bid.([]interface{})[0].(float64)
    bidQty := bid.([]interface{})[1].(float64)
    ask := ctx["data"].(map[string]interface{})["tick"].(map[string]interface{})["asks"].([]interface{})[0]
    askPrice := ask.([]interface{})[0].(float64)
    askQty := ask.([]interface{})[1].(float64)
    st := &Content {
      bidPrice,
      bidQty,
      askPrice,
      askQty,
	}
	b, err := json.Marshal(st)
	if err != nil {
         fmt.Println("encoding faild")
     } else {
         // fmt.Println("encoded data : ")
         // fmt.Println(b)
         // fmt.Println(string(b))
     }
    // fmt.Println(bidPrice, bidQty, askPrice, askQty)
    // fmt.Println(b, st)
    writeToLog(`price`, string(b))
    return st
}

func GetMoney() (money *Money){
	var fileName = "money"
	f, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println("read fail", err)
    }
    stb := &Money{}
    err = json.Unmarshal([]byte(f), &stb)
    return stb
}

func GetLastTransaction() (t *Transaction){
	var fileName = "transaction"
	f, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println("read fail", err)
    }
    stb := &Transaction{}
    err = json.Unmarshal([]byte(f), &stb)
    return stb
}

func writeToLog(fileName string, wireteString string) {

	var f *os.File
	var err1 error
	if checkFileIsExist(fileName) { //如果文件存在
		f, err1 = os.OpenFile(fileName,os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666) //打开文件
		// fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(fileName) //创建文件
		fmt.Println("文件不存在")
	}
	check(err1)
	_, err1 = io.WriteString(f, wireteString) //写入文件(字符串)
	check(err1)
	// fmt.Printf("写入 %d 个字节n", n)
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(fileName string) bool {
	var exist = true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Get(url string, header map[string]string) (s map[string]interface{}, err error) {
	var r http.Request
	r.ParseForm()
	body := strings.NewReader(r.Form.Encode())
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		panic(err)
	}
	for key, val := range header {
		req.Header.Add(key, val)
	}
	req.Close = true
	clt := http.Client{}
	resp, err := clt.Do(req)
	//fmt.Println(resp)

	checkError(err)

	defer resp.Body.Close()    // 绝大多数情况下的正确关闭方式

	HttpCode := resp.StatusCode
	//fmt.Println(HttpCode)
	if HttpCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//fmt.Println(err)
		}
		bodyString := string(bodyBytes)
		// fmt.Println(bodyString)
		// if (viper.GetBool(`log.http`) && connect != nil) {
		// 	connect.Create(&model.Http{Platform: platform, Code: resp.StatusCode, Result: bodyString})
		// }
		//json转map
		var r Requestbody
		r.req = bodyString
		if req2map, err := r.Json2map(); err == nil {
			return req2map, nil
		}
	} else {
		fmt.Println("请求报错")
		fmt.Println(HttpCode)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
	return nil, nil
}

func checkError(err error) {
	if err != nil{
		fmt.Println("请求报错")
		log.Fatalln(err)
	}
}

//把请求包定义成一个结构体
type Requestbody struct {
	req string
}

//以指针的方式传入，但在使用时却可以不用关心
// result 是函数内的临时变量，作为返回值可以直接返回调用层
func (r *Requestbody) Json2map() (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
		return nil, err
	}
	return result, nil
}