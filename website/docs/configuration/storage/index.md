---
sidebar_position: 1
title: "Basic Configuration"
---

# Basic Configuration

Storage is the place where terediX will store the discovered data. You can add multiple storage engines in the
configuration file. You can also add multiple storage engines.

If you add multiple storage engines, you must need to define `default_engine`

| option	           | type   | description	                                                                                                      |
|-------------------|--------|:------------------------------------------------------------------------------------------------------------------|
| 	  batch_size     | number | Number of data to store in a single batch. You can increase the number to speed up the storage process.	          |
| 	  engines        | object | List of storage engines. You can add multiple storage engines.	                                                   |
| 	  default_engine | text   | Name of the default storage engine. If you add multiple storage engines, you must need to define `default_engine` |

### Supported Storage Engines

list of supported storage engines

| storage engine | description	                                                                                    |
|----------------|:------------------------------------------------------------------------------------------------|
| postgresql     | Store data in PostgreSQL database. You can use this storage engine to store data in PostgreSQL. |