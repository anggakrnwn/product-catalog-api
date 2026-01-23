package seed

import (
	"context"
	"log"
	"time"

	"github.com/anggakrnwn/product-catalog-api/internal/domain"
	"github.com/anggakrnwn/product-catalog-api/internal/repository"
	"github.com/google/uuid"
)

func SeedCategories(ctx context.Context, repo repository.CategoryRepository) {
	categories := []domain.Category{
		{
			ID:          uuid.New().String(),
			Name:        "Elektronik",
			Description: "Perangkat elektronik dan gadget",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Fashion",
			Description: "Pakaian dan aksesori fashion",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Rumah Tangga",
			Description: "Perlengkapan dan perabotan rumah",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Olahraga",
			Description: "Peralatan dan aksesori olahraga",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Buku",
			Description: "Berbagai jenis buku dan bacaan",
		},
	}

	seeded := 0
	for _, category := range categories {
		now := time.Now()
		category.CreatedAt = now
		category.UpdatedAt = now

		err := repo.Save(ctx, category)
		if err != nil {
			log.Printf("failed to seed category %s: %v", category.Name, err)
		} else {
			seeded++
		}
	}

	log.Printf("seeded %d categories", seeded)
}
