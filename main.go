package main

import (
	"strconv"
	"os/exec"
 	"time"
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
		Name: "output_voltage",
		Help: "Current output voltage",
	})
	
	outputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "output_voltage_nominal",
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

func recordMetrics(upscBinary string, upsArg string){
	prometheus.MustRegister(batteryCharge)
	prometheus.MustRegister(batteryPacks) 
	prometheus.MustRegister(batteryVoltage)
	prometheus.MustRegister(batteryVoltageNominal)
	prometheus.MustRegister(inputVoltage)
	prometheus.MustRegister(inputVoltageNominal)
	prometheus.MustRegister(outputVoltage)
	prometheus.MustRegister(outputVoltageNominal)
	prometheus.MustRegister(upsPowerNominal)
	prometheus.MustRegister(upsTemp)
	prometheus.MustRegister(upsLoad)
	prometheus.MustRegister(upsStatus)
	
	go func(){
		for {
			upsOutput, err := exec.Command(upscBinary , upsArg).Output()
			
			if err != nil {
				log.Fatal(err)
			}
			
			batteryChargeValue, _ := strconv.ParseFloat(batteryChargeRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			batteryCharge.Set(batteryChargeValue)

			batteryPacksValue, _ := strconv.ParseFloat(batteryPacksRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			batteryPacks.Set(batteryPacksValue)

			batteryVoltageValue, _ := strconv.ParseFloat(batteryVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			batteryVoltage.Set(batteryVoltageValue)

			batteryVoltageNominalValue, _ := strconv.ParseFloat(batteryVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			batteryVoltageNominal.Set(batteryVoltageNominalValue)

			inputVoltageValue, _ := strconv.ParseFloat(inputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			inputVoltage.Set(inputVoltageValue)

			inputVoltageNominalValue, _ := strconv.ParseFloat(inputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			inputVoltageNominal.Set(inputVoltageNominalValue)

			outputVoltageValue, _ := strconv.ParseFloat(outputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			outputVoltage.Set(outputVoltageValue)

			outputVoltageNominalValue, _ := strconv.ParseFloat(outputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			outputVoltageNominal.Set(outputVoltageNominalValue)

			upsPowerNominalValue, _ := strconv.ParseFloat(upsPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			upsPowerNominal.Set(upsPowerNominalValue)

			upsTempValue, _ := strconv.ParseFloat(upsTempRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			upsTemp.Set(upsTempValue)

			upsLoadValue, _ := strconv.ParseFloat(upsLoadRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
			upsLoad.Set(upsLoadValue)

			upsStatusValue := upsStatusRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1]	
			
			switch upsStatusValue {
				case "CAL":
            		upsStatus.Set(0)
				case "TRIM":
					upsStatus.Set(1)
				case "BOOST":
					upsStatus.Set(2)
				case "OL":
					upsStatus.Set(3)
				case "OB":
					upsStatus.Set(4)
				case "OVER":
					upsStatus.Set(5)
				case "LB":
					upsStatus.Set(6)
				case "RB":
					upsStatus.Set(7)
				case "BYPASS":
					upsStatus.Set(8)
				case "OFF":
					upsStatus.Set(9)
				case "CHRG":
					upsStatus.Set(10)
				case "DISCHRG":
					upsStatus.Set(11)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func main() {
	upsArg   := flag.String("ups", "none", "ups name managed by nut")
	portArg  := flag.Int("port", 8100, "port number")
	upscArg  := flag.String("upsc", "/bin/upsc", "upsc path")

	var listenAddr = fmt.Sprintf(":%d", *portArg)

	flag.Parse()
	recordMetrics(*upscArg, *upsArg)
    
	log.Infoln("Starting NUT exporter on ups", *upsArg )	
	http.Handle("/metrics", promhttp.Handler())
    
	log.Infoln("NUT exporter started on port", *portArg)
	http.ListenAndServe(listenAddr, nil)	
}
