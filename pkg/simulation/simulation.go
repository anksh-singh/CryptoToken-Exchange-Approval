package simulation

import (
	"bridge-allowance/config"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Simulation struct {
	env        *config.Config
	logger     *zap.SugaredLogger
	blowFish   *BlowFish
	signAssist *SignAssist
	tenderly   *Tenderly
}

func NewSimulation(env *config.Config, logger *zap.SugaredLogger) *Simulation {
	blowfish := NewBlowFish(env, logger)
	signassist := NewSignAssist(env, logger)
	tenderly := NewTenderlyService(env, logger)
	return &Simulation{env, logger, blowfish, signassist, tenderly}
}

func (s *Simulation) SimulateTx(request SimulateTxRequest) (*SimulateTxResponse, error) {
	if val, ok := simulationSupportedChains[request.Chain]; ok {
		tx, err := s.blowFish.NewSimulate(request)
		if err != nil {
			err = nil
			tx, err = s.signAssist.SimulateTx(request)
			if err != nil {
				err = nil
				tx, err = s.tenderly.SimulateTx(request)
				if err != nil {
					return nil, err
				}
			}
		}
		return tx, err
	} else {
		return nil, status.Errorf(codes.Unavailable, "Unsupported Chain %s", val)
	}
}
