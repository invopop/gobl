import { keygen, build } from "./gobl.js";

// assigning these to the global namespace for cypress tests
window.gobl = {};
window.gobl.keygen = keygen;
window.gobl.build = build;

let goblData = {};

const generateAndDisplayKey = async () => {
  const key = await keygen({ indent: true });
  goblData.key = JSON.parse(key);
  document.getElementById("key").value = key;
};

const processInputFile = async () => {
  const inputFile = document.getElementById("input-file").value;

  const buildData = {
    data: btoa(inputFile),
    privatekey: goblData.key.private,
    indent: true,
  };

  try {
    const buildResult = await build(buildData);
    document.getElementById("output-file").value = buildResult;
    updateStatus("success");
  } catch (e) {
    document.getElementById("output-file").value = "";
    updateStatus("error", e);
  }
};

const displaySuccess = (el) => {
  el.classList.remove("bg-red-200");
  el.classList.add("bg-green-200");
};

const displayError = (el) => {
  el.classList.add("bg-red-200");
  el.classList.remove("bg-green-200");
};

const updateStatus = async (type, message) => {
  const statusEl = document.getElementById("status");
  if (type === "success") {
    statusEl.innerHTML = "Success!";
    displaySuccess(statusEl);
  } else {
    // error case
    statusEl.innerHTML = `Error: ${message}`;
    displayError(statusEl);
  }
};

await generateAndDisplayKey();
await processInputFile();

// process the input file on each keystroke
document.getElementById("input-file").oninput =
  function updateOnInputFileChange() {
    processInputFile();
  };
