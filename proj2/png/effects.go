// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale(YStart int, YEnd int) {
	bounds := img.Out.Bounds()

	// In case we need to use the whole image
	if YStart == -1 {
		YStart = bounds.Min.Y
		YEnd = bounds.Max.Y
	}	
	for y := YStart; y < YEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.In.At(x, y).RGBA()

			// Ensuring that the values are between specified ranges
			greyC := clamp(float64(r+g+b) / 3)
			img.Out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

// Blur applies a blurring effect to the image
func (img *Image) Blur(YStart int, YEnd int) {
	bounds := img.Out.Bounds()
	if YStart == -1 {
		YStart = bounds.Min.Y
		YEnd = bounds.Max.Y
	}
	
	kernel := [9]float64{1 / 9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
	for y := YStart; y < YEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			kernel_and_store(x, y, kernel, img)
		}
	}
}

// Sharpen applies a sharpening effect to the image
func (img *Image) Sharpen(YStart int, YEnd int) {
	bounds := img.Out.Bounds()
	if YStart == -1 {
		YStart = bounds.Min.Y
		YEnd = bounds.Max.Y
	}
	kernel := [9]float64{0, -1, 0, -1, 5, -1, 0, -1, 0}
	for y := YStart; y < YEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			kernel_and_store(x, y, kernel, img)
		}
	}
}

// Edge_Det detects edges in an image
func (img *Image) Edge_Det(YStart int, YEnd int) {
	bounds := img.Out.Bounds()
	if YStart == -1 {
		YStart = bounds.Min.Y
		YEnd = bounds.Max.Y
	}
	
	kernel := [9]float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
	for y := YStart; y < YEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			kernel_and_store(x, y, kernel, img)
		}
	}
}

// The function is a helper function to multiply and store new pixel values
func kernel_and_store(x int, y int, kernel [9]float64, img *Image) {
	var blur_r float64 = 0
	var blur_g float64 = 0
	var blur_b float64 = 0
	var blur_a uint16 = 0
	ctr := -1
	for i := x - 1; i < x+2; i++ {
		for j := y - 1; j < y+2; j++ {

			// Using ctr for kernels
			ctr += 1
			r, g, b, a := img.In.At(i, j).RGBA()
			blur_r += float64(r) * kernel[ctr]
			blur_g += float64(g) * kernel[ctr]
			blur_b += float64(b) * kernel[ctr]
			blur_a = uint16(a)
		}
	}

	ublur_r := clamp(blur_r)
	ublur_g := clamp(blur_r)
	ublur_b := clamp(blur_r)

	img.Out.Set(x, y, color.RGBA64{ublur_r, ublur_g, ublur_b, uint16(blur_a)})
}
