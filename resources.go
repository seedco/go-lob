package lob

// NamedObjectList is used to return the list of countries and states.
type NamedObjectList struct {
	Object string        `json:"object"`
	Data   []NamedObject `json:"data"`
}

// NamedObject is a datum that contains a name and short name for a state or country.
type NamedObject struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Object    string `json:"object"`
}

// GetStates returns a list of US States that Lob recognizes.
func (lob *Lob) GetStates() (*NamedObjectList, error) {
	resp := new(NamedObjectList)
	return resp, Metrics.GetStates.Call(func() error {
		return lob.Get("states/", nil, resp)
	})
}

// GetCountries returns a list of countries that Lob recognizes.
func (lob *Lob) GetCountries() (*NamedObjectList, error) {
	resp := new(NamedObjectList)
	return resp, Metrics.GetCountries.Call(func() error {
		return lob.Get("countries/", nil, resp)
	})
}
