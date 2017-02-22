// http://stackoverflow.com/a/29500100

package main

import (
    "io"
    "io/ioutil"
    "os"
    "strings"
)

func main() {
    fs, _ := ioutil.ReadDir("./assets")
    out, _ := os.Create("textfiles.go")
    out.Write([]byte("package main \n\nconst (\n"))
    for _, f := range fs {
        if strings.HasSuffix(f.Name(), ".html") || strings.HasSuffix(f.Name(), ".css") || strings.HasSuffix(f.Name(), ".js") {
            varname := strings.Replace(f.Name(), ".", "_", -1)
            out.Write([]byte(varname+ " = `"))
            f, _ := os.Open("./assets/"+f.Name())
            io.Copy(out, f)
            out.Write([]byte("`\n"))
        }
    }
    out.Write([]byte(")\n"))
}