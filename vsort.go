package main

import (
	"math/rand"
	"time"

	"github.com/airking05/termui"
)

func main() {

	bubble := func(ops Ops) {
		for i := 0; i < ops.Len(); i++ {
			for j := 0; j < ops.Len(); j++ {
				if ops.LessThen(i, j) {
					ops.Swap(i, j)
				}
			}
		}
	}

	r := Init()
	r.TickTime = time.Millisecond * 1

	r.Run(bubble)

}

type Data struct {
	TickTime time.Duration
	data     []int
	alg      func(ops Ops)
	render   chan bool
	exit     chan bool
}

type Ops struct {
	data *Data
}

func (o *Ops) Swap(i, j int) {
	tmp := o.data.data[j]
	o.data.data[j] = o.data.data[i]
	o.data.data[i] = tmp

	o.data.render <- true
	time.Sleep(o.data.TickTime)
}

func (o *Ops) Len() int {
	return len(o.data.data)
}

func (o *Ops) Equal(i, j int) bool {
	return o.Diff(i, j) == 0
}

func (o *Ops) NotEqual(i, j int) bool {
	return !o.Equal(i, j)
}

func (o *Ops) Diff(i, j int) int {
	return o.data.data[i] - o.data.data[j]
}

func (o *Ops) LessThen(i, j int) bool {
	return o.Diff(i, j) < 0
}

func (o *Ops) GreaterThen(i, j int) bool {
	return o.Diff(i, j) > 0
}

func (o *Ops) LessEqualThen(i, j int) bool {
	return o.Diff(i, j) <= 0
}

func (o *Ops) GreaterEqualThen(i, j int) bool {
	return o.Diff(i, j) >= 0
}

func Init() Data {

	data := make([]int, 150)

	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100) + 1
	}

	return InitWith(data)
}

func InitWith(data []int) Data {

	return Data{
		TickTime: 100 * time.Millisecond,
		data:     data,
		render:   make(chan bool),
		exit:     make(chan bool),
	}
}

func (t *Data) Run(alg func(ops Ops)) {
	t.alg = alg

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	graph := termui.NewSparkline()
	graph.Data = t.data
	graph.Height = 17
	graph.LineColor = termui.ColorYellow

	// border := termui.NewSparklines(graph)
	// border.Height = 20
	// border.Width = 151
	// border.BorderFg = termui.ColorCyan
	// border.X = 0
	// border.BorderLabel = "Sorting"

	// termui.Render(border)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	go func() {
		for {
			select {
			case <-t.render:

				graph.Data = t.data
				border := termui.NewSparklines(graph)
				border.Height = 20
				border.Width = 151
				border.BorderFg = termui.ColorCyan
				border.X = 0
				border.BorderLabel = "Sorting"

				termui.Render(border)
			case <-t.exit:
				return
			}
		}
	}()

	go func() {
		t.alg(Ops{
			data: t,
		})
		t.exit <- true
	}()

	termui.Loop()
}
