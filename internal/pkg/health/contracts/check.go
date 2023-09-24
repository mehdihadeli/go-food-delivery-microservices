package contracts

type Check map[string]Status

func (check Check) AllUp() bool {
	for _, status := range check {
		if !status.IsUp() {
			return false
		}
	}

	return true
}
