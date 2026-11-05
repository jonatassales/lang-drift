package judge

// figureTopics are the topics a figure may speak of: resources only.
var figureTopics = []string{"resource"}

// checkFigure refuses a figure that does not speak of resources, has no fact
// or source, or gives a negative CPU or memory.
func checkFigure(path string, figure Figure) error {
	if err := requireOneOf(path+".topic", figure.Topic, figureTopics); err != nil {
		return err
	}
	if err := requireText(path+".fact", figure.Fact); err != nil {
		return err
	}
	if err := requireText(path+".source", figure.Source); err != nil {
		return err
	}
	if figure.CPUMillis < 0 {
		return refuse(path+".cpuMillis", "must not be negative")
	}
	if figure.MemoryMiB < 0 {
		return refuse(path+".memoryMiB", "must not be negative")
	}

	return nil
}
