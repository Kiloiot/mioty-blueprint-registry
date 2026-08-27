# Swissphone

[Swissphone Wireless AG](https://www.swissphone.com/) builds MIOTY devices and evaluation hardware and
publishes their protocol definitions, blueprint included, on GitHub. Each `spec` below is Swissphone's
own file, unchanged apart from re-indentation.

Retrieved 2026-08-27.

| Model | Type EUI | Source blueprint | Upstream SHA-256 |
|---|---|---|---|
| m.BUTTON | `5C335CF96D627001` | [mbutton_payload_blueprint.json](https://github.com/swionrdhw/mbutton-payload-description/blob/main/mbutton_payload_blueprint.json) | `4b666744ebf1934ef52119ebadf8fea8812da63607ce0270e6616f3d7d24aeac` |
| m3b Demo Board | `5C335CF96D336200` | [m3b_demo_blueprint.txt](https://github.com/swionrdhw/m3b-examples/blob/master/m3b_demo_blueprint.txt) | `36b07052e9f1c883a3752b2283e26a7040069af83fa21ac24c690dc85a01ee47` |

m.BUTTON is the MIOTY button; its blueprint defines seven uplink formats and three downlink formats,
and the repository above also carries the equivalent ASN.1 UPER definition. The m3b demo blueprint
covers the demo applications shipped for the LZE magnolinq M3B maker board.
