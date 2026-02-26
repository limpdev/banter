# Implementation Plan

**Goal**: Transform the static dashboard into a fully reactive, Go-native application using `HTMX`, `Alpine.js`, and `WebSockets`.

## Analysis Summary

The Claude analysis in the README made several assumptions based on typical Go web applications that are **inaccurate** for this specific codebase:

1. **No `bluemonday`**: The project does not use `bluemonday` for sanitization.
2. **No CSP Headers**: The backend does not set strict `Content-Security-Policy` headers blocking inline scripts.
3. **Existing Script Support**: The [widget-html.go](file:///c:/Users/drewg/bin/exe/banter/internal/glance/widget-html.go) component currently allows arbitrary HTML (including `<script>`) to be configured and passed unconditionally as `template.HTML`.

The actual constraint preventing scripts in _other_ widgets is Go's `html/template` standard context-aware escaping.

Since our goal is to fork and spin off if successful, our plan isn't just to "allow" scripts, but to **build a genuinely reactive framework** into the core.

## Proposed Changes

### Phase 1: Foundation (Global Reactivity)

Integrate reactivity libraries directly into the global document template, removing the need for users to manually supply them.

#### [MODIFY] internal/glance/templates/document.html

- Inject `htmx.org` and `alpinejs` via `<script src="...">` tags into the `<head>` block.
- This provides a globally available reactivity foundation for all widgets.

### Phase 2: Widget Partial Rendering (HTMX Support)

HTMX shines when the server can return partial HTML fragments. Currently, Glance builds the whole page at once.

#### [MODIFY] internal/glance/glance.go (or new router file)

- Add a new API endpoint, e.g., `GET /api/widget/{id}` or `GET /api/widget/{type}`.
- This endpoint will render and return ONLY the `<article>` block for a specific widget.

#### [MODIFY] internal/glance/templates/widget-base.html

- Add `hx-get="/api/widget/..."` and `hx-trigger="every XXXs"` attributes to widgets that need periodic polling (e.g., weather, monitor, server-stats).
- This replaces the current need for full-page reloads.

### Phase 3: WebSocket Live Updates

For true real-time widgets, polling isn't enough. We will introduce a WebSocket hub.

#### [NEW] internal/glance/hub.go (or similar)

- Implement a standard Go WebSocket hub (using `gorilla/websocket` or the standard library `/x/net/websocket`).
- The hub will broadcast data payload events to connected clients.

#### [MODIFY] internal/glance/main.go

- Register a `/ws` endpoint that upgrades HTTP requests to WebSocket connections.

#### [MODIFY] internal/glance/templates/document.html

- Add a small globally available JS snippet to establish the WebSocket connection on page load and route incoming messages to widgets via HTMX extension (`htmx-ws`) or Alpine.js event listeners.

### Phase 4: Refactoring Existing Widgets

With the new reactive core, we'll rewrite select widgets to demonstrate the system.

- **Server Stats / Monitor Widgets**: Update to receive data via WebSocket pushes instead of static generation.
- **Interactive Widgets**: Use Alpine.js for local state (tabs, modals, collapsible sections) directly inside the widget templates without needing roundtrips to the Go server.

## Verification Plan

### Automated Tests

- Run `go test ./...` to ensure no existing routing/parsing is broken.

### Manual Verification

1. Verify that HTMX and Alpine.js load correctly in the browser console.
2. Add an HTMX polling attribute to an existing widget (like the clock or weather) and verify network requests in DevTools.
3. Connect to the `/ws` endpoint via a browser console and assert that broadcast messages are received.
4. Render a partial widget via direct `curl` or browser request to the newly created `/api/widget/{id}` endpoint and assert ONLY that widget's HTML is returned.

---

# Walkthrough: Genuinely Reactive Dashboards

We successfully transformed the static Glance fork into a fully reactive application without introducing the massive overhead of a Single Page Application (SPA).

## 1. Global Reactivity Pipeline

Reactivity libraries are now injected globally, freeing users from having to specify them via custom widgets.

[document.html](file:///c:/Users/drewg/bin/exe/banter/internal/glance/templates/document.html) now provides:

- **[HTMX](file:///c:/Users/drewg/bin/exe/banter/internal/glance/widget.go#215-223)**: Allows HTML fragments to dynamically swap without full-page reloads.
- **`Alpine.js`**: A lightweight reactive library used directly inside templates for instantaneous, client-side interactions.

## 2. HTMX Partial Rendering

Widgets used to lock the entire page while rendering. Now, they expose a dedicated render endpoint:

In [glance.go](file:///c:/Users/drewg/bin/exe/banter/internal/glance/glance.go), a new endpoint handles partials:

```go
mux.HandleFunc("/api/widgets/{widget}/{path...}", a.handleWidgetRequest)
// Intercepts `.../render` to execute just the `<article>` widget block.
```

In [widget-base.html](file:///c:/Users/drewg/bin/exe/banter/internal/glance/templates/widget-base.html):
All widgets automatically get `id="widget-XXX"` and `hx-get="/api/widgets/XXX/render"`, along with auto-generated polling schedules based on their Go configured `cacheDuration`.

## 3. WebSockets Hub

For widgets requiring real-time, server-initiated pushes (e.g., live server stats), we implemented a native Go WebSocket hub.

[hub.go](file:///c:/Users/drewg/bin/exe/banter/internal/glance/hub.go) maintains connected clients and runs a broadcast channel. In [document.html](file:///c:/Users/drewg/bin/exe/banter/internal/glance/templates/document.html), a global WebSocket listener waits for incoming pure-HTML payloads from the Go backend, seamlessly swapping updated widget blocks into the DOM in real-time.

## 4. Demonstrating the Framework

### Backend Broadcasts (WebSockets)

We modified the **Server Stats** and **Monitor** widgets to override [setProviders](file:///c:/Users/drewg/bin/exe/banter/internal/glance/widget.go#135-136). They now launch isolated background goroutines that check for updates on their configured ticker schedule, automatically pushing newly minted HTML to all connected browsers.

[widget-server-stats.go](file:///c:/Users/drewg/bin/exe/banter/internal/glance/widget-server-stats.go#L40-L55)
[widget-monitor.go](file:///c:/Users/drewg/bin/exe/banter/internal/glance/widget-monitor.go#L41-L56)

### Local UX Interactivity (Alpine.js)

We fully rewrote the **Group** widget tabs logic. Instead of relying on spaghetti JavaScript event listeners injected remotely, the [group.html](file:///c:/Users/drewg/bin/exe/banter/internal/glance/templates/group.html) template now self-manages state cleanly using Alpine:

```html
<div x-data="{ currentTab: 0 }">
  <button
    @click="currentTab = {{ $i }}"
    :class="{ 'widget-group-title-current': currentTab === {{ $i }} }"
  ></button>
</div>
```

No Go round-trips required for tab switching.

<details><summary> EFFECTED FILES LISTING</summary>

```bash
glance.go
hub.go
widget.go
widget-monitor.go
widget-server-stats.go
group.html
widget-base.html
```

</details>

### What's Next?

The system is now fully prepped for advanced reactive experiments! You can drop an `hx-post` dynamically into a new widget template, or define a new background polling mechanism for any widget by utilizing the established WebSocket hub.
