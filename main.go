package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/wcharczuk/go-chart"
)

// Line is slice of dots on chart
type Line struct {
	XValues []float64 // X axis values
	YValues []float64 // Y axis values
}

func ReadCsvFile(path string) (map[string]Line, error) {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(csvFile)
	requests := make(map[string]Line)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) != 3 {
			return nil, err
		}

		methodName := record[0]
		veCount := record[1]
		maxRPS := record[2]

		line := requests[methodName]
		xValue, err := strconv.ParseFloat(veCount, 64)
		if err != nil {
			return nil, err
		}
		// VE count
		line.XValues = append(line.XValues, xValue)
		yValue, err := strconv.ParseFloat(maxRPS, 64)
		if err != nil {
			return nil, err
		}
		// Max RPS
		line.YValues = append(line.YValues, yValue)
		requests[methodName] = line
	}
	return requests, nil
}

func main() {
	inputCSVPath := flag.String("csv", "requests.csv", " Path to cvs with ReqType,N,RPS.")
	outputPNGPath := flag.String("png", "output.png", " Output png file.")
	flag.Parse()

	lines, err := ReadCsvFile(*inputCSVPath)
	if err != nil {
		log.Fatal("Couldn't read and parse requests", err)
	}

	err = RenderChart(lines, *outputPNGPath)
	if err != nil {
		log.Fatal("Couldn't render chart", err)
	}
}

func RenderChart(requests map[string]Line, fileName string) error {
	var series []chart.Series
	var colorIndex int
	for key, value := range requests {
		series = append(series, chart.ContinuousSeries{
			Name: key,
			Style: chart.Style{
				StrokeColor: chart.GetDefaultColor(colorIndex).WithAlpha(255),
				StrokeWidth: 5,
			},
			XValues: value.XValues,
			YValues: value.YValues,
		})
		colorIndex++
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "VE count",
		},
		YAxis: chart.YAxis{
			Name: "Max RPS",
		},
		Series: series,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	graph.Render(chart.PNG, file)
	err = graph.Render(chart.PNG, file)

	return err
}
