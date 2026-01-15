package aidbox

import "encoding/json"

type Bundle struct {
	Entry []BundleEntry `json:"entry"`
}

type BundleEntry struct {
	Resource json.RawMessage `json:"resource"`
}
