package account

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	//"log"
)

type SignUp struct {
	widget.BaseWidget
	Alias        *widget.Entry
	Password     *widget.Entry
	Confirm      *widget.Entry
	SignUpButton *widget.Button
}

func NewSignUp() *SignUp {
	s := &SignUp{
		Alias:    widget.NewEntry(),
		Password: widget.NewPasswordEntry(),
		Confirm:  widget.NewPasswordEntry(),
		SignUpButton: &widget.Button{
			Style: widget.PrimaryButton,
			Text:  "Sign Up",
		},
	}
	s.Alias.SetPlaceHolder("Alias")
	s.Password.SetPlaceHolder("Password")
	s.Confirm.SetPlaceHolder("Confirm Password")
	s.ExtendBaseWidget(s)
	return s
}

func (s *SignUp) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	return &signUpRenderer{
		signUp: s,
	}
}

type signUpRenderer struct {
	signUp *SignUp
}

// BackgroundColor satisfies the fyne.WidgetRenderer interface.
func (r *signUpRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Destroy satisfies the fyne.WidgetRenderer interface.
func (r *signUpRenderer) Destroy() {
}

func (r *signUpRenderer) Layout(size fyne.Size) {
	//log.Println("signUpRenderer.Layout:", size)
	p := theme.Padding()
	pp := theme.Padding() * 2
	cell := fyne.NewSize(size.Width-pp, size.Height/5-pp)

	x := p
	y := p
	r.signUp.Alias.Move(fyne.NewPos(x, y))
	r.signUp.Alias.Resize(cell)
	y += cell.Height
	r.signUp.Password.Move(fyne.NewPos(x, y))
	r.signUp.Password.Resize(cell)
	y += cell.Height
	r.signUp.Confirm.Move(fyne.NewPos(x, y))
	r.signUp.Confirm.Resize(cell)
	y += cell.Height * 2
	r.signUp.SignUpButton.Move(fyne.NewPos(x, y))
	r.signUp.SignUpButton.Resize(cell)

}

func (r *signUpRenderer) MinSize() (size fyne.Size) {
	cell := fyne.NewSize(minCellWidth, minCellHeight)
	for _, c := range []fyne.CanvasObject{
		r.signUp.Alias,
		r.signUp.Password,
		r.signUp.Confirm,
		r.signUp.SignUpButton,
	} {
		cell = cell.Union(c.MinSize())
	}
	cell.Width += theme.Padding() * 2
	cell.Height += theme.Padding() * 2

	size.Width = cell.Width
	size.Height = cell.Height * 5
	//log.Println("signUpRenderer.MinSize:", size)
	return
}

// Objects satisfies the fyne.WidgetRenderer interface.
func (r *signUpRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		r.signUp.Alias,
		r.signUp.Password,
		r.signUp.Confirm,
		r.signUp.SignUpButton,
	}
}

func (r *signUpRenderer) Refresh() {
	// no-op
}
