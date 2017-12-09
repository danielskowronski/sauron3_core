package main

import (
    "fmt"
    "net/http"
    "net"
    "time"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/yaml.v2"
    "github.com/akamensky/argparse"
    "os"
    "strconv"
)
const retryCount = 5


type Probe struct {
    Title  string    `yaml:"title"`
    Proto  string    `yaml:"proto"`
    Port   int32     `yaml:"port"`
    Alive  bool      `yaml:"alive"`
}
type Host struct {
    Title  string    `yaml:"title"`
    IP     string    `yaml:"ip"`
    Probes []*Probe  `yaml:"probes"`
}
var Database []*Host
var DatabaseJson string

func loadConfig(path string){
    data, err := ioutil.ReadFile(path) 
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    err = yaml.Unmarshal(data, &Database)
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    jsonObj, err := json.Marshal(Database)
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    DatabaseJson = string(jsonObj)
}

func livecheck(ip string, probe Probe) bool{
    if probe.Proto=="tcp" {
    cs := ip+":"+fmt.Sprintf("%v", probe.Port)
        _, err := net.DialTimeout(probe.Proto, cs, TCPOpenTimeout*time.Millisecond)
        return err==nil
    } else if probe.Proto=="udp"{
        return false //non trivial, TODO: implement me
    } else { // assume ICMP PING as everything other for simple fallback
        state, _ := Ping(ip)
        return state
    }
 
}

func checkAll() string{
    LiveStatus := Database

    for _,host := range LiveStatus {
        for _,probe := range host.Probes {
            for i:=0; i<retryCount; i++ {
                probe.Alive=livecheck(host.IP, *probe)
                if probe.Alive { break; }
            }
        }
    }

    jsonObj, err := json.Marshal(LiveStatus)
    if err != nil { return "{\"error!\"}" }
    return string(jsonObj)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path=="/" {
        fmt.Fprintf(w, "%s", index_html)
    } else if r.URL.Path=="/style.css" {
        w.Header().Set("Content-Type", "text/css")
        fmt.Fprintf(w, "%s", style_css)
    } else if r.URL.Path=="/script.js" {
        fmt.Fprintf(w, "%s", script_js)
    } else {
        fmt.Fprintf(w, "Requested URL was not handled by this application.\n%s", r.URL.Path)
    } 
}
func getDefinitions(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, DatabaseJson)
}
func getLivecheck(w http.ResponseWriter, r *http.Request) { 
    fmt.Fprintf(w, checkAll())
}

func main() {
    parser := argparse.NewParser("sauron3", "a real time eye on your network")

    cfg := parser.String("c", "config", 
        &argparse.Options{Required: false, Help: "Path to config file (default is ./config.yml)"})
    port := parser.String("p", "port",  
        &argparse.Options{Required: false, Help: "Port for webgui (default is 8888)"})
    host := parser.String("b", "bind",  
        &argparse.Options{Required: false, Help: "IP to bind (default is 0.0.0.0)"})
    err := parser.Parse(os.Args)
    if err != nil {
        fmt.Print(parser.Usage(err))
        os.Exit(0)
    }

    if *cfg =="" { *cfg="config.yml" }
    if *port=="" { *port="8888" } else {
        portNumeric, err := strconv.Atoi(*port)
        if err!=nil || portNumeric<1 || portNumeric>65535 {
            fmt.Println("[!!] Invalid port number. Exiting")
            os.Exit(1)
        }
    }

    bindString:=*host+":"+*port
    bindStringFull:=bindString
    if *host=="" { 
        bindStringFull="0.0.0.0:"+*port
    }

    fmt.Println("Welcome to sauron runner!")
    fmt.Println("Copyright 2017 Daniel Skowroński <daniel@dsinf.net>")
    fmt.Println("")

    http.HandleFunc("/", staticHandler)
    http.HandleFunc("/definitions/", getDefinitions)
    http.HandleFunc("/probe/", getLivecheck)
    
    fmt.Println("Loading config file "+*cfg)
    loadConfig(*cfg)
    fmt.Println("Starting to listen at http://"+bindStringFull)
    http.ListenAndServe(bindString, nil)
}

//go:generate go run scripts/packTextAssets.go
