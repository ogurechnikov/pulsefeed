package tui

import "pulsefeed/internal/aggregator"

type snapshotMsg aggregator.Snapshot

type snapshotsClosedMsg struct{}
