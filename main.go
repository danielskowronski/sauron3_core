package sauron3

import (
    "fmt"
    "net"
    "time"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/yaml.v2"
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

func LoadConfig(path string){
    data, err := ioutil.ReadFile(path) 
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    err = yaml.Unmarshal(data, &Database)
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    jsonObj, err := json.Marshal(Database)
    if err != nil { panic(fmt.Sprintf("%v",err)) }
    DatabaseJson = string(jsonObj)
}

func Livecheck(ip string, probe Probe) bool{
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

func CheckAll() string{
    LiveStatus := Database

    for _,host := range LiveStatus {
        for _,probe := range host.Probes {
            for i:=0; i<retryCount; i++ {
                probe.Alive=Livecheck(host.IP, *probe)
                if probe.Alive { break; }
            }
        }
    }

    jsonObj, err := json.Marshal(LiveStatus)
    if err != nil { return "{\"error!\"}" }
    return string(jsonObj)
}

