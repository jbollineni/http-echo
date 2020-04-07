package main
  
import (
        "errors"
        "fmt"
        "html/template"
        "net"
        "net/http"
        "os"
        "strings"
)

func main() {
        http.HandleFunc("/", echo)
        http.ListenAndServe(":5000", nil)
}

func echo(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "Close")

        if r.URL.Path == "/" {

                t, err := template.ParseFiles("template.html")
                if err != nil {
                        fmt.Fprintf(w, "Unable to load template")
                }

                interfaceIP, err := interfaceIP()
                if err != nil {
                        fmt.Println(err)
                }

                hostName, _ := os.Hostname()
                httpHost := r.Host
                httpSourceAddr := strings.SplitN(r.RemoteAddr, ":", 2)[0]
                httpSourcePort := strings.SplitN(r.RemoteAddr, ":", 2)[1]
                httpSrcAddrPtrarray, _ := net.LookupAddr(httpSourceAddr)
                httpXLbName := r.Header.Get("X-LB-Name")
                httpXff := r.Header.Get("X-Forwarded-For")
                httpUserAgent := r.Header.Get("User-Agent")
                httpSrcAddrPtr := strings.Join(httpSrcAddrPtrarray, " ") //convert string array to string

                type EchoData struct {
                        HostName      string `json:"hostName"`
                        InterfaceIP   string `json:"interfaceIP"`
                        SourceIP      string `json:"httpSourceAddr"`
                        SourceDNSPtr  string `json:"httpSrcAddrPtr"`
                        SourcePort    string `json:"httpSourcePort"`
                        HostHeader    string `json:"httpHost"`
                        HTTPXLbName   string `json:"httpXLbName"`
                        HTTPXff       string `json:"httpXff"`
                        HTTPUserAgent string `json:"httpUserAgent"`
                }

                data := EchoData{
                        HostName:      hostName,
                        InterfaceIP:   interfaceIP,
                        SourceIP:      httpSourceAddr,
                        SourceDNSPtr:  httpSrcAddrPtr,
                        SourcePort:    httpSourcePort,
                        HostHeader:    httpHost,
			HTTPXLbName:   httpXLbName,
                        HTTPUserAgent: httpUserAgent,
                }
                t.Execute(w, data)

        }
}

//Following code is from https://play.golang.org/p/BDt3qEQ_2H

func interfaceIP() (string, error) {
        ifaces, err := net.Interfaces()
        if err != nil {
                return "", err
        }
        for _, iface := range ifaces {
                if iface.Flags&net.FlagUp == 0 {
                        continue // interface down
                }
                if iface.Flags&net.FlagLoopback != 0 {
                        continue // loopback interface
                }
                addrs, err := iface.Addrs()
                if err != nil {
                        return "", err
                }
                for _, addr := range addrs {
                        var ip net.IP
                        switch v := addr.(type) {
                        case *net.IPNet:
                                ip = v.IP
                        case *net.IPAddr:
                                ip = v.IP
                        }
                        if ip == nil || ip.IsLoopback() {
                                continue
                        }
                        ip = ip.To4()
                        if ip == nil {
                                continue // not an ipv4 address
                        }
                        return ip.String(), nil
                }
        }
        return "", errors.New("Network Connection issue!")
}						