# SDCL Language Specification

|              |            |
| :----------- | :--------- |
| **Version:** | 1.0        |
| **Date:**    | 2025-06-20 |

---

## Abstract

This document specifies the syntax, semantics, and data model for OGATA's Standard Data Character Storage Language (SDCL). SDCL is a human-readable data format designed for configuration files, with core design goals of exceptional readability, simplicity, and a clear, declarative syntax.

Key features of SDCL include:

- **Concise Key-Value Syntax:** Built upon an intuitive `key value` foundation.
- **Structured Data:** Organizes complex data through nested sections (Objects) and lists (Arrays).
- **Explicit Data Types:** Supports fundamental types such as strings, numbers, booleans, and null, with mandatory double quotes for strings to prevent ambiguity.
- **Data Reusability:** Features built-in referencing and inclusion mechanisms to refer to internal or external data sources (e.g., environment variables, other files).
- **Strict Formatting Rules:** Enforces tab-based indentation for multi-line blocks to ensure consistent document style.

This specification provides a conformance standard for SDCL implementations.

---

## Table of Contents

1.  **Introduction**

    - [1.1 Scope](#11-scope)
    - [1.2 Conformance](#12-conformance)
    - [1.3 Terms and Definitions](#13-terms-and-definitions)
    - [1.4 Design Philosophy](#14-design-philosophy)

2.  **Lexical Structure**

    - [2.1 Character Set](#21-character-set)
    - [2.2 Comments](#22-comments)
    - [2.3 Keywords](#23-keywords)
    - [2.4 Identifiers](#24-identifiers)
    - [2.5 Literals](#25-literals)
    - [2.6 Operators and Punctuators](#26-operators-and-punctuators)

3.  **Syntax**

    - [3.1 Document Structure](#31-document-structure)
    - [3.2 Key-Value Pairs](#32-key-value-pairs)
    - [3.3 Sections](#33-sections)
    - [3.4 Lists](#34-lists)
    - [3.5 References](#35-references)
    - [3.6 Front Matter](#36-front-matter)

4.  **Data Model and Semantics**

    - [4.1 Data Types](#41-data-types)
    - [4.2 Type System](#42-type-system)
    - [4.3 Parsing and Evaluation](#43-parsing-and-evaluation)

- **Appendix A: Grammar Summary**

---

## Chapter 1: Introduction

### 1.1 Scope

This specification defines the lexical structure, syntax, data model, and semantics of OGATA's Standard Data Character Storage Language (SDCL). It is intended to be an authoritative reference for developers of SDCL parsers, editors, and other related tooling.

### 1.2 Conformance

A conforming SDCL implementation must satisfy all of the following conditions:

1.  **Parsing:** It must successfully parse all documents that conform to the grammar defined in Appendix A of this specification. It must report an error for documents that do not conform.
2.  **Character Encoding:** It must handle documents encoded in UTF-8.
3.  **Indentation and Line Endings:** It must treat the tab character (`\t`) as the only indentation character for multi-line blocks and the line feed character (`\n`) as the line separator, while ignoring carriage returns (`\r`).
4.  **Data Model:** It must map the parsed SDCL document to an in-memory representation that conforms to the data model described in Chapter 4.
5.  **Reference Resolution:** It must correctly resolve all internal and external references.

### 1.3 Terms and Definitions

- **Document:** A complete SDCL configuration, typically corresponding to a single text file.
- **Key:** An identifier that names a value in a key-value pair or section definition.
- **Value:** The data associated with a key, which can be a scalar or a compound structure.
- **Key-Value Pair:** The fundamental building block of SDCL, consisting of a key and a value.
- **Section:** A nested structure identified by a key, containing a collection of key-value pairs or other sections. It is equivalent to an object or dictionary.
- **List:** An ordered collection of values. It is equivalent to an array.
- **Comment:** A piece of text that is ignored by the parser, used for explanatory purposes.
- **Reference:** A special syntactic construct used to insert or merge a value defined elsewhere.

### 1.4 Design Philosophy

The design of SDCL is based on the following core principles:

- **Human-First Readability:** The syntax is designed to be easily read and written by humans, minimizing syntactic noise.
- **Simplicity:** Avoids complex logical control structures, focusing on the declarative description of data.
- **Unambiguous:** Eliminates common parsing ambiguities through strict rules (e.g., mandatory string quotes, fixed indentation character).
- **Structured:** Provides nested sections and lists to represent hierarchical and collective data naturally.

## Chapter 2: Lexical Structure

### 2.1 Character Set

SDCL documents **must** be encoded in UTF-8.

### 2.2 Comments

Comments begin with a hash symbol (`#`) and extend to the end of the line. Comments **must** be on their own line and are not permitted following any syntactic element.

```sdcl
# This is a valid comment.
key "value" # This is an invalid comment.
```

### 2.3 Keywords

The following words are reserved as keywords and have special meaning. They cannot be used as unquoted string values or identifiers.

- `true`: Represents the boolean value for truth.
- `false`: Represents the boolean value for falsehood.
- `null`: Represents a null value.

### 2.4 Identifiers

Identifiers are used to name keys. Their naming rules are as follows:

1.  Consist of one or more non-whitespace characters.
2.  **Must not** contain space characters.
3.  The dot character (`.`) is a valid part of an identifier and does not itself represent a hierarchy.
4.  Identifiers are case-sensitive. `myKey` and `mykey` are two distinct identifiers.

**Example:**

```sdcl
# Valid identifiers
app.name "My App"
version-1.0 true
_user "admin"
```

### 2.5 Literals

#### 2.5.1 Strings

String literals are sequences of UTF-8 characters enclosed in double quotes (`"`). This is the **only** way to represent textual data, designed to clearly distinguish it from numbers, booleans, nulls, and references.

**Example:**

```sdcl
message "Hello, World!"
path "C:\\Users\\Default"
empty_string ""
```

#### 2.5.2 Numbers

Number literals represent integers or floating-point numbers.

**Example:**

```sdcl
port 5432
version 1.0
negative_integer -10
scientific_notation 6.022e23
```

#### 2.5.3 Booleans

Boolean literals are `true` or `false`.

**Example:**

```sdcl
enabled true
debug_mode false
```

#### 2.5.4 Null

The null literal is `null`, representing the absence of a value.

**Example:**

```sdcl
optional_feature null
```

### 2.6 Operators and Punctuators

- `:` (Colon): Separates the key of a section or list from its body (`{` or `[`).
- `{ }` (Curly Braces): Delimit the start and end of a Section.
- `[ ]` (Square Brackets): Delimit the start and end of a List.
- `( )` (Parentheses): Delimit the path of a Reference.
- `.` (Dot): Acts as a path separator in reference paths or as part of an identifier.
- `#` (Hash): Marks the beginning of a line as a comment.
- `---` (Triple-Dash): Delimits a Front Matter block.

## Chapter 3: Syntax

### 3.1 Document Structure

An SDCL document consists of a series of statements, primarily key-value pairs and sections. The content of the document **must** be indented using the tab character (`\t`). Each level of indentation represents a level of nesting.

### 3.2 Key-Value Pairs

The key-value pair is the basic structure of SDCL.

- **Syntax:** `KEY VALUE`
- The KEY and VALUE are separated by one or more space characters.

**Example:**

```sdcl
app.name "My Awesome App"
app.version "1.0.0"
is_enabled true
```

### 3.3 Sections

Sections are used to organize related key-value pairs into a nested structure. The syntax for sections requires a multi-line block.

- **Syntax:**

  ```sdcl
  key: {
  	inner_key1 "value1"
  	inner_key2 123
  }
  ```

- The opening curly brace `{` **must** be on the same line as the key and colon.
- The closing curly brace `}` **must** be on a new line, and its indentation level must match that of the key defining the section.
- The content within a section **must** be indented by one level.

### 3.4 Lists

Lists are ordered collections of values. They can be defined in two ways:

- **Multi-line Syntax:**

  - The opening square bracket `[` must be on the same line as the key and colon. The closing square bracket `]` must be on a new line, at the same indentation level as the key.
  - Each element must be on its own line and indented by one level.

  ```sdcl
  multi_line_list: [
  	"value1"
  	123
  	true
  ]
  ```

- **Single-line Syntax:**

  - All elements are placed on a single line within the square brackets, separated by one or more spaces.

  ```sdcl
  single_line_list: [1 2 3 "a string" true (a.b.c)]
  ```

Lists can contain values of different types, including anonymous sections (in multi-line lists).

### 3.5 References

The reference mechanism allows for the insertion or merging of values from another location.

#### 3.5.1 Internal Reference

References a value defined within the same document.

- **Value Reference:** `(path.to.value)`

  - Inserts the value of another key at the current position.

- **Shallow Merge (Content Merge):** `(path.to.section)`

  - When used as a direct child of a section, this merges the _contents_ of the referenced section into the current section. This is a shallow merge; only the top-level key-value pairs are copied.

- **Section Insertion (Object Insertion):** `((path.to.section))`

  - This inserts the entire referenced section (including its key) as a new nested section within the current context.

**Example:**

```sdcl
# Base configuration
base: {
	user "guest"
	log_level "info"
}

# Shallow Merge Example
# The contents of 'base' are merged into 'config_shallow'.
config_shallow: {
	(base)
	log_level "debug" # Overrides the log_level from base
}
# Resulting config_shallow is equivalent to:
# { user: "guest", log_level: "debug" }

# Section Insertion Example
# The 'base' section itself is inserted as a child of 'config_insertion'.
config_insertion: {
	((base))
	another_key "value"
}
# Resulting config_insertion is equivalent to:
# { base: { user: "guest", log_level: "info" }, another_key: "value" }
```

#### 3.5.2 External Reference

References a value from a source external to the document. The syntax is an extension of internal references.

- **Syntax:** `.[source].(key)` or `.[source].((key))`
- **Environment Variables:** `.[env].(VAR_NAME)`
- **File Inclusion:** `.[path/to/file.sdcl].(key_in_file)`

### 3.6 Front Matter

An SDCL document may optionally include a "front matter" block at the beginning of the file for metadata.

- This block is enclosed by `---` delimiters.
- The content within the delimiters is parsed as standard SDCL.
- Any content after the second `---` delimiter is ignored by the parser.

**Example:**

```sdcl
---
# This is a front matter section
version "1.0"
author "OG-Open-Source"
---

# The main content of the document starts here, but will be ignored if front matter is present.
main_content: {
	data "this part is ignored"
}
```

## Chapter 4: Data Model and Semantics

### 4.1 Data Types

The SDCL data model is composed of the following types:

- **Scalar Types:**

  - `string`: A sequence of Unicode characters.
  - `boolean`: `true` or `false`.
  - `null`: Represents a null or non-existent value.

- **Numeric Types:**

  - `number`: A generic numeric type.

    - **Integer:** A whole number without a fractional part.
    - **Float:** A number with a fractional part.

- **Compound Types:**

  - `object`: An unordered map from string keys to values. Corresponds to a Section in the syntax.
  - `array`: An ordered sequence of values. Corresponds to a List in the syntax.

### 4.2 Type System

SDCL employs a dynamic type system that is inferred at parse time. The type of a value is determined by its literal representation. For example, `123` is inferred as an Integer, while `"123"` is inferred as a `string`. This specification does not define rules for automatic type coercion.

### 4.3 Parsing and Evaluation

The processing of an SDCL document occurs in two phases: parsing and evaluation.

1.  **Parsing:**

    - The parser reads the document and transforms it into an Abstract Syntax Tree (AST) according to the lexical and syntactic rules.
    - References remain unresolved at this stage.

2.  **Evaluation:**

    - After the document structure is built, the implementation resolves references.
    - The semantics for value assignment and merging are as follows:

      - **Key Uniqueness:** It is illegal for the same key to appear more than once within the same section or at the root level of a document. A compliant parser must report an error in this case.
      - **Override on Merge:** When a section's content is merged using a shallow merge reference `(section)`, keys defined in the current section _after_ the reference will override any keys with the same name from the merged section.

## Appendix A: Grammar Summary

This section provides a formal grammar for SDCL using Extended Backus-Naur Form (EBNF).

```ebnf
document            ::= ( front_matter )? ( statement )*
front_matter        ::= '---' NEWLINE ( statement )* '---' NEWLINE

statement           ::= ( key_value_pair | section | list | comment ) NEWLINE

key_value_pair      ::= IDENTIFIER WHITESPACE value

section             ::= IDENTIFIER ':' WHITESPACE? '{' NEWLINE ( INDENT ( statement )* DEDENT )? '}'
list                ::= IDENTIFIER ':' WHITESPACE? '[' ( single_line_list_items | multi_line_list_items )? ']'
single_line_list_items ::= value ( WHITESPACE value )*
multi_line_list_items ::= NEWLINE ( INDENT ( list_element )+ DEDENT )?

list_element        ::= ( value | anonymous_section ) NEWLINE
anonymous_section   ::= '{' NEWLINE ( INDENT ( statement )* DEDENT )? '}'

value               ::= string_literal | number_literal | boolean_literal | 'null' | reference

string_literal      ::= '"' ( [^"]* ) '"'
number_literal      ::= ( '-' )? [0-9]+ ( '.' [0-9]+ )? ( ( 'e' | 'E' ) ( '+' | '-' )? [0-9]+ )?
boolean_literal     ::= 'true' | 'false'

reference           ::= '((' IDENTIFIER ( '.' IDENTIFIER )* '))' | '(' IDENTIFIER ( '.' IDENTIFIER )* ')' | external_reference
external_reference  ::= '.' '[' ( 'env' | FILE_PATH ) ']' '.' ( '((' IDENTIFIER ( '.' IDENTIFIER )* '))' | '(' IDENTIFIER ( '.' IDENTIFIER )* ')' )

comment             ::= '#' ( .* )

IDENTIFIER          ::= ( [a-zA-Z0-9_.-]+ )
FILE_PATH           ::= ( [^\]]+ )

(* Note: INDENT and DEDENT are conceptual tokens representing an increase/decrease in indentation level, which must be based on TAB characters. *)
```
