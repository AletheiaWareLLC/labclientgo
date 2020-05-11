/*
 * Copyright 2020 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package labclientgo

import (
	"bytes"
	"errors"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	//"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/aliasgo"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/cryptogo"
	"github.com/AletheiaWareLLC/labclientgo/ui"
	"github.com/AletheiaWareLLC/labclientgo/ui/account"
	"github.com/AletheiaWareLLC/labclientgo/ui/data"
	"github.com/AletheiaWareLLC/labclientgo/ui/edit"
	"github.com/AletheiaWareLLC/labclientgo/ui/experiment"
	"github.com/AletheiaWareLLC/labgo"
	"log"
	"os"
)

const (
	MIN_PASSWORD                 = 12
	ERROR_PASSWORD_TOO_SHORT     = "Password Too Short: %d Minimum: %d"
	ERROR_PASSWORDS_DO_NOT_MATCH = "Passwords Do Not Match"
)

type Client struct {
	Root       string
	Cache      bcgo.Cache
	Network    bcgo.Network
	Node       *bcgo.Node
	Experiment *labgo.Experiment
	App        fyne.App
	Window     fyne.Window
	Dialog     dialog.Dialog
}

func (c *Client) ExistingNode(alias string, password []byte, callback func(*bcgo.Node)) {
	// Get key store
	keystore, err := bcgo.GetKeyDirectory(c.Root)
	if err != nil {
		c.ShowError(err)
		return
	}
	// Get private key
	key, err := cryptogo.GetRSAPrivateKey(keystore, alias, password)
	if err != nil {
		c.ShowError(err)
		return
	}
	// Create node
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    c.Cache,
		Network:  c.Network,
		Channels: make(map[string]*bcgo.Channel),
	}

	callback(node)
}

func (c *Client) GetNode() *bcgo.Node {
	if c.Node == nil {
		nc := make(chan *bcgo.Node, 1)
		go c.ShowAccessDialog(func(n *bcgo.Node) {
			nc <- n
		})
		c.Node = <-nc
	}
	return c.Node
}

func (c *Client) GetExperiment() *labgo.Experiment {
	if c.Experiment == nil {
		ec := make(chan *labgo.Experiment, 1)
		go c.ShowExperimentDialog(func(h string, e *labgo.Experiment) {
			if n, ok := c.Network.(*bcgo.TCPNetwork); ok {
				go labgo.Serve(c.Node, c.Cache, n)
				if h != "" && h != "localhost" {
					n.Connect(h, []byte("test"))
				}
			}
			ec <- e
		})
		c.Experiment = <-ec
	}
	return c.Experiment
}

func (c *Client) GetLogo() fyne.CanvasObject {
	return &canvas.Image{
		Resource: data.NewThemedResource(data.LogoUnmasked),
		FillMode: canvas.ImageFillOriginal,
	}
}

func (c *Client) GetOrOpenDeltaChannel(fileId string) *bcgo.Channel {
	node := c.GetNode()
	channel, err := node.GetChannel(labgo.LAB_PREFIX_FILE + fileId)
	if err != nil {
		log.Println(err)
		channel = labgo.OpenFileChannel(fileId)
		// Load channel
		if err := channel.LoadCachedHead(c.Cache); err != nil {
			log.Println(err)
		}
		if c.Network != nil {
			// Pull channel from network
			if err := channel.Pull(c.Cache, c.Network); err != nil {
				log.Println(err)
			}
		}
		// Add channel to node
		node.AddChannel(channel)
	}
	return channel
}

func (c *Client) NewNode(alias string, password []byte, callback func(*bcgo.Node)) {
	// Create Progress Dialog
	progress := dialog.NewProgress("Registering", "message", c.Window)
	defer progress.Hide()
	listener := &ui.ProgressMiningListener{Func: progress.SetValue}

	// Get key store
	keystore, err := bcgo.GetKeyDirectory(c.Root)
	if err != nil {
		c.ShowError(err)
		return
	}
	// Create private key
	key, err := cryptogo.CreateRSAPrivateKey(keystore, alias, password)
	if err != nil {
		c.ShowError(err)
		return
	}
	// Create node
	node := &bcgo.Node{
		Alias:    alias,
		Key:      key,
		Cache:    c.Cache,
		Network:  c.Network,
		Channels: make(map[string]*bcgo.Channel),
	}

	// Register Alias
	if err := aliasgo.Register(node, listener); err != nil {
		c.ShowError(err)
		return
	}

	callback(node)
}

func (c *Client) ShowAccessDialog(callback func(*bcgo.Node)) {
	signIn := account.NewSignIn()
	signUp := account.NewSignUp()

	var content fyne.CanvasObject
	if c.App.Driver().Device().IsMobile() {
		content = widget.NewVBox(
			c.GetLogo(),
			widget.NewTabContainer(
				widget.NewTabItem("Sign In", signIn),
				widget.NewTabItem("Sign Up", signUp),
			),
		)
	} else {
		content = fyne.NewContainerWithLayout(
			layout.NewCenterLayout(),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				widget.NewGroup("Sign In", signIn),
				widget.NewGroup("Sign Up", signUp),
			),
		)
	}
	c.Dialog = dialog.NewCustom("Account Access", "Cancel", content, c.Window)

	if alias, err := bcgo.GetAlias(); err == nil {
		signIn.Alias.SetText(alias)
	}
	if pwd, ok := os.LookupEnv("PASSWORD"); ok {
		signIn.Password.SetText(pwd)
		// TODO if alias was also set auto click
	}
	signIn.SignInButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Sign In Tapped")
		alias := signIn.Alias.Text
		password := []byte(signIn.Password.Text)
		if len(password) < MIN_PASSWORD {
			c.ShowError(errors.New(fmt.Sprintf(ERROR_PASSWORD_TOO_SHORT, len(password), MIN_PASSWORD)))
			return
		}
		c.ExistingNode(alias, password, callback)
	}
	signUp.SignUpButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Sign Up Tapped")
		alias := signUp.Alias.Text
		password := []byte(signUp.Password.Text)
		confirm := []byte(signUp.Confirm.Text)

		err := aliasgo.ValidateAlias(alias)
		if err != nil {
			c.ShowError(err)
			return
		}

		if len(password) < MIN_PASSWORD {
			c.ShowError(errors.New(fmt.Sprintf(ERROR_PASSWORD_TOO_SHORT, len(password), MIN_PASSWORD)))
			return
		}
		if !bytes.Equal(password, confirm) {
			c.ShowError(errors.New(cryptogo.ERROR_PASSWORDS_DO_NOT_MATCH))
			return
		}
		c.NewNode(alias, password, callback)
	}
	c.Dialog.Show()
}

func (c *Client) ShowAccount() {
	//
}

func (c *Client) ShowError(err error) {
	if c.Dialog != nil {
		c.Dialog.Hide()
	}
	c.Dialog = dialog.NewError(err, c.Window)
	c.Dialog.Show()
}

func (c *Client) ShowExperiment() {
	log.Println("ShowExperiment")
	experiment := c.GetExperiment()
	tabber := widget.NewTabContainer()
	center := tabber
	items := make(map[string]*widget.TabItem)
	editors := make(map[string]*edit.ChannelEditor)

	tree := edit.NewTree(experiment.Path, c.Cache, c.Network, func(id string, path ...string) {
		log.Println("Tapped:", id)
		e, ok := editors[id]
		if !ok {
			e = edit.NewChannelEditor(c.GetOrOpenDeltaChannel(id), c.GetNode())
			editors[id] = e
		}
		item, ok := items[id]
		if !ok {
			name := id
			if len(path) > 0 {
				name = path[len(path)-1]
			}
			item = widget.NewTabItem(name, widget.NewVScrollContainer(e))
			items[id] = item
			tabber.Append(item)
		}
		tabber.SelectTab(item)
		if len(items) == 1 {
			// First tab, resize tabber
			tabber.Resize(tabber.MinSize())
		}

	})
	left := widget.NewVScrollContainer(tree)

	chat := widget.NewLabel("Chat")
	status := widget.NewLabel("Ready")
	right := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, nil, nil), status, chat)

	splitter := widget.NewHSplitContainer(left, center)
	splitter.Offset = 0.25
	//splitter = widget.NewVSplitContainer(splitter, terminal)
	//splitter.Offset = 0.75
	splitter = widget.NewHSplitContainer(splitter, right)
	splitter.Offset = 0.75
	c.Window.SetContent(splitter)

	newItem := fyne.NewMenuItem("New", nil)
	newItem.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("File", func() {
			fmt.Println("Menu New->File")
		}),
		fyne.NewMenuItem("Directory", func() {
			fmt.Println("Menu New->Directory")
		}),
	)
	settingsItem := fyne.NewMenuItem("Settings", func() {
		fmt.Println("Menu Settings")
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		ui.ShortcutFocused(&fyne.ShortcutCut{
			Clipboard: c.Window.Clipboard(),
		}, c.Window)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		ui.ShortcutFocused(&fyne.ShortcutCopy{
			Clipboard: c.Window.Clipboard(),
		}, c.Window)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		ui.ShortcutFocused(&fyne.ShortcutPaste{
			Clipboard: c.Window.Clipboard(),
		}, c.Window)
	})
	findItem := fyne.NewMenuItem("Find", func() {
		fmt.Println("Menu Find")
	})

	helpMenu := fyne.NewMenu("Help", fyne.NewMenuItem("Help", func() {
		fmt.Println("Help Menu")
	}))
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		fyne.NewMenu("File", newItem, fyne.NewMenuItemSeparator(), settingsItem),
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)
	c.Window.SetMainMenu(mainMenu)
}

func (c *Client) ShowExperimentDialog(callback func(string, *labgo.Experiment)) {
	log.Println("ShowExperimentDialog")
	create := experiment.NewCreateExperiment(c.Window)
	join := experiment.NewJoinExperiment()

	var content fyne.CanvasObject
	if c.App.Driver().Device().IsMobile() {
		content = widget.NewVBox(
			c.GetLogo(),
			widget.NewTabContainer(
				widget.NewTabItem("Join", join),
				widget.NewTabItem("Create", create),
			),
		)
	} else {
		content = fyne.NewContainerWithLayout(
			layout.NewCenterLayout(),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				widget.NewGroup("Join", join),
				widget.NewGroup("Create", create),
			),
		)
	}
	c.Dialog = dialog.NewCustom("Experiment Access", "Cancel", content, c.Window)

	create.CreateButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Create Tapped")
		uri := create.Path.Text
		go func() {
			progress := dialog.NewProgress("Creating", "message", c.Window)
			defer progress.Hide()
			listener := &ui.ProgressMiningListener{Func: progress.SetValue}
			var reader fyne.FileReadCloser
			if uri != "" {
				r, err := fyne.OpenFileFromURI(uri)
				if err != nil {
					dialog.ShowError(err, c.Window)
					return
				}
				defer r.Close()
				reader = r
			}
			experiment, err := labgo.CreateFromReader(c.GetNode(), listener, uri, reader)
			if err != nil {
				dialog.ShowError(err, c.Window)
				return
			}
			callback("localhost", experiment)
		}()
	}
	join.JoinButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Join Tapped")
		host := join.Host.Text
		id := join.ID.Text
		callback(host, &labgo.Experiment{ID: id})
	}
	c.Dialog.Show()
}

func (c *Client) ShowNode() {
	log.Println("ShowNode")
	node := c.GetNode()
	log.Println("Alias:", node.Alias)
}

/*
func (c *Client) ShowExperimentList(fn func(i int, b binding.Binding)) {
		experimentItems := binding.NewStringList()
		go func() {
			// TODO read experiments from chain
		}()
		var experimentCells int
		experimentList := &widget.List{
			Items: experimentItems,
			OnCreateCell: func() fyne.CanvasObject {
				experimentCells++
				log.Println("Created Label Cell:", experimentCells)
				return &widget.Label{
					Wrapping: fyne.TextWrapBreak,
				}
			},
			OnBindCell: func(object fyne.CanvasObject, data binding.Binding) {
				t, ok := object.(*widget.Label)
				if ok {
					s, ok := data.(binding.String)
					if ok {
						t.Text = s.Get()
					}
					t.Show()
				}
			},
			OnSelected: func(i int, b binding.Binding) {
				log.Println("Selected:", i, b)
				if fn != nil {
					fn(i, b)
				}
			},
		}
		c.Window.SetContent(experimentList)
}
*/
