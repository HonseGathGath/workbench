(function () {
  var themeToggle = document.getElementById("theme-toggle");
  if (themeToggle) {
    themeToggle.addEventListener("click", function () {
      var current = document.documentElement.getAttribute("data-theme");
      var next = current === "dark" ? "light" : "dark";
      document.documentElement.setAttribute("data-theme", next);
      localStorage.setItem("theme", next);
    });
  }

  var titleEl = document.getElementById("w-title");
  var tagsEl = document.getElementById("w-tags");
  var dateEl = document.getElementById("w-date");
  var updatedEl = document.getElementById("w-updated");
  var bodyEl = document.getElementById("w-body");
  var slugEl = document.getElementById("w-slug");
  var statusEl = document.getElementById("w-status");
  var previewEl = document.getElementById("w-preview");
  var deleteEl = document.getElementById("w-delete");
  if (!bodyEl) return; // list page, nothing to wire up

  if (deleteEl) {
    deleteEl.addEventListener("click", function () {
      var title = titleEl.value.trim() || "this post";
      if (!window.confirm('Delete "' + title + '"? This removes the file from /posts — it stays recoverable from git history until you commit, but is gone from the writer and any local build immediately.')) {
        return;
      }
      deleteEl.disabled = true;
      statusEl.textContent = "deleting…";
      fetch("/api/delete", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ slug: slugEl.value })
      })
        .then(function (r) {
          if (!r.ok) throw new Error("delete failed: " + r.status);
          window.location.href = "/";
        })
        .catch(function () {
          statusEl.textContent = "couldn't delete — check the server is still running";
          deleteEl.disabled = false;
        });
    });
  }

  function renderPreview() {
    if (window.marked) {
      var html = marked.parse(bodyEl.value || "");
      previewEl.innerHTML = window.DOMPurify ? DOMPurify.sanitize(html) : html;
    } else {
      previewEl.textContent = bodyEl.value;
    }
  }

  var saveTimer = null;
  var saving = false;

  function scheduleSave() {
    statusEl.textContent = "unsaved changes";
    clearTimeout(saveTimer);
    saveTimer = setTimeout(save, 900);
  }

  function save() {
    if (!titleEl.value.trim()) {
      statusEl.textContent = "add a title to start saving";
      return;
    }
    if (saving) {
      saveTimer = setTimeout(save, 900);
      return;
    }
    saving = true;
    statusEl.textContent = "saving…";
    fetch("/api/save", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        slug: slugEl.value,
        title: titleEl.value,
        tags: tagsEl.value,
        date: dateEl.value,
        updated: updatedEl ? updatedEl.value : "",
        body: bodyEl.value
      })
    })
      .then(function (r) {
        if (!r.ok) throw new Error("save failed: " + r.status);
        return r.json();
      })
      .then(function (data) {
        if (!slugEl.value && data.slug) {
          slugEl.value = data.slug;
          history.replaceState(null, "", "/edit/" + data.slug);
          document.title = titleEl.value + " — Workbench";
        }
        statusEl.textContent = "saved";
      })
      .catch(function () {
        statusEl.textContent = "couldn't save — check the server is still running";
      })
      .finally(function () {
        saving = false;
      });
  }

  [titleEl, tagsEl, dateEl, updatedEl].forEach(function (el) {
    if (el) el.addEventListener("input", scheduleSave);
  });
  bodyEl.addEventListener("input", function () {
    renderPreview();
    scheduleSave();
  });

  renderPreview();
})();
