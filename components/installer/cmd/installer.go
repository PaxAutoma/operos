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
	"bufio"
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"path"
	"strings"

	"github.com/paxautoma/operos/components/common"
	"github.com/paxautoma/operos/components/common/widgets"
	installer "github.com/paxautoma/operos/components/installer/pkg"
	"github.com/paxautoma/operos/components/installer/pkg/network"
	"github.com/paxautoma/operos/components/installer/pkg/screens"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

var gatekeeperAddress = flag.String("gatekeeper", "gatekeeper.paxautoma.com:57345", "address of the Gatkeeper server (host:port)")
var noGatekeeperTLS = flag.Bool("no-gatekeeper-tls", false, "do not use TLS with Gatekeeper")
var versionsFile = flag.String("versions", "versions", "filename with the versions of packages to be installed")
var logFile = flag.String("logfile", "/root/logs/operos-installer.log", "filename of the log file")

// Set through linker args
var operosVersion string

func main() {

	flag.Parse()

	setupLogging(*logFile)
	defer common.LogPanic()

	log.Infof("Installer (Operos v%s) starting", operosVersion)

	context := installer.DefaultContext

	loadVersions(*versionsFile, &context)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.InputEsc = true

	context.G = g
	context.Net = &network.NetworkdConfigurator{}
	context.GatekeeperAddress = *gatekeeperAddress
	context.GatekeeperTLS = !*noGatekeeperTLS
	context.OperosVersion = operosVersion

	screenSet := widgets.NewScreenSet(g, &context)
	screenSet.Screens = []widgets.ScreenCreator{
		screens.IntroScreen,
		screens.EULAScreen,
		screens.HardwareProbeScreen,
		screens.OrgInfoScreen,
		screens.NetworkSettingsPrivateScreen,
		screens.NetworkSettingsPublicIfaceScreen,
		screens.NetworkSettingsPublicIpsScreen,
		screens.NetworkSettingsIpsScreen,
		screens.DiskSelectionScreen,
		screens.StorageSettingsScreen,
		screens.PasswordScreen,
		screens.ConfirmationScreen,
		screens.InstallScreen,
		screens.FinalizeScreen,
		screens.FailScreen,
	}

	screenSet.Start()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	log.Info("Installer finished")
}

func setupLogging(logFile string) {
	logDir := path.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Could not create log directory (%s): %s", logDir, err)
		os.Exit(1)
	}

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open log file (%s): %s", logFile, err)
		os.Exit(1)
	}

	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)

	// Forward the Go "log" module to logrus too
	stdlog.SetOutput(log.StandardLogger().Writer())
}

func loadVersions(versionsFile string, ctx *installer.InstallerContext) {
	file, err := os.Open(versionsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot open 'versions' file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ctx.Versions = []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ctx.Versions = append(ctx.Versions, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Could not read 'versions' file: ", err)
		os.Exit(1)
	}
}
