package mikrotik

type SystemResource struct {
	Uptime       string `json:"uptime"`
	Version      string `json:"version"`
	BuildTime    string `json:"build-time"`
	Platform     string `json:"platform"`
	Board        string `json:"board-name"`
	Architecture string `json:"architecture-name"`
	CPULoad      string `json:"cpu-load"`
	CPUCount     string `json:"cpu-count"`
	FreeMem      string `json:"free-memory"`
	TotalMem     string `json:"total-memory"`
	FreeHDD      string `json:"free-hdd-space"`
	TotalHDD     string `json:"total-hdd-space"`
}

type Interface struct {
	ID         string `json:".id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	MACAddress string `json:"mac-address"`
	MTU        string `json:"mtu"`
	TxByte     string `json:"tx-byte"`
	RxByte     string `json:"rx-byte"`
	TxPacket   string `json:"tx-packet"`
	RxPacket   string `json:"rx-packet"`
	Disabled   string `json:"disabled"`
	Running    string `json:"running"`
	Comment    string `json:"comment"`
}

type IPAddress struct {
	ID        string `json:".id"`
	Address   string `json:"address"`
	Network   string `json:"network"`
	Interface string `json:"interface"`
	Disabled  string `json:"disabled"`
}

type Route struct {
	ID            string `json:".id"`
	DstAddress    string `json:"dst-address"`
	Gateway       string `json:"gateway"`
	GatewayStatus string `json:"gateway-status"`
	Distance      string `json:"distance"`
	Active        string `json:"active"`
	Static        string `json:"static"`
}

type DHCPLease struct {
	ID         string `json:".id"`
	Address    string `json:"address"`
	MACAddress string `json:"mac-address"`
	Hostname   string `json:"host-name"`
	Status     string `json:"status"`
	ExpiresAt  string `json:"expires-after"`
	Server     string `json:"server"`
	LastSeen   string `json:"last-seen"`
	Comment    string `json:"comment"`
}

type Snapshot struct {
	System     *SystemResource   `json:"system"`
	Interfaces []Interface       `json:"interfaces"`
	Addresses  []IPAddress       `json:"addresses"`
	Routes     []Route           `json:"routes"`
	DHCPLeases []DHCPLease       `json:"dhcp_leases"`
	Errors     map[string]string `json:"errors,omitempty"`
}
