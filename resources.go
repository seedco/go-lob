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
func (l *lob) GetStates() (*NamedObjectList, error) {
	resp := new(NamedObjectList)
	if err := l.get("states/", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCountries returns a list of countries that Lob recognizes.
func (l *lob) GetCountries() (*NamedObjectList, error) {
	resp := new(NamedObjectList)
	if err := l.get("countries/", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
