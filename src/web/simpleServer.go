package main

import (
	"net/http"
	"log"
	"fmt"
	"sync"
	"image/color"
	"math/rand"
	"image/gif"
	"image"
	"math"
	"strconv"
)

var mu sync.Mutex
var count int

var palette = []color.Color{
	color.White,
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0xff, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0xff, 0xff},
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/details", requestDetails)
	http.HandleFunc("/gif", lissajous)
	fmt.Println("Server started: http://localhost:8000")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func requestDetails(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "%s %s %s\n", request.Method, request.URL, request.Proto)
	for k, v := range request.Header {
		fmt.Fprintf(writer, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(writer, "Host = %q\n", request.Host)
	fmt.Fprintf(writer, "RemoteAddr = %q\n", request.RemoteAddr)
	if err := request.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range request.Form {
		fmt.Fprintf(writer, "Form[%q] = %q\n", k, v)
	}
}

func counter(writer http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(writer, "Count %d\n", count)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
}

func lissajous(writer http.ResponseWriter, request *http.Request) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: queryParamOrDefaultInt(request,"nframes", nframes)}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		colorIndex := uint8(rand.Int()%3+1)
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		cyclesCalculated := queryParamOrDefaultFloat64(request, "cycles", cycles) * 2 * math.Pi
		for t := 0.0; t < cyclesCalculated; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5),
				size+int(y*size+0.5),
				colorIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(writer, &anim) // NOTE: ignoring encoding errors
}

func queryParamOrDefaultInt(request *http.Request, queryParam string, defaultValue int) int {
	queryParamValue := request.URL.Query().Get(queryParam)
	if queryParamValue == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(queryParamValue)

	if err == nil || result == 0{
		return defaultValue
	} else {
		return result
	}
}

func queryParamOrDefaultFloat64(request *http.Request, queryParam string, defaultValue float64) float64 {
	queryParamValue := request.URL.Query().Get(queryParam)
	if queryParamValue == "" {
		return defaultValue
	}

	result, err := strconv.ParseFloat(queryParamValue, 64)

	if err == nil || result == 0{
		return defaultValue
	} else {
		return result
	}
}
