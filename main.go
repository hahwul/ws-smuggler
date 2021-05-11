package main
import (
        "crypto/tls"
        b64 "encoding/base64"
        "fmt"
        "io"
        //"io/ioutil"
        "log"
        "net"
        "net/http"
        "strings"
        "time"
        "github.com/labstack/echo/v4"
        "github.com/labstack/echo/v4/middleware"
        "github.com/tylerb/graceful"
)
type Response struct {
        Status int `json:"status"`
        Logs []string `json:"logs"`
        Response string `json:"res"`
}
func GetTemplate() string {
        template := "PHNjcmlwdCBzcmM9Imh0dHBzOi8vY29kZS5qcXVlcnkuY29tL2pxdWVyeS0zLjYuMC5taW4uanMiIGludGVncml0eT0ic2hhMjU2LS94VWorM09KVTV5RXhscTZHU1lHU0hrN3RQWGlreW5TN29nRXZEZWovbTQ9IiBjcm9zc29yaWdpbj0iYW5vbnltb3VzIj48L3NjcmlwdD4KPHN0eWxlPgpAaW1wb3J0IHVybChodHRwczovL2ZvbnRzLmdvb2dsZWFwaXMuY29tL2Nzcz9mYW1pbHk9TWVycml3ZWF0aGVyKTsKKiwKKjpiZWZvcmUsCio6YWZ0ZXIgewogIC1tb3otYm94LXNpemluZzogYm9yZGVyLWJveDsKICAtd2Via2l0LWJveC1zaXppbmc6IGJvcmRlci1ib3g7CiAgYm94LXNpemluZzogYm9yZGVyLWJveDsKfQoKaHRtbCwgYm9keSB7CiAgYmFja2dyb3VuZDogIzFhMWExYTsKICBmb250LWZhbWlseTogJ01lcnJpd2VhdGhlcicsIHNhbnMtc2VyaWY7CiAgcGFkZGluZzogMWVtOwp9CgpsYWJlbCB7CiAgZm9udC1mYW1pbHk6ICdNZXJyaXdlYXRoZXInLCBzYW5zLXNlcmlmOwogIGZvbnQtc2l6ZTogMTJweDsKICBjb2xvcjogd2hpdGU7Cn0KCmgxIHsKICB0ZXh0LWFsaWduOiBjZW50ZXI7CiAgY29sb3I6ICM4ODk7Cn0KCmZvcm0gewogIG1heC13aWR0aDogMTAwMHB4OwogIHRleHQtYWxpZ246IGNlbnRlcjsKICBtYXJnaW46IDIwcHggYXV0bzsKfQpmb3JtIGlucHV0LCBmb3JtIHRleHRhcmVhIHsKICBib3JkZXI6IDA7CiAgZm9udC1zaXplOiAxMnB4OwogIG91dGxpbmU6IDA7CiAgcGFkZGluZzogMWVtOwogIGNvbG9yOiAjYmJiOwogIGJhY2tncm91bmQtY29sb3I6ICM0NDQ7CiAgLW1vei1ib3JkZXItcmFkaXVzOiA4cHg7CiAgLXdlYmtpdC1ib3JkZXItcmFkaXVzOiA4cHg7CiAgYm9yZGVyLXJhZGl1czogOHB4OwogIGRpc3BsYXk6IGJsb2NrOwogIHdpZHRoOiAxMDAlOwogIG1hcmdpbi10b3A6IDFlbTsKICBmb250LWZhbWlseTogJ01lcnJpd2VhdGhlcicsIHNhbnMtc2VyaWY7CiAgLW1vei1ib3gtc2hhZG93OiAwIDFweCAxcHggcmdiYSgwLCAwLCAwLCAwLjEpOwogIC13ZWJraXQtYm94LXNoYWRvdzogMCAxcHggMXB4IHJnYmEoMCwgMCwgMCwgMC4xKTsKICBib3gtc2hhZG93OiAwIDFweCAxcHggcmdiYSgwLCAwLCAwLCAwLjEpOwogIHJlc2l6ZTogbm9uZTsKfQpmb3JtIGlucHV0OmZvY3VzLCBmb3JtIHRleHRhcmVhOmZvY3VzIHsKICAtbW96LWJveC1zaGFkb3c6IDAgMHB4IDJweCAjZTc0YzNjICFpbXBvcnRhbnQ7CiAgLXdlYmtpdC1ib3gtc2hhZG93OiAwIDBweCAycHggI2U3NGMzYyAhaW1wb3J0YW50OwogIGJveC1zaGFkb3c6IDAgMHB4IDJweCAjZTc0YzNjICFpbXBvcnRhbnQ7Cn0KZm9ybSAjaW5wdXQtc3VibWl0IHsKICBjb2xvcjogd2hpdGU7CiAgYmFja2dyb3VuZDogIzMzNDsKICBjdXJzb3I6IHBvaW50ZXI7Cn0KZm9ybSAjaW5wdXQtc3VibWl0OmhvdmVyIHsKICAtbW96LWJveC1zaGFkb3c6IDAgMXB4IDFweCAxcHggcmdiYSgxNzAsIDE3MCwgMTcwLCAwLjYpOwogIC13ZWJraXQtYm94LXNoYWRvdzogMCAxcHggMXB4IDFweCByZ2JhKDE3MCwgMTcwLCAxNzAsIDAuNik7CiAgYm94LXNoYWRvdzogMCAxcHggMXB4IDFweCByZ2JhKDE3MCwgMTcwLCAxNzAsIDAuNik7Cn0KZm9ybSB0ZXh0YXJlYSB7CiAgaGVpZ2h0OiAxMjZweDsKfQoKLmhhbGYgewogIGZsb2F0OiBsZWZ0OwogIHdpZHRoOiA0OCU7CiAgbWFyZ2luLWJvdHRvbTogMWVtOwp9CgoucmlnaHQgewogIHdpZHRoOiA1MCU7Cn0KCi5sZWZ0IHsKICBtYXJnaW4tcmlnaHQ6IDIlOwp9CgpAbWVkaWEgKG1heC13aWR0aDogNDgwcHgpIHsKICAuaGFsZiB7CiAgICB3aWR0aDogMTAwJTsKICAgIGZsb2F0OiBub25lOwogICAgbWFyZ2luLWJvdHRvbTogMDsKICB9Cn0KLyogQ2xlYXJmaXggKi8KLmNmOmJlZm9yZSwKLmNmOmFmdGVyIHsKICBjb250ZW50OiAiICI7CiAgLyogMSAqLwogIGRpc3BsYXk6IHRhYmxlOwogIC8qIDIgKi8KfQoKLmNmOmFmdGVyIHsKICBjbGVhcjogYm90aDsKfQoKCjwvc3R5bGU+CjxoMT53cy1zbXVnZ2xlcjwvaDE+Cjxmb3JtIGNsYXNzPSJjZiIgbWV0aG9kPSJwb3N0IiBhY3Rpb249Ii9zZW5kIiB0YXJnZXQ9ImhkbiI+CiAgPGRpdiBjbGFzcz0iaGFsZiBsZWZ0IGNmIj4KICAgIDxsYWJlbD5SZXF1ZXN0PGxhYmVsPgogICAgPGlucHV0IHR5cGU9InRleHQiIG5hbWU9InRhcmdldCIgcGxhY2Vob2xkZXI9IlRhcmdldCBkb21haW46cG9ydCAoZS5nIGhhaHd1bC5jb206MzMxNDUpIj4KICAgIDx0ZXh0YXJlYSBuYW1lPSJvX2RhdGEiIHR5cGU9InRleHQiIHBsYWNlaG9sZGVyPSJPcmlnaW5hbCBIVFRQIFJlcXVlc3QiPjwvdGV4dGFyZWE+CiAgICAKICAgIDx0ZXh0YXJlYSBuYW1lPSJzX2RhdGEiIHR5cGU9InRleHQiIHBsYWNlaG9sZGVyPSJTbXVnZ2xlZCBIVFRQIFJlcXVlc3QiPjwvdGV4dGFyZWE+CiAgICA8bGFiZWwgc3R5bGU9ImRpc3BsYXk6aW5saW5lIj48aW5wdXQgdHlwZT0iY2hlY2tib3giIGlkPSJzc2wiIG5hbWU9InNzbCIgc3R5bGU9ImRpc3BsYXk6aW5saW5lOyB3aWR0aDogYXV0bzsiPiBTU0w8L2xhYmVsPgogICAgICA8aW5wdXQgdHlwZT0ic3VibWl0IiB2YWx1ZT0iU2VuZCIgc3R5bGU9ImN1cnNvcjpwb2ludGVyOyI+CiAgPC9kaXY+CiAgPGRpdiBjbGFzcz0iaGFsZiByaWdodCBjZiI+CiAgICA8bGFiZWw+UmVzcG9uc2U8bGFiZWw+CiAgICA8dGV4dGFyZWEgaWQ9InJlcyIgdHlwZT0idGV4dCIgc3R5bGU9ImhlaWdodDozODNweDsiIHJlYWRvbmx5PjwvdGV4dGFyZWE+CiAgPC9kaXY+CiAgPGRpdj4KCTx0ZXh0YXJlYSBpZD0ibG9ncyIgdHlwZT0idGV4dCIgc3R5bGU9ImhlaWdodDoyMDBweDsgYmFja2dyb3VuZDogdHJhbnNwYXJlbnQ7IGJvcmRlcjogc29saWQgMXB4OyIgcmVhZG9ubHk+PC90ZXh0YXJlYT4KICA8L2Rpdj4KICA8aWZyYW1lIG5hbWU9ImhkbiIgc3R5bGU9InZpc2liaWxpdHk6aGlkZGVuIj48L2lmcmFtZT4KPC9mb3JtPgo8c2NyaXB0PgoJZnVuY3Rpb24gdXBkYXRlTG9ncygpewoJCSQuZ2V0KCcvaW5mbycsIGZ1bmN0aW9uKGRhdGEpIHsKCQkJbG9nU3RyaW5nID0gIiIKCQkJZGF0YS5sb2dzLmZvckVhY2goZnVuY3Rpb24oZCl7CgkJCQlsb2dTdHJpbmcgPSBsb2dTdHJpbmcgKyBkICsgIlxuIgoJCQl9KQoJCQlkb2N1bWVudC5nZXRFbGVtZW50QnlJZCgnbG9ncycpLnZhbHVlID0gbG9nU3RyaW5nCgkJCWRvY3VtZW50LmdldEVsZW1lbnRCeUlkKCdyZXMnKS52YWx1ZSA9IGRhdGEucmVzCgkJfSk7Cgl9CgoJc2V0SW50ZXJ2YWwoZnVuY3Rpb24oKSB7CgkJdXBkYXRlTG9ncygpCgl9LCAxMDAwKTsKPC9zY3JpcHQ+Cg=="
        //dat, err := ioutil.ReadFile("./index.html")
        dat, err := b64.StdEncoding.DecodeString(template)
        _= err
        return string(dat)
}
func Send(target, odata, sdata string, ssl bool) ([]string, string) {
        var logs []string
        var response string
        ssls := "false"
        if ssl {
                ssls = "true"
        }
        logs = append(logs, "Sending connection smuggling request")
        logs = append(logs, "[target: "+target+"][ssl: "+ssls+"]")
        if ssl {
                conf := &tls.Config{
                        InsecureSkipVerify: true,
                }
                conn, err := tls.Dial("tcp", target, conf)
                if nil != err {
                        logs = append(logs, "Failed to connect to server")
                        return logs, response
                }
                recvBuf := make([]byte, 4096)
                req1 := odata
                req2 := sdata
                conn.Write([]byte(req1))
                conn.Read(recvBuf)
                conn.Write([]byte(req2))
                conn.Read(recvBuf)
                bufString := strings.ReplaceAll(string(recvBuf),"\u0000","")
                response = bufString
                if nil != err {
                        if io.EOF == err {
                                log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
                                return logs, response
                        }
                        log.Printf("fail to receive data; err: %v", err)
                        return logs, response
                }
                conn.Close()
        } else {
                conn, err := net.Dial("tcp", target)
                if nil != err {
                        logs = append(logs, "Failed to connect to server")
                        return logs, response
                }
                req1 := odata
                req2 := sdata
                recvBuf := make([]byte, 4096)
                conn.Write([]byte(req1))
                conn.Read(recvBuf)
                conn.Write([]byte(req2))
                conn.Read(recvBuf)
                bufString := strings.ReplaceAll(string(recvBuf),"\u0000","")
                response = bufString
                log.Printf("%s", recvBuf)
                if nil != err {
                        if io.EOF == err {
                                log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
                                return logs, response
                        }
                        log.Printf("fail to receive data; err: %v", err)
                        return logs, response
                }
                conn.Close()
        }
        return logs, response
}
func main() {
        e := echo.New()
        var logs []string
        var response string
        e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
                XSSProtection:         "",
                ContentTypeNosniff:    "",
                XFrameOptions:         "",
                HSTSMaxAge:            3600,
        }))
        e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
                Format: "method=${method}, uri=${uri}, status=${status}\n",
        }))
        e.GET("/", func(c echo.Context) error {
                return c.HTML(http.StatusOK,GetTemplate())
        })
        e.GET("/info", func(c echo.Context) error {
                c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
                r := &Response {
                        Status: 0,
                        Logs: logs,
                        Response: response,
                }
                return c.JSON(http.StatusOK,r)
        })
        e.POST("/send", func(c echo.Context) error {
                c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
                logs = append(logs,"Sending Request..")
                target := c.FormValue("target")
                ssl := c.FormValue("ssl")
                bssl := false
                if ssl == "on" {
                        bssl = true
                }
                odata := c.FormValue("o_data")
                sdata := c.FormValue("s_data")
                fmt.Println(target)
                fmt.Println(odata)
                fmt.Println(sdata)
                fmt.Println(bssl)
                var tlogs []string
                tlogs, response = Send(target,odata,sdata,bssl)
                for _, v := range tlogs {
                        logs = append(logs, v)
                }
                return c.HTML(http.StatusOK,GetTemplate())
        })
        e.Server.Addr = ":4556"
        // Serve it like a boss
        fmt.Println("Start WS Server")
        fmt.Println("Open http://localhost"+e.Server.Addr+" using web bror")
        graceful.ListenAndServe(e.Server, 5*time.Second)
}
