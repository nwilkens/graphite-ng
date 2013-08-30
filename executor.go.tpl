package main
import (
    "fmt"
            "github.com/Dieterbe/graphite-ng/chains"
)

func main () {
    from := int32({{.From}})
    until := int32({{.Until}})
    var dep_el chains.ChainEl

fmt.Print("[")
{{range .Targets}}
    dep_el = {{.Cmd}}
    dep_el.Settings <- from
    dep_el.Settings <- until
    fmt.Printf("{\"target\": \"{{.Query}}\", \"datapoints\": [")
    for {
         d := <- dep_el.Link
         fmt.Printf("[%f, %d]", d.Value, d.Ts)
         if d.Ts >= until {
             break
         } else {
            fmt.Printf(", ")
        }
    }
    fmt.Printf("]},\n") // last shouldn't have extra comma.
{{end}}
fmt.Printf("]")
}
