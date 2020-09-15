import smoothscroll from "smoothscroll-polyfill";

smoothscroll.polyfill();

function scrollToHash(hash) {
  if (hash.indexOf("#") !== 0) {
    return;
  }

  const target = document.querySelector(`[data-id="` + hash.slice(1) + `"]`);

  if (hash !== "#top" && target) {
    document.querySelectorAll(".mt-content-item").forEach((el) => {
      el.classList.remove("visible");
    });
    target.classList.add("visible");
  }

  if (target) {
    target.scrollIntoView({
      behavior: "smooth",
    });
  }
}

document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
  anchor.addEventListener("click", (e) => {
    e.preventDefault();
    scrollToHash(anchor.getAttribute("href"));
  });
});
