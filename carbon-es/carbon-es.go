package main

import (
	"bufio"
	"fmt"
	"github.com/graphite-ng/graphite-ng/es"
	"github.com/graphite-ng/graphite-ng/util"
	"github.com/mattbaird/elastigo/api"
	"github.com/mattbaird/elastigo/core"
	"github.com/stvp/go-toml-config"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		es_host        = config.String("elasticsearch.host", "undefined")
		es_port        = config.Int("elasticsearch.port", 9200)
		es_max_pending = config.Int("elasticsearch.max_pending", 1000000)
		in_port        = config.Int("in.port", 2003)
	)
	fmt.Println(*es_max_pending)
	err := config.Parse("graphite-ng.conf")
	util.DieIfError(err)

	api.Domain = *es_host
	api.Port = strconv.Itoa(*es_port)
	done := make(chan bool)
	core.BulkIndexorGlobalRun(4, done)

	// listen for incoming metrics
	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", *in_port))
	util.DieIfError(err)
	listener, err := net.ListenTCP("tcp", addr)
	util.DieIfError(err)
	defer listener.Close()

	fmt.Printf("carbon-es ready to serve on %d\n", *in_port)
	for {
		conn_in, err := listener.Accept()
		util.DieIfError(err)
		go handleClient(conn_in)
	}
}

func handleClient(conn_in net.Conn) {
	defer conn_in.Close()
	reader := bufio.NewReader(conn_in)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Printf("WARN connection closed uncleanly/broken: %s\n", err.Error())
				return
			}
		}
		str := strings.TrimSpace(string(buf))
		//fmt.Println(str)
		dp_str := strings.Split(str, " ")
		metric_name := dp_str[0]
		value, err := strconv.ParseFloat(dp_str[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse value out of metric '%s': %s\n", str, err.Error())
			continue
		}
		ts, err := strconv.ParseInt(dp_str[2], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse timestamp out of metric '%s': %s\n", str, err.Error())
			continue
		}
		// for some reason IndexBulk needs an id set.
		// seems a little redundant but i guess we can use it to avoid
		// duplicate values
		id := fmt.Sprintf("%s_%d", metric_name, ts)
		dp := es.Datapoint{metric_name, int32(ts), value}
		date := time.Now()
		err = core.IndexBulk("carbon-es", "datapoint", id, &date, &dp)
		util.DieIfError(err)
	}
}
