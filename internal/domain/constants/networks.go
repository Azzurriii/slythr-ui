package constants

// Supported blockchain networks
const (
	NetworkEthereum  = "ethereum"
	NetworkPolygon   = "polygon"
	NetworkBSC       = "bsc"
	NetworkBase      = "base"
	NetworkArbitrum  = "arbitrum"
	NetworkAvalanche = "avalanche"
	NetworkOptimism  = "optimism"
	NetworkGnosis    = "gnosis"
	NetworkFantom    = "fantom"
	NetworkCelo      = "celo"
)

// NetworkChainIDs maps network names to their chain IDs
var NetworkChainIDs = map[string]string{
	NetworkEthereum:  "1",
	NetworkPolygon:   "137",
	NetworkBSC:       "56",
	NetworkBase:      "8453",
	NetworkArbitrum:  "42161",
	NetworkAvalanche: "43114",
	NetworkOptimism:  "10",
	NetworkGnosis:    "100",
	NetworkFantom:    "250",
	NetworkCelo:      "42220",
}

// SupportedNetworks returns a slice of all supported networks
func SupportedNetworks() []string {
	return []string{
		NetworkEthereum,
		NetworkPolygon,
		NetworkBSC,
		NetworkBase,
		NetworkArbitrum,
		NetworkAvalanche,
		NetworkOptimism,
		NetworkGnosis,
		NetworkFantom,
		NetworkCelo,
	}
}

// IsValidNetwork checks if the provided network is supported
func IsValidNetwork(network string) bool {
	_, exists := NetworkChainIDs[network]
	return exists
}

// GetChainID returns the chain ID for the given network
func GetChainID(network string) (string, bool) {
	chainID, exists := NetworkChainIDs[network]
	return chainID, exists
}
