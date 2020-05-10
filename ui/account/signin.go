package account

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	//"log"
)

type SignIn struct {
	widget.BaseWidget
	Alias        *widget.Entry
	Password     *widget.Entry
	SignInButton *widget.Button
}

func NewSignIn() *SignIn {
	s := &SignIn{
		Alias:    widget.NewEntry(),
		Password: widget.NewPasswordEntry(),
		SignInButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Sign In",
		},
	}
	s.Alias.SetPlaceHolder("Alias")
	s.Password.SetPlaceHolder("Password")
	// TODO Alias is single line, handle enter key by moving to password
	// TODO Password is single line, handle enter key by moving to button/auto click
	s.ExtendBaseWidget(s)
	return s
}

func (s *SignIn) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	return &signInRenderer{
		signIn: s,
	}
}

var (
	minCellWidth  = 200
	minCellHeight = 20
)

type signInRenderer struct {
	signIn *SignIn
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *signInRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *signInRenderer) Destroy() {
}

func (r *signInRenderer) Layout(size fyne.Size) {
	//log.Println("signInRenderer.Layout:", size)
	p := theme.Padding()
	pp := theme.Padding() * 2
	cell := fyne.NewSize(size.Width-pp, size.Height/5-pp)

	x := p
	y := p
	r.signIn.Alias.Move(fyne.NewPos(x, y))
	r.signIn.Alias.Resize(cell)
	y += cell.Height
	r.signIn.Password.Move(fyne.NewPos(x, y))
	r.signIn.Password.Resize(cell)
	y += cell.Height * 3
	r.signIn.SignInButton.Move(fyne.NewPos(x, y))
	r.signIn.SignInButton.Resize(cell)

}

func (r *signInRenderer) MinSize() (size fyne.Size) {
	cell := fyne.NewSize(minCellWidth, minCellHeight)
	for _, c := range []fyne.CanvasObject{
		r.signIn.Alias,
		r.signIn.Password,
		r.signIn.SignInButton,
	} {
		cell = cell.Union(c.MinSize())
	}
	cell.Width += theme.Padding() * 2
	cell.Height += theme.Padding() * 2

	size.Width = cell.Width
	size.Height = cell.Height * 5
	//log.Println("signInRenderer.MinSize:", size)
	return
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *signInRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.signIn.Alias,
		r.signIn.Password,
		r.signIn.SignInButton,
	}
}

func (r *signInRenderer) Refresh() {
	// no-op
}
