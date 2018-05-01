package vsort

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/airking05/termui"
)

type Data struct {
	TickTime time.Duration
	data     []int
	alg      func(ops *Ops)
	render   chan bool
	exit     chan bool
}

type Ops struct {
	render bool
	data   *Data
	comps  int
	swaps  int
}

func (o *Ops) Swap(i, j int) {
	tmp := o.data.data[j]
	o.data.data[j] = o.data.data[i]
	o.data.data[i] = tmp
	o.swaps += 1

	if o.render {

		time.Sleep(o.data.TickTime)
		go func() { o.data.render <- true }()
	}

}
func (o *Ops) Diff(i, j int) int {
	o.comps += 1
	if o.render {
		time.Sleep(o.data.TickTime)
	}
	return o.data.data[i] - o.data.data[j]
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

func (o *Ops) Get(i int) int {
	return o.data.data[i]
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

func (t *Data) Run(algs ...func(ops *Ops)) {

	height := 23
	err := termui.Init()
	defer termui.Close()

	if err != nil {
		panic(err)
	}

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	graphs := []*Data{}

	for i, alg := range algs {
		ddata := make([]int, len(t.data))
		copy(ddata, t.data)
		tt := Data{
			alg:      alg,
			data:     ddata,
			TickTime: t.TickTime,
			render:   t.render,
		}
		graphs = append(graphs, &tt)

		// Running
		go func(d *Data) {
			d.alg(&Ops{
				data:   d,
				render: true,
			})
		}(&tt)

		// Timing
		go func(offset int, d *Data) {
			copyData := make([]int, len(d.data))
			copy(copyData, d.data)
			ops := Ops{
				data: &Data{
					data: copyData,
				},
				render: false,
			}

			start := time.Now()
			d.alg(&ops)
			elapsed := time.Now().UnixNano() - start.UnixNano()

			stats := termui.NewPar("Elements: " + strconv.Itoa(len(ops.data.data)) +
				", Comparisons: " + strconv.Itoa(ops.comps) +
				", Swaps: " + strconv.Itoa(ops.swaps) +
				", Time: " + strconv.Itoa(int(elapsed)) + "ns")
			stats.Height = 3
			stats.Width = 151
			stats.Y = offset*height + 20
			stats.BorderLabel = "Stats"
			termui.Render(stats)
		}(i, &tt)

	}

	// rendering
	go func() {
		for {
			select {
			case <-t.render:

				toRender := []termui.Bufferer{}
				for i, d := range graphs {
					graph := termui.NewSparkline()
					graph.Height = 17
					graph.LineColor = termui.ColorYellow
					graph.Data = d.data
					border := termui.NewSparklines(graph)
					border.Height = 20
					border.Width = 151
					border.BorderFg = termui.ColorCyan
					border.X = 0
					border.Y = i * height
					border.BorderLabel = "Sorting"
					toRender = append(toRender, border)
				}

				termui.Render(toRender...)
			case <-t.exit:
				return
			}
		}
	}()

	termui.Loop()
}
