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
    <button @click="currentTab = {{ $i }}" :class="{ 'widget-group-title-current': currentTab === {{ $i }} }">
```
No Go round-trips required for tab switching.

---

### What's Next?
The system is now fully prepped for advanced reactive experiments! You can drop an `hx-post` dynamically into a new widget template, or define a new background polling mechanism for any widget by utilizing the established WebSocket hub.
