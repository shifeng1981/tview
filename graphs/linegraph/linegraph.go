package linegraph

import (
	"github.com/stephenlyu/tview/model"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tview/transform"
	"github.com/therecipe/qt/widgets"
	"github.com/stephenlyu/tview/constants"
	"github.com/therecipe/qt/gui"
	"github.com/stephenlyu/goformula/function"
	"github.com/therecipe/qt/core"
	"github.com/stephenlyu/tview/graphs"
	"fmt"
)

type LineGraph struct {
	Model model.Model
	ValueIndex int
	Scene *widgets.QGraphicsScene
	xTransformer transform.ScaleTransformer

	color *gui.QColor
	startIndex, endIndex int
	PathItem *widgets.QGraphicsPathItem
}

func NewLineGraph(model model.Model, valueIndex int, color *gui.QColor, scene *widgets.QGraphicsScene, xTransformer transform.ScaleTransformer) *LineGraph {
	util.Assert(model != nil, "model != nil")
	util.Assert(model.VarCount() > valueIndex, "len(model.GetGraphTypes()) > valueIndex")
	util.Assert(model.GraphType(valueIndex) == constants.GraphTypeLine, "model.GetGraphTypes()[valueIndex] == constants.GraphTypeLine")

	this := &LineGraph{
		Model: model,
		ValueIndex: valueIndex,
		Scene: scene,
		xTransformer: xTransformer,
		color: color,
	}
	this.init()
	return this
}

func (this *LineGraph) init() {
	this.Model.AddListener(this)
}

func (this *LineGraph) OnDataChanged() {
	this.buildLine()
}

func (this *LineGraph) OnLastDataChanged() {
	if this.Model.Count() <= 0 {
		return
	}

	this.buildLine()
}

func (this *LineGraph) GetValueRange(startIndex int, endIndex int) (float64, float64) {
	if this.Model.Count() == 0 {
		return 0, 0
	}

	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > this.Model.Count() {
		endIndex = this.Model.Count()
	}

	startIndex, endIndex = this.adjustIndices(startIndex, endIndex)

	values := this.Model.GetRaw(startIndex)
	util.Assert(len(values) > this.ValueIndex, "len(values) > this.ValueIndex")

	high := values[this.ValueIndex]
	low := values[this.ValueIndex]

	for i := startIndex + 1; i < endIndex; i++ {
		values := this.Model.GetRaw(i)
		v := values[this.ValueIndex]
		if v > high {
			high = v
		}
		if v < low {
			low = v
		}
	}

	return low, high
}

func (this *LineGraph) buildLine() {
	this.Clear()

	if this.Model.Count() == 0 {
		return
	}

	path := gui.NewQPainterPath()

	needMove := true
	for i := this.startIndex; i < this.endIndex; i++ {
		x := (this.xTransformer.To(float64(i)) + this.xTransformer.To(float64(i + 1))) / 2
		values := this.Model.Get(i)
		v := values[this.ValueIndex]

		if function.IsNaN(v) {
			needMove = true
			continue
		}

		if needMove {
			path.MoveTo2(x, v)
			needMove = false
		} else {
			path.LineTo2(x, v)
		}
	}

	brush := gui.NewQBrush3(this.color, core.Qt__NoBrush)
	pen := gui.NewQPen3(this.color)
	graphs.SetPenStyle(pen, this.Model.LineStyle(this.ValueIndex))
	graphs.SetPenWidth(pen, this.xTransformer, this.Model.LineThick(this.ValueIndex))

	this.PathItem = this.Scene.AddPath(path, pen, brush)
}

func (this *LineGraph) adjustIndices(startIndex int, endIndex int) (int, int) {
	if startIndex > 0 {
		startIndex--
	}
	if endIndex < this.Model.Count() {
		endIndex++
	}
	return startIndex, endIndex
}

// 更新当前显示的K线
func (this *LineGraph) Update(startIndex int, endIndex int) {
	if this.Model.NoDraw(this.ValueIndex) {
		return
	}
	if startIndex < 0 {
		startIndex = 0
	}

	if endIndex > this.Model.Count() {
		endIndex = this.Model.Count()
	}

	startIndex, endIndex = this.adjustIndices(startIndex, endIndex)
	this.startIndex = startIndex
	this.endIndex = endIndex

	// 更新需要显示的K线
	this.buildLine()
}

// 清除所有的K线
func (this *LineGraph) Clear() {
	if this.Model.NoDraw(this.ValueIndex) {
		return
	}
	if this.PathItem != nil {
		this.Scene.RemoveItem(this.PathItem)
		this.PathItem = nil
	}
}

func (this *LineGraph) ShowInfo(index int, display graphs.InfoDisplay) {
	if this.Model.NoText(this.ValueIndex) {
		return
	}
	if index < 0 || index >= this.Model.Count() {
		return
	}

	name := this.Model.GetNames()[this.ValueIndex]
	v := this.Model.GetRaw(index)[this.ValueIndex]

	display.Add(fmt.Sprintf("%s: %s", name, graphs.FormatValue(v, 2)), this.color)
}