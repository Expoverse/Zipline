package mkdir

func mkdir(directory string)  {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		_ = os.Mkdir(directory, 0655)
	}
}
