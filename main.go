package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	inputFile := flag.String("input", "input.png", "Path to the input image file")
	resizeOutput := flag.String("resize_output", "resize.png", "Path to save the resized image")
	grayOutput := flag.String("gray_output", "gray.png", "Path to save the grayscale image")
	asciiOutput := flag.String("ascii_output", "result.txt", "Path to save the ASCII art file")
	htmlOutput := flag.String("html_output", "result.html", "Path to save the ASCII art as HTML")
	imageOutput := flag.String("image_output", "output.png", "Path to save the ASCII art as an image")

	flag.Parse()

	image := LoadImage(*inputFile)
	resizeImage := ResizeImage(image, 200)

	saveImage(resizeImage, *resizeOutput)
	grayImage := ConvGrayScale(resizeImage)
	saveImage(grayImage, *grayOutput)

	resultStr := MapAscii(grayImage)
	saveToFile(resultStr, *asciiOutput)

	AsciiToHTML(resultStr, *htmlOutput)
	//AsciiToHTMLWithColor(resizeImage, resultStr, *htmlOutput)
	AsciiToImage(resultStr, *imageOutput)
}

func LoadImage(inputPath string) image.Image {
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error while opening file %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Printf("Error while decoding image %v\n", err)
		os.Exit(1)
	}
	return img
}

func saveImage(img image.Image, outputPath string) {
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error while creating file %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Error while encoding image %v\n", err)
		os.Exit(1)
	}
}

func ResizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	height := (bounds.Dy() * width) / bounds.Dx()
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, bounds, draw.Over, nil)
	return newImage
}

func ConvGrayScale(img image.Image) image.Image {
	bound := img.Bounds()
	grayImage := image.NewRGBA(bound)

	for i := bound.Min.X; i < bound.Max.X; i++ {
		for j := bound.Min.Y; j < bound.Max.Y; j++ {
			oldPixel := img.At(i, j)
			color := color.GrayModel.Convert(oldPixel)
			grayImage.Set(i, j, color)
		}
	}
	return grayImage
}

func MapAscii(img image.Image) []string {
	asciiChar := "$@B%#*+=,. "

	bound := img.Bounds()
	width := bound.Max.X
	height := bound.Max.Y

	imageAspectRatio := float64(width) / float64(height)

	charAspectRatio := 2.0

	adjustedHeight := int(float64(height) / (charAspectRatio / imageAspectRatio))

	result := make([]string, adjustedHeight)

	for y := 0; y < adjustedHeight; y++ {
		line := ""
		for x := 0; x < width; x++ {
			imgY := int(float64(y) * (float64(height) / float64(adjustedHeight)))

			pixelValue := color.GrayModel.Convert(img.At(x, imgY)).(color.Gray)
			pixel := pixelValue.Y
			asciiIndex := int(pixel) * (len(asciiChar) - 1) / 255
			line += string(asciiChar[asciiIndex])
		}
		result[y] = line
	}
	return result
}


func saveToFile(asciiArt []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range asciiArt {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func AsciiToHTML(ascii []string, outputPath string) {
	HtmlFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error while creating HTML file")
		os.Exit(1)
	}
	defer HtmlFile.Close()

	for lin, lines := range ascii {
		htmlString := `<!DOCTYPE html>
		<html lang="en"><head>
   	 	<meta charset="UTF-8">
    	<meta name="viewport" content="width=device-width, initial-scale=0.8">
    	<title>AsciiImage</title>
		</head>
		<body>
			<code>
		 		<span class="ascii" style="color: black;
		  		background: white;
		  		display:inline-block;
		  		white-space:pre;
		  		letter-spacing:0;
		  		line-height:0.9;
		  		font-family:'Consolas','BitstreamVeraSansMono','CourierNew',Courier,monospace;
		  		font-size:10px;
		  		border-width:1px;
		  		border-style:solid;
		  		border-color:lightgray;">`
		if lin == 0 {
			_, err := HtmlFile.WriteString(htmlString)
			if err != nil {
				fmt.Println("Error while start writing into HTML file")
				os.Exit(1)
			}
		}

		for _, char := range lines {
			_, err := HtmlFile.WriteString(fmt.Sprintf("<span>%v</span>", string(char)))
			if err != nil {
				fmt.Println("Error while writing into HTML file")
				os.Exit(1)
			}
		}
		_, err := HtmlFile.WriteString("<br>")
		if err != nil {
			fmt.Println("Error while writing into HTML file")
			os.Exit(1)
		}
		if lin == len(ascii)-1 {
			_, err := HtmlFile.WriteString("</code></body></html>")
			if err != nil {
				fmt.Println("Error while end writing into HTML file")
				os.Exit(1)
			}

		}
	}
}

func AsciiToHTMLWithColor(img image.Image, ascii []string, outputPath string) {
	HtmlFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error while creating HTML file")
		os.Exit(1)
	}
	defer HtmlFile.Close()

	_, err = HtmlFile.WriteString(`<!DOCTYPE html>
		<html lang="en"><head>
    	<meta charset="UTF-8">
    	<meta name="viewport" content="width=device-width, initial-scale=1.0">
    	<title>Color ASCII Art</title>
		</head>
		<body style="font-family: monospace; line-height: 1; white-space: pre;">`)
	if err != nil {
		fmt.Println("Error while writing to HTML file")
		os.Exit(1)
	}

	bounds := img.Bounds()

	for y, line := range ascii {
		for x := 0; x < len(line); x++ {
			imgX := x
			imgY := y
			if imgY >= bounds.Max.Y || imgX >= bounds.Max.X {
				continue
			}

			r, g, b, _ := img.At(imgX, imgY).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			char := string(line[x])
			htmlString := fmt.Sprintf(`<span style="color: rgb(%d,%d,%d);">%s</span>`, r8, g8, b8, char)
			_, err := HtmlFile.WriteString(htmlString)
			if err != nil {
				fmt.Println("Error while writing to HTML file")
				os.Exit(1)
			}
		}
		_, err = HtmlFile.WriteString("<br>")
		if err != nil {
			fmt.Println("Error while writing to HTML file")
			os.Exit(1)
		}
	}

	_, err = HtmlFile.WriteString(`</body></html>`)
	if err != nil {
		fmt.Println("Error while closing HTML file")
		os.Exit(1)
	}
}


func AsciiToImage(strArray []string, outputPath string) {
	fontImage := image.NewRGBA(image.Rect(0, 0, 1400, len(strArray)*11))
	draw.Draw(fontImage, fontImage.Bounds(), image.White, image.Point{}, draw.Src)

	drawconf := &font.Drawer{
		Dst:  fontImage,
		Src:  image.Black,
		Face: basicfont.Face7x13,
	}

	for i, line := range strArray {
		drawconf.Dot = fixed.Point26_6{
			X: fixed.Int26_6(10 * 64),
			Y: fixed.Int26_6((20 + i*11) * 64),
		}

		drawconf.DrawString(line)

	}

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error while creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	err = png.Encode(file, fontImage)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		os.Exit(1)
	}
}
