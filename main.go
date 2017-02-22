package main

import (
    "fmt"
    "net/http"
    "net"
    "time"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/yaml.v2"
)

type Probe struct {
    Title string    `yaml:"title"`
    Proto string    `yaml:"proto"`
    Port int32      `yaml:"port"`
    Alive bool      `yaml:"alive"`
}
type Host struct {
    Title string    `yaml:"title"`
    IP string       `yaml:"ip"`
    Probes []*Probe `yaml:"probes"`
}
var Database []*Host
var DatabaseJson string

func loadConfig(){
    data, err := ioutil.ReadFile("config.yml") 
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    err = yaml.Unmarshal(data, &Database)
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    jsonObj, err := json.Marshal(Database)
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    DatabaseJson = string(jsonObj)
}

func livecheck(ip string, probe Probe) bool{
    //mock ^.^
    if probe.Proto=="tcp" {
        cs := ip+":"+fmt.Sprintf("%v", probe.Port)
        conn, err := net.DialTimeout(probe.Proto, cs, 250*time.Millisecond)
        return err==nil
    } else if probe.Proto=="udp"{
        return false
    } else { // assume ICMP PING as everything other for simple fallback
        state, _ := Ping(ip)
        return state
    }
 
}

func checkAll() string{
    LiveStatus := Database

    for _,host := range LiveStatus {
        for _,probe := range host.Probes {
            probe.Alive=livecheck(host.IP, *probe)
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
    loadConfig()

    fmt.Println("Welcome to sauron runner!")

    http.HandleFunc("/", staticHandler)
    http.HandleFunc("/definitions/", getDefinitions)
    http.HandleFunc("/probe/", getLivecheck)
    http.ListenAndServe(":8888", nil)
}

//go:generate go run scripts/packTextAssets.go
