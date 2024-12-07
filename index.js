const unitNameSelect = document.getElementById("unit-name");
const unitBanner = document.getElementById("current-unit-banner");
const unitName = document.getElementById("current-unit-name");

const WPN_SELECT = document.getElementById("weapon");
const ASSIST_SELECT= document.getElementById("assist");
const SPECIAL_SELECT= document.getElementById("special");
const PASSIVE_A = document.getElementById("A-passive");
const PASSIVE_B = document.getElementById("B-passive");
const PASSIVE_C = document.getElementById("C-passive");
const PASSIVE_S = document.getElementById("S-passive");
const PASSIVE_X = document.getElementById("X-passive");

const TILE_TO_UNIT_MAP = {};

unitNameSelect.onchange = function(e) {
    const { target: { value } } = e;
    unitName.innerHTML = value;
    unitBanner.src = `https://feheroes.fandom.com/Special:Filepath/${UNITS[value].wikiName.replace(" ENEMY", "").replace(/ /g, "_")}_BtlFace_BU.webp`;

    loadLearnset(value);
};

Array.from(document.getElementsByTagName("a")).forEach((tabAnchor) => {
    tabAnchor.onclick = function() {
        const { id } = this;
        const tabContentContainers = Array.from(document.getElementsByClassName("tab-content"));
        for (let element of tabContentContainers) {
            if (element.id === `tab-${id}`) {
                element.classList.remove("hide");
            } else {
                element.classList.add("hide");
            }
        }
        const tabs = Array.from(document.querySelectorAll("li.tab"));
        for (let tab of tabs) {
            if (tab === this.parentNode) {
                tab.classList.add("active");
            } else {
                tab.classList.remove("active");
            }
        }
    };
});

const SLOTS_SELECTS = {
    "weapon": WPN_SELECT,
    "assist": ASSIST_SELECT,
    "special": SPECIAL_SELECT,
    "passivea": PASSIVE_A,
    "passiveb": PASSIVE_B,
    "passivec": PASSIVE_C,
    "passivex": PASSIVE_X,
};

function loadLearnset(unit) {
    const learnset = LEARNSETS[unit];
    WPN_SELECT.innerHTML = "<option></option>";
    ASSIST_SELECT.innerHTML = "<option></option>";
    SPECIAL_SELECT.innerHTML = "<option></option>";
    PASSIVE_A.innerHTML = "<option></option>";
    PASSIVE_B.innerHTML = "<option></option>";
    PASSIVE_C.innerHTML = "<option></option>";
    PASSIVE_S.innerHTML = "<option></option>";
    PASSIVE_X.innerHTML = "<option></option>";

    for (let skill of learnset) {
        const skillData = SKILLS[skill];
        const targetSelect = SLOTS_SELECTS[skillData.slot];
        if (targetSelect) {
            const opt = document.createElement("option");
            opt.innerHTML = skill;
            targetSelect.appendChild(opt);
        }
    }
};

for (let unit in UNITS) {
    const option = document.createElement("option");
    option.innerHTML = unit;
    unitNameSelect.appendChild(option);
}

fetch("http://localhost:3535/map?filename=S8084C.bin").then((r) => r.json()).then((response) => {
    for (let unit of response.Units) {
        const { X, Y, Name } = unit;
        console.log({ X, Y })
        const img = document.createElement("img");
        img.classList.add("mini", "enemy");
        img.loading = "lazy";
        img.src = `https://feheroes.fandom.com/Special:Filepath/${UNITS[Name].wikiName.replace(" ENEMY", "").replace(/ /g, "_")}_Mini_Unit_Idle.png`;
        TILE_TO_UNIT_MAP[`${Y}-${X}`] = unit;
        document.getElementById(`${Y}-${X}`).appendChild(img);
    }
    document.getElementById("bg").style.background = `url(https://feheroes.fandom.com/Special:Filepath/Map_${response.Id}.png)`;
});

function loadCompleteLearnsets() {
    
}

// var i = document.createElement("img");
// i.src = `https://feheroes.fandom.com/Special:Filepath/Micaiah_Dawning_Maiden_Mini_Unit_Idle.png`;
// i.classList.add("mini");
// document.body.appendChild(i);
