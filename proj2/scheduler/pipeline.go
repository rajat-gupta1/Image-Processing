package scheduler

import (
	"strings"
	"os"
	"proj2/png"
	"encoding/json"
)

func RunPipeline(config Config) {

	threads := config.ThreadCount
	done := make(chan interface{})
	defer close(done)
	
	workerReturn := make(chan bool)
	defer close(workerReturn)
	
	images := LoadImage(done, config)

	for i:= 0; i < threads; i++ {
		go Worker(done, images, threads, config, workerReturn)
	}
	
	finished := 0
	
	for{
		select {
		case <- workerReturn:
			finished += 1
			if finished == threads {
				return
			}
		}
	}
}

func LoadImage(done <- chan interface{}, config Config) <- chan *png.Image {
	file, err := os.Open("data/effects.txt")
	if err != nil {
		panic(err)
	}
	reader := json.NewDecoder(file)
	// defer file.Close()

	directories := strings.Split(config.DataDirs, "+")

	images := make(chan *png.Image)
	go func () {
		defer close(images)
		for reader.More() {
			
			req := Request{}
			err := reader.Decode(&req)
			if err != nil {
				panic(err)
			}

			for _, directory := range directories {
				
				filePath := "data/in/" + directory + "/" + req.InPath
				pngImg, err := png.Load(filePath)
				if err != nil {
					panic(err)
				}
				outfilePath := "data/out/" + directory + "_" + req.OutPath
				pngImg.Effects = req.Effects
				pngImg.OutPath = outfilePath

				select {
				case <- done:
					return 
				case images <- pngImg:
				}
			}
		}
	}()

	return images
}

func Worker(done <- chan interface{}, images <- chan *png.Image, threads int, config Config, workerReturn chan bool) {

	for thisImg := range images {
		push := make (chan string)
		signal := make (chan bool)
		stopGo := make (chan interface{})
		imgSlices := make (chan *png.Image)

		sliceSize := thisImg.Bounds.Max.Y / threads

		for i:= 0; i < threads; i++ {
			YStart := i * sliceSize
			YEnd := YStart + sliceSize

			if i == threads - 1 {
				YEnd = thisImg.Bounds.Max.Y
			}
			
			go miniWorker(done, push, signal, imgSlices, stopGo, YStart, YEnd)
			imgSlices <- thisImg
		}

		for j,  effect := range thisImg.Effects {
			if j > 0 {
				thisImg.Inout()
			}
			for i:= 0; i < threads; i++ {
				push <- effect
			}

			finished := 0
			for i:= 0; i < threads; i++ {
				select {
				case <- done:
					return 
				case <- signal:
					finished += 1
					if finished == threads {
						break
					}
				}
			}
		}
		thisImg.Save(thisImg.OutPath)
		for i:= 0; i < threads; i++ {
			stopGo <- true
		}
		close (push)
		close (signal)
	}
	workerReturn <- true
}

func miniWorker(done <- chan interface{}, push <- chan string, signal chan <- bool, imgSlices <- chan *png.Image, stopGo chan interface{}, YStart int, YEnd int) {

	imgSlice := <-imgSlices

	for {
		select {
		case <- done:
			return
		case effect := <-push:
			if effect == "G" {
				imgSlice.Grayscale(-1, -1)
			} else if effect == "E" {
				imgSlice.Edge_Det(-1, -1)
			} else if effect == "S" {
				imgSlice.Sharpen(-1, -1)
			} else if effect == "B" {
				imgSlice.Blur(-1, -1)
			}
			signal <- true
		case <- stopGo:
			return 
		}
	}
}