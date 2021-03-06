package model

import (
	"github.com/stephenlyu/tview/transform"
	"github.com/stephenlyu/goformula/formulalibrary/base/formula"
)

type ModelListener interface {
	OnDataChanged()
	OnLastDataChanged()
}

type Model interface {
	Count() int									// data count
	Get(index int) []float64					// Get record at index
	GetRaw(index int) []float64					// Get un-scaled record at index
	GetNames() []string

	VarCount() int
	DrawActionCount() int 						// DrawAction数量
	NoDraw(index int) bool 						// 是否绘制图形
	NoText(index int) bool 						// 是否绘制文本
	DrawAbove(index int) bool
	NoFrame(index int) bool
	Color(index int) *formula.Color				// 变量颜色, 形如black或FFFFFF
	LineThick(index int) int 					// 线宽，1-9
	LineStyle(index int) int 					// 线宽，1-9
	GraphType(index int) int
	DrawAction(index int) formula.DrawAction

	SetValueTransformer(transformer transform.Transformer)
	SetScaleTransformer(transformer transform.ScaleTransformer)

	TransformRaw(v float64) float64
	Transform(v float64) float64

	TransformRawFrom(v float64) float64
	TransformFrom(v float64) float64

	AddListener(listener ModelListener)
	RemoveListener(listener ModelListener)

	NotifyDataChanged()
	NotifyLastDataChanged()
}
