package firewall

type FilterRule struct {
	ID           string `json:".id"`
	Chain        string `json:"chain"`
	Action       string `json:"action"`
	SrcAddress   string `json:"src-address"`
	DstAddress   string `json:"dst-address"`
	Protocol     string `json:"protocol"`
	SrcPort      string `json:"src-port"`
	DstPort      string `json:"dst-port"`
	InInterface  string `json:"in-interface"`
	OutInterface string `json:"out-interface"`
	Comment      string `json:"comment"`
	Disabled     string `json:"disabled"`
	Bytes        string `json:"bytes"`
	Packets      string `json:"packets"`
}

type NATRule struct {
	ID          string `json:".id"`
	Chain       string `json:"chain"`
	Action      string `json:"action"`
	SrcAddress  string `json:"src-address"`
	DstAddress  string `json:"dst-address"`
	Protocol    string `json:"protocol"`
	SrcPort     string `json:"src-port"`
	DstPort     string `json:"dst-port"`
	ToAddresses string `json:"to-addresses"`
	ToPorts     string `json:"to-ports"`
	Comment     string `json:"comment"`
	Disabled    string `json:"disabled"`
	Bytes       string `json:"bytes"`
	Packets     string `json:"packets"`
}

type AdressListEntry struct {
	ID       string `json:".id"`
	List     string `json:"list"`
	Address  string `json:"address"`
	Comment  string `json:"comment"`
	Disabled string `json:"disabled"`
	Timeout  string `json:"timeout"`
}
