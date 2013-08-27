package main
import (
    "fmt"
            "github.com/Dieterbe/graphite-ng/chains"
)

func main () {
    from := int32({{.From}})
    until := int32({{.Until}})
    var dep_el chains.ChainEl

{{range .Targets}}
    fmt.Printf("{'target': {{.Query}})")
    dep_el = {{.Cmd}}
    dep_el.Settings <- from
    dep_el.Settings <- until
    for {
         d := <- dep_el.Link
         fmt.Println(d)
         if d.Ts >= until {
             break
         }
    }
{{end}}
}
