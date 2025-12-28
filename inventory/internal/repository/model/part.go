package model

import (
	"time"
)

type PartCategory int

const (
	PartCategoryUnknown PartCategory = iota
	PartCategoryEngine
	PartCategoryFuel
	PartCategoryPorthole
	PartCategoryWing
)

type PartDimensions struct {
	Length float64 `json:"length,omitempty" bson:"length,omitempty"`
	Width  float64 `json:"width,omitempty" bson:"width,omitempty"`
	Height float64 `json:"height,omitempty" bson:"height,omitempty"`
	Weight float64 `json:"weight,omitempty" bson:"weight,omitempty"`
}

type PartManufacturer struct {
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Country string `json:"country,omitempty" bson:"country,omitempty"`
	Website string `json:"website,omitempty" bson:"website,omitempty"`
}

type PartMetadataKind uint8

const (
	MetadataKindUnknown PartMetadataKind = iota
	MetadataKindString
	MetadataKindInt64
	MetadataKindDouble
	MetadataKindBool
)

type Part struct {
	Uuid          string            `json:"uuid,omitempty" bson:"uuid"`
	Name          string            `json:"name,omitempty" bson:"name"`
	Description   string            `json:"description,omitempty" bson:"description"`
	Price         float64           `json:"price,omitempty" bson:"price"`
	StockQuantity int64             `json:"stock_quantity,omitempty" bson:"stock_quantity"`
	Category      PartCategory      `json:"category,omitempty" bson:"category"`
	Dimensions    *PartDimensions   `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Manufacturer  *PartManufacturer `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	Tags          []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	Metadata      map[string]any    `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt     *time.Time        `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     *time.Time        `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
