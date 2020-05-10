package experiment

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
)

type CreateExperiment struct {
	widget.BaseWidget
	Path         *widget.Entry
	CreateButton *widget.Button
}

func NewCreateExperiment(window fyne.Window) *CreateExperiment {
	c := &CreateExperiment{
		Path: widget.NewEntry(),
		CreateButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Create Experiment",
		},
	}
	c.Path.SetPlaceHolder("Path")
	// TODO Path is single line, handle enter key by moving to button/auto click
	c.Path.ActionItem = newFilePicker(window, c.Path)
	c.ExtendBaseWidget(c)
	return c
}

func (c *CreateExperiment) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	return &createExperimentRenderer{
		createExperiment: c,
	}
}

type createExperimentRenderer struct {
	createExperiment *CreateExperiment
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *createExperimentRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *createExperimentRenderer) Destroy() {
}

func (r *createExperimentRenderer) Layout(size fyne.Size) {
	//log.Println("createExperimentRenderer.Layout:", size)
	p := theme.Padding()
	pp := theme.Padding() * 2
	cell := fyne.NewSize(size.Width-pp, size.Height/4-pp)

	x := p
	y := p
	r.createExperiment.Path.Move(fyne.NewPos(x, y))
	r.createExperiment.Path.Resize(cell)
	y += cell.Height * 3
	r.createExperiment.CreateButton.Move(fyne.NewPos(x, y))
	r.createExperiment.CreateButton.Resize(cell)

}

func (r *createExperimentRenderer) MinSize() (size fyne.Size) {
	cell := fyne.NewSize(minCellWidth, minCellHeight)
	for _, c := range []fyne.CanvasObject{
		r.createExperiment.Path,
		r.createExperiment.CreateButton,
	} {
		cell = cell.Union(c.MinSize())
	}
	cell.Width += theme.Padding() * 2
	cell.Height += theme.Padding() * 2

	size.Width = cell.Width
	size.Height = cell.Height * 4
	//log.Println("createExperimentRenderer.MinSize:", size)
	return
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *createExperimentRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.createExperiment.Path,
		r.createExperiment.CreateButton,
	}
}

func (r *createExperimentRenderer) Refresh() {
	// no-op
}
