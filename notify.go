package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "time"
    "io/ioutil"
    "os"
    "io"
    "strconv"
)

type Transaction struct {
    Types       int
    Price       float64
    Status      int
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

type Money struct {
    UsdtQty     float64
    BtcQty      float64
    LastPrice   float64
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

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
        fmt.Println(time.Now())
        fmt.Println("接收到异步请求")
        request.ParseForm()//获取请求参数
        
        uri := request.URL.String()
        method := request.Method
        
        fmt.Println(uri, method)

        mjson,_ :=json.Marshal(request.PostForm)
        mString :=string(mjson)
        fmt.Println("异步数据")
        fmt.Println(mString)

        data := GetLastTransaction()
        
        noticeType := request.Form.Get("noticeType")
        switch noticeType {
        case "DEAL_CREATE":
            data.Status = 3

            money := GetMoney()
            money.LastPrice, _ = strconv.ParseFloat(request.Form.Get("price"), 64)
            break
        case "ORDER_FINISH":
            data.Status = 2
            break
        case "ORDER_CANCEL":
            data.Status = 1
            break
        case "ORDER_CREATE":
            data.Status = 0
            break
        default:
            data.Status = 0
        }

        b, err := json.Marshal(data)
        if err != nil {
             fmt.Println("encoding faild")
        }
        writeToLog(`transaction`, string(b))

        w.Write([]byte(`success`))
    })
    http.ListenAndServe(":3777", nil)
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