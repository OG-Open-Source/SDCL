# SDCL (OGATA's Standard Data Character Storage Language) Specification

## 1. Scope

This document defines the syntax and core features of SDCL (OGATA's Standard Data Character Storage Language). SDCL is a human-readable data format designed for simplicity and clarity in configuration files.

## 2. Document Structure

- **Character Encoding:** SDCL files **must** be encoded in UTF-8.
- **Line Endings:** The newline character (`\n`, LF) is the standard line separator. Parsers **must** ignore carriage return characters (`\r`, CR).
- **Indentation:** Indentation **must** be done using tabs (`\t`). The use of spaces for indentation is **not allowed**.

## 3. Core Components

### 3.1. Key-Value Pairs

- The basic structure of SDCL is the key-value pair.
- **Syntax:** `KEY VALUE`
- A key and its value are separated by one or more spaces.
- **Keys:**
  - **Must not** contain spaces.
  - The dot (`.`) character is a valid part of a key and does not represent a hierarchy.
  - Keys are case-sensitive.
- **Values:**
  - Values are parsed literally from the first non-space character after the key to the end of the line.
  - A value is treated as a `string` and **must** be enclosed in double quotes (`""`), unless it is a `number`, `boolean`, `null`, or a reference.
  - For example, `"My Awesome App"` is a single string value.
  - Double quotes (`""`) are mandatory for string values to avoid ambiguity with other data types.

### 3.2. Data Types

SDCL defines several data types for clarity and structure:

- **`string`:** A sequence of characters enclosed in double quotes (`""`). This is required for all textual data to distinguish it from other types like `number`, `boolean`, `null`, or references.
- **`number`:** Represents both integers and floating-point numbers.
- **`boolean`:** `true` or `false`.
- **`null`:** Represents the absence of a value.
- **`object`:** A collection of key-value pairs, enclosed in curly braces (`{}`).
- **`array`:** An ordered list of values, enclosed in square brackets (`[]`).
- **`insertion`:** A mechanism to include or merge data, denoted by parentheses (`()`).

### 3.3. Comments

- Comments start with a hash symbol (`#`) and extend to the end of the line.
- Comments **must** be on their own line.

## 4. Advanced Features

### 4.1. Objects (Dictionaries/Maps)

- Objects are used to group related key-value pairs.
- **Syntax:**
  ```sdcl
  key: {
  	inner_key1 "value1"
  	inner_key2 "value2"
  }
  ```
- The opening brace `{` must be on the same line as the key and colon, and the closing brace `}` must be on a new line.

### 4.2. Arrays (Lists)

- Arrays are ordered collections of values.
- **Syntax:**
  ```sdcl
  key: [
  	"value1"
  	"a string value"
  	123
  	true
  ]
  ```
- The opening bracket `[` must be on the same line as the key and colon, and the closing bracket `]` must be on a new line.

### 4.3. Referencing and Inclusion

SDCL supports referencing to reuse data.

- **Internal Reference (Insertion):**

  - **Syntax:** `(path.to.value)`
  - **Purpose:** Inserts the value of another key. If used at the top level of an object, it merges the referenced object's content.

- **External Reference:**
  - **Syntax:** `.[source].(key)`
  - **Environment Variables:** `.[env].(VAR_NAME)`
  - **File Inclusion:** `.[path/to/file.sdcl].(key_in_file)`

## 5. Front Matter

- SDCL supports an optional "front matter" block for metadata, enclosed by `---` at the beginning of a file.
- The content within the `---` is parsed as SDCL. The content after is ignored.

## 6. Syntax Examples

```sdcl
# Basic key-value pairs
app.name "My Awesome App"
app.version "1.0.0"
is_enabled true

# Object definition
database: {
	host "localhost"
	port 5432
	user "admin"
	password .[env].(DB_PASS) # External reference
}

# Array definition
features: [
	"User Authentication"
	"Data Processing"
	"reporting"
]

# Inclusion
base_settings: {
	timeout 30
	retries 3
}

production_settings: {
	(base_settings) # Merges base_settings here
	timeout 60      # Overrides timeout
}

# Referencing a value
admin_user (database.user)
```
