package action

func parseCreate(json map[string]any) (Create, error) {
	action, err := parseAction(json)
	if err != nil {
		return Create{}, err
	}
	return Create{action: action}, nil
}
