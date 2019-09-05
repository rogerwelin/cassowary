package main

func (c *cassowary) coordinate() error {

	if c.fileMode {
		urlSuffixes, err := readFile(c.inputFile)
		if err != nil {
			return err
		}
	}

	return nil
}
