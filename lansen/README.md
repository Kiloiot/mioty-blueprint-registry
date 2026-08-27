# Lansen

[Lansen Systems AB](https://www.lansensystems.com/) builds wireless M-Bus and MIOTY sensors in Sweden
and links a MIOTY blueprint from each MIOTY product page. The entries below mirror those files; each
`spec` is Lansen's own file, unchanged apart from re-indentation. Every blueprint carries the product's
Lansen article number in `meta.article`.

Retrieved 2026-08-27.

| Model | Article | Type EUI | Source blueprint | Upstream SHA-256 |
|---|---|---|---|---|
| [LAN-MIOTY-C-TH](https://www.lansensystems.com/products/sensors/air-quality/lan-mioty-c-th/) | LAN-920-0012 | `A0412D0920001201` | [blueprint](https://www.lansensystems.com/media/1655/lan-920-0012_lan-mioty-c-th_blueprint.json) | `1ee06e4c6c0ad8e1d97064c63c47c384cebfa06e20af5454ad564f62470af5c2` |
| [LAN-MIOTY-E2-CO2](https://www.lansensystems.com/products/sensors/air-quality/lan-mioty-e2-co2/) | LAN-920-0011 | `A0412D0920001101` | [blueprint](https://www.lansensystems.com/media/1656/lan-920-0011_lan-mioty-e2-co2_blueprint.json) | `67c55f756afab1a60a3995b8c3ab792c0f25738c38f27415283cb2362e880ce8` |
| [LAN-MIOTY-E2-CO2-I](https://www.lansensystems.com/products/sensors/air-quality/lan-mioty-e2-co2-i/) | LAN-920-0003 | `A0412D0920000301` | [blueprint](https://www.lansensystems.com/media/1657/lan-920-0003_lan-mioty-e2-co2-i_blueprint.json) | `8e65a60f875926d90883f535bdfedb3d58d224f15f96e47ca5182c5d33b68766` |
| [LAN-MIOTY-E2-CO2-S](https://www.lansensystems.com/products/sensors/air-quality/lan-mioty-e2-co2-s/) | LAN-920-0072 | `A0412D0920007201` | [blueprint](https://www.lansensystems.com/media/1658/lan-920-0072_lan-mioty-e2-co2_s_blueprint.json) | `49e9efab01aeb60a5ca7866ccf46bdd8f7bdfbf3187abe4161e249ab59019eee` |
| [LAN-MIOTY-E2-CO2-S-I](https://www.lansensystems.com/products/sensors/air-quality/lan-mioty-e2-co2-s-i/) | LAN-920-0054 | `A0412D0920005401` | [blueprint](https://www.lansensystems.com/media/1659/lan-920-0054_lan-mioty-e2-co2_s_i_blueprint.json) | `5760769e3c87be9c7b4f2391e1173865f9b73b1305997cf699eae48684ba5932` |
| [LAN-MIOTY-O-TH](https://www.lansensystems.com/products/sensors/air-quality/lan-mioty-o-th/) | LAN-920-0074 | `A0412D0920007401` | [blueprint](https://www.lansensystems.com/media/1754/lan-920-0074_lan-mioty-o-th_blueprint.json) | `ddc6c325b9553c6ba4f251525e759e85b03b0c261f17ffd6e5a84d4c525aeded` |
| [LAN-MIOTY-G2-LDP](https://www.lansensystems.com/products/sensors/protection/lan-mioty-g2-ldp-kit-1/) | LAN-920-0061 | `A0412D0920006101` | [blueprint](https://www.lansensystems.com/media/1830/lan-920-0061_lan-mioty-g2-ldp_blueprint.json) | `7dd0122733ab4269a65074e397f3849c782b69b037f1e82bd7c3afa2b6939afe` |
| [LAN-MIOTY-M2](https://www.lansensystems.com/products/sensors/protection/lan-mioty-m2/) | LAN-920-0056 | `A0412D0920005601` | [blueprint](https://www.lansensystems.com/media/1662/lan-920-0056_lan-mioty-m2_blueprint.json) | `129896f865d1f7b7740ab8796816dc1fbdd0c9ca3e0b28d0d0b4388c38b3e92c` |
| [LAN-MIOTY-G2-DC-NO](https://www.lansensystems.com/products/sensors/metering/lan-mioty-g2-dc-no/) | LAN-920-0066 | `A0412D0920006601` | [blueprint](https://www.lansensystems.com/media/1661/lan-920-0066_lan-mioty-g2-dc-no_blueprint.json) | `2c284bcf78c4df2dad4760ec3abbc4a4bc7fc50c8924c7ffc843f1660705bd90` |
| [LAN-MIOTY-G2-DC-NC](https://www.lansensystems.com/products/sensors/metering/lan-mioty-g2-dc-nc/) | LAN-920-0067 | `A0412D0920006701` | [blueprint](https://www.lansensystems.com/media/1660/lan-920-0067_lan-mioty-g2-dc-nc_blueprint.json) | `ea45ee7fa7a7ba458177b87498ba10e6bec2169804af1993fb4f768224540fc5` |

The kit product LAN-MIOTY-G2-LDP-KIT-1 ships the LAN-MIOTY-G2-LDP transmitter; the directory follows
the blueprint's own `meta.name`.

## MIOTY products without a published blueprint

LAN-MIOTY-O-T and LAN-MIOTY-OD-PIR are MIOTY products whose product pages carry a datasheet and
installation guide but no blueprint file, so they are absent here.
