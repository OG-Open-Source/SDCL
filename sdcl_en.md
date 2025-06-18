# SDCL (OGATA's Standard Data Character Storage Language) Specification

## 1. Scope

This document defines the syntax and core features of SDCL (OGATA's Standard Data Character Storage Language). SDCL is designed as a simple, human-readable, and machine-parsable data storage format, specifically for configuration management and data referencing. This specification covers the structure of SDCL documents, assignment syntax, and referencing mechanisms.

This specification is intended for SDCL parser developers and application developers who need a straightforward format for configuration. SDCL's design prioritizes simplicity and strictness to ensure unambiguous parsing.

This specification does not cover the implementation details of SDCL parsers in specific programming languages.

## 2. Terms and Definitions

This document uses the following terms and definitions:

- **Object:** An unordered collection of key-value pairs enclosed in curly braces `{}`.
- **Array:** An ordered collection of values enclosed in square brackets `[]`.
- **Key:** The name used to identify a data item.
- **Value:** The data associated with a key. All values are treated as strings.
- **Path:** A sequence of keys separated by dots (`.`), used to uniquely identify a nested value or structure.
- **Reference:** A mechanism that allows a value to link to another value or content within the same file or from an external source.
- **Parser:** A software component that reads an SDCL document and converts it into an internal data structure.

## 3. Symbols and Abbreviations

- **SDCL:** OGATA's Standard Data Character Storage Language.
- **UTF-8:** A widely used Unicode character encoding scheme.

## 4. Specification

### 4.1. Document Structure

- **Character Encoding:** SDCL documents **must** be encoded in UTF-8.
- **Line Endings:** The newline character (`\n`, LF) is used to separate entries. Parsers **must** ignore the carriage return character (`\r`, CR) when present.
- **Indentation:** Indentation **must** be done using tabs (`\t`). The use of spaces for indentation is **not allowed**. This ensures consistent formatting across all documents.

### 4.2. Syntax Rules

- **Assignment Syntax:** SDCL uses distinct operators for assigning simple values versus structured types.

  - **Simple Value Assignment (`=`):** To assign a value to a key, the equals sign (`=`) **must** be used.
    - `key = value`
  - **Structured Type Assignment (`:`):** To assign an object (`{...}`) or an array (`[...]`) to a key, the colon (`:`) **must** be used.
    - `key: { ... }`
    - `key: [ ... ]`

- **Keys:**

  - Key names **must not** contain spaces, dots (`.`), or be enclosed in quotes.
  - Full paths (e.g., `XXX.YYY.ZZZ`) are used to represent nested structures, where each dot indicates a deeper level of nesting.
  - Within the same object, duplicate key names are **not allowed**. A parser **must** report an error if duplicates are found.

- **Values:**

  - All values are treated as literal strings. They **must not** be enclosed in quotes.
  - A value consists of all characters from the first non-whitespace character after the `=` operator to the end of the line (before any comment).
  - Leading and trailing whitespace in a value **should** be trimmed by the parser.

- **Objects and Arrays:**

  - Objects are defined by `key: { ... }`.
  - Arrays are defined by `key: [ ... ]`.
  - Elements within objects and arrays are separated by newlines.

- **Comments:**
  - Comments begin with a hash symbol (`#`) and continue to the end of the line.
  - Comments **must** appear on their own dedicated line. End-of-line comments are invalid.

### 4.3. Referencing

SDCL supports referencing to enhance data reusability.

- **Value Reference:**

  - **Syntax:** `key = (path.to.value)`
  - **Purpose:** Assigns the value of another key to the current key.

- **Content Inclusion (Objects/Arrays):**

  - **Syntax (without key):** `(path.to.structure)`
  - **Purpose:** Embeds the content of a referenced object or array into the current scope.
  - **Syntax (with key):** `((path.to.structure))`
  - **Purpose:** Embeds an entire structure, including its key, into the current scope.

- **External References:**
  - **Environment Variable:** `key = .env.VAR_NAME`
  - **File Reference:** `key = .path/to/file.sdcl.key`

### 4.4. Front Matter (Metadata Block)

SDCL supports an optional "front matter" block for embedding metadata at the beginning of a file.

- **Syntax:** A front matter block **must** begin on the first line of a file and be enclosed by lines containing only three hyphens (`---`).
- **Parsing:** The content between the `---` delimiters is parsed as a standard SDCL document. All content following the closing `---` is ignored by the SDCL parser.

## 5. SDCL Syntax Examples

```sdcl
# This is a comment.

# --- Basic Assignments ---
application.name = My SDCL App
application.version = 1.2.0
application.enabled = true

# --- Object Assignment ---
server: {
	port = 8080
	host = localhost
}

# --- Array Assignment ---
features: [
	userManagement
	reporting
	notifications
]

# --- Nested Structures ---
database.connection: {
	user = admin
	password = (database.credentials.password) # Value Reference
}
database.credentials.password = .env.DB_PASSWORD # External Reference

# --- Content Inclusion ---
default_settings: {
	timeout = 30
	retries = 3
}

service_config: {
	(default_settings) # Includes timeout and retries
	retries = 5        # Overrides the included value
}

# --- Front Matter Example in a Markdown file ---
# ---
# title = My Document
# author = Kilo Code
# tags: [ sdcl specs example ]
# ---
# # Main Content
# This part is ignored by the SDCL parser.
```

## 6. Conclusion

This document provides the complete specification for the SDCL data format. Its minimalist design, combining unquoted values with powerful referencing, offers a simple yet effective solution for configuration management.
