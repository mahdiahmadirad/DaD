package version

import "strings"

func Normalize(value string) (string, bool) {
	if strings.HasPrefix(value, "v") {
		value = strings.TrimPrefix(value, "v")
	}
	if value == "" || strings.TrimSpace(value) != value {
		return "", false
	}

	coreAndPre := value
	if index := strings.IndexByte(coreAndPre, '+'); index >= 0 {
		if !validIdentifiers(coreAndPre[index+1:], false) {
			return "", false
		}
		coreAndPre = coreAndPre[:index]
	}

	core := coreAndPre
	if index := strings.IndexByte(coreAndPre, '-'); index >= 0 {
		if !validIdentifiers(coreAndPre[index+1:], true) {
			return "", false
		}
		core = coreAndPre[:index]
	}

	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return "", false
	}
	for _, part := range parts {
		if !validCoreNumber(part) {
			return "", false
		}
	}
	return value, true
}

func validCoreNumber(value string) bool {
	if value == "" || len(value) > 1 && value[0] == '0' {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func validIdentifiers(value string, rejectLeadingZeroNumeric bool) bool {
	if value == "" {
		return false
	}
	for _, identifier := range strings.Split(value, ".") {
		if identifier == "" {
			return false
		}
		numeric := true
		for _, character := range identifier {
			switch {
			case character >= '0' && character <= '9':
			case character >= 'A' && character <= 'Z':
				numeric = false
			case character >= 'a' && character <= 'z':
				numeric = false
			case character == '-':
				numeric = false
			default:
				return false
			}
		}
		if rejectLeadingZeroNumeric && numeric &&
			len(identifier) > 1 && identifier[0] == '0' {
			return false
		}
	}
	return true
}
