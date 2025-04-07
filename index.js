const mapEditorState = {
    mapData: [], // Will contain unit data for each cell [x][y]
    currentUnit: null, // Currently edited unit
    selectedCell: null, // Currently selected cell coordinates
};

const unitNameSelect = document.getElementById("name");
const unitBanner = document.getElementById("current-unit-banner");
const unitName = document.getElementById("current-unit-name");

const WPN_SELECT = document.getElementById("weapon");
const ASSIST_SELECT= document.getElementById("assist");
const SPECIAL_SELECT= document.getElementById("special");
const PASSIVE_A = document.getElementById("a-skill");
const PASSIVE_B = document.getElementById("b-skill");
const PASSIVE_C = document.getElementById("c-skill");
const PASSIVE_S = document.getElementById("s-skill");
const PASSIVE_X = document.getElementById("x-skill");


const SLOTS_SELECTS = {
    "weapon": WPN_SELECT,
    "assist": ASSIST_SELECT,
    "special": SPECIAL_SELECT,
    "passivea": PASSIVE_A,
    "passiveb": PASSIVE_B,
    "passivec": PASSIVE_C,
    "sacredseal": PASSIVE_S,
    "passivex": PASSIVE_X,
};

const FORM_ELEMENTS = {
    ...SLOTS_SELECTS,
    name: document.getElementById("name"),
    stats: {
        rarity: {
            slider: document.getElementById("rarity"),
            rangeText: document.getElementsByClassName("range-value")[0],
            stars: document.getElementsByClassName("star-display")[0]
        },
        displayLevel: document.getElementById("display-level"),
        trueLevel: document.getElementById("true-level"),
        hp: {
            input: document.getElementById("hp-stat"),
            slider: document.getElementById("hp-slider"),
        },
        atk: {
            input: document.getElementById("atk-stat"),
            slider: document.getElementById("atk-slider")
        },
        spd: {
            input: document.getElementById("spd-stat"),
            slider: document.getElementById("spd-slider")
        },
        def: {
            input: document.getElementById("def-stat"),
            slider: document.getElementById("def-slider")
        },
        res: {
            input: document.getElementById("res-stat"),
            slider: document.getElementById("res-slider")
        },
    },
    specialControls: {
        customCooldownRange: document.getElementById("custom-cooldown-range"),
        customCooldownInput: document.getElementById("custom-cooldown-input"),
        defaultCooldownRadio: document.getElementById("default-cooldown"),
        customCooldownRadio: document.getElementById("custom-cooldown"),
    },
    allyRadio: document.querySelector('input[value="ally"]'),
    enemyRadio: document.querySelector('input[value="enemy"]'),
}

loadSkills();

function initializeMapData(cols, rows) {
    mapEditorState.mapData = Array(cols + 1).fill().map(() => Array(rows + 1).fill(null));
    console.log(mapEditorState.mapData);
};

initializeMapData(6, 8);

function loadSampleData() {
    // Initialize empty 6x8 grid
    
    // Add some sample units for testing
    mapEditorState.mapData[2][3] = {
        name: "Sigurd: Of the Holy War",
        weapon: "Holy-War Spear",
        assist: "Rally Atk/Spd",
        special: "Override",
        a: "Distant Counter",
        b: "Null Follow-Up 3",
        c: "Atk/Spd Menace",
        s: "Deflect Missile 3",
        x: "",
        stats: {
            hp: 10,
            atk: 90,
            spd: 50,
            def: 40,
            res: 30
        },
        rarity: 5,
        displayLevel: 40,
        trueLevel: 40,
        isEnemy: true,
    };
    
    mapEditorState.mapData[4][6] = {
        name: "Marth: Of Beginnings",
        weapon: "Hero-King Sword",
        assist: "Reposition",
        special: "Aether",
        a: "Atk/Spd Push 4",
        b: "Lull Atk/Spd 3",
        c: "Joint Drive Atk",
        s: "Swift Sparrow 2",
        x: "Canto",
        rarity: 5,
        stats: {
            hp: 36,
            atk: 72,
            spd: 15,
            def: 11,
            res: 80
        },
        trueLevel: 40,
        displayLevel: 40,
        isEnemy: false,
    };
}


document.addEventListener('DOMContentLoaded', function() {
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            tabs.forEach(t => t.classList.remove('active'));
            this.classList.add('active');
            // Add actual tab switching logic here
        });
    });
    
    // File input label update
    const fileInput = document.getElementById('map-file');
    const fileLabel = document.querySelector('.file-input-label span');
    
    fileInput.addEventListener('change', function() {
        if (this.files.length > 0) {
            fileLabel.textContent = this.files[0].name;
        } else {
            fileLabel.textContent = 'Choose Map File';
        }
    });

     // Load sample data
    //  loadSampleData();
    
     // Set up the grid
     setupGrid();
     
     // Add form change listeners to update currentUnit
     const formInputs = document.querySelectorAll('.form-control');
     formInputs.forEach(input => {
         input.addEventListener('change', function() {
             if (mapEditorState.selectedCell) {
                 mapEditorState.currentUnit = getUnitFromForm();
             }
         });
     });
     
     // Add file loading functionality
     const loadButton = document.querySelector('.action-button');
     loadButton.addEventListener('click', function() {
         const fileInput = document.getElementById('map-file');
         if (fileInput.files.length > 0) {
             const file = fileInput.files[0];
             const reader = new FileReader();
             
             reader.onload = function(e) {
                 try {
                     const mapData = JSON.parse(e.target.result);
                     
                     // Validate and convert if needed
                     if (Array.isArray(mapData) && mapData.length > 0) {
                         mapEditorState.mapData = mapData;
                         setupGrid(); // Refresh grid with new data
                         alert('Map loaded successfully');
                     } else {
                         alert('Invalid map data format');
                     }
                 } catch (error) {
                     console.error('Error parsing map file:', error);
                     alert('Error loading map: ' + error.message);
                 }
             };
             
             reader.readAsText(file);
         } else {
             alert('Please select a file first');
         }
     });
});

function loadUnitToForm(unit) {
    if (!unit) {
        // Clear form if no unit
        clearUnitForm();
        return;
    }
    
    // Update each form field with unit data
    FORM_ELEMENTS.name.value = unit.name || '';
    FORM_ELEMENTS.weapon.value = unit.weapon || '';
    FORM_ELEMENTS.assist.value = unit.assist || '';
    FORM_ELEMENTS.special.value = unit.special || '';
    FORM_ELEMENTS.passivea.value = unit.a || '';
    FORM_ELEMENTS.passiveb.value = unit.b || '';
    FORM_ELEMENTS.passivec.value = unit.c || '';
    FORM_ELEMENTS.sacredseal.value = unit.s || '';
    FORM_ELEMENTS.passivex.value = unit.x || '';
    FORM_ELEMENTS.stats.trueLevel.value = +unit.trueLevel;
    FORM_ELEMENTS.stats.displayLevel.value = +unit.displayLevel;
    FORM_ELEMENTS.allyRadio.checked = !unit.isEnemy;
    FORM_ELEMENTS.enemyRadio.checked = unit.isEnemy;
    for (let stat of ["hp", "atk", "spd", "def", "res"]) {
        FORM_ELEMENTS.stats[stat].input.value = unit.stats[stat];
        FORM_ELEMENTS.stats[stat].slider.value = unit.stats[stat];
    }
    console.log(unit);
    if (unit.cooldown === -1) { // default: max cooldown of equipped special
        FORM_ELEMENTS.specialControls.defaultCooldownRadio.click();
    } else {
        FORM_ELEMENTS.specialControls.customCooldownRadio.click();

        FORM_ELEMENTS.specialControls.customCooldownRange.value = unit.cooldown;
        FORM_ELEMENTS.specialControls.customCooldownInput.value = unit.cooldown;
    }

    // document.getEl
    
    // Update character preview
    document.querySelector('.character-title').textContent = `${unit.name}`;
    // In a real app, you would also update the character image here
    unitBanner.src = `https://feheroes.fandom.com/Special:Redirect/file/${UNITS[unit.name].wikiName.replace(" ENEMY", "").replace(/ /g, "_")}_Face_FC.webp`;
}

function clearUnitForm() {
    const selectElements = document.querySelectorAll('.form-control');
    selectElements.forEach(el => {
        if (el.tagName === 'SELECT') {
            el.selectedIndex = 0;
        }
    });
    document.querySelector('.character-title').textContent = 'No Unit Selected';
}

function setupGrid() {
    const grid = document.querySelector('.grid');
    
    // Clear existing grid
    grid.innerHTML = '';
    
    // Create cells for a 6x8 grid
    for (let row = 1; row <= 8; row++) {
        for (let col = 1; col <= 6; col++) {
            const cell = document.createElement('div');
            cell.className = 'grid-cell';
            cell.setAttribute('data-coords', `${col},${row}`);
            
            // Left click - load unit
            cell.addEventListener('click', function(e) {
                e.preventDefault();
                
                // Visual selection
                document.querySelectorAll('.grid-cell').forEach(c => c.classList.remove('selected'));
                cell.classList.add('selected');
                
                // Update selected cell reference
                mapEditorState.selectedCell = [col, row];
                
                // Load unit data if exists
                const unit = mapEditorState.mapData[col][row];
                mapEditorState.currentUnit = unit ? {...unit} : null;
                loadUnitToForm(mapEditorState.currentUnit);
            });
            
            // Right click - save unit
            cell.addEventListener('contextmenu', function(e) {
                e.preventDefault();
                
                // Get current unit from form
                const unitData = getUnitFromForm();
                
                // Save to map data
                mapEditorState.mapData[col][row] = {...unitData};
                
                // Visual indicator of successful save
                const originalBg = cell.style.backgroundColor;
                cell.style.backgroundColor = 'rgba(100, 255, 100, 0.5)';
                setTimeout(() => {
                    cell.style.backgroundColor = originalBg;
                }, 300);
                
                setupGrid();
            });
            
            // Add visual indicator if cell has unit
            if (mapEditorState.mapData[col] && mapEditorState.mapData[col][row]) {
                const unitData = mapEditorState.mapData[col][row];
                console.log(UNITS[unitData.name]);
                const unitMarker = document.createElement('div');
                unitMarker.className = 'unit-indicator';
                unitMarker.style.position = 'absolute';
                unitMarker.style.top = '50%';
                unitMarker.style.left = '50%';
                unitMarker.style.transform = 'translate(-50%, -50%)';
                unitMarker.style.width = '70%';
                unitMarker.style.height = '70%';
                unitMarker.style.borderRadius = '50%';
                unitMarker.style.backgroundColor = unitData.isEnemy ? "rgb(224 97 97)" : 'rgb(97 224 201)';
                cell.appendChild(unitMarker);
                const img = document.createElement("img");
                img.loading = "lazy";
                img.style.position = "absolute";
                img.style.zIndex = 1;
                img.style.inset = "0px";
                img.style.height = "100%";
                img.style.width = "100%";
                img.src = `https://feheroes.fandom.com/Special:Filepath/${UNITS[unitData.name].wikiName.replace(" ENEMY", "").replace(/ /g, "_")}_Mini_Unit_Idle.png`;
                cell.appendChild(img);
                
                // Add unit name as tooltip
                cell.title = mapEditorState.mapData[col][row].name;
            } else {
                cell.innerHTML = "";
            }
            
            grid.appendChild(cell);
        }
    }
}

function getUnitFromForm() {
    return {
        name: FORM_ELEMENTS.name.value,
        weapon: FORM_ELEMENTS.weapon.value,
        assist: FORM_ELEMENTS.assist.value,
        special: FORM_ELEMENTS.special.value,
        a: FORM_ELEMENTS.passivea.value,
        b: FORM_ELEMENTS.passiveb.value,
        c: FORM_ELEMENTS.passivec.value,
        s: FORM_ELEMENTS.sacredseal.value,
        x: FORM_ELEMENTS.passivex.value,
        displayLevel: +FORM_ELEMENTS.stats.displayLevel.value,
        trueLevel: +FORM_ELEMENTS.stats.trueLevel.value,
        rarity: +FORM_ELEMENTS.stats.rarity.value,
        hp: +FORM_ELEMENTS.stats.hp.input.value,
        atk: +FORM_ELEMENTS.stats.atk.input.value,
        spd: +FORM_ELEMENTS.stats.spd.input.value,
        def: +FORM_ELEMENTS.stats.def.input.value,
        res: +FORM_ELEMENTS.stats.res.input.value,
        isEnemy: FORM_ELEMENTS.enemyRadio.checked
    };
}

unitNameSelect.onchange = function(e) {
    const { target: { value } } = e;
    unitName.innerHTML = value;
    unitBanner.src = `https://feheroes.fandom.com/Special:Redirect/file/${UNITS[value].wikiName.replace(" ENEMY", "").replace(/ /g, "_")}_Face_FC.webp`;
};

Array.from(document.getElementsByTagName("a")).forEach((tabAnchor) => {
    tabAnchor.onclick = function() {
        const { innerHTML } = this;
        const tabContentContainers = Array.from(document.getElementsByClassName("tab-content"));
        for (let element of tabContentContainers) {
            if (element.id === innerHTML) {
                element.classList.remove("hide");
            } else {
                element.classList.add("hide");
            }
        }
    };
});

function loadSkills() {
    WPN_SELECT.innerHTML = "<option></option>";
    ASSIST_SELECT.innerHTML = "<option></option>";
    SPECIAL_SELECT.innerHTML = "<option></option>";
    PASSIVE_A.innerHTML = "<option></option>";
    PASSIVE_B.innerHTML = "<option></option>";
    PASSIVE_C.innerHTML = "<option></option>";
    PASSIVE_S.innerHTML = "<option></option>";
    PASSIVE_X.innerHTML = "<option></option>";

    for (let skill in SKILLS) {
        const skillData = SKILLS[skill];
        if (skillData.slot.includes("passive") || skillData.slot.includes("seal")) {
            for (let slot of [PASSIVE_A, PASSIVE_B, PASSIVE_C, PASSIVE_C, PASSIVE_S, PASSIVE_X]) {
                const opt = document.createElement("option");
                opt.innerHTML = skill;
                slot.appendChild(opt);
            }
        } else {
            const targetSelect = SLOTS_SELECTS[skillData.slot];
            if (targetSelect) {
                const opt = document.createElement("option");
                opt.innerHTML = skill;
                targetSelect.appendChild(opt);
            }
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
        const img = document.createElement("img");
        img.classList.add("mini", "enemy");
        img.loading = "lazy";
        img.src = `https://feheroes.fandom.com/Special:Filepath/${UNITS[Name].wikiName.replace(" ENEMY", "").replace(/ /g, "_")}_Mini_Unit_Idle.png`;
        const groupedSkills = Object.groupBy(unit.Skills, skill => skill.slot);
        console.log(groupedSkills)
        
        mapEditorState.mapData[X][Y] = {
            name: unit.Name,
            trueLevel: unit.TrueLevel,
            displayLevel: unit.Level,
            rarity: unit.Rarity,
            weapon: groupedSkills.weapon?.[0].name,
            assist: groupedSkills.assist?.[0].name,
            special: groupedSkills.special?.[0].name,
            a: groupedSkills.passivea?.[0].name,
            b: groupedSkills.passiveb?.[0].name,
            c: groupedSkills.passivec?.[0].name,
            x: groupedSkills.passivex?.[0].name,
            isEnemy: unit.IsEnemy,
            stats: {
                hp: unit.Stats.HP,
                atk: unit.Stats.Atk,
                spd: unit.Stats.Spd,
                def: unit.Stats.Def,
                res: unit.Stats.Res
            },
            cooldown: unit.SpecialCooldown,
                // name: "Sigurd: Of the Holy War",
                // weapon: "Holy-War Spear",
                // assist: "Rally Atk/Spd",
                // special: "Override",
                // a: "Distant Counter",
                // b: "Null Follow-Up 3",
                // c: "Atk/Spd Menace",
                // s: "Deflect Missile 3",
                // x: "",
                // rarity: 5,
                // displayLevel: 40,
                // trueLevel: 40,
                // isEnemy: true,
        };
    }
    setupGrid();
    document.getElementById("grid").style.background = `url(https://feheroes.fandom.com/Special:Filepath/Map_${response.Id}.png)`;
});

const rarityInput = document.getElementById('rarity');
const rarityValue = rarityInput.nextElementSibling;
const starDisplay = document.querySelector('.star-display');

FORM_ELEMENTS.stats.rarity.slider.addEventListener('input', function() {
    const value = this.value;
    FORM_ELEMENTS.stats.rarity.rangeText.textContent = `${value} ★`;
    // Update star display
    FORM_ELEMENTS.stats.rarity.stars.textContent = '★'.repeat(parseInt(value));
});

const statSliders = document.querySelectorAll('.stat-slider');
statSliders.forEach(slider => {
    const inputField = slider.previousElementSibling;
    
    // Sync slider to input
    inputField.addEventListener('input', function() {
        slider.value = this.value;
        mapEditorState.currentUnit[this.id.replace("-stat", "")] = +slider.value;
        inputField.value = this.value;
    });
    
    // Sync input to slider
    slider.addEventListener('input', function() {
        inputField.value = this.value;
        slider.value = this.value;
        mapEditorState.currentUnit[inputField.id.replace("-stat", "")] = +slider.value;
    });
});

const aiRadioButtons = document.querySelectorAll('[name="unit-type"]');

aiRadioButtons.forEach((el) => {
    el.onchange = function(e) {
        const { value } = this;
        mapEditorState.currentUnit.isEnemy = value === "enemy";
    }
});

document.querySelectorAll('input[name="cooldown-type"]').forEach(radio => {
    radio.addEventListener('change', function() {
        const customCooldownContainer = document.getElementById('custom-cooldown-container');
        customCooldownContainer.style.display = this.value === 'custom' ? 'block' : 'none';
    });
});

// Sync the range slider and number input
const cooldownSlider = document.getElementById('custom-cooldown-slider');
const cooldownInput = document.getElementById('custom-cooldown-input');
const cooldownValue = document.querySelector('#custom-cooldown-container .range-value');

cooldownSlider.addEventListener('input', function() {
    console.log("bondour", this.value)
    cooldownInput.value = this.value;
    cooldownValue.textContent = this.value;
});

cooldownInput.addEventListener('input', function() {
    cooldownSlider.value = this.value;
    cooldownValue.textContent = this.value;
});