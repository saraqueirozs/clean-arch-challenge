package usecase

import (
	"github.com/saraqueirozs/clean-arch-challenge/internal/dto"
	"github.com/saraqueirozs/clean-arch-challenge/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrdersUseCase) Execute() ([]dto.OrderOutputDTO, error) {
	orders, err := l.OrderRepository.ListAll()
	if err != nil {
		return []dto.OrderOutputDTO{}, err
	}
	dtos := []dto.OrderOutputDTO{}
	for _, order := range orders {
		dto := dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
