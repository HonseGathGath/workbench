(function () {
  var list = document.getElementById("post-list");
  if (!list) return;

  var cards = Array.prototype.slice.call(list.querySelectorAll(".post-summary"));
  var noResults = document.getElementById("no-results");
  var input = document.getElementById("search-input");
  var tagCloud = document.getElementById("tag-cloud");
  var activeTag = new URLSearchParams(window.location.search).get("tag") || "";
  var textIndex = null; // filled in once search-index.json loads

  if (tagCloud) {
    Array.prototype.slice.call(tagCloud.querySelectorAll(".tag-chip")).forEach(function (chip) {
      if (chip.getAttribute("data-tag") === activeTag) chip.classList.add("active");
    });
  }

  function applyFilter() {
    var q = ((input && input.value) || "").toLowerCase().trim();
    var visible = 0;
    cards.forEach(function (card) {
      var tags = (card.getAttribute("data-tags") || "").split(",");
      var tagOK = !activeTag || tags.indexOf(activeTag) !== -1;
      var textOK = true;
      if (q) {
        var slug = card.getAttribute("data-slug");
        var haystack = (textIndex && textIndex[slug]) || card.textContent.toLowerCase();
        textOK = haystack.toLowerCase().indexOf(q) !== -1;
      }
      var show = tagOK && textOK;
      card.style.display = show ? "" : "none";
      if (show) visible++;
    });
    if (noResults) noResults.hidden = visible > 0;
  }

  if (input) input.addEventListener("input", applyFilter);
  applyFilter();

  // Full post text isn't in the DOM (only the excerpt is), so pull it
  // from the generated index for free-text search once it's available.
  fetch("search-index.json")
    .then(function (r) { return r.json(); })
    .then(function (entries) {
      textIndex = {};
      entries.forEach(function (p) {
        textIndex[p.slug] = (p.title + " " + p.tags.join(" ") + " " + p.text).toLowerCase();
      });
      applyFilter();
    })
    .catch(function () {
      // Title/tag/excerpt filtering above still works without it.
    });
})();
