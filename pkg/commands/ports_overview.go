package commands

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/hunchulchoi/lazydocker/pkg/i18n"
)

type PortOverviewSortColumn int

const (
	PortSortContainer PortOverviewSortColumn = iota
	PortSortService
	PortSortState
	PortSortInternal
	PortSortExternal
	PortSortHost
	PortSortNetworks
)

type PortOverviewRow struct {
	Container string
	Service   string
	State     string
	Internal  string
	External  string
	Host      string
	Networks  string
	Conflict  bool
}

func CollectAllPorts(containers []*Container) []PortOverviewRow {
	rows := make([]PortOverviewRow, 0)

	for _, container := range containers {
		if !container.DetailsLoaded() {
			continue
		}

		rows = append(rows, collectContainerOverviewRows(container)...)
	}

	markPortConflicts(rows)

	return rows
}

func collectContainerOverviewRows(container *Container) []PortOverviewRow {
	networks := containerNetworkNames(container.Details.NetworkSettings.Networks)
	networksLabel := strings.Join(networks, ", ")
	if networksLabel == "" {
		networksLabel = "-"
	}

	service := container.ServiceName
	if service == "" {
		service = "-"
	}

	portRows := publishedPortRows(container.Details.NetworkSettings.Ports, networksLabel)
	if len(portRows) == 0 {
		portRows = exposedPortRows(container.Details.Config.ExposedPorts, networksLabel)
	}

	rows := make([]PortOverviewRow, 0, len(portRows))
	for _, portRow := range portRows {
		rows = append(rows, PortOverviewRow{
			Container: container.Name,
			Service:   service,
			State:     container.Container.State,
			Internal:  portRow[0],
			External:  portRow[1],
			Host:      portRow[2],
			Networks:  portRow[3],
		})
	}

	return rows
}

func markPortConflicts(rows []PortOverviewRow) {
	counts := map[string]int{}
	for i := range rows {
		if rows[i].External == "-" {
			continue
		}

		key := rows[i].Host + ":" + rows[i].External
		counts[key]++
	}

	for i := range rows {
		if rows[i].External == "-" {
			continue
		}

		key := rows[i].Host + ":" + rows[i].External
		if counts[key] > 1 {
			rows[i].Conflict = true
		}
	}
}

func SortPortOverviewRows(rows []PortOverviewRow, column PortOverviewSortColumn, asc bool) {
	sort.Slice(rows, func(i, j int) bool {
		less := ComparePortOverviewRows(rows[i], rows[j], column)
		if asc {
			return less
		}
		return !less
	})
}

func ComparePortOverviewRows(a, b PortOverviewRow, column PortOverviewSortColumn) bool {
	switch column {
	case PortSortService:
		return lessWithNumeric(a.Service, b.Service)
	case PortSortState:
		return lessWithNumeric(a.State, b.State)
	case PortSortInternal:
		return lessWithNumeric(a.Internal, b.Internal)
	case PortSortExternal:
		return lessPortNumber(a.External, b.External)
	case PortSortHost:
		return lessWithNumeric(a.Host, b.Host)
	case PortSortNetworks:
		return lessWithNumeric(a.Networks, b.Networks)
	default:
		return lessWithNumeric(a.Container, b.Container)
	}
}

func lessPortNumber(a, b string) bool {
	aNum, aErr := strconv.Atoi(a)
	bNum, bErr := strconv.Atoi(b)
	if aErr == nil && bErr == nil {
		if aNum != bNum {
			return aNum < bNum
		}
	}
	return lessWithNumeric(a, b)
}

func lessWithNumeric(a, b string) bool {
	if a == b {
		return false
	}
	return a < b
}

func PortOverviewHeaderLabels(tr *i18n.TranslationSet) []string {
	return []string{
		tr.PortsContainerHeader,
		tr.PortsServiceHeader,
		tr.PortsStateHeader,
		tr.PortsInternalHeader,
		tr.PortsExternalHeader,
		tr.PortsHostHeader,
		tr.PortsNetworkHeader,
	}
}

func PortOverviewHeaders(headers []string, column PortOverviewSortColumn, asc bool) []string {
	labeled := append([]string(nil), headers...)
	if int(column) < len(labeled) {
		if asc {
			labeled[column] = labeled[column] + " ↑"
		} else {
			labeled[column] = labeled[column] + " ↓"
		}
	}
	return labeled
}

func PortOverviewDisplayCells(row PortOverviewRow) []string {
	return []string{
		row.Container,
		row.Service,
		row.State,
		row.Internal,
		row.External,
		row.Host,
		row.Networks,
	}
}

func countContainersWithoutDetails(containers []*Container) int {
	count := 0
	for _, container := range containers {
		if !container.DetailsLoaded() {
			count++
		}
	}

	return count
}

func (c *DockerCommand) LoadAllContainerDetails(containers []*Container) error {
	c.ContainerMutex.Lock()
	defer c.ContainerMutex.Unlock()

	sem := make(chan struct{}, 10)
	var wg sync.WaitGroup
	var firstErr error
	var errMutex sync.Mutex

	for _, container := range containers {
		if container.DetailsLoaded() {
			continue
		}

		container := container
		wg.Add(1)
		go func() {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()

			details, err := c.Client.ContainerInspect(context.Background(), container.ID)
			if err != nil {
				c.Log.Error(err)
				errMutex.Lock()
				if firstErr == nil {
					firstErr = err
				}
				errMutex.Unlock()
				return
			}

			container.Details = details
		}()
	}

	wg.Wait()
	return firstErr
}
