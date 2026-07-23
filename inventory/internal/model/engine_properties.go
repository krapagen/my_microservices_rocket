package model

import (
	"fmt"

	errs "github.com/krapagen/my_microservices_rocket/inventory/internal/errors"
)

// EngineProperties — свойства двигателя (Value Object).
type EngineProperties struct {
	requiredStrength int
	class            EngineClass
}

type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

func (e *EngineProperties) RequiredStrength() int {
	return e.requiredStrength
}

func (e *EngineProperties) Class() EngineClass {
	return e.class
}

// requiredStrengthByClass возвращает требуемую прочность корпуса для класса двигателя.
func requiredStrengthByClass(class EngineClass) (int, error) {
	switch class {
	case EngineClassC:
		return 30, nil
	case EngineClassB:
		return 70, nil
	case EngineClassA:
		return 100, nil
	default:
		return 0, fmt.Errorf("неизвестный класс двигателя %q: %w", class, errs.ErrInvalidProperties)
	}
}

// NewEngineProperties создаёт свойства двигателя.
// Класс двигателя определяет требуемую прочность корпуса:
//
//	C -> 30, B -> 70, A -> 100.
func NewEngineProperties(class EngineClass, requiredStrength int) (PartProperties, error) {
	expected, err := requiredStrengthByClass(class)
	if err != nil {
		return PartProperties{}, err
	}

	if requiredStrength != expected {
		return PartProperties{}, fmt.Errorf(
			"для класса %q required_strength должен быть %d, получено %d: %w",
			class, expected, requiredStrength, errs.ErrInvalidProperties,
		)
	}

	return PartProperties{
		engine: &EngineProperties{
			class:            class,
			requiredStrength: requiredStrength,
		},
	}, nil
}
