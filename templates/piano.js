window.Piano = function (container) {
  const order = [
    ["white", "A"], // a
    ["black", "As"],
    ["white", "B"], // b
    ["white", "C"], // c
    ["black", "Cs"],
    ["white", "D"], // d
    ["black", "Ds"],
    ["white", "E"], // e
    ["white", "F"], // f
    ["black", "Fs"],
    ["white", "G"], // g
    ["black", "Gs"],
  ];

  function sleep(ms) {
    return new Promise((res) => {
      setTimeout(() => res(), ms);
    });
  }

  // build piano
  for (let i = 0; i < 88; i++) {
    const noteno = i + 21;
    const [color, notename] = order[i % order.length];

    const el = document.createElement("div");
    el.dataset.noteno = noteno;
    el.dataset.notename = notename;
    el.dataset.color = color;

    el.classList.add("key");
    el.classList.add(color);
    el.classList.add(notename);

    container.appendChild(el);
  }

  async function transition(notes, handleTouch) {
    const keep = new Set(notes.split(",").sort((a, b) => a - b));

    // de-emphasize non-button keys
    for (const key of Array.from(container.children)) {
      if (!keep.has(key.dataset.noteno)) {
        key.style.setProperty("opacity", "0");
      }
    }
    await sleep(1000);

    // add buttons
    const buttons = {};
    let i = 0;
    for (const key of Array.from(container.children)) {
      if (!keep.has(key.dataset.noteno)) {
        continue;
      }
      const el = document.createElement("div");
      const rect = key.getBoundingClientRect();
      el.classList.add("button");
      el.dataset.noteno = key.dataset.noteno;
      el.dataset.notename = key.dataset.notename;
      el.style.setProperty("left", `${rect.x}px`);
      el.style.setProperty("width", `${rect.width}px`);
      el.style.setProperty("height", `${rect.height}px`);
      el.style.setProperty("background-color", key.dataset.color);
      el.style.setProperty("transition-duration", "1s");
      container.appendChild(el);
      buttons[key.dataset.noteno] = el;
    }
    // remove keys
    for (const key of Array.from(container.children)) {
      if (key.classList.contains("key")) {
        container.removeChild(key);
      }
    }
    await sleep(0);

    // move buttons
    let j = 0;
    for (const noteno of keep) {
      const el = buttons[noteno];
      el.style.setProperty("left", `calc(${j} * 100% / ${keep.size})`);
      el.style.setProperty("height", "var(--key-height)");
      el.style.setProperty("width", `calc(100% / ${keep.size})`);
      j++;
    }

    await sleep(1000);

    // add event listeners
    function makeEventListener(el, name, noteno) {
      return function (ev) {
        ev.preventDefault();
        if (name === "on") {
          el.style.setProperty("background-color", "#dddddd");
        } else {
          el.style.setProperty("background-color", "white");
        }
        handleTouch(name, noteno);
      };
    }
    for (const el of Object.values(buttons)) {
      const noteno = Number(el.dataset.noteno);
      const on = makeEventListener(el, "on", noteno);
      const off = makeEventListener(el, "off", noteno);
      el.addEventListener("touchstart", on);
      el.addEventListener("mousedown", on);
      el.addEventListener("touchend", off);
      el.addEventListener("mouseup", off);
    }

    // de-emphasize buttons
    for (const el of Object.values(buttons)) {
      el.style.setProperty("background-color", "white");
    }

    await sleep(1000);
    for (const el of Object.values(buttons)) {
      el.style.setProperty("transition-property", "none");
    }

    handleTouch("ready", 0);
  }

  return { transition };
};
