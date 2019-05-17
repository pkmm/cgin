package zcmu

import (
	"fmt"
	"github.com/ewalker544/libsvm-go"
	"image"
	"os"
	"path"
)

var (
	svmModel *libSvm.Model
)

func init() {
	wd, _ := os.Getwd()
	modelPath := path.Join(wd, "zcmu", "zf.model")
	fmt.Println("zcmu model load in ", modelPath)
	svmModel = libSvm.NewModelFromFile(modelPath)
}

// recognize verify code.
func crop(src image.Image, name string) map[string][]float64 {
	vec := make(map[string][]float64, 0)
	rgbImg := src
	index := 0
	for i := 2; i < 50; i += 12 {
		var tmp []float64
		for y := 1; y < 22; y++ {
			for x := 0; x <= 16; x++ {
				pixel := rgbImg.At(x+i, y)
				r, g, b, _ := pixel.RGBA()
				y := float64(0.3*float64(r)+0.59*float64(g)+0.11*float64(b)) / 257.0
				tmp = append(tmp, y/255.0)
			}
		}
		vec[fmt.Sprintf("%s-%d", name, index)] = tmp
		index++
	}
	return vec
}

func Predict(im image.Image) (string, error) {
	vec := crop(im, "loc")
	ret := make([]byte, 0)
	x := make(map[int]float64)
	for ind := 0; ind < 4; ind++ {
		for index, value := range vec[fmt.Sprintf("loc-%d", ind)] {
			x[index+1] = value
		}
		predictLabel := svmModel.Predict(x)
		ans := byte(predictLabel)
		ret = append(ret, ans)
	}
	return string(ret), nil
}
