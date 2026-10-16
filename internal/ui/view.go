package ui

import (
	"fmt"
	"strings"

	"wax/internal/client"
)

// View renders the full TUI.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("  " + m.spinner.View() + " loading…\n")
	} else {
		switch m.mode {
		case viewCollection:
			b.WriteString(m.renderCollection())
		case viewSearch:
			b.WriteString(m.renderSearchResults())
		case viewMarketplace:
			b.WriteString(m.renderMarketplace())
		case viewDetail:
			b.WriteString(m.renderDetail())
		}
	}

	b.WriteString("\n")
	b.WriteString(m.renderFooter())
	return b.String()
}

func (m Model) renderHeader() string {
	title := styleTitle.Render("wax")
	var sub string
	switch m.mode {
	case viewCollection:
		sub = styleStatus.Render(fmt.Sprintf("collection  page %d/%d  (%d releases)", m.collectionPage, m.collectionPages, m.collectionTotal))
	case viewSearch:
		sub = styleStatus.Render(fmt.Sprintf("search: %q  page %d/%d  (%d results)", m.lastQuery, m.searchPage, m.searchPages, m.searchTotal))
	case viewDetail:
		if m.detail != nil {
			sub = styleStatus.Render(fmt.Sprintf("%s — %s", artistNames(m.detail.Artists), m.detail.Title))
		}
	}
	return title + "  " + sub
}

func (m Model) renderFooter() string {
	if m.searching {
		return styleFilter.Render("/") + " " + m.searchInput.View()
	}
	if m.status != "" {
		s := m.status
		if m.statusErr {
			return styleError.Render(s)
		}
		return styleStatus.Render(s)
	}
	help := "↑/↓ k/j move  n/p page  / search  enter detail  esc back  q quit"
	if m.mode == viewDetail {
		help = "↑/↓ k/j scroll  esc back  q quit"
	}
	return styleHelp.Render(help)
}

// ── collection view ──────────────────────────────────────────────────────────

const (
	colArtist = 32
	colTitle  = 56
	colYear   = 6
	colFormat = 14
	colLabel  = 38
)

func (m Model) renderCollection() string {
	var b strings.Builder

	// header row
	b.WriteString(renderHeaderRow())
	b.WriteString("\n")

	visible := visibleRows(m.height - 6)
	start, end := scrollWindow(m.collectionCursor, len(m.collectionReleases), visible)

	for i := start; i < end; i++ {
		r := m.collectionReleases[i]
		row := collectionRow(r.BasicInfo)
		b.WriteString(renderRow(row, i, m.collectionCursor))
		b.WriteString("\n")
	}
	return b.String()
}

func renderHeaderRow() string {
	cols := []struct{ title string; width int }{
		{"Artist", colArtist},
		{"Title", colTitle},
		{"Year", colYear},
		{"Format", colFormat},
		{"Label", colLabel},
	}
	var cells []string
	for _, c := range cols {
		cells = append(cells, styleHeader.Width(c.width).Render(truncate(c.title, c.width-2)))
	}
	return strings.Join(cells, "")
}

func renderRow(cells []string, idx, cursor int) string {
	widths := []int{colArtist, colTitle, colYear, colFormat, colLabel}
	style := styleRowEven
	if idx%2 == 1 {
		style = styleRowOdd
	}
	if idx == cursor {
		style = styleSelected
	}

	var parts []string
	for i, w := range widths {
		val := ""
		if i < len(cells) {
			val = cells[i]
		}
		parts = append(parts, style.Width(w).Render(truncate(val, w-2)))
	}
	return strings.Join(parts, "")
}

func collectionRow(r client.ReleaseBasicInfo) []string {
	year := fmt.Sprintf("%d", r.Year)
	if r.Year == 0 {
		year = "-"
	}
	return []string{
		artistNames(r.Artists),
		r.Title,
		year,
		formatName(r.Formats),
		labelName(r.Labels),
	}
}

// ── search view ──────────────────────────────────────────────────────────────

func (m Model) renderSearchResults() string {
	var b strings.Builder
	b.WriteString(renderHeaderRow())
	b.WriteString("\n")

	visible := visibleRows(m.height - 6)
	start, end := scrollWindow(m.searchCursor, len(m.searchResults), visible)

	for i := start; i < end; i++ {
		r := m.searchResults[i]
		row := []string{
			fmt.Sprintf("%d", r.ID),
			r.Title,
			r.Year,
			strings.Join(r.Format, ", "),
			strings.Join(r.Label, ", "),
		}
		b.WriteString(renderRow(row, i, m.searchCursor))
		b.WriteString("\n")
	}
	return b.String()
}

// ── marketplace view ─────────────────────────────────────────────────────────

func (m Model) renderMarketplace() string {
	var b strings.Builder

	// marketplace has different columns
	header := strings.Join([]string{
		styleHeader.Width(10).Render("Condition"),
		styleHeader.Width(20).Render("Sleeve"),
		styleHeader.Width(10).Render("Price"),
		styleHeader.Width(10).Render("Currency"),
		styleHeader.Width(24).Render("Seller"),
		styleHeader.Width(72).Render("Notes"),
	}, "")
	b.WriteString(header)
	b.WriteString("\n")

	visible := visibleRows(m.height - 6)
	start, end := scrollWindow(m.marketplaceCursor, len(m.marketplaceListings), visible)

	for i := start; i < end; i++ {
		l := m.marketplaceListings[i]
		style := styleRowEven
		if i%2 == 1 {
			style = styleRowOdd
		}
		if i == m.marketplaceCursor {
			style = styleSelected
		}
		row := strings.Join([]string{
			style.Width(10).Render(truncate(l.Condition, 8)),
			style.Width(20).Render(truncate(l.SleeveCondition, 18)),
			style.Width(10).Render(fmt.Sprintf("%.2f", l.Price.Value)),
			style.Width(10).Render(l.Price.Currency),
			style.Width(24).Render(truncate(l.Seller.Username, 22)),
			style.Width(72).Render(truncate(l.Comments, 70)),
		}, "")
		b.WriteString(row)
		b.WriteString("\n")
	}
	return b.String()
}

// ── detail view ──────────────────────────────────────────────────────────────

func (m Model) renderDetail() string {
	if m.detail == nil {
		return ""
	}
	d := m.detail
	var lines []string

	// metadata block
	lines = append(lines, styleTitle.Render(artistNames(d.Artists)+" — "+d.Title))
	lines = append(lines, "")

	year := fmt.Sprintf("%d", d.Year)
	if d.Year == 0 {
		year = "—"
	}
	lines = append(lines, styleStatus.Render(fmt.Sprintf("Year: %s   Country: %s", year, d.Country)))
	lines = append(lines, styleStatus.Render(fmt.Sprintf("Label: %s", labelName(d.Labels))))
	lines = append(lines, styleStatus.Render(fmt.Sprintf("Format: %s", formatName(d.Formats))))
	if len(d.Genres) > 0 {
		lines = append(lines, styleStatus.Render(fmt.Sprintf("Genre: %s", strings.Join(d.Genres, ", "))))
	}
	if len(d.Styles) > 0 {
		lines = append(lines, styleStatus.Render(fmt.Sprintf("Style: %s", strings.Join(d.Styles, ", "))))
	}

	// community stats
	lines = append(lines, "")
	lowestPrice := "—"
	if d.Community.LowestPrice != nil {
		lowestPrice = fmt.Sprintf("%.2f %s", d.Community.LowestPrice.Value, d.Community.LowestPrice.Currency)
	}
	lines = append(lines, styleDim.Render(fmt.Sprintf(
		"Have: %d   Want: %d   Lowest: %s   Rating: %.2f/5 (%d votes)",
		d.Community.Have, d.Community.Want, lowestPrice,
		d.Community.Rating.Average, d.Community.Rating.Count,
	)))

	// tracklist
	if len(d.Tracklist) > 0 {
		lines = append(lines, "")
		lines = append(lines, styleTitle.Render("Tracklist"))
		lines = append(lines, "")
		for _, t := range d.Tracklist {
			dur := ""
			if t.Duration != "" {
				dur = "  " + styleDim.Render(t.Duration)
			}
			pos := styleFilter.Render(fmt.Sprintf("%-4s", t.Position))
			lines = append(lines, pos+" "+t.Title+dur)
		}
	}

	// notes
	if d.Notes != "" {
		lines = append(lines, "")
		lines = append(lines, styleTitle.Render("Notes"))
		lines = append(lines, "")
		// word-wrap notes to terminal width
		for _, l := range wrapText(d.Notes, m.width-4) {
			lines = append(lines, styleDim.Render(l))
		}
	}

	// apply scroll
	visible := m.height - 6
	start := m.detailScroll
	if start >= len(lines) {
		start = len(lines) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + visible
	if end > len(lines) {
		end = len(lines)
	}

	return strings.Join(lines[start:end], "\n")
}

// wrapText wraps s to max width at word boundaries.
func wrapText(s string, width int) []string {
	if width <= 0 {
		return []string{s}
	}
	words := strings.Fields(s)
	var lines []string
	var current strings.Builder
	for _, w := range words {
		if current.Len()+len(w)+1 > width && current.Len() > 0 {
			lines = append(lines, current.String())
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteByte(' ')
		}
		current.WriteString(w)
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	return lines
}

// ── helpers ──────────────────────────────────────────────────────────────────

func visibleRows(h int) int {
	if h < 5 {
		return 5
	}
	return h
}

func scrollWindow(cursor, total, visible int) (start, end int) {
	start = cursor - visible/2
	if start < 0 {
		start = 0
	}
	end = start + visible
	if end > total {
		end = total
		start = end - visible
		if start < 0 {
			start = 0
		}
	}
	return
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 3 {
		return string(r[:max])
	}
	return string(r[:max-1]) + "…"
}

func artistNames(artists []client.Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

func labelName(labels []client.Label) string {
	seen := map[string]bool{}
	var out []string
	for _, l := range labels {
		if !seen[l.Name] {
			seen[l.Name] = true
			out = append(out, l.Name)
		}
	}
	return strings.Join(out, ", ")
}

func formatName(formats []client.Format) string {
	names := make([]string, len(formats))
	for i, f := range formats {
		names[i] = f.Name
	}
	return strings.Join(names, ", ")
}
