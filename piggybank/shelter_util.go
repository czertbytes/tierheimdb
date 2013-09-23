package piggybank

func intAnimalTypes(animalTypes []string) int {
	val := 0
	for _, animalType := range animalTypes {
		switch animalType {
		case "cat":
			val |= (1 << 0)
		case "dog":
			val |= (1 << 1)
		}
	}

	return val
}

func animalTypes(intAnimalTypes int) []string {
	animalTypes := []string{}
	if intAnimalTypes&(1<<0) != 0 {
		animalTypes = append(animalTypes, "cat")
	}

	if intAnimalTypes&(1<<1) != 0 {
		animalTypes = append(animalTypes, "dog")
	}

	return animalTypes
}
