package experiment

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	//"log"
)

type JoinExperiment struct {
	widget.BaseWidget
	Host       *widget.Entry
	ID         *widget.Entry
	JoinButton *widget.Button
}

func NewJoinExperiment() *JoinExperiment {
	j := &JoinExperiment{
		Host: widget.NewEntry(),
		ID:   widget.NewEntry(),
		JoinButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Join Experiment",
		},
	}
	j.Host.SetPlaceHolder("Host")
	j.ID.SetPlaceHolder("ID")
	// TODO Host is single line, handle enter key by moving to id
	// TODO ID is single line, handle enter key by moving to button/auto click
	j.ExtendBaseWidget(j)
	return j
}

func (j *JoinExperiment) CreateRenderer() fyne.WidgetRenderer {
	j.ExtendBaseWidget(j)
	return &joinExperimentRenderer{
		joinExperiment: j,
	}
}

var (
	minCellWidth  = 200
	minCellHeight = 20
)

type joinExperimentRenderer struct {
	joinExperiment *JoinExperiment
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *joinExperimentRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *joinExperimentRenderer) Destroy() {
}

func (r *joinExperimentRenderer) Layout(size fyne.Size) {
	//log.Println("joinExperimentRenderer.Layout:", size)
	p := theme.Padding()
	pp := theme.Padding() * 2
	cell := fyne.NewSize(size.Width-pp, size.Height/4-pp)

	x := p
	y := p
	r.joinExperiment.Host.Move(fyne.NewPos(x, y))
	r.joinExperiment.Host.Resize(cell)
	y += cell.Height
	r.joinExperiment.ID.Move(fyne.NewPos(x, y))
	r.joinExperiment.ID.Resize(cell)
	y += cell.Height * 2
	r.joinExperiment.JoinButton.Move(fyne.NewPos(x, y))
	r.joinExperiment.JoinButton.Resize(cell)

}

func (r *joinExperimentRenderer) MinSize() (size fyne.Size) {
	cell := fyne.NewSize(minCellWidth, minCellHeight)
	for _, c := range []fyne.CanvasObject{
		r.joinExperiment.Host,
		r.joinExperiment.ID,
		r.joinExperiment.JoinButton,
	} {
		cell = cell.Union(c.MinSize())
	}
	cell.Width += theme.Padding() * 2
	cell.Height += theme.Padding() * 2

	size.Width = cell.Width
	size.Height = cell.Height * 4
	//log.Println("joinExperimentRenderer.MinSize:", size)
	return
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *joinExperimentRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.joinExperiment.Host,
		r.joinExperiment.ID,
		r.joinExperiment.JoinButton,
	}
}

func (r *joinExperimentRenderer) Refresh() {
	// no-op
}
