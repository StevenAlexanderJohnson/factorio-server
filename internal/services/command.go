package services

type CommandService struct {
	factorioService *FactorioService
}

func NewCommandService(factorioService *FactorioService) *CommandService {
	return &CommandService{
		factorioService: factorioService,
	}
}

func (r *CommandService) SendCommand(cmd string) (string, error) {
	if !r.factorioService.IsRunning() {
		return "", ErrFactorioServerNotRunning
	}
	return r.factorioService.SendCommand(cmd)
}
