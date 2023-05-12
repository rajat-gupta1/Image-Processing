package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
	"proj2/png"
	"strings"
)

type Request struct {
	InPath  string
	OutPath string
	Effects []string
}

func RunSequential(config Config) {
	file, err := os.Open("data/effects.txt")
	if err != nil {
		fmt.Println(err)
	}
	reader := json.NewDecoder(file)
	defer file.Close()

	directories := strings.Split(config.DataDirs, "+")
	for reader.More() {
		req := Request{}
		err := reader.Decode(&req)

		if err != nil {
			print(err)
			return
		}

		for _, directory := range directories {
			filePath := "data/in/" + directory + "/" + req.InPath
			pngImg, err := png.Load(filePath)
			if err != nil {
				print(err)
			}

			// Applying effects
			for j, effect := range req.Effects {
				if j > 0 {
					pngImg.Inout()
				}
				if effect == "G" {
					pngImg.Grayscale(-1, -1)
				} else if effect == "E" {
					pngImg.Edge_Det(-1, -1)
				} else if effect == "S" {
					pngImg.Sharpen(-1, -1)
				} else if effect == "B" {
					pngImg.Blur(-1, -1)
				}
				
			}

			outfilePath := "data/out/" + directory + "_" + req.OutPath
			err = pngImg.Save(outfilePath)

			//Checks to see if there were any errors when saving.
			if err != nil {
				panic(err)
			}
		}
	}
}
