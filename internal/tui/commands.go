package tui

import (
	tea "charm.land/bubbletea/v2"
	"pulsefeed/internal/aggregator"
)

func waitForSnapshot(snapshots <-chan aggregator.Snapshot) tea.Cmd {
	return func() tea.Msg {
		snapshot, ok := <-snapshots
		if !ok {
			return snapshotsClosedMsg{}
		}
		return snapshotMsg(snapshot)
	}
}
