#!/usr/bin/env python3
"""Fetch Open5e's SRD 5.1 data for the pin-test fixtures in this directory:
srd-2014-armor.json, srd-2014-items.json and srd-2014-weapons.json.

Usage:
    python3 fetch-srd-2014.py armor > srd-2014-armor.json
    python3 fetch-srd-2014.py items > srd-2014-items.json
    python3 fetch-srd-2014.py weapons > srd-2014-weapons.json

Open5e's v1 API 403s a request with no User-Agent header; every endpoint gets
one here. See content/items_srd_test.go for how each fixture is consumed:
TestArmorMatchesTheSrd for srd-2014-armor.json, TestGearMatchesTheSrd for
srd-2014-items.json, TestWeaponsMatchTheSrd for srd-2014-weapons.json.
"""
import json
import sys
import urllib.request

USER_AGENT = "Mozilla/5.0"


def pages(url):
    while url:
        req = urllib.request.Request(url, headers={"User-Agent": USER_AGENT})
        doc = json.load(urllib.request.urlopen(req, timeout=60))
        yield from doc["results"]
        url = doc.get("next")


def pounds(raw):
    # Render a v2 weight ("8.000", "0.250") the way the fixtures spell it,
    # keeping a fractional weight rather than truncating it to a whole pound.
    return f"{float(raw):g} lb."


def normalize_armor_name(name):
    # Open5e's v2 items API appends "Armor" to some names ("Hide Armor")
    # that v1's armor API gives bare ("Hide"); v1 is this fixture's name of
    # record, so strip the suffix before matching the two sources by name.
    name = name.strip()
    if name.lower().endswith(" armor"):
        name = name[: -len(" armor")]
    return name.lower()


def fetch_armor():
    # v1's armor endpoint is the SRD 5.1 armour table (base AC, Dex cap,
    # Strength requirement, stealth, cost) but carries no weight.
    v1 = list(pages("https://api.open5e.com/v1/armor/?document__slug=wotc-srd&limit=100"))
    wanted_categories = {"Light Armor", "Medium Armor", "Heavy Armor", "Shield"}
    v1 = [r for r in v1 if r["category"] in wanted_categories]

    # v2's items endpoint, filtered to the SRD 2014 document and the armor/
    # shield categories, supplies the weight v1 lacks.
    weight_by_name = {}
    for r in pages("https://api.open5e.com/v2/items/?limit=400"):
        if r["document"]["key"] != "srd-2014":
            continue
        cat = r.get("category")
        if not cat or cat["key"] not in ("armor", "shield"):
            continue
        weight_by_name[normalize_armor_name(r["name"])] = r["weight"]

    rows = []
    for r in sorted(v1, key=lambda r: r["slug"]):
        raw_weight = weight_by_name.get(normalize_armor_name(r["name"]))
        weight = pounds(raw_weight) if raw_weight is not None else None
        rows.append(
            {
                "slug": r["slug"],
                "name": r["name"],
                "category": r["category"],
                "base_ac": r["base_ac"],
                "plus_dex_mod": r["plus_dex_mod"],
                "plus_max": r["plus_max"],
                "plus_flat_mod": r["plus_flat_mod"],
                "strength_requirement": r["strength_requirement"],
                "stealth_disadvantage": r["stealth_disadvantage"],
                "cost": r["cost"],
                "weight": weight,
            }
        )
    return rows


def fetch_items():
    # The full SRD 2014 item table, Open5e's raw per-unit cost and weight,
    # untouched: content/items_srd_test.go scales and corrects it.
    items = [r for r in pages("https://api.open5e.com/v2/items/?limit=400") if r["document"]["key"] == "srd-2014"]
    items.sort(key=lambda r: r["key"])
    return [
        {
            "key": r["key"],
            "name": r["name"],
            "category": r["category"]["key"] if r.get("category") else None,
            "cost": r.get("cost"),
            "weight": r.get("weight"),
        }
        for r in items
    ]


def fetch_weapons():
    # v1's weapons endpoint is the SRD 5.1 weapon table, complete: v2's
    # weapon data has gaps (missing properties, one wrong damage type).
    # Only the eight fields the fixture pins are projected.
    weapons = list(pages("https://api.open5e.com/v1/weapons/?document__slug=wotc-srd&limit=100"))
    weapons.sort(key=lambda r: r["slug"])
    return [
        {
            "slug": r["slug"],
            "name": r["name"],
            "category": r["category"],
            "cost": r["cost"],
            "weight": r["weight"],
            "damage_dice": r["damage_dice"],
            "damage_type": r["damage_type"],
            # Open5e serves a weapon with no properties as null (the
            # morningstar); the fixture spells that as an empty list.
            "properties": r["properties"] or [],
        }
        for r in weapons
    ]


FETCHERS = {"armor": fetch_armor, "items": fetch_items, "weapons": fetch_weapons}


def main():
    if len(sys.argv) != 2 or sys.argv[1] not in FETCHERS:
        print("usage: fetch-srd-2014.py armor|items|weapons", file=sys.stderr)
        sys.exit(1)
    rows = FETCHERS[sys.argv[1]]()
    json.dump(rows, sys.stdout, indent=1)
    print()


if __name__ == "__main__":
    main()
