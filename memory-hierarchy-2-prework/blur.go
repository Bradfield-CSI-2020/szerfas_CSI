package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

var kernel = [][]float64{
	{0.006693, 0.00832, 0.00972, 0.01067, 0.011007, 0.01067, 0.00972, 0.00832, 0.006693},
	{0.00832, 0.010343, 0.012082, 0.013264, 0.013682, 0.013264, 0.012082, 0.010343, 0.00832},
	{0.00972, 0.012082, 0.014114, 0.015494, 0.015983, 0.015494, 0.014114, 0.012082, 0.00972},
	{0.01067, 0.013264, 0.015494, 0.017009, 0.017546, 0.017009, 0.015494, 0.013264, 0.01067},
	{0.011007, 0.013682, 0.015983, 0.017546, 0.0181, 0.017546, 0.015983, 0.013682, 0.011007},
	{0.01067, 0.013264, 0.015494, 0.017009, 0.017546, 0.017009, 0.015494, 0.013264, 0.01067},
	{0.00972, 0.012082, 0.014114, 0.015494, 0.015983, 0.015494, 0.014114, 0.012082, 0.00972},
	{0.00832, 0.010343, 0.012082, 0.013264, 0.013682, 0.013264, 0.012082, 0.010343, 0.00832},
	{0.006693, 0.00832, 0.00972, 0.01067, 0.011007, 0.01067, 0.00972, 0.00832, 0.006693},
}

// ideas
// go into library and figure out how we're accessing the pixels in the image, see if can store/retrieve more efficiently
// tile the movement through the nested outer loops given we can fit the entier kernal in L1 several times over
// I think we'd want to tile it so that it fits in (size of L1 cache - size of kernal)

//func blur(in image.Image) image.Image {
//	out := image.NewRGBA(in.Bounds())
//
//	for j := 0; j < out.Rect.Max.Y; j++ {
//		for i := 0; i < out.Rect.Max.X; i++ {
//			var rr, gg, bb float64
//			for kern_i, kern_row := range kernel {
//				for kern_j, weight := range kern_row {
//					ii := i - len(kernel)/2 + kern_i
//					jj := j - len(kern_row)/2 + kern_j
//					if 0 <= ii && ii < out.Rect.Max.X && 0 <= jj && jj < out.Rect.Max.Y {
//						r, g, b, _ := in.At(ii, jj).RGBA()
//						rr += float64(r>>8) * weight
//						gg += float64(g>>8) * weight
//						bb += float64(b>>8) * weight
//					}
//				}
//			}
//			out.Set(i, j, color.RGBA{uint8(rr), uint8(gg), uint8(bb), 255})
//		}
//	}
//	return out
//}

var TILE_SIZE int = 31 * 1024 / 8	// Estimating 31kb remaining in L1 cache after pulling in kernal; convert to bytes and divide by unit of increment for our loops, which is float64 or 8 bytes
//var TILE_LENGTH int = int(math.Sqrt(float64(TILE_SIZE))) // square root to get the length of one side <-- actually don't think we want this because not working with squares

func blur(in image.Image) image.Image {
	out := image.NewRGBA(in.Bounds())

	var ti, tj int      	// indexes of the tile
	var i, j int            // indexes within a tile
	var	i_end, j_end int// end when matrix dim not a multiple of tile_size

	var image_rows int = out.Rect.Max.Y
	var image_cols int = out.Rect.Max.X

	fmt.Printf("tile length is %d\n", TILE_SIZE)
	fmt.Printf("ti is %d\n", ti)
	fmt.Printf("tj is %d\n", tj)
	for ti = 0; ti < image_rows; ti += TILE_SIZE {
		i_end = int(math.Min(float64(ti) + float64(TILE_SIZE), float64(image_rows)))
		fmt.Printf("i_end is %d\n", i_end)
		for tj = 0; tj < image_cols; tj += TILE_SIZE {
			j_end = int(math.Min(float64(tj) + float64(TILE_SIZE), float64(image_cols)))
			fmt.Printf("j_end is %d\n", j_end)


			for i = 0; i < i_end; i++ {
				fmt.Printf("i is %d\n", i)
				for j = 0; j < j_end; j++ {
					fmt.Printf("j is %d\n", j)
					fmt.Printf("j end within j loop is %d\n", j_end)
					var rr, gg, bb float64
					for kern_i, kern_row := range kernel {
						fmt.Printf("kern_i is %d\n", kern_i)
						for kern_j, weight := range kern_row {
							i2 := i - len(kernel)/2 + kern_i
							j2 := j - len(kern_row)/2 + kern_j
							if 0 <= i2 && i2 < out.Rect.Max.Y && 0 <= j2 && j2 < out.Rect.Max.X {
								fmt.Printf("i2 is %d\n", i2)
								fmt.Printf("j2 is %d\n", j2)
								// might want to change this to j2, i2
								r, g, b, _ := in.At(i2, j2).RGBA()
								rr += float64(r>>8) * weight
								gg += float64(g>>8) * weight
								bb += float64(b>>8) * weight
							}
						}
					}
					out.Set(i, j, color.RGBA{uint8(rr), uint8(gg), uint8(bb), 255})
				}
			}

		}
	}

	return out
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Extract BMP pixels from file
	f, err := os.Open("input.jpeg")
	check(err)
	in, err := jpeg.Decode(f)
	check(err)

	// Perform blur
	out := blur(in)

	// Write to output
	fout, err := os.Create("output.jpeg")
	check(err)
	jpeg.Encode(fout, out, &jpeg.Options{Quality: 75})

	fmt.Printf("ok\n")
}
