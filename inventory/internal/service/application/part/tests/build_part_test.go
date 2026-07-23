package tests

import (
	"time"

	"github.com/google/uuid"

	"github.com/krapagen/my_microservices_rocket/inventory/internal/model"
)

func newFakePart(partType model.PartType) model.Part {
	return model.RestorePart(
		uuid.New(),
		"Fake Part",
		"description",
		partType,
		1000,
		10,
		0,
		model.PartProperties{},
		time.Now().UTC(),
	)
}

func newFakeHull(strength int) model.Part {
	props, _ := model.NewHullProperties(strength)
	return model.RestorePart(
		uuid.New(),
		"Hull",
		"hull description",
		model.PartTypeHull,
		5000,
		10,
		0,
		props,
		time.Now().UTC(),
	)
}

func newFakeEngine(class model.EngineClass, requiredStrength int) model.Part {
	props, _ := model.NewEngineProperties(class, requiredStrength)
	return model.RestorePart(
		uuid.New(),
		"Engine",
		"engine description",
		model.PartTypeEngine,
		3000,
		10,
		0,
		props,
		time.Now().UTC(),
	)
}

func newFakeShield(shieldType model.ShieldType) model.Part {
	props, _ := model.NewShieldProperties(shieldType)
	return model.RestorePart(
		uuid.New(),
		"Shield",
		"shield description",
		model.PartTypeShield,
		2000,
		10,
		0,
		props,
		time.Now().UTC(),
	)
}

func newFakeWeapon(weaponType model.WeaponType) model.Part {
	props, _ := model.NewWeaponProperties(weaponType)
	return model.RestorePart(
		uuid.New(),
		"Weapon",
		"weapon description",
		model.PartTypeWeapon,
		2000,
		10,
		0,
		props,
		time.Now().UTC(),
	)
}
