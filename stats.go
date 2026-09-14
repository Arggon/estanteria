package main

import "time"

// Stats is the aggregated view of the shelf. Computed by a pure function:
// it never mutates the input nor touches the ledger.
type Stats struct {
	Total         int
	ByStatus      map[Status]int
	AvgRating     float64
	PagesPerMonth map[string]int // clave "2006-01", ventana de 12 meses
}

// ComputeStats aggregates books read-only. Books without Pages > 0 or
// without Finished still count in Total/ByStatus but contribute no pages.
// The pages window is the 12 months ending in now's month (inclusive).
func ComputeStats(books []Book, now time.Time) Stats {
	st := Stats{
		ByStatus:      make(map[Status]int),
		PagesPerMonth: make(map[string]int),
	}
	windowStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -11, 0)
	startMonth := windowStart.Year()*12 + int(windowStart.Month())
	endMonth := now.Year()*12 + int(now.Month())

	var ratingSum, ratingCount int
	for _, b := range books {
		st.Total++
		st.ByStatus[b.Status]++
		if b.Rating != 0 {
			ratingSum += b.Rating
			ratingCount++
		}
		if b.Pages > 0 && b.Finished != nil {
			f := *b.Finished
			fm := f.Year()*12 + int(f.Month())
			if fm >= startMonth && fm <= endMonth {
				st.PagesPerMonth[f.Format("2006-01")] += b.Pages
			}
		}
	}
	if ratingCount > 0 {
		st.AvgRating = float64(ratingSum) / float64(ratingCount)
		st.AvgRating = float64(int(st.AvgRating*100+0.5)) / 100
	}
	return st
}
