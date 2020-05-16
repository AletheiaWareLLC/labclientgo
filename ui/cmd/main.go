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

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"github.com/AletheiaWareLLC/bcfynego/ui/account"
)

func main() {
	app := app.New()
	window := app.NewWindow("LAB UI")
	window.SetContent(account.NewSignIn().CanvasObject())
	window.Resize(fyne.NewSize(800, 600))
	window.CenterOnScreen()
	window.ShowAndRun()
}
