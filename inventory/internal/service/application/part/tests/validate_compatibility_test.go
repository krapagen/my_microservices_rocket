package tests

import (
	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

func (s *ServiceSuite) TestValidateCompatibility_Success() {
	hull := newFakeHull(100)
	engine := newFakeEngine(model.EngineClassC, 30)
	slots := model.ShipSlots{
		HullUUID:   hull.UUID(),
		EngineUUID: engine.UUID(),
	}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: []uuid.UUID{hull.UUID(), engine.UUID()}}).Return([]model.Part{hull, engine}, nil)

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.NoError(err)
}

func (s *ServiceSuite) TestValidateCompatibility_SuccessFullSlots() {
	hull := newFakeHull(100)
	engine := newFakeEngine(model.EngineClassC, 30)
	shield := newFakeShield(model.ShieldTypeEnergy)
	weapon := newFakeWeapon(model.WeaponTypeMissile)
	slots := model.ShipSlots{
		HullUUID:   hull.UUID(),
		EngineUUID: engine.UUID(),
		ShieldUUID: shield.UUID(),
		WeaponUUID: weapon.UUID(),
	}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: []uuid.UUID{hull.UUID(), engine.UUID(), shield.UUID(), weapon.UUID()}}).Return([]model.Part{hull, engine, shield, weapon}, nil)

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.NoError(err)
}

func (s *ServiceSuite) TestValidateCompatibility_MissingHull() {
	engine := newFakeEngine(model.EngineClassC, 30)
	slots := model.ShipSlots{EngineUUID: engine.UUID()}

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.Error(err)
	s.ErrorIs(err, errs.ErrPartTypeMismatch)
	s.partRepository.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestValidateCompatibility_MissingEngine() {
	hull := newFakeHull(100)
	slots := model.ShipSlots{HullUUID: hull.UUID()}

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.Error(err)
	s.ErrorIs(err, errs.ErrPartTypeMismatch)
	s.partRepository.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestValidateCompatibility_DuplicateUUID() {
	shared := uuid.New()
	slots := model.ShipSlots{
		HullUUID:   shared,
		EngineUUID: shared,
	}

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.Error(err)
	s.ErrorIs(err, errs.ErrPartTypeMismatch)
	s.partRepository.AssertNotCalled(s.T(), "List")
}

func (s *ServiceSuite) TestValidateCompatibility_PartNotFound() {
	hull := newFakeHull(100)
	engineUUID := uuid.New()
	slots := model.ShipSlots{
		HullUUID:   hull.UUID(),
		EngineUUID: engineUUID,
	}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: []uuid.UUID{hull.UUID(), engineUUID}}).Return(nil, errs.ErrPartNotFound)

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestValidateCompatibility_WrongPartType() {
	weapon := newFakeWeapon(model.WeaponTypeLaser)
	engine := newFakeEngine(model.EngineClassC, 30)
	slots := model.ShipSlots{
		HullUUID:   weapon.UUID(),
		EngineUUID: engine.UUID(),
	}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: []uuid.UUID{weapon.UUID(), engine.UUID()}}).Return([]model.Part{weapon, engine}, nil)

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.ErrorIs(err, errs.ErrPartTypeMismatch)
}

func (s *ServiceSuite) TestValidateCompatibility_HullTooWeak() {
	hull := newFakeHull(50)
	engine := newFakeEngine(model.EngineClassA, 100)
	slots := model.ShipSlots{
		HullUUID:   hull.UUID(),
		EngineUUID: engine.UUID(),
	}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: []uuid.UUID{hull.UUID(), engine.UUID()}}).Return([]model.Part{hull, engine}, nil)

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.ErrorIs(err, errs.ErrIncompatibleParts)
}

func (s *ServiceSuite) TestValidateCompatibility_PlasmaLaserConflict() {
	hull := newFakeHull(100)
	engine := newFakeEngine(model.EngineClassC, 30)
	shield := newFakeShield(model.ShieldTypePlasma)
	weapon := newFakeWeapon(model.WeaponTypeLaser)
	slots := model.ShipSlots{
		HullUUID:   hull.UUID(),
		EngineUUID: engine.UUID(),
		ShieldUUID: shield.UUID(),
		WeaponUUID: weapon.UUID(),
	}

	s.partRepository.EXPECT().List(s.ctx, input.PartFilter{UUIDs: []uuid.UUID{hull.UUID(), engine.UUID(), shield.UUID(), weapon.UUID()}}).Return([]model.Part{hull, engine, shield, weapon}, nil)

	err := s.service.ValidateCompatibility(s.ctx, slots)
	s.ErrorIs(err, errs.ErrIncompatibleParts)
}
