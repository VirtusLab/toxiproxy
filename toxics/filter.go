package toxics

import "github.com/sirupsen/logrus"

// The FilterToxic drops traffic from the addresses that are filtered.
type FilterToxic struct{
	FilteredAddresses []string `json:"filtered_addresses"`
}

func (t *FilterToxic) Pipe(stub *ToxicStub) {
	for {
		select {
		case <-stub.Interrupt:
			return
		case c := <-stub.Input:
			if c == nil {
				stub.Close()
				return
			}

			for _, filteredAddress := range t.FilteredAddresses {
				if stub.Meta.DownstreamAddress == filteredAddress {
					logrus.Infof("Encountered connection from filtered address %s, dropping the connection.", filteredAddress)
					stub.Close()
					return
				}
			}

			stub.Output <- c
		}
	}
}

func init() {
	Register("filter", new(FilterToxic))
}

