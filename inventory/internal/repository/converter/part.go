package converter

import (
	"encoding/json"
	"fmt"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
	"github.com/krapagen/my_microservices_rocket/inventory/internal/repository/record"
)

func PartRecordToModel(rec record.PartRecord) (model.Part, error) {
	var propsRec record.PartPropertiesRecord
	if err := json.Unmarshal(rec.Properties, &propsRec); err != nil {
		return model.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
	}

	props, err := partPropertiesFromRecord(propsRec)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать свойства: %w", err)
	}

	partType, err := model.NewPartType(rec.PartType)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	return model.RestorePart(
		rec.UUID,
		rec.Name,
		rec.Description,
		partType,
		rec.Price,
		rec.StockQuantity,
		rec.Reserved,
		props,
		rec.CreatedAt,
	), nil
}

func PartsRecordToModel(rec []record.PartRecord) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(rec))
	for _, r := range rec {
		part, err := PartRecordToModel(r)
		if err != nil {
			return nil, fmt.Errorf("конвертировать запись детали: %w", err)
		}
		parts = append(parts, part)
	}
	return parts, nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (model.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return model.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return model.NewEngineProperties(model.EngineClass(rec.Engine.Class), rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return model.NewShieldProperties(model.ShieldType(rec.Shield.ShieldType))
	case rec.Weapon != nil:
		return model.NewWeaponProperties(model.WeaponType(rec.Weapon.WeaponType))
	default:
		return model.PartProperties{}, nil
	}
}
