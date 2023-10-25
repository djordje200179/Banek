package results

type none struct{}

func (none none) String() string { return "--none--" }

var None Result = none{}
