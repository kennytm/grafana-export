import Clipboard from "clipboard";

const scriptContent = require("!!raw-loader?esModule=false!../scripts/dist/index.min.js");
document.querySelector("#scriptContentBox").value = scriptContent;

const copy = new Clipboard("#copyButton");
copy.on("success", () => {
  document.querySelector("#copiedNote").classList.remove("is-invisible");
  setTimeout(() => {
    document.querySelector("#copiedNote").classList.add("is-invisible");
  }, 1000);
});
