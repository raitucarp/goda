# Changelog

## [0.1.0] — 2026-07-02

### Added

- Core Flexbox layout engine (pure Go port of Facebook Yoga)
- Full CSS Flexbox property support: flex-direction, justify-content, align-items, align-self, flex-wrap, gap, aspect-ratio, box-sizing
- CSS Grid line and track types (`GridLine`, `GridTrackSize`, `GridTrackList`)
- Dimension properties: width, height, min/max width/height with point, percent, auto, max-content, fit-content, stretch units
- Spacing properties: margin, padding, border with edge-specific and shorthand support (all, horizontal, vertical, top, right, bottom, left, start, end)
- Fluent builder pattern — all setters return `*Node` for chaining
- CSS string parsing API (`ApplyStyleString`) with semicolon/newline-separated declarations
- CSS map API (`ApplyStyle`) for programmatic property setting
- Layout output struct (`LayoutOut`) with `Rect`, `Edges`, direction, and overflow for GUI consumption
- `rem`/`em` unit support with `SetFontSizeEstimate` / `GetFontSizeEstimate`
- Configurable point scale factor for pixel grid rounding
- Errata bitmask for legacy behavior compatibility
- Experimental feature flags
- Custom measure functions, baseline functions, and clone node callbacks
- Measurement caching for optimized relayout
- Absolute positioning support
- Baseline alignment support
- Display: contents and display: none support
- Direction-aware (LTR/RTL) property resolution
- Comprehensive test suite (102 tests)
- E-commerce example using `fogleman/gg` for rendering

### Changed

- N/A (initial release)

### Removed

- C++ memory management artifacts: `Free()`, `FreeRecursive()`, `Finalize()`, `Config.Free()`, `Reset()`
- C++ event subsystem: `EventType`, `EventData`, `EventSubscriber`, allocation/deallocation events
- `YG` prefix from all public types (`YGValue` → `Value`, `YGSize` → `Size`, `YGUndefined` → `Undefined`, `YGFloatIsUndefined` → `IsUndefinedFloat`)
- `Yoga` and `YN` naming prefixes

### Fixed

- N/A (initial release)
