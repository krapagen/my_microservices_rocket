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
// Класс двигателя — метка, а requiredStrength задаётся явно.
func NewEngineProperties(class EngineClass, requiredStrength int) (PartProperties, error) {
	if _, err := requiredStrengthByClass(class); err != nil {
		return PartProperties{}, err
	}

	return PartProperties{
		engine: &EngineProperties{
			class:            class,
			requiredStrength: requiredStrength,
		},
	}, nil
}
