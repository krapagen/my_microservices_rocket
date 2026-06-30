package service

import (
	"context"
	"log/slog"
	"sort"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь космического корабля
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

// Server реализует gRPC сервис
type Server struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
}

// NewServer создаёт сервер с предзагруженными seed-данными
func NewServer() *Server {
	now := timestamppb.Now()

	return &Server{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Экспериментальный корпус (нет на складе)",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

// GetPart возвращает деталь по UUID
func (s *Server) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	op := "Функция inventory/pkg/service/GetPart"
	log := slog.With("op", op)
	// 1. Проверить, что uuid не пустой → INVALID_ARGUMENT

	if req.GetUuid() == "" {
		log.ErrorContext(ctx, "uuid обязателен")
		return &inventoryv1.GetPartResponse{}, status.Error(codes.InvalidArgument, "uuid обязателен")
	}

	// 2. Валидировать формат UUID → INVALID_ARGUMENT

	partUuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		log.ErrorContext(ctx, "неверный формат uuid", "error", err)
		return &inventoryv1.GetPartResponse{}, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", req.GetUuid())
	}
	log.InfoContext(ctx, "валидный формат uuid", "uuid", partUuid.String())
	// 3. Найти деталь в map

	part, ok := s.parts[partUuid]

	// 4. Если не найдена → NOT_FOUND

	if !ok {
		log.ErrorContext(ctx, "деталь не найдена", "uuid", partUuid.String())
		return &inventoryv1.GetPartResponse{}, status.Errorf(codes.NotFound, "деталь %s не найдена", req.GetUuid())
	}
	// 5. Преобразовать в inventoryv1.Part

	retPart := &inventoryv1.Part{
		Uuid:          partUuid.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      part.PartType,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}

	// 6. Вернуть деталь
	log.InfoContext(ctx, "деталь найдена", "uuid", partUuid.String(), "name", part.Name)
	return &inventoryv1.GetPartResponse{
		Part: retPart,
	}, nil
}

// ListParts возвращает список деталей с опциональной фильтрацией по типу
func (s *Server) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	op := "Функция inventory/pkg/service/ListParts"
	log := slog.With("op", op)
	// 1. Если передан список uuids → найти детали по UUID (сохраняя порядок запроса)

	if len(req.GetUuids()) > 0 {
		var parts []*inventoryv1.Part
		for _, strUuid := range req.GetUuids() {
			//    - Проверить формат каждого UUID → INVALID_ARGUMENT
			partUuid, err := uuid.Parse(strUuid)
			if err != nil {
				log.ErrorContext(ctx, "неверный формат uuid", "error", err, "uuid", strUuid)
				return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", strUuid)
			}

			part, ok := s.parts[partUuid]
			//    - Если хоть один UUID не найден → NOT_FOUND
			if !ok {
				log.ErrorContext(ctx, "деталь не найдена", "uuid", strUuid)
				return nil, status.Errorf(codes.NotFound, "деталь %s не найдена", strUuid)
			}
			parts = append(parts, &inventoryv1.Part{
				Uuid:          partUuid.String(),
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				PartType:      part.PartType,
				StockQuantity: part.StockQuantity,
				CreatedAt:     part.CreatedAt,
			})
		}
		log.InfoContext(ctx, "детали найдены по UUID", "count", len(parts))
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}

	// 2. Иначе если part_type == UNSPECIFIED → вернуть все детали

	if req.GetPartType() == inventoryv1.PartType_PART_TYPE_UNSPECIFIED {
		var parts []*inventoryv1.Part
		for _, part := range s.parts {
			parts = append(parts, &inventoryv1.Part{
				Uuid:          part.UUID,
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				PartType:      part.PartType,
				StockQuantity: part.StockQuantity,
				CreatedAt:     part.CreatedAt,
			})
		}
		// Сортируем по имени
		sort.Slice(parts, func(i, j int) bool {
			return parts[i].Name < parts[j].Name
		})
		log.InfoContext(ctx, "возвращены все детали", "part_type", inventoryv1.PartType_PART_TYPE_UNSPECIFIED, "count", len(parts))
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}
	// 3. Иначе → фильтровать по типу

	var parts []*inventoryv1.Part
	for _, part := range s.parts {
		if part.PartType == req.GetPartType() {
			parts = append(parts, &inventoryv1.Part{
				Uuid:          part.UUID,
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				PartType:      part.PartType,
				StockQuantity: part.StockQuantity,
				CreatedAt:     part.CreatedAt,
			})
		}
	}
	log.InfoContext(ctx, "детали отфильтрованы по типу", "part_type", req.GetPartType(), "count", len(parts))
	// 4. Отсортировать по имени (для фильтрации по типу и UNSPECIFIED, не для uuids)

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Name < parts[j].Name
	})

	return &inventoryv1.ListPartsResponse{Parts: parts}, nil
}
