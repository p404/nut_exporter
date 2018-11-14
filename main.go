package main

import (
 	"time"
	"math/rand"
	"fmt"
	"flag"
	"regexp"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)


// Regex
var(

	batteryChargeRegex          =   regexp.MustCompile(`(?:battery[.]charge:(?:\s)(.*))`)
    batteryPacksRegex           =   regexp.MustCompile(`(?:battery[.]packs:(?:\s)(.*))`)
	batteryVoltageRegex         =   regexp.MustCompile(`(?:battery[.]voltage:(?:\s)(.*))`)
	batteryVoltageNominalRegex  =   regexp.MustCompile(`(?:battery[.]voltage[.]nominal:(?:\s)(.*))`)
    inputVoltageRegex           =   regexp.MustCompile(`(?:input[.]voltage:(?:\s)(.*))`)
	inputVoltageNominalRegex    =   regexp.MustCompile(`(?:input[.]voltage[.]nominal:(?:\s)(.*))`)
    outputVoltageRegex          =   regexp.MustCompile(`(?:output[.]voltage:(?:\s)(.*))`)
    outputVoltageNominalRegex   =   regexp.MustCompile(`(?:output[.]voltage[.]nominal:(?:\s)(.*))`)
    upsPowerNominalRegex        =   regexp.MustCompile(`(?:ups[.]power[.]nominal:(?:\s)(.*))`)
    upsTempRegex                =   regexp.MustCompile(`(?:ups[.]temperature:(?:\s)(.*))`)
    upsLoadRegex                =   regexp.MustCompile(`(?:ups[.]load:(?:\s)(.*))`)
	upsStatusRegex              =   regexp.MustCompile(`(?:ups[.]status:(?:\s)(.*))`)
)

// NUT Gauges
var (

	batteryCharge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_charge",
		Help: "Current battery charge (percent)",
	})
	
	batteryPacks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_pack",
		Help: "Number of battery packs on the UPS",
	})

	batteryVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_voltage",
		Help: "Current battery voltage",
	})

	batteryVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_voltage_nominal",
		Help: "Nominal battery voltage",
	})

	inputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "input_voltage",
		Help: "Current input voltage",
	})

	inputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "input_voltage_nominal",
		Help: "Nominal input voltage",
	})

	outputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ouput_voltage",
		Help: "Current output voltage",
	})
	
	outputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ouput_voltage_nominal",
		Help: "Nominal output voltage",
	})
	
	upsPowerNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_power_nominal",
		Help: "Nominal ups power",
	})
	
	upsTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_temp",
		Help: "UPS Temperature (degrees C)",
	})
	
	upsLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_load",
		Help: "Current UPS load (percent)",
	})

	upsStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_status",
		Help: "Current UPS Status (0=Calibration, 1=SmartTrim, 2=SmartBoost, 3=Online, 4=OnBattery, 5=Overloaded, 6=LowBattery, 7=ReplaceBattery, 8=OnBypass, 9=Off, 10=Charging, 11=Discharging)",
	})
)

func recordMetrics(){
	prometheus.MustRegister(upsLoad)
//	out, err := exec.Command("ls", "-lah").Output()
//	if err != nil {
//		log.Fatal(err)
//	}
	go func(){
		for {
			upsLoad.Set(rand.Float64())
			time.Sleep(2 * time.Second)
		}
	}()

}


func main() {
	upsArg   := flag.String("ups", "none", "ups name managed by nut")
	portArg  := flag.Int("port", 8100, "port number")
	
	var listenAddr = fmt.Sprintf(":%d", *portArg)
    //upscArg  := flag.String("upsc", "/bin/upsc", "upsc path")

	flag.Parse()
	recordMetrics()
    //fmt.Printf("testing cmd %s", out)
    log.Infoln("Starting NUT exporter on ups", *upsArg )	
	http.Handle("/metrics", promhttp.Handler())
    
	log.Infoln("NUT exporter started on port", *portArg)
    http.ListenAndServe(listenAddr, nil)	
}

