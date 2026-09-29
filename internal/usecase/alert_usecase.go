package usecase

import (
	"acoustic-sensor-backend/internal/domain"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/Syipmong/acoustic-sensor-backend/internal/domain"
)

type AlertRepository interface{
	LogAlert(ctx context.Context, alert *domain.Alert) error
}

type AlertUsecase struct{
	repo AlertRepository
	publicKey *ecdsa.PublicKey
}

func NewAlertUsecase(repo AlertRepository, pubKey *ecdsa.PublicKey) *AlertUsecase{
	return &AlertUsecase{
		repo: repo,
		publicKey: pubKey,
	}
}

func (u *AlertUsecase) ProcessIncomingFrame(ctx context.Context, rawPayload []byte, alert *domain.Alert) error{
	digest := sha256.Sum256(rawPayload[:22])
}
