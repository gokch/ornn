package config

type Global struct {
	Import    []*Import `json:"import"`
	FieldType []string  `json:"field_type"`
}

func (t *Global) ConvFieldType(fldType string) string {
	for _, fld := range t.FieldType {
		if fld == fldType {
			return fldType
		}
	}
	return ""
}

type Import struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
