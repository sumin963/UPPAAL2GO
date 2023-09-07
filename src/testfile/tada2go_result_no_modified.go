package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// #define N  6
// typedef  int id_t;
// typedef struct Train{
// }Train;
// typedef struct Gate{
//         id_t list[N+1];
//         int len;
// }Gate;
// void enqueue(Gate *Gate, id_t element)
// {
//         Gate->list[Gate-> len++] = element;
// }
// void dequeue(Gate *Gate )
// {
//         int i = 0;
//         Gate-> len -= 1;
//         while (i < Gate-> len)
//         {
//                 Gate->list[i] = Gate->list[i + 1];
//                 i++;
//         }
//         Gate->list[i] = 0;
// }
// id_t front(Gate *Gate )
// {
//    return Gate->list[0];
// }
// id_t tail(Gate *Gate )
// {
//    return Gate->list[Gate-> len - 1];
// }
//
import "C"

type loc_info struct {
	id         string
	loc_circle *canvas.Circle
	xPox       int
	yPox       int
}

func main() {
	eps := time.Millisecond * 10
	appr_chan := make([]chan bool, C.N)
	for i := range appr_chan {
		appr_chan[i] = make(chan bool)
	}
	stop_chan := make([]chan bool, C.N)
	for i := range stop_chan {
		stop_chan[i] = make(chan bool)
	}
	leave_chan := make([]chan bool, C.N)
	for i := range leave_chan {
		leave_chan[i] = make(chan bool)
	}
	go_chan := make([]chan bool, C.N)
	for i := range go_chan {
		go_chan[i] = make(chan bool)
	}
	//
	//
	app := app.New()
	window := app.NewWindow("Diagonal")
	var loc_all []loc_info
	loc_id0 := &loc_info{"id0", &canvas.Circle{StrokeColor: color.RGBA{220, 20, 60, 255}, StrokeWidth: 15}, 100, 100}
	loc_id1 := &loc_info{"id1", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 200, 300}
	loc_id2 := &loc_info{"id2", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 300, 100}
	loc_id2p := &loc_info{"id2p", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 350, 100}
	loc_id2pp := &loc_info{"id2pp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 400, 100}
	loc_id3 := &loc_info{"id3", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 100, 200}
	loc_id3p := &loc_info{"id3p", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 150, 200}
	loc_id3pp := &loc_info{"id3pp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 200, 200}
	loc_id4 := &loc_info{"id4", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 300, 200}
	loc_id4p := &loc_info{"id4p", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 350, 200}
	loc_id4pp := &loc_info{"id4pp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 400, 200}
	loc_exp := &loc_info{"exp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 200, 500}

	loc_all = append(loc_all, *loc_id0, *loc_id1, *loc_id2, *loc_id2p, *loc_id2pp, *loc_id3, *loc_id3p, *loc_id3pp, *loc_id4, *loc_id4p, *loc_id4pp, *loc_exp)

	for _, val := range loc_all {
		val.loc_circle.Resize(fyne.NewSize(15, 15))
		val.loc_circle.Move(fyne.NewPos(float32(val.xPox), float32(val.yPox)))
	}

	tick := time.NewTicker(time.Second * 1)

	ch := make(chan string, 1)

	go func() {
		current_loc := "id0"
		past_loc := "id0"
		for {
			window.SetContent(container.NewWithoutLayout(
				loc_all[0].loc_circle,
				loc_all[1].loc_circle,
				loc_all[2].loc_circle,
				loc_all[3].loc_circle,
				loc_all[4].loc_circle,
				loc_all[5].loc_circle,
				loc_all[6].loc_circle,
				loc_all[7].loc_circle,
				loc_all[8].loc_circle,
				loc_all[9].loc_circle,
				loc_all[10].loc_circle,
				loc_all[11].loc_circle))
			<-tick.C
			var _loc string
			select {
			case _loc = <-ch:
			default:
				_loc = ""
			}
			for i, val := range loc_all {
				if val.id == _loc {
					for _, val2 := range loc_all {
						if val2.id == past_loc {
							val2.loc_circle = &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}
							val2.loc_circle.Resize(fyne.NewSize(15, 15))
							val2.loc_circle.Move(fyne.NewPos(float32(val2.xPox), float32(val2.yPox)))
							val2.loc_circle.Refresh()
							break
						}
					}
					past_loc = current_loc
					current_loc = _loc
					loc_all[i].loc_circle = &canvas.Circle{StrokeColor: color.RGBA{220, 20, 60, 255}, StrokeWidth: 15}
					val.loc_circle.Resize(fyne.NewSize(15, 15))
					val.loc_circle.Move(fyne.NewPos(float32(val.xPox), float32(val.yPox)))
					val.loc_circle.Refresh()
					break
				}
			}

		}
	}()
	//
	//
	Train := func(id C.id_t, ch chan string) {
		x_now := time.Now()
		x := time.Since(x_now)
		var id2_passage []string
		var id3_passage []string
		var id4_passage []string
		//
		//
		// app := app.New()
		// window := app.NewWindow("Diagonal")
		// var loc_all []loc_info
		// loc_id0 := loc_info{"id0", &canvas.Circle{StrokeColor: color.RGBA{220, 20, 60, 255}, StrokeWidth: 15}, 100, 100}
		// loc_id1 := loc_info{"id1", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 150, 300}
		// loc_id2 := loc_info{"id2", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 300, 100}
		// loc_id2p := loc_info{"id2p", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 350, 100}
		// loc_id2pp := loc_info{"id2pp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 400, 100}
		// loc_id3 := loc_info{"id3", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 100, 200}
		// loc_id3p := loc_info{"id3p", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 150, 200}
		// loc_id3pp := loc_info{"id3pp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 200, 200}
		// loc_id4 := loc_info{"id4", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 300, 200}
		// loc_id4p := loc_info{"id4p", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 350, 200}
		// loc_id4pp := loc_info{"id4pp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 400, 200}
		// loc_exp := loc_info{"exp", &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}, 0, 0}
		// loc_all = append(loc_all, loc_id0, loc_id1, loc_id2, loc_id2p, loc_id2pp, loc_id3, loc_id3p, loc_id3p, loc_id3pp, loc_id4, loc_id4p, loc_id4pp, loc_exp)

		// for _, val := range loc_all {
		// 	val.loc_circle.Resize(fyne.NewSize(15, 15))
		// 	val.loc_circle.Move(fyne.NewPos(float32(val.xPox), float32(val.yPox)))
		// }

		// tick := time.NewTicker(time.Second * 1)

		// //var current_loc *canvas.Circle

		// ch := make(chan string, 1)

		// go func() {
		// 	current_loc := "id0"
		// 	past_loc := "id0"
		// 	for {
		// 		window.SetContent(container.NewWithoutLayout(loc_all[0].loc_circle,
		// 			loc_all[1].loc_circle,
		// 			loc_all[2].loc_circle,
		// 			loc_all[3].loc_circle,
		// 			loc_all[4].loc_circle,
		// 			loc_all[5].loc_circle,
		// 			loc_all[6].loc_circle,
		// 			loc_all[7].loc_circle,
		// 			loc_all[8].loc_circle,
		// 			loc_all[9].loc_circle,
		// 			loc_all[10].loc_circle,
		// 			loc_all[11].loc_circle))
		// 		<-tick.C
		// 		var _loc string
		// 		select {
		// 		case _loc = <-ch:
		// 		default:
		// 			_loc = ""
		// 		}
		// 		for i, val := range loc_all {
		// 			if val.id == _loc {
		// 				for _, val2 := range loc_all {
		// 					if val2.id == past_loc {
		// 						val2.loc_circle = &canvas.Circle{StrokeColor: color.RGBA{39, 112, 180, 255}, StrokeWidth: 15}
		// 						val2.loc_circle.Resize(fyne.NewSize(15, 15))
		// 						val2.loc_circle.Move(fyne.NewPos(float32(val2.xPox), float32(val2.yPox)))
		// 						val2.loc_circle.Refresh()
		// 						break
		// 					}
		// 				}
		// 				past_loc = current_loc
		// 				current_loc = _loc
		// 				loc_all[i].loc_circle = &canvas.Circle{StrokeColor: color.RGBA{220, 20, 60, 255}, StrokeWidth: 15}
		// 				val.loc_circle.Resize(fyne.NewSize(15, 15))
		// 				val.loc_circle.Move(fyne.NewPos(float32(val.xPox), float32(val.yPox)))
		// 				val.loc_circle.Refresh()
		// 				break
		// 			}
		// 		}

		// 	}
		// }()
		//
		//
		goto id0
	id0:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Safe", "location", "x", ":", x)
		select {
		case appr_chan[id] <- true:
			x_now = time.Now()
			ch <- "id3"
			goto id3
		}
	id1:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Stop", "location", "x", ":", x)
		select {
		case <-go_chan[id]:
			x_now = time.Now()
			ch <- "id4"

			goto id4
		}
	id2:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross_0", "location", "x", ":", x)
		id2_passage = []string{"x==3", "x>3", "x>5"}
		switch time_passage(id2_passage, x) {
		case 0:
		case 1:
			ch <- "id2p"

			goto id2p
		case 2:
			ch <- "id2pp"

			goto id2pp
		case 3:
			ch <- "exp"

			goto exp
		}
		select {
		case <-time.After(time.Second*3 - x - eps):
			ch <- "id2p"

			goto id2p
		}
	id3:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr_0", "location", "x", ":", x)
		id3_passage = []string{"x==10", "x>10", "x>20"}
		switch time_passage(id3_passage, x) {
		case 0:
		case 1:
			ch <- "id3p"

			goto id3p
		case 2:
			ch <- "id3pp"

			goto id3pp
		case 3:
			ch <- "exp"

			goto exp
		}
		select {
		case <-time.After(time.Second*10 - x - eps):
			ch <- "id3p"

			goto id3p
		case <-stop_chan[id]:
			ch <- "id1"

			goto id1
		}
	id4:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start_0", "location", "x", ":", x)
		id4_passage = []string{"x==7", "x>7", "x>15"}
		switch time_passage(id4_passage, x) {
		case 0:
		case 1:
			ch <- "id4p"

			goto id4p
		case 2:
			ch <- "id4pp"

			goto id4pp
		case 3:
			ch <- "exp"

			goto exp
		}
		select {
		case <-time.After(time.Second*7 - x - eps):
			ch <- "id4p"

			goto id4p
		}
	exp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "exp", "location", "x", ":", x)
		select {}
	id2p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross_1", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*3 - x):
			ch <- "id2pp"

			goto id2pp
		case leave_chan[id] <- true:
			ch <- "id0"

			goto id0
		}
	id2pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross_2", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*5 - x):
			ch <- "exp"

			goto exp
		case leave_chan[id] <- true:
			ch <- "id0"

			goto id0
		}
	id3p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr_1", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*10 - x):
			ch <- "id3pp"

			goto id3pp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			ch <- "id2"

			goto id2
		case <-stop_chan[id]:
			ch <- "id1"

			goto id1
		}
	id3pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr_2", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*20 - x):
			ch <- "exp"

			goto exp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			ch <- "id2"

			goto id2
		}
	id4p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start_1", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*7 - x):
			ch <- "id4pp"

			goto id4pp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			ch <- "id2"

			goto id2
		}
	id4pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start_2", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*15 - x):
			ch <- "exp"

			goto exp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			ch <- "id2"

			goto id2
		}
	}
	Gate := func() {
		local_val := C.Gate{list: [C.N + 1]C.id_t{}, len: 0}
		goto id7
	id5:
		select {
		case stop_chan[C.tail(&local_val)] <- true:
			goto id6
		}
	id6:
		select {
		case <-appr_chan[0]:
			C.enqueue(&local_val, 0)
			goto id5
		case <-leave_chan[C.front(&local_val)]:
			C.dequeue(&local_val)
			goto id7
		}
	id7:
		select {
		case when(local_val.len > 0, go_chan[C.front(&local_val)]) <- true:
			goto id6
		case <-when(local_val.len == 0, appr_chan[0]):
			C.enqueue(&local_val, 0)
			goto id6
		}
	}
	go Train(0, ch)
	go Gate()
	window.ShowAndRun()
	<-time.After(time.Second * 20)
}
func when(guard bool, channel chan bool) chan bool {
	if !guard {
		return nil
	}
	return channel
}
func when_guard(guard bool) <-chan time.Time {
	if !guard {
		return nil
	}
	return time.After(time.Second * 0)
}
func time_passage(time_passage []string, ctime time.Duration) int {
	for i, val := range time_passage {
		if strings.Contains(val, "==") {
			num, _ := strconv.Atoi(val[strings.Index(val, "==")+2:])
			if time.Second*time.Duration(num) > ctime {
				return i
			}
		} else if strings.Contains(val, "<") {
			num, _ := strconv.Atoi(val[strings.Index(val, "==")+1:])
			if time.Second*time.Duration(num) == ctime {
				return i
			}
		}
	}
	return len(time_passage)
}
