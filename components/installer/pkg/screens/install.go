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
	"bytes"
	"fmt"
	"html/template"
	"image"
	"os"
	"os/exec"
	"strings"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"github.com/paxautoma/operos/components/common"
	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
)

func ConfirmationScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Confirm"

	t := template.Must(template.New("confirmation").Parse(`You're about to install the Operos Controller on this system with the following
settings:

    Cluster name:                 {{.OrgInfo.Cluster}}
    Organization:                 {{.OrgInfo.Organization}}{{if .OrgInfo.Department}} ({{.OrgInfo.Department}}){{end}}
    Location:                     {{.OrgInfo.City}}{{if .OrgInfo.Province}}, {{.OrgInfo.Province}}{{end}}, {{.OrgInfo.Country}}

    Private interface:            {{.PrivateInterface}}, {{.PrivateSubnet}}
    Public interface:             {{.PublicNetwork.Interface}}, {{.PublicIPInfo}}
    Pod subnet:                   {{.PodSubnet}}
    Service subnet:               {{.ServiceSubnet}}

    Install to controller disk:   {{.ControllerDisk}}
    Worker code/data storage:     {{.StorageSystemPercentage}}% / {{.StorageDataPercentage}}%

To change these settings, please use the Back button.`))

	var message bytes.Buffer
	err := t.Execute(&message, ctx.Responses)
	if err != nil {
		panic(err)
	}

	screen.Message = message.String()
	screen.ShowNext(false)

	screen.OnPrev = func() error {
		screenSet.Back(1)
		return nil
	}

	buttonStart := widgets.NewButton("button-start", "Install now", 30, 0, 19, 3)
	buttonStart.Handler = func(string) error {
		screenSet.Forward(1)
		return nil
	}

	screen.Content = buttonStart

	return screen
}

func InstallScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)

	screen := widgets.NewScreen()
	screen.Title = "Installing"
	screen.ShowPrev(false)
	screen.ShowNext(false)

	cmd := exec.Command("./install.sh")
	cmd.Env = os.Environ()

	privateSubnetParts := strings.Split(ctx.Responses.PrivateSubnet, "/")

	cmd.Env = append(cmd.Env, ctx.Versions...)
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("CONTROLLER_PRIVATE_IF=%s", ctx.Responses.PrivateInterface),
		fmt.Sprintf("CONTROLLER_PUBLIC_IF=%s", ctx.Responses.PublicNetwork.Interface),
		fmt.Sprintf("CONTROLLER_PUBLIC_IF_MODE=%s", ctx.Responses.PublicNetwork.Mode),
		fmt.Sprintf("CONTROLLER_DISK=%s", ctx.Responses.ControllerDisk),
		fmt.Sprintf("OPEROS_VERSION=%s", ctx.OperosVersion),
		fmt.Sprintf("OPEROS_INSTALL_ID=%s", ctx.InstallID),
		fmt.Sprintf("OPEROS_CONTROLLER_IP=%s", ctx.Responses.ControllerIP),
		fmt.Sprintf("OPEROS_KUBE_API_INSECURE_PORT=%d", 8080),
		fmt.Sprintf("OPEROS_KUBE_API_SECURE_PORT=%d", 8443),
		fmt.Sprintf("OPEROS_NODE_MASK=/%s", privateSubnetParts[1]),
		fmt.Sprintf("OPEROS_POD_CIDR=%s", ctx.Responses.PodSubnet),
		fmt.Sprintf("OPEROS_SERVICE_CIDR=%s", ctx.Responses.ServiceSubnet),
		fmt.Sprintf("OPEROS_DNS_SERVICE_IP=%s", ctx.Responses.DNSIP),
		fmt.Sprintf("OPEROS_DNS_DOMAIN=%s", ctx.Responses.DNSDomain),
		fmt.Sprintf("OPEROS_WORKER_STORAGE_PERCENTAGE=%d", ctx.Responses.StorageSystemPercentage),
		fmt.Sprintf("OPEROS_CLUSTER_NAME=%s", ctx.Responses.OrgInfo.Cluster),
		fmt.Sprintf("OPEROS_CLUSTER_ORG=%s", ctx.Responses.OrgInfo.Organization),
		fmt.Sprintf("OPEROS_CLUSTER_DEPARTMENT=%s", ctx.Responses.OrgInfo.Department),
		fmt.Sprintf("OPEROS_CLUSTER_CITY=%s", ctx.Responses.OrgInfo.City),
		fmt.Sprintf("OPEROS_CLUSTER_PROVINCE=%s", ctx.Responses.OrgInfo.Province),
		fmt.Sprintf("OPEROS_CLUSTER_COUNTRY=%s", ctx.Responses.OrgInfo.Country),
		fmt.Sprintf("INSTALLER_ROOT_PASSWD=%s", ctx.Responses.RootPassword),
	)

	if ctx.Responses.PublicNetwork.Mode == "static" {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("CONTROLLER_PUBLIC_IF_IPNET=%s", ctx.Responses.PublicNetwork.Subnet),
			fmt.Sprintf("CONTROLLER_PUBLIC_IF_GW=%s", ctx.Responses.PublicNetwork.Gateway),
		)
	}

	output := widgets.NewPar("install-sh", "")
	output.Bounds = image.Rect(1, 1, 78, 18)

	executor := common.NewCmdExecutor(cmd)
	executor.SuccessMessage = widgets.ColorString(gocui.ColorGreen, "Installation complete")
	executor.FailMessage = widgets.ColorString(gocui.ColorRed, "Installation failed")

	installSuccessful := false

	executor.OnFinish = func(success bool) {
		ctx.G.Update(func(g *gocui.Gui) error {
			installSuccessful = success
			screen.ShowNext(true)
			screen.FocusableSet.Next()

			return nil
		})
	}

	screen.Content = output

	screen.OnInitialize = func(g *gocui.Gui) {
		go func() {
			fmt.Fprintln(output, "> Generating controller certificates")

			err := installer.CreateControllerCerts(ctx)
			if err != nil {
				fmt.Fprintf(output, "Failed to create controller certificates")
				fmt.Fprintf(output, err.Error())

				installSuccessful = false
				screen.ShowNext(true)
				screen.FocusableSet.Next()

				return
			}

			cmd.Env = append(cmd.Env,
				fmt.Sprintf("INSTALLER_CONTROLLER_KEY=%s", ctx.ControllerKey),
				fmt.Sprintf("INSTALLER_CONTROLLER_CERT=%s", ctx.ControllerCert),
				fmt.Sprintf("INSTALLER_SERVER_KEY=%s", ctx.ServerKey),
				fmt.Sprintf("INSTALLER_SERVER_CERT=%s", ctx.ServerCert),
			)

			err = executor.Start(output)
			if err != nil {
				panic(err)
			}
		}()
	}

	screen.OnNext = func() error {
		if installSuccessful {
			screenSet.Forward(1)
			return nil
		}

		message := fmt.Sprintf(`
	One of the installation steps has failed.

	In order to allow our engineers to debug this issue, you can
	submit the installation logs to Pax Automa. Would you like to
	do this now?`)

		dialog := widgets.NewDialog("confirm-upload", message, 70)
		dialog.OnAccept = func() error {
			dialog.Close()
			screenSet.Forward(2)
			return nil
		}

		dialog.OnCancel = func() error {
			dialog.Close()
			screenSet.Restart()
			return nil
		}

		screen.Blur()
		return dialog.Layout(ctx.G)
	}

	return screen
}

func FinalizeScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	screen := widgets.NewScreen()
	screen.Title = "Installation complete"
	screen.Message = `
The Operos controller has been installed.
The machine must reboot to continue.`

	screen.ShowNext(false)
	screen.ShowPrev(false)

	butReboot := widgets.NewButton("but-reboot", "Reboot", 35, 10, 10, 3)
	butReboot.Handler = func(string) error {
		log.Debug("Rebooting")
		return installer.Reboot()
	}

	screen.Content = butReboot

	return screen
}

func FailScreen(screenSet *widgets.ScreenSet, context interface{}) *widgets.Screen {
	ctx := context.(*installer.InstallerContext)
	pubNet := ctx.Responses.PublicNetwork

	screen := widgets.NewScreen()
	screen.Title = "Upload install logs"

	screen.ShowNext(false)
	screen.ShowPrev(false)

	screen.OnNext = func() error {
		screenSet.Restart()
		return nil
	}

	cmd := exec.Command("/root/installerdiag.sh", "--install", ctx.InstallID)

	output := widgets.NewPar("diagnostics", "")
	output.Bounds = image.Rect(0, 0, 78, 18)

	executor := common.NewCmdExecutor(cmd)
	executor.SuccessMessage = widgets.BoldString(gocui.ColorGreen, "Logs have been uploaded to Pax Automa.")
	executor.FailMessage = widgets.BoldString(gocui.ColorRed, "Failed uploading logs to Pax Automa.")

	screen.Content = output

	screen.OnInitialize = func(g *gocui.Gui) {
		go func() {
			defer common.LogPanic()

			fmt.Fprintf(output, "Activating %s\n", pubNet.Interface)

			if err := ctx.Net.ConfigureInterface(pubNet); err != nil {
				fmt.Fprintf(output, "Failed to configure network interface %s: %s", pubNet.Interface, err)
				log.Errorf("Failed to configure network interface: %+v", err.Error())
				return
			}

			if err := executor.Start(output); err != nil {
				panic(err)
			}
		}()
	}

	executor.OnFinish = func(success bool) {
		ctx.G.Update(func(g *gocui.Gui) error {
			screen.ShowNext(true)
			screen.FocusableSet.Next()
			return nil
		})
	}

	return screen
}
