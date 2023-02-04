package main

//this program should perform a padding-operation on images (jpg- and png-)
//here this means adding padding in one axis, rendering the original image in the center
//it looks for narrow pictures (stretched out in format like i.e. 1:2)
//the goal is to give us/me a square picture, from this narrow picture

//TODO: if the narrow image is png and seems to have a transparent background, make the new image likewise
import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var files []string

	/*settings concerning images*/
	//this is the factor by which the (image) long side should exceed its short side in order to trigger the padding-operation
	var maxRatio float32 = 2

	var minRatio float32 = 1 / maxRatio // this is the ratio from an inverse division of the dimensions x y

	var userMaxRatioDefault float64 = 2.0
	userMaxRatio := flag.Float64("ratiolimit", userMaxRatioDefault, "A decimal value that defines the upper limit of an images' ratio")

	flag.Parse()

	fmt.Printf("Images ratio fix. Command line arguments:\n\n")
	fmt.Printf("'-ratiolimit' \tfloat\texample: -ratiolimit=3.1\n")
	fmt.Printf("deefault value: %.2f\n", userMaxRatioDefault)
	fmt.Printf("----------------------\n")
	fmt.Printf("userMaxRatio: %.2f\n", *userMaxRatio)

	root := "./"
	fmt.Printf("Images and their ratio between long and short side, maxRatio alowed is %.2f\nLook in folder: %s\n----------------------\n", maxRatio, root)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".dat" { //example code
			return nil
		}
		files = append(files, path)
		//files = append(files, info.Name()) //example code
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)

		//image handling
		existingImageFile, err := os.Open(file)
		if err != nil {
			panic(err)
		} else {

			//fmt.Printf("\nexistingImageFile: %v\n", existingImageFile)

			_, err := existingImageFile.Seek(0, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", existingImageFile.Name(), err)
			}

			//fmt.Println(http.DetectContentType(existingImageFile)) //didnt work, cant get http package

			// Calling the generic image.Decode() will tell give us the data
			// and type of image it is as a string. We expect "png"
			imageData, imageType, err := image.Decode(existingImageFile)
			if err != nil {
				//fmt.Printf("error, this file (%s) could not be handled as an image", existingImageFile.Name())
				//fmt.Fprintf(os.Stderr, "%s: %v\n", existingImageFile.Name(), err)
				defer existingImageFile.Close()
				continue
			}
			fmt.Println(imageData.ColorModel())
			fmt.Println(imageType)

			existingImageFile.Seek(0, 0)

			im, _, err := image.DecodeConfig(existingImageFile)
			if err != nil {
				panic(err)
			}
			fmt.Printf("width: %d height: %d\n", im.Width, im.Height)

			var ratio float32 = float32(im.Width) / float32(im.Height)
			fmt.Printf("ratio: %.2f\n", ratio)
			if ratio > maxRatio || ratio < minRatio {
				fmt.Printf(">>>>>>>>>Image is exceeding maxRatio! (in either width or height)<<<<<<<<<\n")

				var biggest int = 0
				var paddX bool
				if ratio > 1 {
					biggest = im.Width
					paddX = false
				} else {
					biggest = im.Height
					paddX = true
				}
				m := image.NewRGBA(image.Rect(0, 0, biggest, biggest)) //a square image
				//blue := color.RGBA{0, 0, 255, 255}
				white := color.White
				draw.Draw(m, m.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

				var dp image.Point //have a start point in the destination
				dp.X = 0
				dp.Y = 0
				if paddX {
					//starting x must be away from destination img left edge
					dp.X = (biggest / 2) - im.Width/2
				} else {
					//away from dest. img top edge
					dp.Y = (biggest / 2) - im.Height/2
				}
				var sr image.Rectangle
				sr.Min.X = 0
				sr.Min.Y = 0
				sr.Max.X = im.Width
				sr.Max.Y = im.Height
				r := image.Rectangle{dp, dp.Add(sr.Size())}
				draw.Draw(m, r, imageData, sr.Min, draw.Src) //draw to in-memory image

				newFileName := root + fileNameWithoutExtension(existingImageFile.Name()) + "_squared"

				if imageType == "jpeg" {
					newFileName += ".jpg"
				} else if imageType == "png" {
					newFileName += ".png"
				} else {
					continue
				}

				f, err := os.Create(newFileName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %v\n", existingImageFile.Name(), err)
				}
				defer f.Close()

				// Specify the quality, between 0-100
				// Higher is better

				if imageType == "jpeg" {
					opt := jpeg.Options{
						Quality: 90,
					}
					saveJpeg(f, m, opt)
				} else if imageType == "png" {
					savePng(f, m)
				}
			}
		}
		defer existingImageFile.Close()
		fmt.Printf("--------------------------------------\n")
	}

}

func saveJpeg(osFile *os.File, imageRGBA *image.RGBA, options jpeg.Options) error {
	fmt.Printf("saveJpeg\n")
	er := jpeg.Encode(osFile, imageRGBA, &options)
	if er != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", osFile.Name(), er)
		return er
	}
	return nil
}

//credits: https://riptutorial.com/go/example/31686/loading-and-saving-image

func savePng(osFile *os.File, imageRGBA *image.RGBA) error {
	fmt.Printf("savePng\n")
	er := png.Encode(osFile, imageRGBA)
	if er != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", osFile.Name(), er)
		return er
	}
	return nil
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

//credits: https://gist.github.com/ivanzoid
