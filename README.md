# teredix
`TeReDiX` (Tech Resource Discover &amp; Exploration) is a tool to discover tech resources for an organization and explore them as a resource graph

## Technical Architecture

```
+-------------------+       +------------------------+        +------------------------+
| Source 1          |       | Scanner 1              |        | Storage Engine         |
| - Scanner         +------>+ - Scan()               |        | - Prepare()            |
|                   |       |                        |        | - Persist()            |
+-------------------+       +------------------------+        | - Find()               |
                                                                   +------------------------+
+-------------------+       +------------------------+
| Source 2          |       | Scanner 2              |
| - Scanner         +------>+ - Scan()               |
|                   |       |                        |
+-------------------+       +------------------------+

                           +------------------------+
                           | Processor              |
                           | - Process()            |
                           +------------------------+
```

The architecture consists of several components:

**Sources:** represent the different sources from which resources can be discovered. Each source has its own scanner to discover resources

**Scanners** These are responsible for scanning a particular source and returning a list of resources.

**Storage Engine** This component stores the discovered resources. It is responsible for preparing the storage schema, 
persisting resources, and finding resources based on a filter.

**Processor** This component orchestrates the discovery process. It starts all the scanners in parallel and processes 
the resources as they become available.


## Getting Started

