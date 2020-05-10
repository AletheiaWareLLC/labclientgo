package experiment

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"log"
)

var _ desktop.Cursorable = (*filePicker)(nil)
var _ fyne.Tappable = (*filePicker)(nil)
var _ fyne.Widget = (*filePicker)(nil)

type filePicker struct {
	widget.BaseWidget
	window fyne.Window
	icon   *canvas.Image
	entry  *widget.Entry
}

func newFilePicker(w fyne.Window, e *widget.Entry) *filePicker {
	pr := &filePicker{
		window: w,
		icon:   canvas.NewImageFromResource(theme.FolderIcon()),
		entry:  e,
	}
	pr.ExtendBaseWidget(pr)
	return pr
}

// CreateRenderer satisfies the fyne.Widget interface.
func (r *filePicker) CreateRenderer() fyne.WidgetRenderer {
	return &filePickerRenderer{
		icon:  r.icon,
		entry: r.entry,
	}
}

// Cursor satisfies the desktop.Cursorable interface.
func (r *filePicker) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// Tapped satisfies the fyne.Tappable interface.
func (r *filePicker) Tapped(*fyne.PointEvent) {
	log.Println("filePicker.Tapped")
	// Show open file dialog
	dialog.ShowFileOpen(func(reader fyne.FileReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, r.window)
			return
		}
		if reader == nil {
			return
		}
		// Set entry text to file uri
		r.entry.SetText(reader.URI())
		reader.Close()
	}, r.window)
}

var _ fyne.WidgetRenderer = (*filePickerRenderer)(nil)

type filePickerRenderer struct {
	entry *widget.Entry
	icon  *canvas.Image
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *filePickerRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *filePickerRenderer) Destroy() {
}

// Layout satisfies the fyne.WidgetRenderer interface.
func (r *filePickerRenderer) Layout(size fyne.Size) {
	r.icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
	r.icon.Move(fyne.NewPos((size.Width-theme.IconInlineSize())/2, (size.Height-theme.IconInlineSize())/2))
}

// MinSize satisfies the fyne.WidgetRenderer interface.
func (r *filePickerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *filePickerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.icon}
}

// Refresh satisfies the fyne.WidgetRenderer interface.
func (r *filePickerRenderer) Refresh() {
	canvas.Refresh(r.icon)
}
