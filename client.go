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
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
	"github.com/AletheiaWareLLC/bcfynego"
	bcui "github.com/AletheiaWareLLC/bcfynego/ui"
	bcdata "github.com/AletheiaWareLLC/bcfynego/ui/data"
	"github.com/AletheiaWareLLC/bcgo"
	"github.com/AletheiaWareLLC/labclientgo/ui/data"
	"github.com/AletheiaWareLLC/labclientgo/ui/edit"
	"github.com/AletheiaWareLLC/labclientgo/ui/experiment"
	"github.com/AletheiaWareLLC/labgo"
	"log"
	"os"
)

type Client struct {
	bcfynego.Client
	Experiment *labgo.Experiment
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
		c.Window.SetTitle("LAB - " + c.Experiment.ID)
	}
	return c.Experiment
}

func (c *Client) GetLogo() fyne.CanvasObject {
	return &canvas.Image{
		Resource: bcdata.NewThemedResource(data.LogoUnmasked),
		//FillMode: canvas.ImageFillContain,
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

func (c *Client) ShowExperiment() {
	log.Println("ShowExperiment")
	experiment := c.GetExperiment()
	tabber := widget.NewTabContainer()
	center := tabber
	items := make(map[string]*widget.TabItem)
	editors := make(map[string]*edit.ChannelEditor)
	listener := &bcgo.PrintingMiningListener{Output: os.Stdout}

	selectPath := func(id string, path ...string) {
		log.Println("Selected:", id, path)
		e, ok := editors[id]
		if !ok {
			e = edit.NewChannelEditor(c.GetNode(), listener, c.GetOrOpenDeltaChannel(id))
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

	}
	tree := edit.NewTree(experiment.Path, c.Cache, c.Network, selectPath)
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
			fmt.Println("Show file open dialog")
			fmt.Println("Call selectPath(id, path)")
		}),
		fyne.NewMenuItem("Import", func() {
			fmt.Println("Menu New->Import")
			fmt.Println("Show file open dialog")
		}),
		fyne.NewMenuItem("Export", func() {
			fmt.Println("Menu New->Export")
			fmt.Println("Show file save dialog")
		}),
	)
	settingsItem := fyne.NewMenuItem("Settings", func() {
		fmt.Println("Menu Settings")
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		bcui.ShortcutFocused(&fyne.ShortcutCut{
			Clipboard: c.Window.Clipboard(),
		}, c.Window)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		bcui.ShortcutFocused(&fyne.ShortcutCopy{
			Clipboard: c.Window.Clipboard(),
		}, c.Window)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		bcui.ShortcutFocused(&fyne.ShortcutPaste{
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

	c.Dialog = dialog.NewCustom("Experiment Access", "Cancel",
		widget.NewAccordionContainer(
			&widget.AccordionItem{Title: "Create", Detail: create.CanvasObject(), Open: true},
			widget.NewAccordionItem("Join", join.CanvasObject()),
		), c.Window)

	create.CreateButton.OnTapped = func() {
		c.Dialog.Hide()
		log.Println("Create Tapped")
		uri := create.Path.Text
		go func() {
			progress := dialog.NewProgress("Creating", "message", c.Window)
			defer progress.Hide()
			listener := &bcui.ProgressMiningListener{Func: progress.SetValue}
			var reader fyne.FileReadCloser
			if uri != "" {
				r, err := storage.OpenFileFromURI(storage.NewURI(uri))
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
