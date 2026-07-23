package test

import (
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestValidateCompatibility_Success() {
	hull := uuid.New().String()
	engine := uuid.New().String()
	req := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   hull,
		EngineUuid: engine,
	}

	s.partService.EXPECT().ValidateCompatibility(s.ctx, model.ShipSlots{
		HullUUID:   uuid.MustParse(hull),
		EngineUUID: uuid.MustParse(engine),
	}).Return(nil)

	resp, err := s.api.ValidateCompatibility(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestValidateCompatibility_FullSlots() {
	hull := uuid.New().String()
	engine := uuid.New().String()
	shield := uuid.New().String()
	weapon := uuid.New().String()
	req := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   hull,
		EngineUuid: engine,
		ShieldUuid: shield,
		WeaponUuid: weapon,
	}

	s.partService.EXPECT().ValidateCompatibility(s.ctx, model.ShipSlots{
		HullUUID:   uuid.MustParse(hull),
		EngineUUID: uuid.MustParse(engine),
		ShieldUUID: uuid.MustParse(shield),
		WeaponUUID: uuid.MustParse(weapon),
	}).Return(nil)

	resp, err := s.api.ValidateCompatibility(s.ctx, req)
	s.NoError(err)
	s.NotNil(resp)
}

func (s *APISuite) TestValidateCompatibility_MissingHull() {
	engine := uuid.New().String()
	req := &inventoryv1.ValidateCompatibilityRequest{EngineUuid: engine}

	s.partService.EXPECT().ValidateCompatibility(s.ctx, model.ShipSlots{
		HullUUID:   uuid.Nil,
		EngineUUID: uuid.MustParse(engine),
	}).Return(errs.ErrPartTypeMismatch)

	resp, err := s.api.ValidateCompatibility(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrPartTypeMismatch)
	s.Nil(resp)
}

func (s *APISuite) TestValidateCompatibility_InvalidEngineUUID() {
	hull := uuid.New().String()
	req := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   hull,
		EngineUuid: "not-a-uuid",
	}

	s.partService.EXPECT().ValidateCompatibility(s.ctx, model.ShipSlots{
		HullUUID:   uuid.MustParse(hull),
		EngineUUID: uuid.Nil,
	}).Return(errs.ErrPartTypeMismatch)

	resp, err := s.api.ValidateCompatibility(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrPartTypeMismatch)
	s.Nil(resp)
}

func (s *APISuite) TestValidateCompatibility_ServiceError() {
	hull := uuid.New().String()
	engine := uuid.New().String()
	req := &inventoryv1.ValidateCompatibilityRequest{
		HullUuid:   hull,
		EngineUuid: engine,
	}

	s.partService.EXPECT().ValidateCompatibility(s.ctx, model.ShipSlots{
		HullUUID:   uuid.MustParse(hull),
		EngineUUID: uuid.MustParse(engine),
	}).Return(errs.ErrIncompatibleParts)

	resp, err := s.api.ValidateCompatibility(s.ctx, req)
	s.Error(err)
	s.ErrorIs(err, errs.ErrIncompatibleParts)
	s.Nil(resp)
}
