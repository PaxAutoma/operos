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

package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	stdlog "log"
	"net/http"
	_ "net/http/pprof"
	"os/exec"

	"github.com/paxautoma/operos/components/common"
	"github.com/paxautoma/operos/components/common/widgets"
	statustty "github.com/paxautoma/operos/components/statustty/pkg"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"github.com/wercker/journalhook"
)

var canExit = flag.Bool("e", false, "allow user to exit via Ctrl-C")
var useBootIface = flag.Bool("boot-if", false, "instead of public/private settings, use the interface indicated by the BOOT_IF in the kernel command line")
var publicIface = flag.String("public-if", "", "public interface")
var privateIface = flag.String("private-if", "", "private interface")
var kubeURL = flag.String("kube-url", "", "kubenetes API server URL")
var nodeType = flag.String("node-type", "Controller", "type of node (Controller/Worker)")
var debugAddr = flag.String("debug-addr", "", "enable debug server on this address")

// Set through linker flags
var operosVersion string

func main() {
	defer common.LogPanic()

	// Log directly to journald
	journalhook.Enable()
	log.SetLevel(log.DebugLevel)
	log.SetOutput(ioutil.Discard)
	stdlog.SetOutput(log.StandardLogger().Writer())

	flag.Parse()

	log.Info("starting statustty")

	s, err := statustty.NewSystemd()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	gui, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Fatalf("error initializing gui: %v", err)
	}
	defer gui.Close()

	setupGui(gui, *nodeType)

	closer := make(chan struct{})
	defer close(closer)

	var ifaces []statustty.IfaceSpec
	if *useBootIface {
		bootIf, err := statustty.GetBootIface()
		if err != nil {
			log.Fatalf("could not obtain boot interface: %v", err)
		}

		ifaces = append(ifaces, statustty.IfaceSpec{
			Title:  "Private",
			Device: bootIf,
		})
	} else {
		ifaces = append(ifaces, statustty.IfaceSpec{
			Title:  "Private",
			Device: *privateIface,
		})
	}

	if *publicIface != "" {
		ifaces = append(ifaces, statustty.IfaceSpec{
			Title:  "Public",
			Device: *publicIface,
		})
	}

	go func() {
		chBoot := s.SubscribeBootProgress(closer)
		chHost := statustty.SubscribeHostname(closer)
		chUnit, chErr := s.SubscribeUnitStats(closer)
		chNet := statustty.SubscribeNetStatus(ifaces, closer)
		chKube := statustty.SubscribeKubeStatus(*kubeURL, closer)

		var unitStats *statustty.UnitStats
		var netStatus *statustty.NetStatus
		var kubeStatus *statustty.KubeStatus
		var hostname *string

		for {
			select {
			case prog := <-chBoot:
				gui.Update(func(g *gocui.Gui) error {
					if prog.Error != nil {
						log.Errorf("error receiving boot status: %v", err)
					} else {
						showProgressBar(g, prog.Progress)
					}
					return nil
				})
				continue
			case err := <-chErr:
				log.Errorf("error receiving unit stats update: %v", err)
			case unitStats = <-chUnit:
			case netStatus = <-chNet:
				if netStatus.Error != nil {
					log.Errorf("error obtaining network status: %v", err)
					continue
				}
			case kubeStatus = <-chKube:
			case hostname = <-chHost:
			}

			gui.Update(func(g *gocui.Gui) error {
				showSystemStatus(g, unitStats, netStatus, kubeStatus, hostname)
				return nil
			})
		}
	}()

	if *debugAddr != "" {
		go func() {
			http.ListenAndServe(*debugAddr, nil)
		}()
	}

	gui.MainLoop()
}

func setupGui(g *gocui.Gui, nodeType string) {
	g.InputEsc = true

	g.SetManagerFunc(func(gui *gocui.Gui) error {
		w, h := g.Size()

		if v, err := g.SetView("main", -1, -1, w, 3); err != nil {
			if err != gocui.ErrUnknownView {
				log.Fatalf("error setting main view: %v", err)
			}

			if *canExit {
				g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
					return gocui.ErrQuit
				})
			}

			g.SetKeybinding("main", gocui.KeyCtrlD, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				sendDiagnostics(g)
				return nil
			})

			v.BgColor = gocui.ColorWhite
			v.FgColor = gocui.ColorBlack
			v.Frame = false
			fmt.Fprintf(v, "\n\n Pax Automa Operos \033[34;1m%s\033[0m (v%s)\n", nodeType, operosVersion)

			if _, err := g.SetCurrentView("main"); err != nil {
				log.Fatalf("Cannot activate main view: %s", err.Error())
			}
		}

		if v, err := g.SetView("footer", -1, h-2, w, h); err != nil {
			if err != gocui.ErrUnknownView {
				log.Fatalf("error setting footer view: %v", err)
			}

			v.BgColor = gocui.ColorWhite
			v.FgColor = gocui.ColorBlack
			v.Frame = false
			fmt.Fprint(v, " <Alt-Right> Log in to system    <Ctrl-D> Send diagnostics")
		}

		return nil
	})
}

func showProgressBar(g *gocui.Gui, progress float64) {
	screenW, _ := g.Size()
	barW := int(float64(screenW) * progress)

	v, err := g.SetView("progressbar", -1, 2, barW, 4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			log.Fatalf("cannot set progress bar view: %v", err)
		}

		v.Frame = false
	}

	if progress >= 1 {
		v.BgColor = gocui.ColorGreen
	} else {
		v.BgColor = gocui.ColorYellow
	}
}

func showSystemStatus(g *gocui.Gui, unitStats *statustty.UnitStats, netStatus *statustty.NetStatus, kubeStatus *statustty.KubeStatus, hostname *string) {
	w, h := g.Size()
	v, err := g.SetView("status", 0, 4, w, h-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			log.Fatalf("cannot set progress bar view: %v", err)
		}

		v.Frame = false
		v.Wrap = false
		v.BgColor = gocui.ColorBlack
	}

	v.Clear()

	fmt.Fprintf(v, "\033[37;1mSystem status:\033[0m\n\n")

	if hostname != nil {
		fmt.Fprintf(v, "• Hostname: %s\n\n", *hostname)
	}

	if netStatus != nil {
		if netStatus.IsOK() {
			fmt.Fprintf(v, "\033[32;1m• Network ok\033[0m\n\n")
		} else {
			fmt.Fprintf(v, "\033[31;1m• Network degraded\033[0m\n\n")
		}

		for _, status := range netStatus.Statuses {
			fmt.Fprintf(v, "    %s interface (%s):\n", status.Spec.Title, status.Spec.Device)
			fmt.Fprintf(v, "        IP: %s\n", status.NiceIP())

			if *nodeType == "Controller" && status.IP != nil {
				fmt.Fprintf(v, "        Operos UI: \033[4mhttp://%s\033[0m\n", status.IP.IP)
			}

			fmt.Fprintln(v)
		}
	}

	if unitStats != nil {
		switch {
		case len(unitStats.Stopping) > 0:
			fmt.Fprintf(v, "\033[33;1m• Services stopping\033[0m (%d running)\n", unitStats.Active.CountServices())
		case len(unitStats.Starting) > 0:
			fmt.Fprintf(v, "\033[33;1m• Services starting\033[0m (%d running)\n", unitStats.Active.CountServices())
		case len(unitStats.Failed) > 0:
			fmt.Fprintf(v, "\033[31;1m• Services degraded\033[0m (%d running, %d failed)\n", unitStats.Active.CountServices(), len(unitStats.Failed))
		default:
			fmt.Fprintf(v, "\033[32;1m• Services ok\033[0m (%d running)\n", unitStats.Active.CountServices())
		}

		if len(unitStats.Starting) > 0 || len(unitStats.Failed) > 0 || len(unitStats.Stopping) > 0 {
			fmt.Fprintln(v)

			for _, desc := range unitStats.Starting.GetDescriptions() {
				fmt.Fprintf(v, "    > Starting %s\n", desc)
			}
			for _, desc := range unitStats.Stopping.GetDescriptions() {
				fmt.Fprintf(v, "    > Stopping %s\n", desc)
			}
			for _, desc := range unitStats.Failed.GetDescriptions() {
				fmt.Fprintf(v, "    > \033[31;1mFailed\033[0m %s\n", desc)
			}
		}

		fmt.Fprintln(v)
	}

	if kubeStatus != nil {
		if kubeStatus.Reachable {
			fmt.Fprintf(v, "\033[32;1m• Kubernetes API is reachable\033[0m\n")
		} else {
			fmt.Fprintf(v, "\033[33;1m• Waiting for Kubernetes API\n")
		}
		fmt.Fprintln(v)
	}
}

func sendDiagnostics(g *gocui.Gui) {
	screenW, screenH := g.Size()
	windowW, windowH := int(float32(screenW)*0.8), int(float32(screenH)*0.8)
	windowL := int(float32(screenW) * 0.1)
	windowT := int(float32(screenH) * 0.1)

	if v, err := g.SetView("diag-window", windowL, windowT, windowL+windowW, windowT+windowH); err != nil {
		if err != gocui.ErrUnknownView {
			log.Fatalf("Failed to create diagnostics window: %s", err.Error())
		}

		v.Frame = true
		v.Wrap = true
		v.Autoscroll = true

		fmt.Fprintln(v, widgets.BoldString(gocui.ColorWhite, " Diagnostics"))
		fmt.Fprintln(v)
		fmt.Fprintln(v, " This tool will collect diagnostic information from this machine and upload it")
		fmt.Fprintln(v, " to Pax Automa.")
		fmt.Fprintln(v)
		fmt.Fprintln(v, " Please press Enter to begin or Esc to cancel")
	}

	status := "confirm"

	if _, err := g.SetCurrentView("diag-window"); err != nil {
		log.Fatalf("Failed to activate diagnostics window")
	}

	output := widgets.NewPar("diagnostics", "")
	output.Bounds = image.Rect(1, 1, windowW-3, windowH-3)

	cmd := exec.Command("operosdiag")
	executor := common.NewCmdExecutor(cmd)
	executor.SuccessMessage = widgets.BoldString(gocui.ColorGreen,
		"\nThe diagnostic package has been submitted.\nPress Enter to continue.")
	executor.FailMessage = widgets.BoldString(gocui.ColorRed,
		"\nThe diagnostic package could not be submitted.\nPlease contact Pax Automa for support.\nPress Enter to continue.")

	executor.OnFinish = func(success bool) {
		status = "done"
	}

	cleanup := func() {
		output.Destroy(g)
		g.DeleteKeybinding("diag-window", gocui.KeyEnter, gocui.ModNone)
		g.DeleteKeybinding("diag-window", gocui.KeyEsc, gocui.ModNone)
		g.DeleteView("diag-window")

		if _, err := g.SetCurrentView("main"); err != nil {
			log.Fatalf("Failed to activate main view: %s", err.Error())
		}
	}

	g.SetKeybinding("diag-window", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		switch status {
		case "confirm":
			status = "running"

			if err := output.Render(g, image.Point{X: windowL + 1, Y: windowT + 2}, nil); err != nil {
				log.Fatalf("Failed to render diagnostics window: %s", err.Error())
			}

			if err := executor.Start(output); err != nil {
				log.Fatalf("Failed to start diagnostic utility: %s", err.Error())
			}
		case "done":
			cleanup()
		}

		return nil
	})

	g.SetKeybinding("diag-window", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if status == "confirm" {
			cleanup()
		}
		return nil
	})
}
