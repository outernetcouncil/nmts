# Entity-Relationship Model

## Entities

### Entity Categorization

Entities can be categorized according to many different criteria, and no
single categorization framework may prove to be perfect nor fit for all
purposes--an exception or category error might always be found.
Nevertheless, some guidelines for categorizing Entity types are as follows.

#### Physical vs Logical

At a high level, Entities may be divided into *physical* and *logical*
types.  **Physical Entities** have attributes and relationships to other
Entities that characterize their tangible existence, including:
* a fixed position on Earth or a description of their orbital motion or
flight trajectory, or
* a description of a Field of View/Field of Regard and orienation relative
to another physical Entity to which they might be attached.

**Physical Entities** may also have attributes or relationships to Entities
that are less tangible but integral to understanding their performance in
the context of modeling a communications network, including:
* antenna gain patterns,
* signal processing pipelines (e.g., in an RF or optical chain), and
* other aspects that contribute to modeling Layer 1 of a traditional
network layer model.

The attributes ascribed to defined Physical Entities take inspiration from
a variety of existing modeling strategies tailored to the relevant physical
domain (RF systems, optical systems, orbiting platforms, ...).

**Logical Entities** are generally anything else essential to modeling the
network which has no relevant physical characteristics, or elements for
which the modeling-relevant physical characteristics have been ascribed to
a Physical Entity with which a Logical Entity has a defined Relationship.

These are often virtual elements for which there might be a software
construct or implementation, e.g. an abstract Network Node, a logical packet
link (at Layer 2, Layer 3, ...), an IPsec tunnel, or an IP/MPLS forwarding
table.

Attributes ascribed to Logical Entities take inspiration from existing model
and API work, including YANG models from [OpenConfig](https://github.com/openconfig/public)
efforts and [other published standards bodies](https://github.com/YangModels/yang).

#### Network vs Service -related

**Logical Entities** may themselves be further categorized according to a
variety of criteria, but a useful distinction can be made between
*network-related* and *service-related* Logical Entities. One example of
this is the logical entity representing a network interface as distinct
from the UNI service entity that might by assigned/configured to traverse
the network interface.

### Entity Attributes

Each entity kind has at most one attributes message type that may be
included with each instance. The attributes typically comprise:
* elements that are defined within the scope of the entity attributes
  definition itself, and
* elements influenced by or mapped from other models, e.g. IEEE, IETF,
  and MEF YANG models.

Not every aspect of models defined elsewhere may be represented in these
attributes, as they may not be directly relevant to the connectivity and
service modeling goals of NMTS.

## Relationships

As in MALT, relationships do not have attributes.  They are directed edges
that have a tail (`a:`), a head (`z:`), and a kind.
