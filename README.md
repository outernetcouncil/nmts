# Network Model for Temporospatial Systems (NMTS)

<div style="text-align: right">
<br/>
"All models are wrong, but some are useful"<br/>
<a href="https://en.wikipedia.org/wiki/All_models_are_wrong">wikipedia</a><br/>
&mdash; <a href="https://en.wikipedia.org/wiki/George_E._P._Box">George E. P. Box</a><br/>
<br/>
</div>

Inspired by Google's Multi-Abstraction-Layer Topology (MALT; see
[NSDI 2020 site](https://www.usenix.org/conference/nsdi20/presentation/mogul)
for paper, slides, and recorded presentation), this collection of protobuf
definitions and supporting code aims to make feasible richer communications
network modeling of time-dynamic elements.

As in the MALT paper, this is an Entity-Relationship system where each
have "kinds" but only Entities have attributes.

The primary purpose of this data model is to support constructing a graph
of a network with its temporo-spatial dynamic aspects suitable for:
1. evaluating whether and how a requested service might be built atop the network, and
1. issuing instructions to selected elements necessary to enact the service.

In other words, this data model is in service of a Network Digital Twin
of the style depicted in [draft-irtf-nmrg-network-digital-twin-arch](https://datatracker.ietf.org/doc/html/draft-irtf-nmrg-network-digital-twin-arch-04#section-6):
```
        +---------------------------------------------------------+
        |   +-------+   +-------+          +-------+              |
        |   | App 1 |   | App 2 |   ...    | App n |   Application|
        |   +-------+   +-------+          +-------+              |
        +-------------^-------------------+-----------------------+
                      |Capability Exposure| Intent Input
                      |                   |
        +-------------+-------------------v-----------------------+
        |                        Instance of Digital Twin Network |
        |  +--------+   +------------------------+   +--------+   |
        |  |        |   | Service Mapping Models |   |        |   |
        |  |        |   |  +------------------+  |   |        |   |
        |  | Data   +--->  |Functional Models |  +---> Digital|   |
        |  | Repo-  |   |  +-----+-----^------+  |   | Twin   |   |
        |  | sitory |   |        |     |         |   | Network|   |
        |  |        |   |  +-----v-----+------+  |   |  Mgmt  |   |
        |  |        <---+  |  Basic Models    |  <---+        |   |
        |  |        |   |  +------------------+  |   |        |   |
        |  +--------+   +------------------------+   +--------+   |
        +--------^----------------------------+-------------------+
                 |                            |
                 | data collection            | control
        +--------+----------------------------v-------------------+
        |                      Real Network                       |
        |                                                         |
        +---------------------------------------------------------+

          Figure 2: Reference Architecture of Digital Twin Network
```
It aims to facilitate use by an [SDN](https://www.rfc-editor.org/rfc/rfc7426.html)
Controller.

## Documentation

* [Entity-Relationship model](docs/entity_relationship.md)
