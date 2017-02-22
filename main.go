package main

import (
    "fmt"
    "net/http"
)

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
func getHostsList(w http.ResponseWriter, r *http.Request) {
    //mock
    fmt.Fprintf(w, "[{\"id\":0,\"name\":\"yuggoth\",\"ip\":\"10.64.73.7\"},{\"id\":1,\"name\":\"vm\",\"ip\":\"10.64.73.8\"}]")
}
func getLivecheck(w http.ResponseWriter, r *http.Request) { 
    //mock
    fmt.Fprintf(w, "[{\"host_id\":0,\"check_id\":0,\"name\":\"ping\",\"alive\":true},{\"host_id\":1,\"check_id\":1,\"name\":\"ping\",\"alive\":true},{\"host_id\":1,\"check_id\":2,\"name\":\"tcp/80\",\"alive\":false}]")
}
func getLivecheckList(w http.ResponseWriter, r *http.Request) { 
    getLivecheck(w,r)
}

func main() {
    fmt.Println("Welcome to sauron runner!")

    http.HandleFunc("/", staticHandler)
    http.HandleFunc("/hosts/", getHostsList)
    http.HandleFunc("/livechecks/", getLivecheckList)
    http.HandleFunc("/probe/", getLivecheck)
    http.ListenAndServe(":8888", nil)
}

//go:generate go run scripts/packTextAssets.go
