# SDCL (OGATA's Standard Data Character Storage Language) Specification

## 1. Scope

This document defines the syntax, semantics, and core features of SDCL (OGATA's Standard Data Character Storage Language). SDCL is designed as a human-readable and machine-parsable data storage format, specifically for configuration management and data referencing. This specification covers the structure of SDCL documents, data types, referencing mechanisms, comment rules, and the representation of expanded and compact forms.

This specification is intended for SDCL parser developers, data architects, and application developers who need to use SDCL for data exchange or configuration management. SDCL's design draws upon the simplicity of JSON, the readability of YAML, and TOML's friendliness for configuration files, while introducing powerful internal and external referencing mechanisms to promote modularity and reusability of configurations, all while maintaining a strict syntax to simplify parsing and reduce ambiguity.

This specification does not cover the implementation details of SDCL parsers in specific programming languages, nor does it cover network transmission protocols.

## 2. Terms and Definitions

This document uses the following terms and definitions, arranged alphabetically:

- **Array:** An ordered collection of values enclosed in square brackets `[]`.
- **Compact Form:** An SDCL format that allows objects and arrays to be defined on a single line, using spaces as delimiters between elements.
- **Content Inclusion:** A referencing mechanism used to embed the content (without its key name) or the entire structure (including its key name) of a specified object or array from a given path into the current scope.
- **Expanded Form:** An SDCL format that uses newlines and indentation to clearly define the hierarchical structure of objects and arrays.
- **Environment Variable:** A named variable provided by the operating system or execution environment, external to the SDCL document.
- **Key:** The name used to identify a data item.
- **Literal:** A fixed representation of a value in source code; for example, `123` is a numeric literal, `"hello"` is a string literal.
- **Object:** An unordered collection of key-value pairs enclosed in curly braces `{}`.
- **Parser:** A software component responsible for reading an SDCL document, validating its syntax, and converting it into an internal data structure usable by an application.
- **Path:** A sequence of keys separated by dots (`.`), used to uniquely identify a nested value or structure within an SDCL document.
- **Reference:** A mechanism that allows a value to link directly to another value or content within the same file or from an external source (such as an environment variable or another SDCL file).
- **Scalar Value:** A single, indivisible unit of data, such as a string, number, boolean, null, date, time, datetime, country code, or Base64 encoded data.
- **Structured Type:** Refers to objects or arrays, which can contain other values (including scalar values or other structured types).
- **Value:** The data associated with a key.

## 3. Symbols and Abbreviations

- **ISO 3166-1 alpha-2:** International Organization for Standardization standard for country codes.
- **ISO 8601:** International Organization for Standardization standard for date and time representation.
- **RFC 4648:** Related standard for Base64 encoding.
- **SDCL:** OGATA's Standard Data Character Storage Language.
- **UTF-8:** A widely used Unicode character encoding scheme.

## 4. Specification

### 4.1. Data Types and Referencing Rules

SDCL defines a precise set of basic data types, each with clear referencing requirements to ensure unambiguous parsing and data integrity.

- **Strings:** A sequence of Unicode characters. All string values **must** be enclosed in double quotes (`""`). Multi-line strings are not supported; all string content must reside on a single logical line.

  - Example: `"hello world"`, `"user@example.com"`, `"123 Main St."`

- **Numbers:** Represent integer and floating-point numerical values. Numbers **must not** be enclosed in quotes.

  - **Integers:** Whole numbers (e.g., `123`, `-45`).
  - **Floating-Point Numbers:** Numbers with a decimal part (e.g., `3.14`, `-0.001`).
  - Example: `123`, `3.14`, `-100`, `0.5`

- **Booleans:** Represent logical truth values. Only `true` or `false` are recognized, and they **must not** be enclosed in quotes.

  - Example: `true`, `false`

- **Null:** Represents the absence of a value. Only `null` is recognized, and it **must not** be enclosed in quotes.

  - Example: `null`

- **Dates:** Represent calendar dates. Values **must** follow the `YYYY-MM-DD` format and **must not** be enclosed in quotes.

  - Example: `2025-05-27`, `1999-12-31`

- **Times:** Represent a time of day. Values **must** follow the `HH:MM:SS` format (24-hour clock) and **must not** be enclosed in quotes.

  - Example: `14:30:00`, `09:15:05`

- **Datetimes:** Represent a specific point in time, combining date and time. Values **must** follow the ISO 8601 format (e.g., `YYYY-MM-DDTHH:MM:SSZ` for UTC or `YYYY-MM-DDTHH:MM:SS+HH:MM` with offset) and **must not** be enclosed in quotes.

  - Example: `2025-05-27T14:30:00Z`, `2023-10-27T10:00:00+08:00`

- **Countries:** Represent ISO 3166-1 alpha-2 country codes. Values **must** consist of exactly two uppercase ASCII letters and **must not** be enclosed in quotes.

  - Example: `TW`, `US`, `JP`

- **Base64:** Represents Base64 encoded binary data. Values **must** follow Base64 encoding rules and **must not** be enclosed in quotes. Commonly used for embedding small binary data, such as images or cryptographic keys.
  - Example: `SGVsbG8gV29ybGQ=` (Base64 encoded "Hello World")

**Referencing Rules Summary:**

- **Unquoted Types:** `int`, `float`, `bool`, `null`, `country`, `date`, `time`, `datetime`, and `base64` values **must not** be enclosed in double quotes (`""`).
- **Quoted Types:** Only `string` values **must** be enclosed in double quotes (`""`).

### 4.1.10 Character Encoding

SDCL documents should default to UTF-8 encoding. All string data is considered Unicode strings.

### 4.2. Strictness and Syntax

SDCL enforces a strict set of syntax rules to ensure unambiguous parsing and consistent data representation. Adherence to these rules is crucial for valid SDCL documents.

- **Key Naming Conventions:**

  - Key names, especially the last part in a dot-separated path (e.g., `YYY` in `XXX.YYY`), **must not** contain spaces or dot (`.`) symbols. This ensures a clear hierarchy. Key names themselves **must not** be enclosed in quotes.
  - Full paths (e.g., `XXX.YYY.ZZZ`) are used to represent nested structures, where each dot indicates a deeper level of nesting.
  - Within the same object scope, identical key names **must not** appear. If a parser encounters duplicate key names within the same object scope, it **must** throw a parsing error and terminate processing.

- **String Value Delimitation:**

  - All literal string values **must** be enclosed in double quotes (`""`).
  - Single quotes (`''`) are not supported as string delimiters.
  - Multi-line string formats (e.g., Python's triple quotes `"""` or YAML's block scalars `|`/`>`) are not supported. All string content must be defined on a single logical line.

- **Structural Elements:**

  - **Objects (Tables):** Represented by key-value pairs enclosed in curly braces (`{}`). Each key-value pair within an object should be on a new line in expanded form, or separated by at least one space in compact form.
  - **Arrays:** Represented by a sequence of values enclosed in square brackets (`[]`). Elements within an array should be on a new line (for complex elements like objects) or separated by at least one space (for simple scalar values) in expanded form, and separated by at least one space in compact form.

- **Basic Key-Value Pair Syntax:**

  - The basic syntax for defining a data entry is `key: value`. A colon (`:`) separates the key from its corresponding value.

- **Object Scope Definition:**
  - Key-value pairs that logically belong to an object **must** be explicitly enclosed in curly braces (`{}`). This clearly defines the scope and hierarchy of the data.

### 4.3. Delimiters and Comments

- **Delimiters:** SDCL employs specific delimiters to demarcate data elements and structures:

  - **Newlines:** Primarily used to separate different key-value pairs or elements in multi-line contexts (e.g., within objects or arrays). Each new line indicates a new entry or element within a structured block.
  - **Spaces:** Used to separate elements in single-line contexts (especially within arrays). Multiple spaces are treated as a single delimiter.
  - **Commas (`,`):** Explicitly **forbidden** as delimiters for array elements or key-value pairs. The prohibition of commas simplifies the syntax, reduces confusion with formats like JSON, and promotes a clearer newline/space-separated style. This strict rule enhances clarity and avoids ambiguity with other data formats.

- **Comments:** SDCL supports the inclusion of comments for documentation and clarity:
  - **Single-line Comments:** Begin with a hash symbol (`#`). All text from `#` to the end of the line is considered a comment and is ignored by the parser.
  - **Placement:** Comments **must** appear on their own dedicated line.
  - **Forbidden End-of-line Comments:** Comments are not allowed to follow a key-value pair or any other data element on the same line (e.g., `key: "value" # comment` is invalid). This rule ensures consistent formatting and simplifies parsing.

### 4.4. SDCL Unique Features

SDCL stands out through several unique features designed to enhance data management, modularity, and reusability within configuration and data files.

- **Explicit Type Tags (Optional):**

  - SDCL provides optional `<type></type>` tags to explicitly declare the expected format or type of a value. While parsers can generally infer types based on literal representation and quoting rules, explicit type tags are particularly useful in the following situations:
    1.  Enhancing readability, clearly indicating the expected type of a value.
    2.  In some cases, avoiding potential ambiguity (e.g., an unquoted value might be misinterpreted as a specific type).
    3.  Serving as a validation mechanism, ensuring the value conforms to the expected format.
  - **Supported Tags:** `<int></int>`, `<float></float>`, `<bool></bool>`, `<str></str>`, `<date></date>`, `<time></time>`, `<datetime></datetime>`, `<country></country>`, `<base64></base64>`.
  - **Usage:** For primitive types (e.g., `int`, `float`, `bool`, `null`, `string`, `country`, `base64`), these tags are optional. If present, the parser will prioritize the explicit tag; otherwise, it infers the type.
  - **Example:**

  ```sdcl
  # Explicitly tagged datetime
  startDate: <datetime>2024-05-26T18:30:00Z</datetime>
  # Explicitly tagged country
  userCountry: <country>TW</country>
  # Explicitly tagged Base64
  imageData: <base64>T0dBVEEncyBTdGFuZGFyZCBEYXRhIENoYXJhY3RlciBTdG9yYWdlIExhbmd1YWdl</base64>
  # These are functionally equivalent to their untagged, correctly quoted counterparts (e.g., `startDate: 2024-05-26T18:30:00Z`).
  ```

- **Path-Based Referencing and Content Inclusion:**

  - SDCL utilizes dot notation (`.`) to represent hierarchical paths within a document, where `XXX.YYY: value` is parsed as `XXX: { YYY: value }`. This notation is central to SDCL's powerful referencing and inclusion mechanisms:

  1.  **Value Reference (within key-value pairs):**

      - **Syntax:** `key: (path.to.value)`
      - **Purpose:** This syntax allows the value of a key to dynamically reference another value located at a specified path _within the same SDCL file_. This creates a live link: if the referenced value changes, the value of the referencing key will automatically update.
      - **Example:**

      ```sdcl
      base.config.defaultLogLevel: "DEBUG"
      # service.logging.level will be "DEBUG"
      service.logging.level: (base.config.defaultLogLevel)
      ```

  2.  **Content Inclusion (Objects/Arrays, without key name):**

      - **Syntax:** `(path.to.object_or_array)` (on its own line, without a colon and value)
      - **Purpose:** This mechanism is used to directly embed the _content_ (nested elements) of an already defined object or array into the current scope. It functions similarly to YAML's merge key (`<<: *alias`) but for direct content injection without including the source key name itself. If the referenced path does not point to an existing object or array, or if the resolved value is not the expected type (e.g., for object inclusion, the path points to a scalar value), the parser **must** throw an error.
      - **Example:**

      ```sdcl
      common_resource_limits: {
        cpu: "500m"
        memory: "512Mi"
      }
      service.resources.limits: {
        # Inserts cpu: "500m" and memory: "512Mi"
        (common_resource_limits)
        # Overrides the cpu value from common_resource_limits
        cpu: "250m"
      }
      # Resulting structure for service.resources.limits: { cpu: "250m", memory: "512Mi" }
      # This is JSON syntax; SDCL does not use commas as separators.
      ```

  3.  **Content Inclusion (Objects/Arrays, with key name):**
      - **Syntax:** `((path.to.object_or_array))`
      - **Purpose:** This syntax is used to include the _entire_ object or array, including its key name and all content, from a specified path into the current scope. This is suitable for embedding complete, named data structures. This mechanism only supports objects and arrays, not simple key-value pairs.
      - **Example:**
      ```sdcl
      users: [
        { id: 1 name: "Alice" }
        { id: 2 name: "Bob" }
      ]
      user_data_container: {
        # Inserts the entire users array
        ((users))
      }
      # Resulting structure for user_data_container: { users: [ {id:1, name:"Alice"}, {id:2, name:"Bob"} ] }
      # This is JSON syntax; SDCL does not use commas as separators.
      ```

- **External References:**

  - **Syntax:** `.env.KEY`
  - **Purpose:** Values starting with a dot (`.`) and not enclosed in double quotes (`""`) are specifically interpreted as references to external environment variables. This allows SDCL configurations to dynamically pull values from the execution environment. Parsers should provide mechanisms to read environment variables from the execution environment. It should be noted that storing sensitive information in environment variables may pose security risks and should be carefully evaluated.
  - **Supported Formats:** Currently, only `.env.KEY` and `.XXX.sdcl.KEY` formats are supported, where `KEY` is the name of the environment variable to be referenced, and `XXX` is the name of the SDCL file to be referenced.
  - **Example:**

  ```sdcl
  # Refers to an environment variable named API_KEY
  api.key: .env.API_KEY

  # Refers to a key named 'database.port' in an external SDCL file named 'config.sdcl'
  my.external.port: .config.sdcl.database.port
  ```

  - **`.XXX.sdcl.KEY` Reference:** This syntax is used to reference a value in another SDCL file (`XXX.sdcl`). Parsers need to define a file resolution strategy. Typically, this might involve:
    - Paths relative to the current SDCL file.
    - A predefined list of search paths, for example: Parsers can support configuring one or more base directories as search paths for referencing external SDCL files.
      Parsers **must** prevent circular references (e.g., A references B, B references A), and if a circular reference is detected, an error **should** be thrown. If the referenced external file does not exist or cannot be read, or if the specified path within the file does not exist, the parser **must** throw an error.

- **Content Overriding/Merging:**

  - The "last definition wins" principle in SDCL applies only to content introduced through **content inclusion** or **value referencing** mechanisms. This means that when a key's value or content is introduced into the current scope via a reference, if that key already exists in the current scope, the newly introduced value will override the old one.
  - **Duplicate Directly Defined Keys:** Within the same scope, duplicate key names that are directly defined are **not allowed**. If a parser encounters duplicate directly defined key names within the same scope, it **must** throw a parsing error and terminate processing.
  - Example:

  ```sdcl
  # Overriding via content inclusion
  default_settings: {
      timeout: 1000
      retries: 3
  }
  service_config: {
      (default_settings)
      retries: 5 # Overrides retries from default_settings
  }
  # Result: service_config: { timeout: 1000, retries: 5 }

  # Overriding via value reference
  base.url: "http://example.com"
  api.endpoint: (base.url)
  base.url: "http://new-example.com" # Overrides base.url, api.endpoint will also update
  # Result: api.endpoint: "http://new-example.com"

  # Invalid: Duplicate directly defined key
  # invalid_config: {
  #   key1: "value1"
  #   key1: "value2" # Error: Duplicate key name
  # }
  ```

### 4.5. Expanded and Compact Forms

SDCL supports two primary forms for representing data structures: Expanded Form and Compact Form. These two forms provide flexibility to accommodate different readability and space efficiency needs.

- **Expanded Form:**

  - Expanded form uses newlines and indentation to clearly define the hierarchical structure of objects and arrays. Each key-value pair or array element typically resides on its own line, enhancing readability, especially for complex or nested data structures.
  - **Indentation:** In expanded form, indentation is primarily used to enhance human readability and is not syntactically enforced. Parsers should primarily rely on curly braces `{}` and square brackets `[]`, as well as newlines (or spaces between elements within objects/arrays) to determine structure and separation. However, consistent indentation style (e.g., 2 spaces or 4 spaces) is recommended to improve maintainability.
  - **Example:**

  ```sdcl
  application: {
    name: "My App"
    version: 1.0
    settings: {
      debug: true
      logLevel: "INFO"
    }
  }
  ```

- **Compact Form:**
  - Compact form allows objects and arrays to be defined on a single line, using spaces as delimiters between elements. This form is more compact in terms of space efficiency and is suitable for simple data structures or when minimizing file size is desired.
  - **Example:**
  ```sdcl
  application:{name:"My App" version:1.0 settings:{debug:true logLevel:"INFO"}}
  features:["userManagement" "reporting" "notifications"]
  ```
- **Non-Interchangeable:**
  - Within a single SDCL document, expanded and compact forms **must not** be mixed. A document **must** fully adhere to one form to ensure consistent and predictable parsing.
  - **Choice:** Expanded form is recommended for complex, nested data structures or situations requiring high readability. Compact form is suitable for simple data structures, single-line configurations, or scenarios where minimizing file size is desired.

## 5. SDCL Syntax Examples

This section provides comprehensive examples of SDCL syntax, covering various data types, structures, and referencing mechanisms.

```sdcl
# This is a single-line comment in SDCL.
# All comments must reside on their own dedicated line.

# SDCL Application Example Configuration

# Basic key-value pairs within an object (key: value)
application: {
  name: "My SDCL App"
  version: 1.0
  enabled: true
  debugMode: false
}

# Numbers and Null values (these type tags are optional)
server: {
  port: 8080
  timeout: 30.5
  maxConnections: 100
  logLevel: "INFO"
  adminEmail: null
}

# Date, Time, Datetime, Country, and Base64 types (explicit type tags are optional for inference)
event: {
  startDate: <datetime>2024-05-26T18:30:00Z</datetime>
  # Datetime without explicit tag
  endDate: 2024-05-27T09:00:00Z
  eventDate: <date>2024-05-26</date>
  # Time without explicit tag
  eventTime: 18:30:00
  # Country without explicit tag
  originCountry: TW
  # Country with explicit tag
  destinationCountry: <country>US</country>
  # Base64 with explicit tag
  profileImage: <base64>T0dBVEEncyBTdGFuZGFyZCBEYXRhIENoYXJhY3RlciBTdG9yYWdlIExhbmd1YWdl</base64>
  # Base64 without explicit tag
  documentContent: VGhpcyBpcyBhIHRlc3QgZG9jdW1lbnQu
}

# Nested configuration using dot notation (as single-line key-values)
database.type: "PostgreSQL" # Equivalent to { "database": { "type": "PostgreSQL" } }
database.host: "localhost" # Equivalent to { "database": { "host": "localhost" } }
database.port: 5432 # Equivalent to { "database": { "port": 5432 } }
database.user: "admin" # Equivalent to { "database": { "user": "admin" } }
database.password: "secure_password_123" # Equivalent to { "database": { "password": "secure_password_123" } }

# Defining nested objects directly using key: { ... }
database.connectionPool: {
  maxSize: 20
  idleTimeout: 60000
}

# Array of strings (elements separated by spaces or newlines)
features.enabledFeatures: [
  "userManagement"
  "reporting"
  "notifications"
]

# Array of objects (objects separated by spaces or newlines)
users: [
  {
    id: 1
    name: "Alice"
    email: "alice@example.com"
  }
  {
    id: 2
    name: "Bob"
    email: "bob@example.com"
  }
]

# Value reference example (internal reference)
# First, define a base value
base.config.defaultLogLevel: "DEBUG"

# Now, reference it
# This value will be "DEBUG"; if base.config.defaultLogLevel changes, it will automatically update.
service.logging.level: (base.config.defaultLogLevel)

# External reference example using dot prefix
# Refers to an environment variable named API_KEY
api.key: .env.API_KEY

# Refers to a key named 'database.port' in an external SDCL file named 'config.sdcl'
my.external.port: .config.sdcl.database.port

# Configuration for a specific service
service.api: {
  baseUrl: "https://api.example.com/v1"
  apiKey: "your_api_key_here"
  rateLimit: {
    requestsPerMinute: 100
    burst: 10
  }
}

# Content inclusion example (object, without key name) using (path.to.object)
# This includes the content of database.connectionPool. If keys conflict, later definitions will override earlier ones.
db_settings: {
  (database.connectionPool)
  # Overrides the 'maxSize' value from database.connectionPool
  maxSize: 25
}

# Content inclusion example (array, without key name) using (path.to.array)
# This demonstrates merging elements of one array into another.
additional_features: ["adminPanel" "analyticsDashboard"]
all_features: [
  "userManagement"
  # Includes elements from the 'additional_features' array.
  (additional_features)
]

# Content inclusion example (object, with key name) using ((path.to.object))
user_data_container: {
  # Includes the entire 'users' array, including its key 'users'.
  ((users))
}

# Content inclusion example (object, with key name) using ((path.to.object)) for nested objects
service_limits: {
  # Includes the entire 'rateLimit' object, including its key 'rateLimit'.
  ((service.api.rateLimit))
}

# Value reference example (standalone) using (path.to.value)
# This demonstrates a value reference.
# If 'database.port' is 5432, then 'my.referenced.port' will effectively be 5432.
my.referenced.port: (database.port)

# Value reference example (standalone) for object content inclusion using (path.to.object)
# This demonstrates object content inclusion without a key name.
direct_db_settings: {
  (database.connectionPool)
  # Overrides the 'idleTimeout' value from 'database.connectionPool'.
  idleTimeout: 70000
}

# Value reference example (standalone) for array content inclusion using (path.to.array)
# This demonstrates array content inclusion without a key name.
direct_user_list: [
  (users)
]

# Expanded and Compact Forms Examples
expanded_example: {
  key1: "value1"
  nested_object: {
    nested_key: 123
    another_key: true
  }
  array_example: [
    "item1"
    "item2"
  ]
}

compact_example:{name:"My App" version:1.0 settings:{debug:true logLevel:"INFO"}}
features:["userManagement" "reporting" "notifications"]

# Invalid Syntax Examples
# Invalid: End-of-line comment
key: "value" # This is an invalid end-of-line comment

# Invalid: Comma used in array
# invalid_array: [ "a", "b" ]

# Invalid: Single quotes used for string
# invalid_string: 'test'

# Invalid: Key name contains spaces (in the last part)
# "invalid key name": "value"

# Invalid: Duplicate key within object scope
# duplicate_key_object: {
#   myKey: "value1"
#   myKey: "value2"
# }

# Invalid: Mixing expanded and compact forms
# mixed_form_example: { key1: "value1"
#   key2: "value2" }
```

## 6. Conclusion

This document serves as the comprehensive design specification for SDCL, encompassing all agreed-upon features, syntax rules, and examples. It reflects the collaborative effort made to define a robust and intuitive data storage language. Through its unique combination of features, SDCL aims to be a powerful, flexible, and easy-to-use data storage language, particularly suitable for modern application configuration management and data representation needs.
