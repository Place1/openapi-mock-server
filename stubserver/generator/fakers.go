package generator

func FakeString() string {
	values := []string{
		"lorem ipsum",
		"neque porro",
		"dolorem ipsum",
	}
	return values[randInt(0, len(values)-1)]
}
