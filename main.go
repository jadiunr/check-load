package main

import (
	"fmt"
    "strconv"
    "strings"

	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Critical	string
	Warning		string
    PerCPU      bool
	MetricsOnly bool
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-load",
			Short:    "Check load averages and provide metrics",
			Keyspace: "sensu.io/plugins/check-load/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "critical",
			Argument:  "critical",
			Shorthand: "c",
			Default:   "0.85,0.8,0.75",
			Usage:     "Critical threshold for load averages",
			Value:     &plugin.Critical,
		},
		{
			Path:      "warning",
			Argument:  "warning",
			Shorthand: "w",
			Default:   "0.75,0.7,0.65",
			Usage:     "Warning threshold for load averages",
			Value:     &plugin.Warning,
		},
        {
            Path:      "percpu",
            Argument:  "percpu",
            Shorthand: "r",
            Usage:     "Divide the load averages by the number of CPUs",
            Value:     &plugin.PerCPU,
        },
		{
			Path:	   "metricsonly",
			Argument:  "metricsonly",
			Shorthand: "m",
			Default:   false,
			Usage:     "Outputs only the metrics without checking the threshold.",
            Value:     &plugin.MetricsOnly,
		},
	}

    cload [3]float64
    wload [3]float64
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
    for i, thStr := range strings.Split(plugin.Critical, ",") {
        th, err := strconv.ParseFloat(thStr, 64)
        if err != nil {
            return sensu.CheckStateWarning, fmt.Errorf("Could not parse the critical threshold.")
        }
        cload[i] = th
    }
    for i, thStr := range strings.Split(plugin.Warning, ",") {
        th, err := strconv.ParseFloat(thStr, 64)
        if err != nil {
            return sensu.CheckStateWarning, fmt.Errorf("Could not parse the warning threshold.")
        }
        wload[i] = th
    }
	for i, w := range wload {
        if w > cload[i] {
            return sensu.CheckStateWarning, fmt.Errorf("--warning cannot be greater than --critical")
        }
    }
    return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	info, err := cpu.Info()
    if err != nil {
        return sensu.CheckStateCritical, fmt.Errorf("Error: obtaining CPU info: %v", err)
    } else if len(info) < 1 {
        return sensu.CheckStateCritical, fmt.Errorf("Error: no CPU info found")
    }
    cores := float64(0)
    for _, i := range info {
        cores += float64(i.Cores)
    }
    var loadStats *load.AvgStat
    for i := 1; i < 10; i++ {
        loadStats, err = load.Avg()
    }
    if err != nil {
        return sensu.CheckStateCritical, fmt.Errorf("Error: obtaining load stats: %v", err)
    }

    load := []float64{loadStats.Load1, loadStats.Load5, loadStats.Load15}
    loadPerCPU := []float64{load[0] / cores, load[1] / cores, load[2] / cores}
    perfData := fmt.Sprintf(
        "load1=%.2f, load5=%.2f, load15=%.2f, load1_per_cpu=%.2f, load5_per_cpu=%.2f, load15_per_cpu=%.2f",
        load[0], load[1], load[2],
        loadPerCPU[0], loadPerCPU[1], loadPerCPU[2],
    )
    var targetLoad = load
    if plugin.PerCPU {
        targetLoad = loadPerCPU
    }
    if !plugin.MetricsOnly {
        for i, l := range targetLoad {
            if l > cload[i] {
                fmt.Printf("%s Critical - load average: %s | %s\n", plugin.PluginConfig.Name, fmt.Sprintf("%.2f, %.2f, %.2f", targetLoad[0], targetLoad[1], targetLoad[2]), perfData)
                return sensu.CheckStateCritical, nil
            }
            if l > wload[i] {
                fmt.Printf("%s Warning - load average: %s | %s\n", plugin.PluginConfig.Name, fmt.Sprintf("%.2f, %.2f, %.2f", targetLoad[0], targetLoad[1], targetLoad[2]), perfData)
                return sensu.CheckStateWarning, nil
            }
        }
    }

    fmt.Printf("%s OK - load average: %s | %s\n", plugin.PluginConfig.Name, fmt.Sprintf("%.2f, %.2f, %.2f", targetLoad[0], targetLoad[1], targetLoad[2]), perfData)
    return sensu.CheckStateOK, nil
}
