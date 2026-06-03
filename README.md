# go-discovery

## Project Description

The core functionality of this project is to discover network devices within a Local Area Network (LAN). Currently, it only supports the ONVIF protocol.

Parts of the modules in this project are derived from [use-go/onvif](https://github.com/use-go/onvif). Since `use-go/onvif` has not been updated for a long time, issues arose when referencing it in practical projects. Therefore, we extracted the device discovery functionality code from `use-go/onvif`, updated its dependencies, and further encapsulated and refined its features. We would like to express our sincere gratitude to the authors of `use-go/onvif`.
