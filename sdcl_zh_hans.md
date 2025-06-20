# SDCL 语言规格书 (SDCL Language Specification)

|                      |            |
| :------------------- | :--------- |
| **版本 (Version)：** | 1.0        |
| **日期 (Date)：**    | 2025-06-20 |

---

## 摘要 (Abstract)

本文件定义了 SDCL (OGATA's Standard Data Character Storage Language) 的语法、语义与数据模型。SDCL 是一种专为配置文件设计的人类可读数据格式，其核心设计目标是实现卓越的可读性、简洁性与清晰的声明式语法。

SDCL 的主要特性包含：

- **简洁的键值语法：** 以直观的 `key value` 形式作为基础。
- **结构化数据：** 通过嵌套的节 (Sections/Objects) 与列表 (Lists/Arrays) 来组织复杂的数据。
- **明确的数据类型：** 支持字符串、数字、布尔值与空值等基本类型，并强制要求字符串使用双引号，以避免歧义。
- **数据重用机制：** 内建引用 (Referencing) 与插入 (Inclusion) 机制，可引用内部或外部（如环境变量、其他文件）的数据。
- **严格的格式化规则：** 对于多行区块采用基于制表符 (`\t`) 的缩进，确保文件风格的一致性。

本规格书为 SDCL 的实现提供了一致性的标准。

---

## 目录 (Table of Contents)

1.  **第 1 章：导论 (Introduction)**

    - [1.1 范围 (Scope)](#11-范围-scope)
    - [1.2 合规性 (Conformance](#12-合规性-conformance)
    - [1.3 术语与定义 (Terms and Definitions](#13-术语与定义-terms-and-definitions)
    - [1.4 设计哲学 (Design Philosophy](#14-设计哲学-design-philosophy)

2.  **第 2 章：词法结构 (Lexical Structure)**

    - [2.1 字符集 (Character Set](#21-字符集-character-set)
    - [2.2 注释 (Comments](#22-注释-comments)
    - [2.3 关键字 (Keywords](#23-关键字-keywords)
    - [2.4 标识符 (Identifiers](#24-标识符-identifiers)
    - [2.5 字面量 (Literals](#25-字面量-literals)
    - [2.6 运算符与标点符号 (Operators and Punctuators](#26-运算符与标点符号-operators-and-punctuators)

3.  **第 3 章：语法 (Syntax)**

    - [3.1 文件结构 (Document Structure](#31-文件结构-document-structure)
    - [3.2 键值对 (Key-Value Pairs](#32-键值对-key-value-pairs)
    - [3.3 节 (Sections](#33-节-sections)
    - [3.4 列表 (Lists](#34-列表-lists)
    - [3.5 引用 (References](#35-引用-references)
    - [3.6 前导数据 (Front Matter](#36-前导数据-front-matter)

4.  **第 4 章：数据模型与语义 (Data Model and Semantics)**

    - [4.1 数据类型 (Data Types](#41-数据类型-data-types)
    - [4.2 类型系统 (Type System](#42-类型系统-type-system)
    - [4.3 解析与求值 (Parsing and Evaluation](#43-解析与求值-parsing-and-evaluation)

- **附录 A：语法摘要 (Grammar Summary)**

---

## 第 1 章：导论 (Introduction)

### 1.1 范围 (Scope)

本规格定义了 SDCL (OGATA's Standard Data Character Storage Language) 的词法结构、语法、数据模型及语义。本文件旨在为 SDCL 解析器、编辑器及其他相关工具的开发者提供一份权威性的参考。

### 1.2 合规性 (Conformance)

一个合规的 SDCL 实现必须满足以下所有条件：

1.  **解析 (Parsing)：** 必须能成功解析所有符合本规格附录 A 中定义的语法之文件。对于不符合语法的文件，必须回报错误。
2.  **字符编码 (Character Encoding)：** 必须能处理以 UTF-8 编码的文件。
3.  **缩进与换行 (Indentation and Line Endings)：** 必须将制表符 (`\t`) 视为唯一的缩进字符，并将换行符 (`\n`) 视为行分隔符，同时忽略回车符 (`\r`)。
4.  **数据模型 (Data Model)：** 必须将解析后的 SDCL 文件映射到一个与第 4 章所描述的数据模型相符的内存中表示。
5.  **引用解析 (Reference Resolution)：** 必须能正确解析所有内部及外部引用。

### 1.3 术语与定义 (Terms and Definitions)

- **文件 (Document)：** 一个完整的 SDCL 配置，通常对应一个文本文件。
- **键 (Key)：** 一个标识符，在键值对或节的定义中，作为值的名称。
- **值 (Value)：** 与键关联的数据，可以是标量值 (scalar) 或复合结构 (compound structure)。
- **键值对 (Key-Value Pair)：** SDCL 的基本组成单元，由一个键和一个值组成。
- **节 (Section)：** 一个由键标识的嵌套结构，包含一组键值对或其他的节。其作用等同于一个对象 (object) 或字典 (dictionary)。
- **列表 (List)：** 一个值的有序集合。其作用等同于一个数组 (array)。
- **注释 (Comment)：** 一段会被解析器忽略的文本，用于提供说明。
- **引用 (Reference)：** 一种特殊的语法结构，用于插入或合并其他地方定义的值。

### 1.4 设计哲学 (Design Philosophy)

SDCL 的设计基于以下核心原则：

- **人类优先的可读性 (Human-First Readability)：** 语法设计旨在让用户能轻松阅读和撰写，尽可能减少语法噪音。
- **简洁性 (Simplicity)：** 避免引入复杂的逻辑控制结构，专注于数据的声明式描述。
- **无歧义 (Unambiguous)：** 通过严格的规则（如强制字符串引号、固定的缩进字符）来消除常见的解析歧义。
- **结构化 (Structured)：** 提供嵌套的节和列表，以自然的方式表示层级式和集合式数据。

## 第 2 章：词法结构 (Lexical Structure)

### 2.1 字符集 (Character Set)

SDCL 文件 **必须** 使用 UTF-8 字符编码。

### 2.2 注释 (Comments)

注释以井号 (`#`) 开始，并延伸至该行结束。注释 **必须** 独占一行，不允许出现在任何语法元素的后面。

```sdcl
# 这是一个有效的注释。
key "value" # 这是一个无效的注释。
```

### 2.3 关键字 (Keywords)

以下单词被保留为关键字，具有特殊意义，不可作为未加引号的字符串值或标识符。

- `true`：代表布尔值的 "真"。
- `false`：代表布尔值的 "假"。
- `null`：代表空值。

### 2.4 标识符 (Identifiers)

标识符用于命名键 (keys)。其命名规则如下：

1.  由一个或多个非空白字符组成。
2.  **不得** 包含空格 (space characters)。
3.  点号 (`.`) 是标识符的有效部分，本身不表示层级关系。
4.  标识符区分大小写。`myKey` 和 `mykey` 是两个不同的标识符。

**示例 (Example)：**

```sdcl
# 有效的标识符
app.name "My App"
version-1.0 true
_user "admin"
```

### 2.5 字面量 (Literals)

#### 2.5.1 字符串 (Strings)

字符串字面量是由双引号 (`"`) 包围的 UTF-8 字符序列。这是表示文本数据的 **唯一** 方式，旨在将其与数字、布尔值、空值和引用明确区分开来。

**示例 (Example)：**

```sdcl
message "Hello, World!"
path "C:\\Users\\Default"
empty_string ""
```

#### 2.5.2 数字 (Numbers)

数字字面量表示整数 (integers) 或浮点数 (floating-point numbers)。

**示例 (Example)：**

```sdcl
port 5432
version 1.0
negative_integer -10
scientific_notation 6.022e23
```

#### 2.5.3 布尔值 (Booleans)

布尔字面量是 `true` 或 `false`。

**示例 (Example)：**

```sdcl
enabled true
debug_mode false
```

#### 2.5.4 空值 (Null)

空值字面量是 `null`，表示值的缺失。

**示例 (Example)：**

```sdcl
optional_feature null
```

### 2.6 运算符与标点符号 (Operators and Punctuators)

- `:` (冒号)：分隔节或列表的键与其主体 (`{` 或 `[`)。
- `{ }` (花括号)：界定节 (Section) 的开始与结束。
- `[ ]` (方括号)：界定列表 (List) 的开始与结束。
- `( )` (圆括号)：界定引用 (Reference) 的路径。
- `.` (点号)：在引用路径中作为路径分隔符，或作为标识符的一部分。
- `#` (井号)：标识一行的开始为注释。
- `---` (三连横线)：界定前导数据 (Front Matter) 区块。

## 第 3 章：语法 (Syntax)

### 3.1 文件结构 (Document Structure)

一个 SDCL 文件由一系列的语句 (statements) 组成，主要是键值对和节。文件内容的缩进 **必须** 使用制表符 (`\t`)。每个缩进层级代表一个嵌套层级。

### 3.2 键值对 (Key-Value Pairs)

键值对是 SDCL 的基本结构。

- **语法：** `KEY VALUE`
- 键 (KEY) 和值 (VALUE) 之间由一个或多个空格分隔。

**示例 (Example)：**

```sdcl
app.name "My Awesome App"
app.version "1.0.0"
is_enabled true
```

### 3.3 节 (Sections)

节用于将相关的键值对组织成嵌套结构。节的语法要求使用多行区块。

- **语法：**

  ```sdcl
  key: {
  	inner_key1 "value1"
  	inner_key2 123
  }
  ```

- 起始花括号 `{` **必须** 与键和冒号在同一行。
- 结束花括号 `}` **必须** 位于一个新行，且其缩进层级必须与定义该节的键相同。
- 节内的内容 **必须** 缩进一个层级。

### 3.4 列表 (Lists)

列表是值的有序集合。它们可以透过两种方式定义：

- **多行语法 (Multi-line Syntax)：**

  - 起始方括号 `[` 必须与键和冒号在同一行。结束方括号 `]` 必须位于一个新行，且其缩进层级必须与定义该列表的键相同。
  - 列表中的每个元素 **必须** 独占一行，并缩进一个层级。

  ```sdcl
  multi_line_list: [
  	"value1"
  	123
  	true
  ]
  ```

- **单行语法 (Single-line Syntax)：**

  - 所有元素皆放置于方括号内的同一行，并以一个或多个空格分隔。

  ```sdcl
  single_line_list: [1 2 3 "a string" true (a.b.c)]
  ```

列表可以包含不同类型的值，包含匿名的节 (仅限于多行语法)。

### 3.5 引用 (References)

引用机制允许在文件的一个位置插入或合并另一个位置的值。

#### 3.5.1 内部引用 (Internal Reference)

引用同一文件内已定义的值。

- **值引用 (Value Reference)：** `(path.to.value)`

  - 在当前位置插入另一个键的值。

- **浅层合并 (Shallow Merge / Content Merge)：** `(path.to.section)`

  - 当此语法作为一个节的直接子项目时，它会将所引用节的 _内容_ 合并到当前节中。此为浅层合并，仅复制顶层的键值对。

- **节插入 (Section Insertion / Object Insertion)：** `((path.to.section))`

  - 此语法会将所引用的 _整个节_（包含其键名）作为一个新的嵌套节插入到当前内容中。

**示例 (Example)：**

```sdcl
# 基础配置
base: {
	user "guest"
	log_level "info"
}

# 浅层合并示例
# 'base' 的内容被合并到 'config_shallow' 中。
config_shallow: {
	(base)
	log_level "debug" # 覆写来自 base 的 log_level
}
# 结果等同于:
# { user: "guest", log_level: "debug" }

# 节插入示例
# 'base' 节本身被作为一个子节插入到 'config_insertion' 中。
config_insertion: {
	((base))
	another_key "value"
}
# 结果等同于:
# { base: { user: "guest", log_level: "info" }, another_key: "value" }
```

#### 3.5.2 外部引用 (External Reference)

引用来自文件外部来源的值。其语法是内部引用的延伸。

- **语法：** `.[source].(key)` 或 `.[source].((key))`
- **环境变量：** `.[env].(VAR_NAME)`
- **文件包含：** `.[path/to/file.sdcl].(key_in_file)`

### 3.6 前导数据 (Front Matter)

SDCL 文件可以选择性地在文件开头包含一个 "front matter" 区块，用于存放元数据。

- 此区块由 `---` 分隔符包围。
- 分隔符内的内容被解析为标准的 SDCL。
- 第二个 `---` 分隔符之后的任何内容都将被解析器忽略。

**示例 (Example)：**

```sdcl
---
# 这是一个 front matter 区块
version "1.0"
author "OG-Open-Source"
---

# 文件的主要内容从这里开始，但如果存在 front matter，这部分将被忽略。
main_content: {
	data "this part is ignored"
}
```

## 第 4 章：数据模型与语义 (Data Model and Semantics)

### 4.1 数据类型 (Data Types)

SDCL 的数据模型由以下类型构成：

- **标量类型 (Scalar Types)：**

  - `string`：一个 Unicode 字符序列。
  - `boolean`：`true` 或 `false`。
  - `null`：代表空或不存在的值。

- **数值类型 (Numeric Types)：**

  - `number`：一个泛称的数值类型。

    - **整数 (Integer)：** 不含小数部分的整数。
    - **浮点数 (Float)：** 包含小数部分的数字。

- **复合类型 (Compound Types)：**

  - `object`：一个从字符串键到值的无序映射 (map)。对应语法中的节 (Section)。
  - `array`：一个值的有序序列 (sequence)。对应语法中的列表 (List)。

### 4.2 类型系统 (Type System)

SDCL 采用一个在解析时推断的动态类型系统。值的类型由其字面表示法决定。例如，`123` 被推断为整数，而 `"123"` 被推断为 `string`。本规格不定义自动类型转换规则。

### 4.3 解析与求值 (Parsing and Evaluation)

SDCL 文件的处理分为解析与求值两个阶段。

1.  **解析 (Parsing)：**

    - 解析器读取文件，并根据词法和语法规则将其转换为一个抽象语法树 (AST)。
    - 在此阶段，引用 (references) 保持未解析状态。

2.  **求值 (Evaluation)：**

    - 文件结构建立后，实现开始解析引用。
    - 值的赋值和合并语义如下：

      - **键的唯一性 (Key Uniqueness)：** 在同一个节或文件的根层级中，相同的键不允许多次出现。合规的解析器必须在此情况下回报错误。
      - **合并时覆写 (Override on Merge)：** 当使用浅层合并引用 `(section)` 合并一个节的内容时，在引用语句 _之后_ 定义的键，将会覆写被合并节中的同名键。

## 附录 A：语法摘要 (Grammar Summary)

此处使用扩展巴科斯范式 (EBNF) 提供 SDCL 的形式语法。

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

(* 注：INDENT 和 DEDENT 是概念性的 token，代表基于 TAB 字符的缩进层级增加/减少。*)
```
