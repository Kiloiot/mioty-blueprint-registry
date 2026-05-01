# MIOTY Central Blueprint Registry

Community-contributed MIOTY device blueprints (payload decoders) for use with [KiloCenter](https://github.com/Kiloiot/KiloServiceCenter).

## Structure

```
blueprints/
  {manufacturer}/
    {device-model}/
      {version}.json
```

Each blueprint JSON file contains the decoder specification, device metadata, and Type EUI mapping.

## Contributing

Blueprints are submitted directly from KiloCenter's web interface. When you create a blueprint and click "Submit to Registry", a pull request is automatically created in this repository.

## Usage

Import blueprints from this registry into your KiloCenter instance through the device catalog.

## License

Blueprint specifications in this repository are contributed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
