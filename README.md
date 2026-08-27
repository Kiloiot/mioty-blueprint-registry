# MIOTY Central Blueprint Registry

A community central registry of MIOTY blueprints — the JSON Type Format Descriptions defined by the
MIOTY Application Layer Specification that tell a service center how to turn a device's raw uplink
bytes into named, typed, unit-carrying values.

A blueprint is matched to a device by its **Type EUI**, so one blueprint serves every unit of that
device model on every network. That is exactly the kind of thing worth publishing once, in the open,
instead of re-deriving from a datasheet in every deployment.

## Using the registry

Import blueprints from this registry into your [KiloCenter](https://github.com/Kiloiot/kilo-service-center)
instance through the device catalog. The files follow the standard blueprint format, so they are equally
readable by any other MIOTY backend — application center, cloud or base station.

## Manufacturers: this is your registry too

If you build MIOTY hardware, you are invited to publish and maintain your device blueprints here.

Today most MIOTY blueprints live wherever each vendor happens to put them — a link on a product page,
an appendix in a PDF, a file attached to a support ticket, or nowhere public at all. Integrators end
up hand-transcribing payload tables, and a firmware revision that changes the payload silently breaks
decoding downstream. A shared registry fixes both ends of that: one canonical location per device, and
a version history anyone can diff.

What you get by publishing here:

- Your blueprint is importable directly into KiloCenter and reusable by any other MIOTY backend that
  reads the standard blueprint format.
- Versioned files, so a payload change ships as a new version rather than silently replacing the old one.
- Attribution: entries record the manufacturer and, for entries mirrored from a vendor site, the exact
  upstream URL they came from (see the `README.md` in each manufacturer directory).

You keep ownership of your device's blueprint. Open a pull request to add a model, publish a new
version, or correct an entry that mirrors your published file.

## Layout

Manufacturers sit at the root of the repository:

```
{manufacturer}/{device-model}/v{version}.json
```

for example `sentinum/febris-co2/v1.0.json`. Both directory names are lowercase kebab-case; the
filename carries a leading `v` and otherwise reproduces the blueprint's own `version` string. This is
the layout KiloCenter writes when a blueprint is submitted from its web interface, so hand-written
entries and submitted entries land in the same place.

Each manufacturer directory also carries a `README.md` recording, for every entry mirrored from a
vendor site, its source URL, the retrieval date and the SHA-256 of the upstream bytes.

## File format

Each file wraps the blueprint `spec` in a small envelope of catalog metadata:

| Key | Meaning |
|---|---|
| `spec` | the blueprint itself: `version`, `typeEui`, `meta`, `component`, `uplink`, `downlink` |
| `type_eui` | the spec's `typeEui` as 16 uppercase hex characters |
| `version` | the spec's `version`, repeated for indexing |
| `manufacturer_name` | manufacturer as it is written on the product |
| `device_model_name` | product name as the manufacturer markets it |
| `device_model_code` | the model directory slug |
| `id`, `device_model_id` | UUIDs identifying this blueprint and its device model |
| `is_default` | whether this is the default blueprint for the model |
| `created_at`, `updated_at` | RFC3339 timestamps |

`spec` is kept exactly as the manufacturer publishes it — key order and all — so an entry can be
compared against the vendor's own file byte for byte after re-indenting.

### UUIDs for hand-added entries

A blueprint submitted from a running KiloCenter carries UUIDs from that instance's database. An entry
added by hand has no instance behind it, so it derives both UUIDs deterministically instead, and any
tool can reproduce them:

```
namespace       = uuidv5(URL namespace, "https://github.com/Kiloiot/mioty-blueprint-registry")
device_model_id = uuidv5(namespace, "{manufacturer}/{model}")
id              = uuidv5(namespace, "{manufacturer}/{model}/{version}")
```

Deriving `device_model_id` without the version means every version of a model correctly shares one
model id.

## Contributing

Two routes, both ending in a pull request against this repository:

1. **From KiloCenter** — create the blueprint in the device catalog and use "Submit to Registry". The
   branch, file path and envelope are generated for you.
2. **By hand** — add `{manufacturer}/{model}/v{version}.json` following the layout and format above. If
   the file mirrors one published on a vendor site, add it to that manufacturer's `README.md` with its
   source URL, the retrieval date and the SHA-256 of the upstream bytes.

Please make sure a new entry: uses the device's real Type EUI (a placeholder Type EUI will never match
a real device); does not collide with a Type EUI already in the registry; and parses and validates
against the MIOTY Application Layer Specification.

## License

Blueprint specifications contributed to this repository are licensed [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
Entries that mirror a manufacturer's published file remain the work of that manufacturer and are
credited in the manufacturer's `README.md`.
