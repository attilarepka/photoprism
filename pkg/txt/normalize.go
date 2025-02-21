package txt

import "strings"

// NormalizeName sanitizes and capitalizes names.
func NormalizeName(name string) string {
	if name == "" {
		return ""
	}

	// Remove double quotes and other special characters.
	name = strings.Map(func(r rune) rune {
		switch r {
		case '"', '`', '~', '\\', '/', '*', '%', '&', '|', '+', '=', '$', '@', '!', '?', ':', ';', '<', '>', '{', '}':
			return -1
		}
		return r
	}, name)

	name = strings.TrimSpace(name)

	if name == "" {
		return ""
	}

	// Shorten and capitalize.
	return Clip(Title(name), ClipDefault)
}

// NormalizeState returns the full, normalized state name.
func NormalizeState(stateName, countryCode string) string {
	// Remove whitespace from name.
	stateName = strings.TrimSpace(stateName)

	// Empty?
	if stateName == "" || stateName == UnknownStateCode {
		// State doesn't have a name.
		return ""
	}

	// Normalize country code.
	countryCode = strings.ToLower(strings.TrimSpace(countryCode))

	// Is the name an abbreviation that should be normalized?
	if states, found := StatesByCountry[countryCode]; !found {
		// Unknown country.
	} else if normalized, found := states[stateName]; !found {
		// Unknown abbreviation.
	} else if normalized != "" {
		// Yes, use normalized name.
		stateName = normalized
	}

	// Return normalized state name.
	return stateName
}

// NormalizeQuery replaces search operator with default symbols.
func NormalizeQuery(s string) string {
	s = strings.ToLower(Clip(s, ClipQuery))
	s = strings.ReplaceAll(s, Spaced(EnOr), Or)
	s = strings.ReplaceAll(s, Spaced(EnAnd), And)
	s = strings.ReplaceAll(s, Spaced(EnWith), And)
	s = strings.ReplaceAll(s, Spaced(EnIn), And)
	s = strings.ReplaceAll(s, Spaced(EnAt), And)
	s = strings.ReplaceAll(s, SpacedPlus, And)
	s = strings.ReplaceAll(s, "%", "*")
	return strings.Trim(s, "+&|_-=!@$%^(){}\\<>,.;: ")
}

// NormalizeUsername returns the normalized username (lowercase, whitespace trimmed).
func NormalizeUsername(s string) string {
	return strings.ToLower(Clip(s, ClipUsername))
}
