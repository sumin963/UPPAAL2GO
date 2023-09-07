package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Canvas")
// 	myCanvas := myWindow.Canvas()

// 	blue := color.NRGBA{R: 0, G: 0, B: 180, A: 255}
// 	rect := canvas.NewRectangle(blue)
// 	myCanvas.SetContent(rect)

// 	go func() {
// 		time.Sleep(time.Second)
// 		green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
// 		rect.FillColor = green
// 		rect.Refresh()
// 		time.Sleep(time.Second)

// 		setContentToCircle(myCanvas)
// 	}()

// 	myWindow.Resize(fyne.NewSize(100, 100))
// 	myWindow.ShowAndRun()
// }
// func setContentToCircle(c fyne.Canvas) {
// 	red := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
// 	circle := canvas.NewCircle(color.White)
// 	circle.StrokeWidth = 4
// 	circle.StrokeColor = red
// 	c.SetContent(circle)
// }

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Box Layout")

// 	text1 := canvas.NewText("Hello", color.White)
// 	text2 := canvas.NewText("There", color.White)
// 	text3 := canvas.NewText("(right)", color.White)
// 	content := container.New(layout.NewHBoxLayout(), text1, text2, layout.NewSpacer(), text3)

// 	text4 := canvas.NewText("centered", color.White)
// 	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())
// 	myWindow.SetContent(container.New(layout.NewVBoxLayout(), content, centered))
// 	myWindow.ShowAndRun()
// }

type diagonal struct {
}

// func (d *diagonal) MinSize(objects []fyne.CanvasObject) fyne.Size {
// 	w, h := float32(0), float32(0)
// 	fmt.Println("2", objects)

// 	for _, o := range objects {
// 		childSize := o.MinSize()

// 		w += childSize.Width
// 		h += childSize.Height
// 		//fmt.Println(w, h)

// 	}
// 	return fyne.NewSize(w, h)
// }
// func (d *diagonal) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
// 	//pos := fyne.NewPos(0, containerSize.Height-d.MinSize(objects).Height)
// 	pos := fyne.NewPos(96, 96)
// 	for _, o := range objects {
// 		size := o.MinSize()
// 		o.Resize(size)
// 		o.Move(pos)

//			pos = pos.Add(fyne.NewPos(200, 96))
//		}
//	}

func main() {
	a := app.New()
	w := a.NewWindow("Diagonal")
	var circle_1 *canvas.Circle
	var circle_2 *canvas.Circle
	circle_1 = &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}
	circle_2 = &canvas.Circle{StrokeColor: color.RGBA{220, 20, 60, 255}, StrokeWidth: 15}

	circle_1.Resize(fyne.NewSize(15, 15))
	circle_1.Move(fyne.NewPos(96, 96))
	circle_2.Resize(fyne.NewSize(15, 15))
	circle_2.Move(fyne.NewPos(200, 96))
	tick := time.NewTicker(time.Second * 5)
	go func() {
		for {
			w.SetContent(container.NewWithoutLayout(circle_1, circle_2))
			circle_2 = &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}
			circle_2.Resize(fyne.NewSize(15, 15))
			circle_2.Move(fyne.NewPos(200, 96))
			circle_2.Refresh()
			<-tick.C
		}
	}()

	//w.SetContent(container.NewWithoutLayout(circle_1, circle_2))

	w.ShowAndRun()
}
