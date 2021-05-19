package i18n

const (
	defaultLanguage = EN
)

// String provides a simple map of locales to texts.
type String map[Lang]string

// String provides a single string from the map using the
// language requested or resorting to the default.
func (s String) String(lang Lang) string {
	if v, ok := s[lang]; ok {
		return v
	}
	if v, ok := s[defaultLanguage]; ok {
		return v
	}
	for _, v := range s {
		return v // provide first entry
	}
	return ""
}
