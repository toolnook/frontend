package generator

// CSS is the stylesheet for the generated documentation site.
const CSS = `:root {
  --bg: #1a1a2e;
  --fg: #e0e0e0;
  --accent: #7c3aed;
  --accent-light: #a78bfa;
  --card: #16213e;
  --border: #2a2a4a;
  --badge: #0f3460;
  --missing: #ef4444;
  --documented: #10b981;
}
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif; background: var(--bg); color: var(--fg); display: flex; min-height: 100vh; }
.sidebar { width: 280px; background: var(--card); padding: 24px 16px; border-right: 1px solid var(--border); position: sticky; top: 0; height: 100vh; overflow-y: auto; }
.logo { font-size: 1.4rem; font-weight: 700; color: var(--accent-light); margin-bottom: 16px; }
.logo a { color: inherit; text-decoration: none; }
.search-input { width: 100%; padding: 8px 12px; border: 1px solid var(--border); border-radius: 6px; background: var(--bg); color: var(--fg); font-size: 0.9rem; margin-bottom: 12px; }
.search-input:focus { outline: none; border-color: var(--accent); }
.module-list { list-style: none; }
.module-list li { padding: 6px 0; }
.module-list a { color: var(--fg); text-decoration: none; font-size: 0.9rem; }
.module-list a:hover { color: var(--accent-light); }
.back-link { display: block; margin-top: 12px; color: var(--accent-light); text-decoration: none; font-size: 0.85rem; }
.content { flex: 1; padding: 32px 40px; max-width: 900px; }
.content h1 { font-size: 2rem; margin-bottom: 8px; }
.content h2 { font-size: 1.4rem; margin: 24px 0 12px; }
.content h3 { font-size: 1.1rem; margin: 16px 0 8px; }
.version { color: var(--accent-light); font-size: 0.9rem; margin-bottom: 20px; }
.language { color: var(--accent-light); font-size: 0.85rem; margin-bottom: 16px; }
.badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 0.75rem; background: var(--badge); color: var(--fg); }
.badge.documented { background: var(--documented); }
.badge.missing { background: var(--missing); }
.kind { font-size: 0.75rem; padding: 1px 6px; border-radius: 3px; background: var(--badge); text-transform: uppercase; }
.doc { background: var(--card); padding: 16px 20px; border-radius: 8px; margin-bottom: 16px; border: 1px solid var(--border); }
.doc.missing-docs { border-left: 3px solid var(--missing); }
.doc h3 { display: flex; align-items: center; gap: 8px; }
.signature { background: var(--bg); padding: 8px 12px; border-radius: 4px; font-family: 'Fira Code', monospace; font-size: 0.85rem; margin: 8px 0; }
.params { list-style: none; margin: 8px 0; }
.params li { padding: 4px 0; font-size: 0.9rem; }
.type { color: var(--accent-light); font-family: monospace; font-size: 0.85rem; }
.returns { margin: 8px 0; font-size: 0.9rem; }
.example { background: var(--bg); padding: 12px; border-radius: 4px; margin: 8px 0; }
.example pre { overflow-x: auto; }
.example code { font-family: 'Fira Code', monospace; font-size: 0.85rem; }
.coverage { background: var(--card); padding: 20px; border-radius: 8px; margin: 20px 0; }
.progress-bar { background: var(--bg); border-radius: 4px; height: 8px; margin: 8px 0; overflow: hidden; }
.progress-fill { background: var(--documented); height: 100%; border-radius: 4px; transition: width 0.3s; }
.module-grid { list-style: none; display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 12px; }
.module-grid li { background: var(--card); padding: 16px; border-radius: 8px; }
.module-grid a { color: var(--accent-light); text-decoration: none; font-weight: 600; }
#search-results { margin-bottom: 12px; }
#search-results .result { padding: 4px 0; }
#search-results .result a { color: var(--fg); text-decoration: none; font-size: 0.85rem; }
#search-results .result a:hover { color: var(--accent-light); }
#graph-container { background: var(--card); border-radius: 8px; padding: 16px; min-height: 400px; }
#graph-container svg { width: 100%; height: 400px; }
`

// SearchJS is the client-side search functionality.
const SearchJS = `(function() {
  let searchIndex = [];
  const searchInput = document.getElementById('search');
  const resultsDiv = document.getElementById('search-results');

  // Load search index
  fetch('search-index.json')
    .then(r => r.json())
    .then(data => { searchIndex = data; })
    .catch(() => {});

  searchInput.addEventListener('input', function() {
    const query = this.value.toLowerCase().trim();
    if (!query) { resultsDiv.innerHTML = ''; return; }

    const results = searchIndex.filter(function(entry) {
      return entry.name.toLowerCase().includes(query) ||
             entry.description.toLowerCase().includes(query) ||
             entry.module.toLowerCase().includes(query);
    }).slice(0, 10);

    resultsDiv.innerHTML = results.map(function(r) {
      return '<div class="result"><a href="' + r.url + '">' +
        r.name + ' <span class="badge">' + r.kind + '</span> ' +
        '<small>' + r.module + '</small></a></div>';
    }).join('');
  });
})();`

// GraphJS renders the dependency graph as an interactive SVG.
const GraphJS = `(function() {
  fetch('graph-data.json')
    .then(r => r.json())
    .then(function(data) {
      const container = document.getElementById('graph-container');
      if (!container || !data.nodes || data.nodes.length === 0) {
        if (container) container.innerHTML = '<p>No dependencies found.</p>';
        return;
      }

      const width = container.clientWidth || 800;
      const height = 400;
      const cx = width / 2;
      const cy = height / 2;
      const radius = Math.min(width, height) / 2 - 60;

      // Position nodes in a circle
      const positions = {};
      data.nodes.forEach(function(node, i) {
        const angle = (i / data.nodes.length) * 2 * Math.PI;
        positions[node.id] = {
          x: cx + radius * Math.cos(angle),
          y: cy + radius * Math.sin(angle)
        };
      });

      // Build SVG
      let svg = '<svg viewBox="0 0 ' + width + ' ' + height + '">';
      // Draw edges
      data.edges.forEach(function(edge) {
        const from = positions[edge.from];
        const to = positions[edge.to];
        if (from && to) {
          svg += '<line x1="' + from.x + '" y1="' + from.y +
            '" x2="' + to.x + '" y2="' + to.y +
            '" stroke="#7c3aed" stroke-opacity="0.3" stroke-width="1"/>';
        }
      });
      // Draw nodes
      data.nodes.forEach(function(node) {
        const pos = positions[node.id];
        const r = Math.max(12, Math.min(24, 8 + node.degree * 3));
        svg += '<circle cx="' + pos.x + '" cy="' + pos.y +
          '" r="' + r + '" fill="#16213e" stroke="#7c3aed" stroke-width="2"/>';
        svg += '<text x="' + pos.x + '" y="' + (pos.y + 4) +
          '" text-anchor="middle" fill="#e0e0e0" font-size="10">' +
          node.label.substring(0, 12) + '</text>';
      });
      svg += '</svg>';
      container.innerHTML = svg;
    })
    .catch(function() {
      const container = document.getElementById('graph-container');
      if (container) container.innerHTML = '<p>Failed to load graph data.</p>';
    });
})();`
