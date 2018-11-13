package main

import (
	"fmt"
	"flag"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (

)

func recordMetrics(){
//	out, err := exec.Command("ls", "-lah").Output()
//	if err != nil {
//		log.Fatal(err)
//	}
}


func main() {
	upsArg       := flag.String("ups", "cyberpower", "ups name managed by nut")
	portArg      := flag.Int("port", 8100, "port number")
	
	var listenAddr = fmt.Sprintf(":%d", *portArg)
    //upscArg  := flag.String("upsc", "/bin/upsc", "upsc path")

	flag.Parse()
    //fmt.Printf("testing cmd %s", out)
    log.Infoln("Starting NUT exporter on ups", *upsArg )	
	http.Handle("/metrics", promhttp.Handler())
    
	log.Infoln("NUT exporter started on", listenAddr)
    http.ListenAndServe(listenAddr, nil)	
}

