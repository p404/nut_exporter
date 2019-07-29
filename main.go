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
	upsArg   = flag.String("ups", "none", "ups name managed by nut")
	portArg  = flag.Int("port", 8100, "port number")
	upscArg  = flag.String("upsc", "/bin/upsc", "upsc path")

	batteryChargeRegex          =   regexp.MustCompile(`(?:battery[.]charge:(?:\s)(.*))`)
	batteryPacksRegex           =   regexp.MustCompile(`(?:battery[.]packs:(?:\s)(.*))`)
	batteryVoltageRegex         =   regexp.MustCompile(`(?:battery[.]voltage:(?:\s)(.*))`)
	batteryVoltageNominalRegex  =   regexp.MustCompile(`(?:battery[.]voltage[.]nominal:(?:\s)(.*))`)
	batteryRuntimeRegex         =   regexp.MustCompile(`(?:battery[.]runtime:(?:\s)(.*))`)
	batteryRuntimeLowRegex      =   regexp.MustCompile(`(?:battery[.]runtime[.]low:(?:\s)(.*))`)
	inputVoltageRegex           =   regexp.MustCompile(`(?:input[.]voltage:(?:\s)(.*))`)
	inputVoltageNominalRegex    =   regexp.MustCompile(`(?:input[.]voltage[.]nominal:(?:\s)(.*))`)
	outputVoltageRegex          =   regexp.MustCompile(`(?:output[.]voltage:(?:\s)(.*))`)
	outputVoltageNominalRegex   =   regexp.MustCompile(`(?:output[.]voltage[.]nominal:(?:\s)(.*))`)
	upsPowerNominalRegex        =   regexp.MustCompile(`(?:ups[.]power[.]nominal:(?:\s)(.*))`)
	upsTempRegex                =   regexp.MustCompile(`(?:ups[.]temperature:(?:\s)(.*))`)
	upsLoadRegex                =   regexp.MustCompile(`(?:ups[.]load:(?:\s)(.*))`)
	upsStatusRegex              =   regexp.MustCompile(`(?:ups[.]status:(?:\s)(.*))`)


	batteryCharge = prometheus.NewGauge(prometheus.GaugeOpts{})
	batteryPacks= prometheus.NewGauge(prometheus.GaugeOpts{})
	batteryVoltage= prometheus.NewGauge(prometheus.GaugeOpts{})
	batteryVoltageNominal= prometheus.NewGauge(prometheus.GaugeOpts{})
	batteryRuntime= prometheus.NewGauge(prometheus.GaugeOpts{})
	batteryRuntimeLow= prometheus.NewGauge(prometheus.GaugeOpts{})
	inputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{})
	inputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{})
	outputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{})
	outputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{})
	upsPowerNominal = prometheus.NewGauge(prometheus.GaugeOpts{})
	upsTemp = prometheus.NewGauge(prometheus.GaugeOpts{})
	upsLoad = prometheus.NewGauge(prometheus.GaugeOpts{})
	upsStatus = prometheus.NewGauge(prometheus.GaugeOpts{})

)



func initMetrics( upsArg string) {

	constLabels := map[string]string{ "ups" : upsArg }
	batteryCharge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_charge",
		Help: "Current battery charge (percent)",
		ConstLabels: constLabels,
	})

	batteryPacks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_pack",
		Help: "Number of battery packs on the UPS",
		ConstLabels: constLabels,
	})

	batteryVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_voltage",
		Help: "Current battery voltage",
		ConstLabels: constLabels,
	})

	batteryVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_voltage_nominal",
		Help: "Nominal battery voltage",
		ConstLabels: constLabels,
	})
	batteryRuntime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_runtime",
		Help: "Battery runtime",
		ConstLabels: constLabels,
	})
	batteryRuntimeLow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "battery_runtime_low",
		Help: "Battery runtime low",
		ConstLabels: constLabels,
	})
	inputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "input_voltage",
		Help: "Current input voltage",
		ConstLabels: constLabels,
	})

	inputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "input_voltage_nominal",
		Help: "Nominal input voltage",
		ConstLabels: constLabels,
	})

	outputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "output_voltage",
		Help: "Current output voltage",
		ConstLabels: constLabels,
	})

	outputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "output_voltage_nominal",
		Help: "Nominal output voltage",
		ConstLabels: constLabels,
	})

	upsPowerNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_power_nominal",
		Help: "Nominal ups power",
		ConstLabels: constLabels,
	})

	upsTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_temp",
		Help: "UPS Temperature (degrees C)",
		ConstLabels: constLabels,
	})

	upsLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_load",
		Help: "Current UPS load (percent)",
		ConstLabels: constLabels,
	})

	upsStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ups_status",
		Help: "Current UPS Status (0=Calibration, 1=SmartTrim, 2=SmartBoost, 3=Online, 4=OnBattery, 5=Overloaded, 6=LowBattery, 7=ReplaceBattery, 8=OnBypass, 9=Off, 10=Charging, 11=Discharging)",
		ConstLabels: constLabels,
	})

}

func recordMetrics(){
	prometheus.MustRegister(batteryCharge)
	prometheus.MustRegister(batteryPacks)
	prometheus.MustRegister(batteryVoltage)
	prometheus.MustRegister(batteryVoltageNominal)
	prometheus.MustRegister(batteryRuntime)
	prometheus.MustRegister(batteryRuntimeLow)
	prometheus.MustRegister(inputVoltage)
	prometheus.MustRegister(inputVoltageNominal)
	prometheus.MustRegister(outputVoltage)
	prometheus.MustRegister(outputVoltageNominal)
	prometheus.MustRegister(upsPowerNominal)
	prometheus.MustRegister(upsTemp)
	prometheus.MustRegister(upsLoad)
	prometheus.MustRegister(upsStatus)
}
func sampleMetrics(upscBinary string, upsArg string) {
	go func(){
		for {
			upsOutput, err := exec.Command(upscBinary , upsArg).Output()

			if err != nil {
				log.Fatal(err)
			}

			if batteryChargeRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryCharge)
			} else {
				batteryChargeValue, _ := strconv.ParseFloat(batteryChargeRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryCharge.Set(batteryChargeValue)
			}

			if batteryPacksRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryPacks)
			} else {
				batteryPacksValue, _ := strconv.ParseFloat(batteryPacksRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryPacks.Set(batteryPacksValue)
			}

			if batteryVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryVoltage)
			} else {
				batteryVoltageValue, _ := strconv.ParseFloat(batteryVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryVoltage.Set(batteryVoltageValue)
			}

			if batteryVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryVoltageNominal)
			} else {
				batteryVoltageNominalValue, _ := strconv.ParseFloat(batteryVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryVoltageNominal.Set(batteryVoltageNominalValue)
			}
			if batteryRuntimeRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryRuntime)
			} else {
				batteryRuntimeValue, _ := strconv.ParseFloat(batteryRuntimeRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryRuntime.Set(batteryRuntimeValue)
			}


			if batteryRuntimeLowRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryRuntimeLow)
			} else {
				batteryRuntimeLowValue, _ := strconv.ParseFloat(batteryRuntimeLowRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryRuntimeLow.Set(batteryRuntimeLowValue)
			}

			if inputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(inputVoltage)
			} else {
				inputVoltageValue, _ := strconv.ParseFloat(inputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				inputVoltage.Set(inputVoltageValue)
			}

			if inputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(inputVoltageNominal)
			} else {
				inputVoltageNominalValue, _ := strconv.ParseFloat(inputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				inputVoltageNominal.Set(inputVoltageNominalValue)
			}

			if outputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(outputVoltage)
			} else {
				outputVoltageValue, _ := strconv.ParseFloat(outputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
  				outputVoltage.Set(outputVoltageValue)
			}

			if outputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(outputVoltageNominal)
			} else {
				outputVoltageNominalValue, _ := strconv.ParseFloat(outputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				outputVoltageNominal.Set(outputVoltageNominalValue)
			}

			if upsPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsPowerNominal)
			} else {
				upsPowerNominalValue, _ := strconv.ParseFloat(upsPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				upsPowerNominal.Set(upsPowerNominalValue)
			}

			if upsTempRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsTemp)
			} else {
				upsTempValue, _ := strconv.ParseFloat(upsTempRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				upsTemp.Set(upsTempValue)
			}

			if upsLoadRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsLoad)
			} else {
				upsLoadValue, _ := strconv.ParseFloat(upsLoadRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				upsLoad.Set(upsLoadValue)
			}

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

  flag.Parse()

	var listenAddr = fmt.Sprintf(":%d", *portArg)
	initMetrics(*upsArg)
	recordMetrics()
	sampleMetrics(*upscArg, *upsArg)

	log.Infoln("Starting NUT exporter on ups", *upsArg )
	http.Handle("/metrics", promhttp.Handler())

	log.Infoln("NUT exporter started on port", *portArg)
	http.ListenAndServe(listenAddr, nil)
}
