package zcmu

import (
	"fmt"
	libSvm "github.com/ewalker544/libsvm-go"
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
	fmt.Println("=== OCR model load in ", modelPath)
	svmModel = libSvm.NewModelFromFile(modelPath)
}

// recognize verify code.
func picture2vector(picture image.Image) *[][]float64 {
	vec := make([][]float64, 4) // 4个字符 分为4个向量
	index := 0
	for i := 2; i < 50; i += 12 {
		for y := 1; y < 22; y++ {
			for x := 0; x <= 16; x++ {
				pixel := picture.At(x+i, y)
				r, g, b, _ := pixel.RGBA()
				y := (0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b)) / 257.0
				vec[index] = append(vec[index], y/255.0)
			}
		}
		index++
	}
	return &vec
}

func Predict(picture image.Image) (string, error) {
	vector := picture2vector(picture)
	result := make([]byte, 0)
	label := make(map[int]float64)
	for i := 0; i < 4; i++ {
		for index, value := range (*vector)[i] {
			label[index+1] = value
		}
		result = append(result, byte(svmModel.Predict(label)))
	}
	return string(result), nil
}
