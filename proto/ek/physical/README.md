# Physical Layer Models

This directory contains NMTS' models for the physical layer (OSI Model Layer 1):
* [antenna.proto](./antenna.proto): Models of RF or optical antennas, including their physical constraints, field of regard limitations, and beamforming and beamhopping capabilities. 
* [common.proto](./common.proto): Messages that are shared throughout the physical layer models. 
* [platform.proto](./platform.proto): Models of physical platforms, such as a ground station, satellite, fixed or mobile user terminal, aircraft, ship, terrestrial vehicle, etc.
* [port.proto](./port.proto): Models of physical ports, such as an Ethernet jack or a physical port on a processor or sensor. NMTS considers these `Port`s to originate and terminate signals in the physical layer.
* [modulation.proto](./modulation.proto): Models of modulators, demodulators, and adpative coding and modulation configurations. The AdaptiveCodingAndModulation message maps thresholds of various measurements of received signal quality (e.g. carrier to noise-plus-interference) to the effective data rate that the link could sustain. Based on wireless propagation analysis, an interested observer can predict which MODCOD the Adaptive Coding and Modulation would select, and use this predicted MODCOD to estimate the capacity of wireless links. 
* [signal_processing_chain.proto](./signal_processing_chain.proto): Models of amplifiers, filters and mixers.
* [transceiver.proto](./transceiver.proto): Models of transmitters and receivers.

A basic example of how one antenna in a user terminal or regenerative payload  
could be modeled is:

```
      --------------------- Port -----------------------
     |                                                  |
     | (RK_ORIGINATES)                                  | (RK_TERMINATES)             
     v                                                  V
 Modulator *                                       Demodulator *
     |                                                  ^
     | (RK_SIGNAL_TRANSITS)                             | (RK_SIGNAL_TRANSITS)             
     v                                                  |
SignalProcessingChain                           SignalProcessingChain
     |                                                  ^
     | (RK_SIGNAL_TRANSITS)                             | (RK_SIGNAL_TRANSITS)             
     v                                                  |
 Transmitter                                         Receiver
     |                                                  ^
     | (RK_SIGNAL_TRANSITS)                             | (RK_SIGNAL_TRANSITS)             
     |                                                  |
      ---------------------> Antenna ------------------- 
                              ^
                              |
                (RK_CONTAINS) |
                              |
                          Platform
```
\* The `Modulator` and `Demodulator` would both have relationships to an `AdaptiveCodingAndModulation`, such as:

```
           ---------------(RK_CHARACTERIZES)--> Modulator
          |
 AdaptiveCodingAndModulation
          |
           ---------------(RK_CHARACTERIZES)--> Deodulator
```

## Developer Notes
The expected configurations of entities and relationships is specified in the comments of the protocol buffers. A summary is below:
* There must be an `RK_CHARACTERIZES` relationship from an `AdaptiveCodingAndModulation` (ACM) to each `Modulator` and `Demodulator` for which this ACM configuration applies: `AdaptiveCodingAndModulation` --`RK_CHARACTERIZES`--> `Modulator` or `AdaptiveCodingAndModulation` --`RK_CHARACTERIZES`--> `Demodulator`. 
* There must be an `RK_CONTAINS` relationship from a `Platform` to each `Antenna` that is attached to the platform: `Platform` --`RK_CONTAINS`--> `Antenna`. The `Platform` can have multiple outgoing `RK_CONTAINS` relationships if it has multiple antennas, such as for a satellite with multiple antennas.
* There must be an `RK_ORIGINATES` relationship from a `Port` to the first entity in the transmit chain, such as: `Port` --`RK_ORIGINATES`--> `Modulator`.  
* There must be an `RK_TERMINATES` relationship from a `Port` to the final entity in the receive chain, such as: `Port` --`RK_TERMINATES`--> `Demodulator`.
* Between the physical layer entities, there must be an `RK_SIGNAL_TRANSITS` relationship to indicate the direction that a signal would flow through the devices.    
* The `Modulator` and `Demodulator` entities can be omitted if no modulation or demodulation occurs.

There must be a 1:1 relationship between a `Port`, a `Transmitter` or `Receiver` or both, and an `Antenna`. In other words, for each `Antenna`, there must exist a single `Transmitter` or a single `Receiver` or both that has an `RK_SIGNAL_TRANSITS` relationship with it, and there must exist a single `Port` in the chain of entities.

All fields are optional unless specified as 'Required' in the comments of the protocol buffers.
