package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/gosimple/slug"

	"github.com/photoprism/photoprism/internal/maps"
	"github.com/photoprism/photoprism/pkg/txt"
)

// Moment contains photo counts per month and year
type Moment struct {
	Label      string `json:"Label"`
	Country    string `json:"Country"`
	State      string `json:"State"`
	Year       int    `json:"Year"`
	Month      int    `json:"Month"`
	PhotoCount int    `json:"PhotoCount"`
}

var MomentLabels = map[string]string{
	"park":             "Parks & Gardens",
	"botanical-garden": "Parks & Gardens",
	"water-park":       "Parks & Gardens",
	"alpine":           "Outdoor Adventures",
	"hiking":           "Outdoor Adventures",
	"mountain":         "Outdoor Adventures",
	"mountains":        "Outdoor Adventures",
	"camping":          "Outdoor Adventures",
	"camper":           "Outdoor Adventures",
	"bench":            "Outdoor Adventures",
	"bunker":           "Outdoor Adventures",
	"castle":           "Outdoor Adventures",
	"viewpoint":        "Outdoor Adventures",
	"nature-reserve":   "Nature & Landscape",
	"landscape":        "Nature & Landscape",
	"nature":           "Nature & Landscape",
	"flower":           "Nature & Landscape",
	"field":            "Nature & Landscape",
	"forest":           "Nature & Landscape",
	"rocks":            "Nature & Landscape",
	"valley":           "Nature & Landscape",
	"bay":              "Bays, Capes & Beaches",
	"beach":            "Bays, Capes & Beaches",
	"seashore":         "Bays, Capes & Beaches",
	"cape":             "Bays, Capes & Beaches",
	"ship":             "Water, Ships & Boats",
	"water":            "Water, Ships & Boats",
	"pier":             "Water, Ships & Boats",
	"boat":             "Water, Ships & Boats",
	"boathouse":        "Water, Ships & Boats",
	"lakeside":         "Water, Ships & Boats",
	"shark":            "Water, Ships & Boats",
	"fish":             "Water, Ships & Boats",
	"jellyfish":        "Water, Ships & Boats",
	"submarine":        "Water, Ships & Boats",
	"diving":           "Water, Ships & Boats",
	"festival":         "Festivals & Entertainment",
	"nightclub":        "Festivals & Entertainment",
	"microphone":       "Festivals & Entertainment",
	"stage":            "Festivals & Entertainment",
	"theater":          "Festivals & Entertainment",
	"theme park":       "Festivals & Entertainment",
	"event":            "Festivals & Entertainment",
	"wine":             "Festivals & Entertainment",
	"cat":              "Pets",
	"dog":              "Pets",
	"rabbit":           "Pets",
	"hamster":          "Pets",
}

// MomentLabelsFilter returns the smart filter string for a moment based on a matching label.
func MomentLabelsFilter(label string) string {
	// TODO: Needs refactoring
	label = strings.SplitN(label, txt.Or, 2)[0]

	title := MomentLabels[label]

	if title == "" {
		return ""
	}

	var l []string

	for i := range MomentLabels {
		if MomentLabels[i] == title {
			l = append(l, i)
		}
	}

	return strings.Join(txt.UniqueWords(l), txt.Or)
}

// CountryName returns the country name if any.
func (m Moment) CountryName() string {
	if m.Country == "" {
		return ""
	}

	return maps.CountryName(m.Country)
}

// Slug returns an identifier string for a moment.
func (m Moment) Slug() (s string) {
	state := txt.NormalizeState(m.State, m.Country)

	if state == "" {
		return m.TitleSlug()
	}

	country := maps.CountryName(m.Country)

	if m.Year > 1900 && m.Month == 0 {
		s = fmt.Sprintf("%s-%s-%04d", country, state, m.Year)
	} else if m.Year > 1900 && m.Month > 0 && m.Month <= 12 {
		s = fmt.Sprintf("%s-%s-%04d-%02d", country, state, m.Year, m.Month)
	} else {
		s = fmt.Sprintf("%s-%s", country, state)
	}

	return slug.Make(s)
}

// TitleSlug returns an identifier string based on the title.
func (m Moment) TitleSlug() string {
	return slug.Make(m.Title())
}

// Title returns an english title for the moment.
func (m Moment) Title() string {
	state := txt.NormalizeState(m.State, m.Country)

	if m.Year == 0 && m.Month == 0 {
		if m.Label != "" {
			return MomentLabels[strings.SplitN(m.Label, txt.Or, 2)[0]]
		}
		if state != "" {
			return state
		}

		return m.CountryName()
	}

	if m.Country != "" && m.Year > 1900 && m.Month == 0 {
		if state != "" {
			return fmt.Sprintf("%s / %d", state, m.Year)
		}

		return fmt.Sprintf("%s %d", m.CountryName(), m.Year)
	}

	if m.Year > 1900 && m.Month > 0 && m.Month <= 12 {
		date := time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.UTC)

		if state != "" {
			return fmt.Sprintf("%s / %s", state, date.Format("January 2006"))
		}

		if m.Country == "" {
			return date.Format("January 2006")
		}

		return fmt.Sprintf("%s / %s", m.CountryName(), date.Format("January 2006"))
	}

	if m.Month > 0 && m.Month <= 12 {
		return time.Month(m.Month).String()
	}

	return m.CountryName()
}

// Moments represents a list of moments.
type Moments []Moment

// MomentsTime counts photos by month and year.
func MomentsTime(threshold int) (results Moments, err error) {
	db := UnscopedDb().Table("photos").
		Select("photos.photo_year AS year, photos.photo_month AS month, COUNT(*) AS photo_count").
		Where("photos.photo_quality >= 3 AND deleted_at IS NULL AND photo_private = 0 AND photos.photo_year > 0 AND photos.photo_month > 0").
		Group("photos.photo_year, photos.photo_month").
		Order("photos.photo_year DESC, photos.photo_month DESC").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

// MomentsCountries returns the most popular countries by year.
func MomentsCountries(threshold int) (results Moments, err error) {
	db := UnscopedDb().Table("photos").
		Select("photo_country AS country, photo_year AS year, COUNT(*) AS photo_count ").
		Where("photos.photo_quality >= 3 AND deleted_at IS NULL AND photo_private = 0 AND photo_country <> 'zz' AND photo_year > 0").
		Group("photo_country, photo_year").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

// MomentsStates returns the most popular states and countries by year.
func MomentsStates(threshold int) (results Moments, err error) {
	db := UnscopedDb().Table("photos").
		Select("p.place_country AS country, p.place_state AS state, COUNT(*) AS photo_count").
		Joins("JOIN places p ON p.id = photos.place_id").
		Where("photos.photo_quality >= 3 AND photos.deleted_at IS NULL AND photo_private = 0 AND p.place_state <> '' AND p.place_country <> 'zz'").
		Group("p.place_country, p.place_state").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&results).Error; err != nil {
		return results, err
	}

	return results, nil
}

// MomentsLabels returns the most popular photo labels.
func MomentsLabels(threshold int) (results Moments, err error) {
	var cats []string

	for cat, _ := range MomentLabels {
		cats = append(cats, cat)
	}

	m := Moments{}

	db := UnscopedDb().Table("photos").
		Select("l.label_slug AS label, COUNT(*) AS photo_count").
		Joins("JOIN photos_labels pl ON pl.photo_id = photos.id AND pl.uncertainty < 100").
		Joins("JOIN labels l ON l.id = pl.label_id").
		Where("photos.photo_quality >= 3 AND photos.deleted_at IS NULL AND photo_private = 0 AND l.label_slug IN (?)", cats).
		Group("l.label_slug").
		Having("photo_count >= ?", threshold)

	if err := db.Scan(&m).Error; err != nil {
		return m, err
	}

	done := make(map[string]bool)

	for i := 0; i < len(m); i++ {
		f := MomentLabelsFilter(m[i].Label)

		if _, ok := done[f]; ok {
			continue
		} else {
			done[f] = true
		}

		m[i].Label = f
		results = append(results, m[i])
	}

	return results, nil
}
