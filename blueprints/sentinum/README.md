# Sentinum

[Sentinum GmbH](https://sentinum.de/) builds MIOTY, LoRaWAN and cellular sensors in Germany. The
entries below mirror the MIOTY blueprint files Sentinum publishes on its documentation site; each
`spec` is Sentinum's own file, unchanged apart from re-indentation.

Retrieved 2026-08-27.

| Model | Type EUI | Source blueprint | Upstream SHA-256 |
|---|---|---|---|
| [Apollon-Q](https://docs.sentinum.de/en/apollon-q-t/r/tr-general-description) | `FCA84A0100000000` | [Apollon_Q_mioty_blueprint_V2.json](https://docs.sentinum.de/hubfs/blueprint%20test%202026/Apollon_Q_mioty_blueprint_V2.json) | `0d2fdd6c23b3786d0873a82d6c30e7e1da79c301dd46d3696c679121755a9f2d` |
| [Febris CO2](https://docs.sentinum.de/en/febris-co2-general-description) | `FCA84A0300000001` | [Febris_CO2_mioty_blueprint_V2.json](https://docs.sentinum.de/hubfs/blueprint%20test%202026/Febris_CO2_mioty_blueprint_V2.json) | `16572eecded94f6cacca9dddbce3cea9bee344ac826e86c090d5bc3fe580a4a1` |
| [Febris TH](https://docs.sentinum.de/en/febris-th-general-description) | `FCA84A0300000000` | [Febris_TH_mioty_blueprint_V2.json](https://docs.sentinum.de/hubfs/blueprint%20test%202026/Febris_TH_mioty_blueprint_V2.json) | `8cdff213facf3e47de3659cc18bde48128f57b8a2491a6e863a1d9c95472d3cf` |
| [Hyperion](https://docs.sentinum.de/en/would-you-like-to-find-out-more-about-the-hyperion-sensor) | `FCA84A0000000007` | [Hyperion_mioty_blueprint_V2.json](https://docs.sentinum.de/hubfs/blueprint%20test%202026/Hyperion_mioty_blueprint_V2.json) | `0532eb201359ac892eb699d90ec7588266bea6ca6a84fac65ef5caa2ecb25ea3` |
| [Juno TH TILT](https://docs.sentinum.de/en/juno-tilit-general-description) | `FCA84A0900000002` | [Juno_TH_TILT_mioty_blueprint_V2.json](https://docs.sentinum.de/hubfs/blueprint%20test%202026/Juno_TH_TILT_mioty_blueprint_V2.json) | `6cb5151baf4f50ed1aea7da0979a26333e57c928baa6618b1f82fd8b665fade7` |
| [Nyx](https://docs.sentinum.de/en/nyx-genral-description) | `FCA84A0700000000` | [Nyx_mioty_blueprint_V2.json](https://docs.sentinum.de/hubfs/blueprint%20test%202026/Nyx_mioty_blueprint_V2.json) | `422a7dfa84b317f18a96406765b8b6cff2cde2cdfadec181ec94f9fa9aea660a` |

The Hyperion blueprint declares `meta.name` as "Hyperion v3" and the Juno TH TILT blueprint declares
`meta.name` as "Juno"; the directory names follow the product names Sentinum markets.

## MIOTY products without a published blueprint

Aion IO-Link MIOTY Adapter, Apollon-ZETA, Febris SCW and Neptun are MIOTY products for which Sentinum
publishes a payload description but no blueprint file. They are absent here rather than transcribed by
hand, because a transcription cannot supply a verified Type EUI.
