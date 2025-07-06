const input = document.getElementById("input");
const output = document.getElementById("output");

input.addEventListener("keydown", async function (e) {
  if (e.key === "Enter") {
    const cmd = input.value.trim();
    output.innerHTML += `<div>&gt; ${cmd}</div>`;
    input.value = "";

    if (cmd === "targets") {
      const res = await fetch("/challenges");
      const data = await res.json();
      data.forEach(entry => {
        output.innerHTML += `<div>  - ${entry.host}:${entry.port}</div>`;
      });
    } else if (cmd.startsWith("submit ")) {
      const flag = cmd.slice(7);
      const res = await fetch("/submit", {
        method: "POST",
        body: flag
      });
      const text = await res.text();
      output.innerHTML += `<div>${text}</div>`;
    } else {
      output.innerHTML += `<div>Unknown command</div>`;
    }

    output.scrollTop = output.scrollHeight;
  }
});
