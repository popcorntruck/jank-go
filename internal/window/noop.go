package window

type NoopWindowService struct {
}

func NewNoopWindowService() (*NoopWindowService, error) {
	return &NoopWindowService{}, nil
}

func (s *NoopWindowService) Close() error {
	return nil
}

func (s *NoopWindowService) ActiveWindow() *WindowInfo {
	return nil
}
