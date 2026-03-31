package acpcli

import "time"

// DefaultPermissionTimeout keeps direct ACP provider approval waits aligned
// with the hub's fail-closed permission window.
const DefaultPermissionTimeout = 2 * time.Hour
