package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/service/input"
)

// inventory/internal/service/application/part/validate_compatibility.go

// ValidateCompatibility разрешает UUID слотов в детали и проверяет их совместимость.
func (s *service) ValidateCompatibility(ctx context.Context, slots model.ShipSlots) error {
	resolved, err := s.resolveShipSlots(ctx, slots)
	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolved)
}

// resolveShipSlots загружает детали по UUID из слотов, проверяет соответствие слот↔тип
// и собирает результат в model.ResolvedShipSlots.
func (s *service) resolveShipSlots(ctx context.Context, slots model.ShipSlots) (model.ResolvedShipSlots, error) {
	// Проверяем, что обязательные слоты (корпус и двигатель) заполнены.
	if slots.HullUUID == uuid.Nil {
		return model.ResolvedShipSlots{}, fmt.Errorf("отсутствует UUID корпуса: %w", errs.ErrPartTypeMismatch)
	}
	if slots.EngineUUID == uuid.Nil {
		return model.ResolvedShipSlots{}, fmt.Errorf("отсутствует UUID двигателя: %w", errs.ErrPartTypeMismatch)
	}

	// Собираем список слотов для разрешения: обязательные + непустые опциональные.
	type slot struct {
		name     string
		uuid     uuid.UUID
		partType model.PartType
	}

	toResolve := []slot{
		{name: "hull", uuid: slots.HullUUID, partType: model.PartTypeHull},
		{name: "engine", uuid: slots.EngineUUID, partType: model.PartTypeEngine},
	}
	if slots.ShieldUUID != uuid.Nil {
		toResolve = append(toResolve, slot{name: "shield", uuid: slots.ShieldUUID, partType: model.PartTypeShield})
	}
	if slots.WeaponUUID != uuid.Nil {
		toResolve = append(toResolve, slot{name: "weapon", uuid: slots.WeaponUUID, partType: model.PartTypeWeapon})
	}

	// Проверяем дублирование UUID между слотами (один UUID в разных слотах запрещён).
	seen := make(map[uuid.UUID]string, len(toResolve))
	uuids := make([]uuid.UUID, 0, len(toResolve))
	for _, sl := range toResolve {
		if prev, ok := seen[sl.uuid]; ok {
			return model.ResolvedShipSlots{}, fmt.Errorf("UUID %s используется в слотах %s и %s: %w", sl.uuid, prev, sl.name, errs.ErrPartTypeMismatch)
		}
		seen[sl.uuid] = sl.name
		uuids = append(uuids, sl.uuid)
	}

	// Загружаем детали по непустым UUID одним вызовом partRepo.List.
	parts, err := s.partRepository.List(ctx, input.PartFilter{UUIDs: uuids})
	if err != nil {
		return model.ResolvedShipSlots{}, err
	}

	// Проверяем, что каждая загруженная деталь соответствует ожидаемому типу слота,
	// и собираем проверенные детали в model.ResolvedShipSlots.
	partsByUUID := make(map[uuid.UUID]*model.Part, len(parts))
	for _, val := range parts {
		partsByUUID[val.UUID()] = new(val)
	}

	var resolved model.ResolvedShipSlots
	for _, sl := range toResolve {
		part, ok := partsByUUID[sl.uuid]
		if !ok {
			return model.ResolvedShipSlots{}, fmt.Errorf("деталь %s (uuid=%s) не найдена: %w", sl.name, sl.uuid, errs.ErrPartNotFound)
		}
		if part.PartType() != sl.partType {
			return model.ResolvedShipSlots{}, fmt.Errorf(
				"%s (uuid=%s) имеет тип %s, ожидался %s: %w",
				sl.name, sl.uuid, part.PartType(), sl.partType, errs.ErrPartTypeMismatch,
			)
		}
		switch sl.partType {
		case model.PartTypeHull:
			resolved.Hull = *part
		case model.PartTypeEngine:
			resolved.Engine = *part
		case model.PartTypeShield:
			resolved.Shield = part
		case model.PartTypeWeapon:
			resolved.Weapon = part
		}
	}

	return resolved, nil
}
