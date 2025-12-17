package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) AddPart(part *inventoryV1.Part) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.parts[part.Uuid]; exists {
		return fmt.Errorf("part %s already exists", part.Uuid)
	}
	s.parts[part.Uuid] = part
	return nil
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	parts := make([]*inventoryV1.Part, 0, len(s.parts))
	for _, p := range s.parts {
		parts = append(parts, p)
	}

	f := getFilter(req)

	if isEmptyFilter(f) {
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	pred, ok := buildPredicate(f)
	if !ok {
		return &inventoryV1.ListPartsResponse{Parts: parts}, nil
	}

	out := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		if pred(p) {
			out = append(out, p)
		}
	}

	return &inventoryV1.ListPartsResponse{Parts: out}, nil
}

type partPredicate func(*inventoryV1.Part) bool

func buildPredicate(f *inventoryV1.PartsFilter) (partPredicate, bool) {
	var preds []partPredicate

	if set := normalizeStringSet(f.GetUuids(), false); len(set) > 0 {
		preds = append(preds, func(p *inventoryV1.Part) bool {
			_, ok := set[p.GetUuid()]
			return ok
		})
	}

	if terms := normalizeStringTerms(f.GetNames()); len(terms) > 0 {
		preds = append(preds, func(p *inventoryV1.Part) bool {
			name := strings.ToLower(p.GetName())
			return anyContains(name, terms)
		})
	}

	if cats := normalizeCategorySet(f.GetCategories()); len(cats) > 0 {
		preds = append(preds, func(p *inventoryV1.Part) bool {
			_, ok := cats[p.GetCategory()]
			return ok
		})
	}

	if countries := normalizeStringSet(f.GetManufacturerCountries(), true); len(countries) > 0 {
		preds = append(preds, func(p *inventoryV1.Part) bool {
			m := p.GetManufacturer()
			if m == nil {
				return false
			}
			country := strings.ToLower(strings.TrimSpace(m.GetCountry()))
			_, ok := countries[country]
			return ok
		})
	}

	if tags := normalizeStringSet(f.GetTags(), true); len(tags) > 0 {
		preds = append(preds, func(p *inventoryV1.Part) bool {
			return hasAnyTag(p.GetTags(), tags)
		})
	}

	if len(preds) == 0 {
		return nil, false
	}

	return func(p *inventoryV1.Part) bool {
		for _, pr := range preds {
			if !pr(p) {
				return false
			}
		}
		return true
	}, true
}

func normalizeStringTerms(src []string) []string {
	out := make([]string, 0, len(src))
	for _, s := range src {
		s = strings.ToLower(strings.TrimSpace(s))
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func normalizeStringSet(src []string, lower bool) map[string]struct{} {
	out := make(map[string]struct{}, len(src))
	for _, s := range src {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if lower {
			s = strings.ToLower(s)
		}
		out[s] = struct{}{}
	}
	return out
}

func isEmptyFilter(f *inventoryV1.PartsFilter) bool {
	if f == nil {
		return true
	}
	return len(f.GetUuids()) == 0 &&
		len(f.GetTags()) == 0 &&
		len(f.GetManufacturerCountries()) == 0 &&
		len(f.GetCategories()) == 0 &&
		len(f.GetNames()) == 0
}

func getFilter(req *inventoryV1.ListPartsRequest) *inventoryV1.PartsFilter {
	if req == nil {
		return nil
	}
	return req.GetFilter()
}

func normalizeCategorySet(src []inventoryV1.Category) map[inventoryV1.Category]struct{} {
	out := make(map[inventoryV1.Category]struct{}, len(src))
	for _, c := range src {
		out[c] = struct{}{}
	}
	return out
}

func anyContains(haystack string, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(haystack, t) {
			return true
		}
	}
	return false
}

func hasAnyTag(partTags []string, filterTags map[string]struct{}) bool {
	for _, t := range partTags {
		t = strings.ToLower(strings.TrimSpace(t))
		if _, ok := filterTags[t]; ok {
			return true
		}
	}
	return false
}

func seedParts() []*inventoryV1.Part {
	now := timestamppb.New(time.Now())

	return []*inventoryV1.Part{
		{
			Uuid:          "11111111-1111-1111-1111-111111111111",
			Name:          "Main Engine X1",
			Description:   "Primary propulsion engine",
			Price:         1200000,
			StockQuantity: 5,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 4.5,
				Width:  2.1,
				Height: 2.0,
				Weight: 1800,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "Orbital Dynamics",
				Country: "Germany",
				Website: "https://orbital-dynamics.de",
			},
			Tags: []string{"main", "engine", "booster"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"fuel_type": {Kind: &inventoryV1.MetadataValue_StringValue{StringValue: "liquid"}},
				"reusable":  {Kind: &inventoryV1.MetadataValue_BoolValue{BoolValue: wrapperspb.Bool(true)}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},

		{
			Uuid:          "22222222-2222-2222-2222-222222222222",
			Name:          "Fuel Tank A9",
			Description:   "High-pressure fuel container",
			Price:         320000,
			StockQuantity: 12,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryV1.Dimensions{
				Length: 3.0,
				Width:  1.8,
				Height: 1.8,
				Weight: 900,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "CosmoFuel",
				Country: "USA",
				Website: "https://cosmofuel.com",
			},
			Tags: []string{"fuel", "tank"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"capacity_liters": {Kind: &inventoryV1.MetadataValue_Int_64Value{Int_64Value: 5000}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},

		{
			Uuid:          "33333333-3333-3333-3333-333333333333",
			Name:          "Observation Porthole",
			Description:   "Reinforced glass porthole",
			Price:         85000,
			StockQuantity: 20,
			Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryV1.Dimensions{
				Length: 1.2,
				Width:  1.2,
				Height: 0.2,
				Weight: 80,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "SpaceGlass",
				Country: "France",
				Website: "https://spaceglass.fr",
			},
			Tags: []string{"window", "glass"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"radiation_protected": {
					Kind: &inventoryV1.MetadataValue_BoolValue{
						BoolValue: wrapperspb.Bool(true),
					},
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},

		{
			Uuid:          "44444444-4444-4444-4444-444444444444",
			Name:          "Wing Module L",
			Description:   "Left aerodynamic wing",
			Price:         210000,
			StockQuantity: 7,
			Category:      inventoryV1.Category_CATEGORY_WING,
			Dimensions: &inventoryV1.Dimensions{
				Length: 6.0,
				Width:  2.5,
				Height: 0.8,
				Weight: 600,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    "AeroSpace Ltd",
				Country: "UK",
				Website: "https://aerospace.co.uk",
			},
			Tags: []string{"wing", "left"},
			Metadata: map[string]*inventoryV1.MetadataValue{
				"material": {Kind: &inventoryV1.MetadataValue_StringValue{StringValue: "carbon"}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:          "55555555-5555-5555-5555-555555555555",
			Name:          "Wing Module R",
			Description:   "Right aerodynamic wing",
			Price:         210000,
			StockQuantity: 7,
			Category:      inventoryV1.Category_CATEGORY_WING,
			Tags:          []string{"wing", "right"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},

		{
			Uuid:          "66666666-6666-6666-6666-666666666666",
			Name:          "Auxiliary Engine B2",
			Description:   "Secondary maneuvering engine",
			Price:         480000,
			StockQuantity: 4,
			Category:      inventoryV1.Category_CATEGORY_ENGINE,
			Tags:          []string{"engine", "aux"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},

		{
			Uuid:          "77777777-7777-7777-7777-777777777777",
			Name:          "Fuel Valve V1",
			Description:   "Fuel flow regulator",
			Price:         15000,
			StockQuantity: 40,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Tags:          []string{"fuel", "valve"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},

		{
			Uuid:          "88888888-8888-8888-8888-888888888888",
			Name:          "Thermal Porthole",
			Description:   "Heat resistant porthole",
			Price:         92000,
			StockQuantity: 10,
			Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
			Tags:          []string{"window", "thermal"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},

		{
			Uuid:          "99999999-9999-9999-9999-999999999999",
			Name:          "Fuel Pump P3",
			Description:   "High efficiency pump",
			Price:         60000,
			StockQuantity: 15,
			Category:      inventoryV1.Category_CATEGORY_FUEL,
			Tags:          []string{"fuel", "pump"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	for _, part := range seedParts() {
		err = service.AddPart(part)
		if err != nil {
			log.Printf("failed to add part: %v", err)
		}
	}
	inventoryV1.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on port %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down gRPC server...")
	s.GracefulStop()
	log.Printf("Server gracefully stopped")
}
