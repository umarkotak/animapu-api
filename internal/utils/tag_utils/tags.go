package tag_utils

func Update(current, incoming []string, mode string) ([]string, bool) {
	switch mode {
	case "replace":
		return incoming, true
	case "add":
		seen := make(map[string]bool, len(current)+len(incoming))
		for _, tag := range current {
			seen[tag] = true
		}
		for _, tag := range incoming {
			if !seen[tag] {
				current = append(current, tag)
				seen[tag] = true
			}
		}
		return current, true
	case "remove":
		remove := make(map[string]bool, len(incoming))
		for _, tag := range incoming {
			remove[tag] = true
		}
		tags := current[:0]
		for _, tag := range current {
			if !remove[tag] {
				tags = append(tags, tag)
			}
		}
		return tags, true
	default:
		return nil, false
	}
}
