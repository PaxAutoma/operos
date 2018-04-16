/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package screens

import (
	"image"

	"github.com/jroimartin/gocui"
	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
)

func OrgInfoScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Organization information"
	screen.Message = `Please tell us about your organization. This information will be used in the
a TLS certificates for this cluster.`

	menu := widgets.NewEditableList("menu-general-info", []*widgets.EditableListItem{
		widgets.NewEditableListItem("Cluster name", "cluster", ctx.Responses.OrgInfo.Cluster, widgets.ValidateNotEmpty),
		widgets.NewEditableListItem("Organization", "organization", ctx.Responses.OrgInfo.Organization, widgets.ValidateNotEmpty),
		widgets.NewEditableListItem("Department", "department", ctx.Responses.OrgInfo.Department, nil),
		widgets.NewEditableListItem("City/town", "city", ctx.Responses.OrgInfo.City, widgets.ValidateNotEmpty),
		widgets.NewEditableListItem("Province/state", "province", ctx.Responses.OrgInfo.Province, nil),
		widgets.NewEditableListItem("Country", "country", ctx.Responses.OrgInfo.Country, widgets.ValidateNotEmpty),
	}, 80, 8)

	errorList := widgets.NewPar("par-errors", "")
	errorList.Bounds = image.Rect(1, 0, 79, 5)
	errorList.FgColor = gocui.ColorRed

	var valid bool
	validate := func() {
		errors := menu.Validate()
		errorList.SetText(widgets.JoinValidationErrors(errors))
		valid = len(errors) == 0
		screen.ShowNext(valid)
	}
	validate()

	menu.OnItemChange = func(item *widgets.EditableListItem) {
		switch item.Key {
		case "cluster":
			ctx.Responses.OrgInfo.Cluster = item.Value
		case "organization":
			ctx.Responses.OrgInfo.Organization = item.Value
		case "department":
			ctx.Responses.OrgInfo.Department = item.Value
		case "city":
			ctx.Responses.OrgInfo.City = item.Value
		case "province":
			ctx.Responses.OrgInfo.Province = item.Value
		case "country":
			ctx.Responses.OrgInfo.Country = item.Value
		}

		validate()
	}

	vl := widgets.NewVerticalLayout()
	vl.Items = []widgets.Renderable{
		menu,
		errorList,
	}

	screen.Content = vl

	screen.OnPrev = func() error {
		screenSet.Back(2)
		return nil
	}

	screen.OnNext = func() error {
		if valid {
			screenSet.Forward(1)
		}
		return nil
	}

	return screen
}
