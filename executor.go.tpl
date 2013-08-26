package main
import (
    "fmt"
)

func main () {
    from := int32({{.From}})
    until := int32({{.Until}})
    out := {{.Cmd}}
    for {
        d := <-out
        fmt.Println(d)
        if d.ts >= until {
            break
        }
    }
}
