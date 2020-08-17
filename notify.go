package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "time"
)

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

        
        // 第二种方式
        // username := request.Form.Get("username")
        // password := request.Form.Get("password")

        // fmt.Printf("POST form-urlencoded: username=%s, password=%s\n", username, password)

        // fmt.Fprintf(writer, `{"code":0}`)

        // w.Write([]byte(`hello world`))
    })
    http.ListenAndServe(":3777", nil)
}