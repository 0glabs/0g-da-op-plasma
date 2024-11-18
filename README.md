<br />
<br />
  <p align="center" href="https://optimism.io/">
  <img alt="Optimism" src="https://raw.githubusercontent.com/ethereum-optimism/brand-kit/main/assets/svg/OPTIMISM-R.svg" width=180>
  </p>
  <h3 align="center">+</h3>
  <p align="center">
  <img src="https://framerusercontent.com/images/JJi9BT4FAjp4W63c3jjNz0eezQ.png" alt="Logo" width="80" height="80">
  </p>
<br />


## Optimism x 0G DA Integration Guide

### Overview:

The Optimism 0G integration allows developers to deploy an OP Stack Rollup using 0G (Zero Gravity) for data availability. This integration marks the first AI-focused external integration to the OP Stack protocol, offering an alternative to Ethereum for high-performance data availability.


### Key Components

1. DA-server Implementation: This repository implements a da-server that runs as a sidecar process alongside the OP Stack rollup node with generic commitments. This server connects to a 0G DA client to securely communicate with the 0G DA network.
2. DA Client and DA Encoder Implementations: These implementations are used to interact with the 0G DA network and to encode/decode data for storage and retrieval.
3. OP Stack: OP Stack implementation with configuration adjustments

### Documentation

To get started with the Optimism 0G implementation, please refer to the [OP Stack on 0G DA Documentation](https://docs.0g.ai/build-with-0g/rollups-and-appchains/op-stack-on-0g-da).

### Learn More About 0G

[0G Website](https://0g.ai/)
[0G Github](https://github.com/0glabs)

### Learn More About Arbitrum Orbit

[OP Stack Docs](https://docs.arbitrum.io/launch-orbit-chain/orbit-gentle-introduction)
