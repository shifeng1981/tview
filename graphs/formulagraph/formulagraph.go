package formulagraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/graphs"
	"github.com/stephenlyu/tview/constants"
	"github.com/stephenlyu/tview/graphs/linegraph"
	"github.com/stephenlyu/tview/graphs/stickgraph"
	"github.com/stephenlyu/tview/graphs/volgraph"
)

type FormulaGraph struct {
	Model model.Model
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	graphs []graphs.Graph
}

func NewFormulaGraph(model model.Model, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *FormulaGraph {
	util.Assert(model != nil, "model != nil")

	this := &FormulaGraph{
		Model: model,
		Scene: scene,
		xTransformer: xTransformer,
	}
	this.init()
	return this
}

func (this *FormulaGraph) init() {
	graphTypes := this.Model.GetGraphTypes()
	this.graphs = make([]graphs.Graph, len(graphTypes))
	j := 0
	for i, graphType := range graphTypes {
		var graph graphs.Graph
		switch graphType {
		case constants.GraphTypeLine:
			color := graphs.COLORS[j % len(graphs.COLORS)]
			j++
			graph = linegraph.NewLineGraph(this.Model, i, color, this.Scene, this.xTransformer)
		case constants.GraphTypeStick:
			graph = stickgraph.NewStickGraph(this.Model, i, this.Scene, this.xTransformer)
		case constants.GraphTypeVolStick:
			graph = volgraph.NewVolStickGraph(this.Model, i, this.Scene, this.xTransformer)
		}
		this.graphs[i] = graph
	}
}

func (this *FormulaGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
	low, high := this.graphs[0].GetValueRange(startIndex, endIndex)

	for i := 1; i < len(this.graphs); i++ {
		low1, high1 := this.graphs[i].GetValueRange(startIndex, endIndex)
		if low1 < low {
			low = low1
		}
		if high1 > high {
			high = high1
		}
	}

	return low, high
}

// 更新当前显示的K线
func (this *FormulaGraph) Update(startIndex int, endIndex int) {
	for _, graph := range this.graphs {
		graph.Update(startIndex, endIndex)
	}
}

// 清除所有的K线
func (this *FormulaGraph) Clear() {
	for _, graph := range this.graphs {
		graph.Clear()
	}
}