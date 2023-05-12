package scheduler

import (
	"sync"
	"encoding/json"
	"proj2/png"
	"os"
	"strings"
)

type bspWorkerContext struct {

	// For locking while performing operations
	Cond *sync.Cond
	Mu sync.Mutex
	Done bool 

	Threads int
	ThreadsWaiting int
	DirectoryIdx int
	EffectCount int
	Directories []string

	Reader *json.Decoder
	Request *Request
	PngImg *png.Image
}

func NewBSPContext(config Config) *bspWorkerContext {
	ctx := &bspWorkerContext{}
	ctx.Cond = sync.NewCond(&ctx.Mu)
	ctx.Done = false

	ctx.Threads = config.ThreadCount
	ctx.ThreadsWaiting = 0
	
	file, _ := os.Open("data/effects.txt")
	ctx.Reader = json.NewDecoder(file)

	ctx.Request = &Request{Effects: []string{}}

	ctx.Directories = strings.Split(config.DataDirs, "+")
	ctx.DirectoryIdx = 0
	ctx.EffectCount = 0

	return ctx
}

func RunBSPWorker(id int, ctx *bspWorkerContext) {
	for {
		ctx.Mu.Lock()
		ctx.ThreadsWaiting += 1

		// If this thread is the last thread, we want this to instruct other threads
		if ctx.ThreadsWaiting == ctx.Threads {
			ctx.ThreadsWaiting = 0

			// If this is the first effect in the image (Its a new image)
			if ctx.EffectCount == 0{

				// If this is the first type in the directory
				if ctx.DirectoryIdx == 0 {

					// In case we are getting more images
					if ctx.Reader.More() {
						ctx.Request = &Request{}
						ctx.Reader.Decode(ctx.Request)
						filePath := "data/in/" + ctx.Directories[ctx.DirectoryIdx] + "/" + ctx.Request.InPath
						ctx.PngImg, _ = png.Load(filePath)

						// Waking up workers to do the tasks
						ctx.Cond.Broadcast()

					// We are done with all the images
					} else {
						ctx.Done = true
						ctx.Cond.Broadcast()
						ctx.Mu.Unlock()
						return
					}

					// In case there were multiple type of files in the directory
				} else if ctx.DirectoryIdx < len(ctx.Directories) {
					filePath := "data/in/" + ctx.Directories[ctx.DirectoryIdx] + "/" + ctx.Request.InPath
					ctx.PngImg, _ = png.Load(filePath)

					ctx.Cond.Broadcast()
				}

			// There were multiple effects for this image, instructing threads to finish the remaining effects
			} else {
				ctx.Cond.Broadcast()
			}
		} else {

			// In case it wasnt the last thread
			ctx.Cond.Wait()

			if ctx.Done{
				ctx.Mu.Unlock()
				return
			}
		}
		ctx.Mu.Unlock()

		// Each of the threads work on a part of each image
		stepSize := ctx.PngImg.Bounds.Max.Y / ctx.Threads
		YStart := id * stepSize
		YEnd := YStart + stepSize
		
		// In case its the last thread
		if id == ctx.Threads - 1 {
			YEnd = ctx.PngImg.Bounds.Max.Y
		}

		// Which function to call depends on the effect requested
		if ctx.Request.Effects[ctx.EffectCount] == "G" {
			ctx.PngImg.Grayscale(YStart, YEnd)
		} else if ctx.Request.Effects[ctx.EffectCount] == "E" {
			ctx.PngImg.Edge_Det(YStart, YEnd)
		} else if ctx.Request.Effects[ctx.EffectCount] == "S" {
			ctx.PngImg.Sharpen(YStart, YEnd)
		} else if ctx.Request.Effects[ctx.EffectCount] == "B" {
			ctx.PngImg.Blur(YStart, YEnd)
		}

		// Locking to update the threads count and writing to the file/getting the image ready for the next effect
		ctx.Mu.Lock()
		ctx.ThreadsWaiting += 1

		// In Case its the last thread
		if ctx.ThreadsWaiting == ctx.Threads {

			// Resetting the counter
			ctx.ThreadsWaiting = 0
			ctx.EffectCount += 1

			// All the effects have been applied
			if ctx.EffectCount == len(ctx.Request.Effects) {
				outfilePath := "data/out/" + ctx.Directories[ctx.DirectoryIdx] + "_" + ctx.Request.OutPath 
				err := ctx.PngImg.Save(outfilePath)

				if err != nil {
					panic(err)
				}

				ctx.EffectCount = 0
				ctx.DirectoryIdx += 1

				if ctx.DirectoryIdx == len(ctx.Directories) {
					ctx.DirectoryIdx = 0
				}
				// Continuing with the same image and different effect
			} else {
				ctx.PngImg.Inout()
			}
			ctx.Cond.Broadcast()
		} else {

			// Waiting till the last thread finishes writing on file or get the image ready for the next effect
			ctx.Cond.Wait()
		}
		ctx.Mu.Unlock()
	}
}
