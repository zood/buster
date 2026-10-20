#!/usr/bin/env python3
"""Generate the dark map style from the light one by recolouring it.

The light style, map-style.json, is edited in Maputnik. Run this afterwards from the
directory that holds it to regenerate map-style-dark.json there, rather than editing the
dark style by hand. If the light style gains a colour this script doesn't know, it lists the
colour and exits; add a dark value for it to GLOBAL (or OVERRIDES) below.

The script works on the text, so the output keeps Maputnik's formatting and a diff
against the light style shows only the changes.

usage: cd assets/static && ./make_dark_map_style.py
"""
import json
import re
import sys

# --- palette -------------------------------------------------------------------------
LAND = "#181b25"  # also map_placeholder in george's app/src/main/res/values-night/colors.xml
LAND_RGB = "24, 27, 37"
ROAD_MAJOR = "#3e465c"
ROAD_MINOR = "#2f3547"
CASING = "#2c3141"
TUNNEL = "#2a2f3f"
TUNNEL_CASING = "#3c4257"
PATH = "#6b6255"
RAIL = "#4b5164"
RAIL_RGB = "75, 81, 100"
WATER = "#14506e"
WATERWAY = "#1a5f80"
BLDG_TOP = "#272b3a"
BLDG_SIDE = "#1f2230"
ICE = "#2e3446"
AERO = "#2b3040"
AERO_CASING = "#3a4054"
BOUNDARY = "#6b7185"

# light colour -> dark colour, for every layer
GLOBAL = {
    "#f8f4f0": LAND,
    "hsla(30, 19%, 90%, 0.4)": "rgba(40, 45, 62, 0.4)",
    "hsla(30, 19%, 90%, 0.2)": "rgba(40, 45, 62, 0.2)",
    "hsla(0, 60%, 87%, 0.23)": "rgba(95, 60, 80, 0.23)",
    "hsla(49, 100%, 88%, 0.34)": "rgba(80, 72, 45, 0.34)",
    "#e0e4dd": "#232e27",
    "#fde": "#33222c",
    "#f0e8f8": "#262239",
    "#6a4": "#6a4",
    "hsla(0, 0%, 0%, 0.03)": "hsla(0, 0%, 0%, 0.03)",
    "#a9ea83": "#294530",
    "#d8e8c8": "#243a2b",
    "rgba(245, 238, 188, 1)": "rgba(48, 46, 36, 1)",
    "#a0c8f0": WATERWAY,
    "#5cd1ff": WATER,
    "hsl(210, 67%, 85%)": "hsl(200, 55%, 26%)",
    "#f2eae2": BLDG_TOP,
    "#dfdbd7": BLDG_SIDE,
    "rgba(200, 147, 102, 1)": TUNNEL_CASING,
    "#e9ac77": TUNNEL_CASING,
    "rgba(244, 209, 158, 1)": TUNNEL,
    "#fff4c6": TUNNEL,
    "#ffdaa6": TUNNEL,
    "#cba": PATH,
    "#bbb": RAIL,
    "hsla(0, 0%, 73%, 0.77)": f"rgba({RAIL_RGB}, 0.77)",
    "hsla(0, 0%, 73%, 0.68)": f"rgba({RAIL_RGB}, 0.68)",
    "#cfcdca": CASING,
    "hsla(0, 0%, 89%, 0.56)": "rgba(46, 51, 69, 0.56)",
    "#cad3de": ROAD_MINOR,
    "#fff": ROAD_MAJOR,
    "#ffffff": ROAD_MAJOR,
    "hsl(0, 0%, 70%)": RAIL,
    "hsl(248, 7%, 70%)": "#5a6074",
    "#74aee9": "#7fb9e0",
    "rgba(255,255,255,0.7)": f"rgba({LAND_RGB},0.7)",
    "rgba(255,255,255,0.8)": f"rgba({LAND_RGB},0.8)",
    "#666": "#a9aebd",
    "hsl(30, 23%, 62%)": "#9a8f80",
    "#765": "#b8ad9f",
    "#633": "#c9a3a3",
    "#333": "#d5d8e0",
    "#334": "#c3c8d6",
}

# per-layer overrides, where a light colour plays a different role
LABEL_HALO = {"#ffffff": LAND}
OVERRIDES = {
    "landcover-glacier": {"#fff": ICE},
    "landcover-ice-shelf": {"#fff": ICE},
    "aeroway-taxiway-casing": {"rgba(153, 153, 153, 1)": AERO_CASING},
    "aeroway-runway-casing": {"rgba(153, 153, 153, 1)": AERO_CASING},
    "aeroway-area": {"rgba(255, 255, 255, 1)": AERO},
    "aeroway-taxiway": {"rgba(255, 255, 255, 1)": AERO},
    "aeroway-runway": {"rgba(255, 255, 255, 1)": AERO},
    "boundary-land-level-2": {"rgba(0, 0, 0, 1)": BOUNDARY},
    # the shields are white sprites, so their text stays black
    "highway-shield-us-interstate": {"rgba(0, 0, 0, 1)": "rgba(0, 0, 0, 1)"},
    "highway-shield-us-other": {"rgba(0, 0, 0, 1)": "rgba(0, 0, 0, 1)"},
    "poi-level-3": LABEL_HALO,
    "poi-level-2": LABEL_HALO,
    "poi-level-1": LABEL_HALO,
    "poi-railway": LABEL_HALO,
    "airport-label-major": LABEL_HALO,
}

# other text edits: layer id -> [(old, new)]
EDITS = {
    # the arrow sprite is grey, which is faint on dark roads at half opacity
    "road_oneway": [('"paint": {"icon-opacity": 0.5}', '"paint": {"icon-opacity": 0.8}')],
    "road_oneway_opposite": [('"paint": {"icon-opacity": 0.5}', '"paint": {"icon-opacity": 0.8}')],
    # airport_11 is a dark grey plane with no badge, invisible on dark land
    "airport-label-major": [('"icon-image": "airport_11"', '"icon-image": "airfield_11"')],
    "highway-name-minor": [('"text-halo-blur": 0.5,', f'"text-halo-blur": 0.5,\n        "text-halo-color": "{LAND}",')],
    "highway-name-major": [('"text-halo-blur": 0.5,', f'"text-halo-blur": 0.5,\n        "text-halo-color": "{LAND}",')],
}

COLOUR = re.compile(r'"(#[0-9a-fA-F]{3,8}|rgba?\([^"]*\)|hsla?\([^"]*\))"')
LAYER_ID = re.compile(r'^      "id": "([^"]+)",$', re.M)


def main(src, dst):
    text = open(src).read()
    assert text.count('"name": "Zood",') == 1
    text = text.replace('"name": "Zood",', '"name": "Zood Dark",')

    starts = [(m.start(), m.group(1)) for m in LAYER_ID.finditer(text)]
    end = text.index('\n  ],\n  "id"')
    out = [text[:starts[0][0]]]
    unmapped = set()
    for i, (pos, layer) in enumerate(starts):
        seg = text[pos:starts[i + 1][0] if i + 1 < len(starts) else end]
        table = {**GLOBAL, **OVERRIDES.get(layer, {})}

        def sub(m):
            if m.group(1) not in table:
                unmapped.add((layer, m.group(1)))
                return m.group(0)
            return f'"{table[m.group(1)]}"'

        seg = COLOUR.sub(sub, seg)
        for old, new in EDITS.get(layer, []):
            assert seg.count(old) == 1, (layer, old)
            seg = seg.replace(old, new)
        out.append(seg)
    out.append(text[end:])
    result = "".join(out)
    if unmapped:
        sys.exit(f"unmapped colours: {sorted(unmapped)}")

    # the output must be valid and differ from the light style only in paint/icon values
    light, dark = json.loads(open(src).read()), json.loads(result)
    assert [l["id"] for l in light["layers"]] == [l["id"] for l in dark["layers"]]
    for a, b in zip(light["layers"], dark["layers"]):
        for k in set(a) | set(b):
            if k not in ("paint", "layout"):
                assert a.get(k) == b.get(k), (a["id"], k)
    open(dst, "w").write(result)


if __name__ == "__main__":
    main("map-style.json", "map-style-dark.json")
