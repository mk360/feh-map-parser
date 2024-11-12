const unitNameSelect = document.getElementById("unit-name");
const unitBanner = document.getElementById("current-unit-banner");
const unitName = document.getElementById("current-unit-name");

unitNameSelect.onchange = function(e) {
    const { target: { value } } = e;

    unitName.innerHTML = value;
    unitBanner.src = `https://feheroes.fandom.com/Special:Filepath/${UNITS[value].wikiName.replace(/ /g, "_").replace(" ENEMY", "")}_BtlFace_BU.webp`;
}

for (let unit in UNITS) {
    const option = document.createElement("option");
    option.innerHTML = unit;
    unitNameSelect.appendChild(option);
}

document.getElementById("allow-illegal").onchange = console.log;
