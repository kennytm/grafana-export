function setLanguage(lang) {
  document.body.setAttribute("lang", lang);
}

document.querySelectorAll("[data-switch-lang]").forEach((item) => {
  item.addEventListener("click", (e) => {
    const lang = e.target.getAttribute("data-switch-lang");
    setLanguage(lang);
    localStorage.setItem("language", lang);
  });
});

function detectLanguage() {
  const availableLanguages = ["en-US", "zh-CN", "zh-TW"];

  const storedLang = localStorage.getItem("language");
  if (availableLanguages.indexOf(storedLang) > -1) {
    setLanguage(storedLang);
    return;
  }

  const fullLang = navigator.language;
  if (availableLanguages.indexOf(fullLang) > -1) {
    setLanguage(fullLang);
    return;
  }

  const lang = fullLang.split("-")[0];
  for (const l of availableLanguages) {
    if (lang === l.split("-")[0]) {
      setLanguage(l);
      return;
    }
  }

  setLanguage(availableLanguages[0]);
}

detectLanguage();

document.querySelector(".mt-nav").classList.add("visible");
