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

package statustty

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/paxautoma/operos/components/common"

	"github.com/coreos/go-systemd/dbus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Systemd struct {
	conn *dbus.Conn
}

func NewSystemd() (*Systemd, error) {
	conn, err := dbus.NewSystemConnection()
	if err != nil {
		return nil, errors.Wrap(err, "could not initialize systemd connection")
	}

	if err := conn.Subscribe(); err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to systemd events")
	}

	return &Systemd{
		conn: conn,
	}, nil
}

func (s *Systemd) Close() {
	s.conn.Close()
}

type UnitStatusList map[string]*dbus.UnitStatus

func (usl UnitStatusList) CountServices() int {
	res := 0
	for name := range usl {
		if strings.HasSuffix(name, ".service") {
			res++
		}
	}
	return res
}

func (usl UnitStatusList) GetDescriptions() []string {
	descrs := make([]string, len(usl))
	idx := 0
	for _, unit := range usl {
		descrs[idx] = unit.Description
		idx++
	}

	sort.Strings(descrs)
	return descrs
}

type UnitStats struct {
	Active   UnitStatusList
	Starting UnitStatusList
	Stopping UnitStatusList
	Inactive UnitStatusList
	Failed   UnitStatusList
}

func (stats *UnitStats) Delete(name string) {
	delete(stats.Active, name)
	delete(stats.Starting, name)
	delete(stats.Stopping, name)
	delete(stats.Inactive, name)
	delete(stats.Failed, name)
}

func (stats *UnitStats) Update(unit *dbus.UnitStatus) {
	//logrus.Debugf("unit update: %s / %s / %s", unit.Name, unit.ActiveState, unit.SubState)

	stats.Delete(unit.Name)

	switch unit.ActiveState {
	case "active":
		stats.Active[unit.Name] = unit
	case "reloading", "activating":
		stats.Starting[unit.Name] = unit
	case "deactivating":
		stats.Stopping[unit.Name] = unit
	case "failed":
		stats.Failed[unit.Name] = unit
	case "inactive":
		stats.Inactive[unit.Name] = unit
	default:
		logrus.Errorf("unknown state: %s, %s", unit.Name, unit.ActiveState)
	}
}

func (s *Systemd) GetUnitStats() (UnitStats, error) {
	units, err := s.conn.ListUnits()
	if err != nil {
		return UnitStats{}, errors.Wrap(err, "could not get list of systemd units")
	}

	stats := UnitStats{
		Active:   make(UnitStatusList),
		Starting: make(UnitStatusList),
		Stopping: make(UnitStatusList),
		Inactive: make(UnitStatusList),
		Failed:   make(UnitStatusList),
	}
	for _, unit := range units {
		stats.Update(&unit)
	}

	return stats, nil
}

func (s *Systemd) GetBootProgress() (float64, error) {
	progress, err := s.conn.GetManagerProperty("Progress")
	if err != nil {
		return 0, errors.Wrap(err, "could not get systemd Progress property")
	}

	if !strings.HasPrefix(progress, "@d ") {
		return 0, errors.Errorf("systemd Progress property returned invalid value: %v", progress)
	}

	f, err := strconv.ParseFloat(progress[3:], 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse Progress value from systemd")
	}

	return f, nil
}

func (s *Systemd) SubscribeUnitStats(closer <-chan struct{}) (<-chan *UnitStats, <-chan error) {
	chStats := make(chan *UnitStats, 1)
	chOutErr := make(chan error, 1)

	stats := UnitStats{
		Active:   make(UnitStatusList),
		Starting: make(UnitStatusList),
		Stopping: make(UnitStatusList),
		Inactive: make(UnitStatusList),
		Failed:   make(UnitStatusList),
	}

	chUnits, chErr := s.conn.SubscribeUnits(time.Second / 2)
	go func() {
		defer common.LogPanic()

		for {
			select {
			case units := <-chUnits:
				for name, unit := range units {
					if unit == nil {
						logrus.Debugf("unit deleted: %s", name)
						stats.Delete(name)
					} else {
						stats.Update(unit)
					}
				}
				chStats <- &stats
			case err := <-chErr:
				chOutErr <- errors.Wrap(err, "unit update failed")
			case _, ok := <-closer:
				if !ok {
					break
				}
			}
		}
	}()

	return chStats, chOutErr
}

type BootProgress struct {
	Progress float64
	Error    error
}

func (b BootProgress) String() string {
	return fmt.Sprintf("%.f%%", b.Progress*100)
}

func (s *Systemd) SubscribeBootProgress(closer <-chan struct{}) <-chan BootProgress {
	ch := make(chan BootProgress)

	go func() {
		defer common.LogPanic()

		for {
			select {
			case _, ok := <-closer:
				if !ok {
					break
				}
			default:
			}

			progress, err := s.GetBootProgress()
			if err != nil {
				ch <- BootProgress{Error: errors.Wrap(err, "failed to obtain boot progress")}
			} else {
				ch <- BootProgress{Progress: progress}
			}

			time.Sleep(time.Second)
		}
	}()

	return ch
}
