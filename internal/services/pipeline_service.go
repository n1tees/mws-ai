package services

type DummyPipeline struct{}

func (p *DummyPipeline) Run(analysisID uint) error {
	return nil
}
