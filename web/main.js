window.onload = () => {
  const loadingElement = document.getElementById("loading");
  const screenElement = document.getElementById("screen");
  const finishedElement = document.getElementById("finished");

  screenElement.addEventListener("contextmenu", (event) => {
    event.preventDefault();
  });

  const hideLoading = () => {
    loadingElement.style.display = "none";
    screenElement.style.display = "block";
  };

  const showFinished = () => {
    screenElement.style.display = "none";
    finishedElement.style.display = "block";
  };

  console.log("Loading WebAssembly executable...");
  const go = new Go();
  WebAssembly.instantiateStreaming(
    fetch("web/main.wasm"),
    go.importObject
  ).then((result) => {
    console.log("Running WebAssembly executable...");
    hideLoading();
    go.run(result.instance).then(() => {
      console.log("Finished WebAssembly executable.");
      showFinished();
    });
  });
};
