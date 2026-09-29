package usecase

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/Syipmong/acoustic-sensor/internal/domain"
)

type AlertRepository interface {
	LogAlert(ctx context.Context, alert *domain.Alert) error
}

type AlertUsecase struct {
	repo      AlertRepository
	publicKey *ecdsa.PublicKey
}

func NewAlertUsecase(repo AlertRepository, pubKey *ecdsa.PublicKey) *AlertUsecase {
	return &AlertUsecase{
		repo:      repo,
		publicKey: pubKey,
	}
}

func (u *AlertUsecase) ProcessIncomingFrame(ctx context.Context, rawPayload []byte, alert *domain.Alert) error {
	if len(rawPayload) < 22 || len(alert.SignatureRS) != 28 {
		return errors.New("invalid frame or signature length")
	}
	digest := sha256.Sum256(rawPayload[:22])

	r := new(big.Int).SetBytes(alert.SignatureRS[:14])
	s := new(big.Int).SetBytes(alert.SignatureRS[14:])

	isValid := ecdsa.Verify(u.publicKey, digest[:], r, s)
	if !isValid {
		return errors.New("Cryptographic Verification Failed, Package Dropped")
	}

	err := u.repo.LogAlert(ctx, alert)
	if err != nil {
		return err
	}
	return nil

}
