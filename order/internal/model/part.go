package model

import (
	"time"
)

type PartsFilter struct {
	UUIDs                 []string       `json:"uuids,omitempty"`
	Names                 []string       `json:"names,omitempty"`
	Categories            []PartCategory `json:"categories,omitempty"`
	ManufacturerCountries []string       `json:"manufacturer_countries,omitempty"`
	Tags                  []string       `json:"tags,omitempty"`
}

type PartCategory int

const (
	PartCategoryUnknown PartCategory = iota
	PartCategoryEngine
	PartCategoryFuel
	PartCategoryPorthole
	PartCategoryWing
)

type PartDimensions struct {
	Length float64 `json:"length,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
	Weight float64 `json:"weight,omitempty"`
}

type PartManufacturer struct {
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
	Website string `json:"website,omitempty"`
}

type PartMetadataKind uint8

const (
	MetadataKindUnknown PartMetadataKind = iota
	MetadataKindString
	MetadataKindInt64
	MetadataKindDouble
	MetadataKindBool
)

type PartMetadataValue struct {
	Kind   PartMetadataKind `json:"kind"`
	String *string          `json:"string_value,omitempty"`
	Int64  *int64           `json:"int_64_value,omitempty"`
	Double *float64         `json:"double_value,omitempty"`
	Bool   *bool            `json:"bool_value,omitempty"`
}

func MetaString(v string) *PartMetadataValue {
	return &PartMetadataValue{Kind: MetadataKindString, String: &v}
}

func MetaInt64(v int64) *PartMetadataValue {
	return &PartMetadataValue{Kind: MetadataKindInt64, Int64: &v}
}

func MetaDouble(v float64) *PartMetadataValue {
	return &PartMetadataValue{Kind: MetadataKindDouble, Double: &v}
}

func MetaBool(v bool) *PartMetadataValue {
	return &PartMetadataValue{Kind: MetadataKindBool, Bool: &v}
}

type Part struct {
	Uuid          string                        `json:"uuid,omitempty"`
	Name          string                        `json:"name,omitempty"`
	Description   string                        `json:"description,omitempty"`
	Price         float64                       `json:"price,omitempty"`
	StockQuantity int64                         `json:"stock_quantity,omitempty"`
	Category      PartCategory                  `json:"category,omitempty"`
	Dimensions    *PartDimensions               `json:"dimensions,omitempty"`
	Manufacturer  *PartManufacturer             `json:"manufacturer,omitempty"`
	Tags          []string                      `json:"tags,omitempty"`
	Metadata      map[string]*PartMetadataValue `json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	CreatedAt     *time.Time                    `json:"created_at,omitempty"`
	UpdatedAt     *time.Time                    `json:"updated_at,omitempty"`
}
