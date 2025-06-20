---
Version "1.0"
Date "2025-06-20"
---
# SDCL (OGATA's Standard Data Character Storage Language) Documentation

This repository contains the official documentation for SDCL (OGATA's Standard Data Character Storage Language), a human-readable and machine-parsable data storage format designed for configuration management and data referencing.

## Key Features

- **Human-Readable Syntax:** Designed for simplicity and clarity, making it easy to write and maintain configuration files.
- **Structured Data:** Supports nested objects (sections) and lists (arrays) to represent complex data hierarchies.
- **Strict Typing:** Enforces explicit data types, such as double-quoted strings, numbers, booleans, and null, to reduce ambiguity.
- **Powerful Referencing:**
  - **Shallow Merge `()`:** Merges the contents of one section into another.
  - **Section Insertion `(())`:** Inserts an entire section as a nested object.
  - **External Data:** Supports referencing environment variables and including content from other files.

## Full Specification

The complete language specification is provided in an ISO standard-style format across multiple languages to ensure broad accessibility and clarity.

- **English:** [`sdcl_en.md`](sdcl_en.md)
- **Simplified Chinese (简体中文):** [`sdcl_zh_hans.md`](sdcl_zh_hans.md)
- **Traditional Chinese (繁體中文):** [`sdcl_zh_hant.md`](sdcl_zh_hant.md)

Each document provides a comprehensive guide to SDCL's syntax, semantics, data types, and referencing mechanisms.

## Purpose

The SDCL documentation serves as a definitive reference for:

- SDCL parser developers
- Data architects
- Application developers utilizing SDCL for data exchange or configuration management

## Contribution

For any contributions, feedback, or suggestions regarding the SDCL specification or its documentation, please open an issue in this repository.
