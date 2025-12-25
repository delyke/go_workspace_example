package part

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
	repoModel "github.com/delyke/go_workspace_example/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filters *model.PartsFilter) ([]*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	parts := make([]*repoModel.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, part)
	}

	if isEmptyFilter(filters) {
		return converter.RepoListPartsToModel(parts), nil
	}

	pred, ok := buildPredicate(filters)
	if !ok {
		return converter.RepoListPartsToModel(parts), nil
	}

	out := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		if pred(part) {
			out = append(out, lo.ToPtr(converter.RepoPartToModel(*part)))
		}
	}
	return out, nil
}

func isEmptyFilter(f *model.PartsFilter) bool {
	if f == nil {
		return true
	}
	return len(f.UUIDs) == 0 &&
		len(f.Tags) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Categories) == 0 &&
		len(f.Names) == 0
}

type partPredicate func(part *repoModel.Part) bool

func buildPredicate(filters *model.PartsFilter) (partPredicate, bool) {
	var preds []partPredicate

	if set := normalizeStringSet(filters.UUIDs, false); len(set) > 0 {
		preds = append(preds, func(part *repoModel.Part) bool {
			_, ok := set[part.Uuid]
			return ok
		})
	}

	if terms := normalizeStringTerms(filters.Names); len(terms) > 0 {
		preds = append(preds, func(part *repoModel.Part) bool {
			name := strings.ToLower(part.Name)
			return anyContains(name, terms)
		})
	}

	if cats := normalizeCategorySet(filters.Categories); len(cats) > 0 {
		preds = append(preds, func(part *repoModel.Part) bool {
			_, ok := cats[part.Category]
			return ok
		})
	}

	if countries := normalizeStringSet(filters.ManufacturerCountries, true); len(countries) > 0 {
		preds = append(preds, func(p *repoModel.Part) bool {
			m := p.Manufacturer
			if m == nil {
				return false
			}
			country := strings.ToLower(strings.TrimSpace(m.Country))
			_, ok := countries[country]
			return ok
		})
	}

	if tags := normalizeStringSet(filters.Tags, true); len(tags) > 0 {
		preds = append(preds, func(p *repoModel.Part) bool {
			return hasAnyTag(filters.Tags, tags)
		})
	}

	if len(preds) == 0 {
		return nil, false
	}
	return func(p *repoModel.Part) bool {
		for _, pr := range preds {
			if !pr(p) {
				return false
			}
		}
		return true
	}, true
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

func anyContains(haystack string, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(haystack, t) {
			return true
		}
	}
	return false
}

func normalizeCategorySet(src []model.PartCategory) map[repoModel.PartCategory]struct{} {
	out := make(map[repoModel.PartCategory]struct{}, len(src))
	for _, c := range src {
		out[converter.ModelPartCategoryToRepo(c)] = struct{}{}
	}
	return out
}
