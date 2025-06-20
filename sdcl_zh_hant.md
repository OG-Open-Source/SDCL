# SDCL 語言規格書 (SDCL Language Specification)

|                      |            |
| :------------------- | :--------- |
| **版本 (Version)：** | 1.0        |
| **日期 (Date)：**    | 2025-06-20 |

---

## 摘要 (Abstract)

本文件定義了 SDCL (OGATA's Standard Data Character Storage Language) 的語法、語義與資料模型。SDCL 是一種專為組態檔設計的人類可讀資料格式，其核心設計目標是實現卓越的可讀性、簡潔性與清晰的聲明式語法。

SDCL 的主要特性包含：

- **簡潔的鍵值語法：** 以直觀的 `key value` 形式作為基礎。
- **結構化資料：** 透過巢狀的節 (Sections/Objects) 與列表 (Lists/Arrays) 來組織複雜的資料。
- **明確的資料型別：** 支援字串、數字、布林值與空值等基本型別，並強制要求字串使用雙引號，以避免歧義。
- **資料重用機制：** 內建引用 (Referencing) 與插入 (Inclusion) 機制，可引用內部或外部（如環境變數、其他檔案）的資料。
- **嚴格的格式化規則：** 對於多行區塊採用基於定位字元 (`\t`) 的縮排，確保文件風格的一致性。

本規格書為 SDCL 的實作提供了一致性的標準。

---

## 目錄 (Table of Contents)

1.  **第 1 章：導論 (Introduction)**

    - [1.1 範圍 (Scope)](#11-範圍-scope)
    - [1.2 合規性 (Conformance)](#12-合規性-conformance)
    - [1.3 術語與定義 (Terms and Definitions)](#13-術語與定義-terms-and-definitions)
    - [1.4 設計哲學 (Design Philosophy)](#14-設計哲學-design-philosophy)

2.  **第 2 章：詞法結構 (Lexical Structure)**

    - [2.1 字元集 (Character Set)](#21-字元集-character-set)
    - [2.2 註解 (Comments)](#22-註解-comments)
    - [2.3 關鍵字 (Keywords)](#23-關鍵字-keywords)
    - [2.4 識別碼 (Identifiers)](#24-識別碼-identifiers)
    - [2.5 字面量 (Literals)](#25-字面量-literals)
    - [2.6 運算子與標點符號 (Operators and Punctuators)](#26-運算子與標點符號-operators-and-punctuators)

3.  **第 3 章：語法 (Syntax)**

    - [3.1 文件結構 (Document Structure)](#31-文件結構-document-structure)
    - [3.2 鍵值對 (Key-Value Pairs)](#32-鍵值對-key-value-pairs)
    - [3.3 節 (Sections)](#33-節-sections)
    - [3.4 列表 (Lists)](#34-列表-lists)
    - [3.5 引用 (References)](#35-引用-references)
    - [3.6 前導資料 (Front Matter)](#36-前導資料-front-matter)

4.  **第 4 章：資料模型與語義 (Data Model and Semantics)**

    - [4.1 資料型別 (Data Types)](#41-資料型別-data-types)
    - [4.2 型別系統 (Type System)](#42-型別系統-type-system)
    - [4.3 解析與求值 (Parsing and Evaluation)](#43-解析與求值-parsing-and-evaluation)

- **附錄 A：文法摘要 (Grammar Summary)**

---

## 第 1 章：導論 (Introduction)

### 1.1 範圍 (Scope)

本規格定義了 SDCL (OGATA's Standard Data Character Storage Language) 的詞法結構、語法、資料模型及語義。本文件旨在為 SDCL 解析器、編輯器及其他相關工具的開發者提供一份權威性的參考。

### 1.2 合規性 (Conformance)

一個合規的 SDCL 實作必須滿足以下所有條件：

1.  **解析 (Parsing)：** 必須能成功解析所有符合本規格附錄 A 中定義的文法之文件。對於不符合文法的文件，必須回報錯誤。
2.  **字元編碼 (Character Encoding)：** 必須能處理以 UTF-8 編碼的文件。
3.  **縮排與換行 (Indentation and Line Endings)：** 必須將定位字元 (`\t`) 視為唯一的縮排字元，並將換行符 (`\n`) 視為行分隔符，同時忽略回車符 (`\r`)。
4.  **資料模型 (Data Model)：** 必須將解析後的 SDCL 文件映射到一個與第 4 章所描述的資料模型相符的記憶體中表示。
5.  **引用解析 (Reference Resolution)：** 必須能正確解析所有內部及外部引用。

### 1.3 術語與定義 (Terms and Definitions)

- **文件 (Document)：** 一個完整的 SDCL 組態，通常對應一個文字檔案。
- **鍵 (Key)：** 一個識別碼，在鍵值對或節的定義中，作為值的名稱。
- **值 (Value)：** 與鍵相關聯的資料，可以是純量值 (scalar) 或複合結構 (compound structure)。
- **鍵值對 (Key-Value Pair)：** SDCL 的基本組成單元，由一個鍵和一個值組成。
- **節 (Section)：** 一個由鍵標識的巢狀結構，包含一組鍵值對或其他的節。其作用等同於一個物件 (object) 或字典 (dictionary)。
- **列表 (List)：** 一個值的有序集合。其作用等同於一個陣列 (array)。
- **註解 (Comment)：** 一段會被解析器忽略的文字，用於提供說明。
- **引用 (Reference)：** 一種特殊的語法結構，用於插入或合併其他地方定義的值。

### 1.4 設計哲學 (Design Philosophy)

SDCL 的設計基於以下核心原則：

- **人類優先的可讀性 (Human-First Readability)：** 語法設計旨在讓人類能輕鬆閱讀和撰寫，盡可能減少語法噪音。
- **簡潔性 (Simplicity)：** 避免引入複雜的邏輯控制結構，專注於資料的聲明式描述。
- **無歧義 (Unambiguous)：** 透過嚴格的規則（如強制字串引號、固定的縮排字元）來消除常見的解析歧義。
- **結構化 (Structured)：** 提供巢狀的節和列表，以自然的方式表示階層式和集合式資料。

## 第 2 章：詞法結構 (Lexical Structure)

### 2.1 字元集 (Character Set)

SDCL 文件 **必須** 使用 UTF-8 字元編碼。

### 2.2 註解 (Comments)

註解以井號 (`#`) 開始，並延伸至該行結束。註解 **必須** 獨佔一行，不允許出現在任何語法元素的後面。

```sdcl
# 這是一個有效的註解。
key "value" # 這是一個無效的註解。
```

### 2.3 關鍵字 (Keywords)

以下單詞被保留為關鍵字，具有特殊意義，不可作為未加引號的字串值或識別碼。

- `true`：代表布林值的 "真"。
- `false`：代表布林值的 "假"。
- `null`：代表空值。

### 2.4 識別碼 (Identifiers)

識別碼用於命名鍵 (keys)。其命名規則如下：

1.  由一個或多個非空白字元組成。
2.  **不得** 包含空格 (space characters)。
3.  點號 (`.`) 是識別碼的有效部分，本身不表示階層關係。
4.  識別碼區分大小寫。`myKey` 和 `mykey` 是兩個不同的識別碼。

**範例 (Example)：**

```sdcl
# 有效的識別碼
app.name "My App"
version-1.0 true
_user "admin"
```

### 2.5 字面量 (Literals)

#### 2.5.1 字串 (Strings)

字串字面量是由雙引號 (`"`) 包圍的 UTF-8 字元序列。這是表示文字資料的 **唯一** 方式，旨在將其與數字、布林值、空值和引用明確區分開來。

**範例 (Example)：**

```sdcl
message "Hello, World!"
path "C:\\Users\\Default"
empty_string ""
```

#### 2.5.2 數字 (Numbers)

數字字面量表示整數 (integers) 或浮點數 (floating-point numbers)。

**範例 (Example)：**

```sdcl
port 5432
version 1.0
negative_integer -10
scientific_notation 6.022e23
```

#### 2.5.3 布林值 (Booleans)

布林字面量是 `true` 或 `false`。

**範例 (Example)：**

```sdcl
enabled true
debug_mode false
```

#### 2.5.4 空值 (Null)

空值字面量是 `null`，表示值的缺失。

**範例 (Example)：**

```sdcl
optional_feature null
```

### 2.6 運算子與標點符號 (Operators and Punctuators)

- `:` (冒號)：分隔節或列表的鍵與其主體 (`{` 或 `[`)。
- `{ }` (花括號)：界定節 (Section) 的開始與結束。
- `[ ]` (方括號)：界定列表 (List) 的開始與結束。
- `( )` (圓括號)：界定引用 (Reference) 的路徑。
- `.` (點號)：在引用路徑中作為路徑分隔符，或作為識別碼的一部分。
- `#` (井號)：標識一行的開始為註解。
- `---` (三連橫線)：界定前導資料 (Front Matter) 區塊。

## 第 3 章：語法 (Syntax)

### 3.1 文件結構 (Document Structure)

一個 SDCL 文件由一系列的語句 (statements) 組成，主要是鍵值對和節。文件內容的縮排 **必須** 使用定位字元 (`\t`)。每個縮排層級代表一個巢狀層級。

### 3.2 鍵值對 (Key-Value Pairs)

鍵值對是 SDCL 的基本結構。

- **語法：** `KEY VALUE`
- 鍵 (KEY) 和值 (VALUE) 之間由一個或多個空格分隔。

**範例 (Example)：**

```sdcl
app.name "My Awesome App"
app.version "1.0.0"
is_enabled true
```

### 3.3 節 (Sections)

節用於將相關的鍵值對組織成巢狀結構。節的語法要求使用多行區塊。

- **語法：**

  ```sdcl
  key: {
  	inner_key1 "value1"
  	inner_key2 123
  }
  ```

- 起始花括號 `{` **必須** 與鍵和冒號在同一行。
- 結束花括號 `}` **必須** 位於一個新行，且其縮排層級必須與定義該節的鍵相同。
- 節內的內容 **必須** 縮排一個層級。

### 3.4 列表 (Lists)

列表是值的有序集合。它們可以透過兩種方式定義：

- **多行語法 (Multi-line Syntax)：**

  - 起始方括號 `[` 必須與鍵和冒號在同一行。結束方括號 `]` 必須位於一個新行，且其縮排層級必須與定義該列表的鍵相同。
  - 列表中的每個元素 **必須** 獨佔一行，並縮排一個層級。

  ```sdcl
  multi_line_list: [
  	"value1"
  	123
  	true
  ]
  ```

- **單行語法 (Single-line Syntax)：**

  - 所有元素皆放置於方括號內的同一行，並以一個或多個空格分隔。

  ```sdcl
  single_line_list: [1 2 3 "a string" true (a.b.c)]
  ```

列表可以包含不同型別的值，包含匿名的節 (僅限於多行語法)。

### 3.5 引用 (References)

引用機制允許在文件的一個位置插入或合併另一個位置的值。

#### 3.5.1 內部引用 (Internal Reference)

引用同一文件內已定義的值。

- **值引用 (Value Reference)：** `(path.to.value)`

  - 在當前位置插入另一個鍵的值。

- **淺層合併 (Shallow Merge / Content Merge)：** `(path.to.section)`

  - 當此語法作為一個節的直接子項目時，它會將被引用節的 _內容_ 合併到當前節中。此為淺層合併，僅複製頂層的鍵值對。

- **節插入 (Section Insertion / Object Insertion)：** `((path.to.section))`

  - 此語法會將被引用的 _整個節_（包含其鍵名）作為一個新的巢狀節插入到當前內容中。

**範例 (Example)：**

```sdcl
# 基礎組態
base: {
	user "guest"
	log_level "info"
}

# 淺層合併範例
# 'base' 的內容被合併到 'config_shallow' 中。
config_shallow: {
	(base)
	log_level "debug" # 覆寫來自 base 的 log_level
}
# 結果等同於:
# { user: "guest", log_level: "debug" }

# 節插入範例
# 'base' 節本身被作為一個子節插入到 'config_insertion' 中。
config_insertion: {
	((base))
	another_key "value"
}
# 結果等同於:
# { base: { user: "guest", log_level: "info" }, another_key: "value" }
```

#### 3.5.2 外部引用 (External Reference)

引用來自文件外部來源的值。其語法是內部引用的延伸。

- **語法：** `.[source].(key)` 或 `.[source].((key))`
- **環境變數：** `.[env].(VAR_NAME)`
- **檔案包含：** `.[path/to/file.sdcl].(key_in_file)`

### 3.6 前導資料 (Front Matter)

SDCL 文件可以選擇性地在文件開頭包含一個 "front matter" 區塊，用於存放元資料。

- 此區塊由 `---` 分隔符包圍。
- 分隔符內的內容被解析為標準的 SDCL。
- 第二個 `---` 分隔符之後的任何內容都將被解析器忽略。

**範例 (Example)：**

```sdcl
---
# 這是一個 front matter 區塊
version "1.0"
author "OG-Open-Source"
---

# 文件的主要內容從這裡開始，但如果存在 front matter，這部分將被忽略。
main_content: {
	data "this part is ignored"
}
```

## 第 4 章：資料模型與語義 (Data Model and Semantics)

### 4.1 資料型別 (Data Types)

SDCL 的資料模型由以下型別構成：

- **純量型別 (Scalar Types)：**

  - `string`：一個 Unicode 字元序列。
  - `boolean`：`true` 或 `false`。
  - `null`：代表空或不存在的值。

- **數值型別 (Numeric Types)：**

  - `number`：一個泛稱的數值型別。

    - **整數 (Integer)：** 不含小數部分的整數。
    - **浮點數 (Float)：** 包含小數部分的數字。

- **複合型別 (Compound Types)：**

  - `object`：一個從字串鍵到值的無序映射 (map)。對應語法中的節 (Section)。
  - `array`：一個值的有序序列 (sequence)。對應語法中的列表 (List)。

### 4.2 型別系統 (Type System)

SDCL 採用一個在解析時推斷的動態型別系統。值的型別由其字面表示法決定。例如，`123` 被推斷為整數，而 `"123"` 被推斷為 `string`。本規格不定義自動型別轉換規則。

### 4.3 解析與求值 (Parsing and Evaluation)

SDCL 文件的處理分為解析與求值兩個階段。

1.  **解析 (Parsing)：**

    - 解析器讀取文件，並根據詞法和語法規則將其轉換為一個抽象語法樹 (AST)。
    - 在此階段，引用 (references) 保持未解析狀態。

2.  **求值 (Evaluation)：**

    - 文件結構建立後，實作開始解析引用。
    - 值的賦值和合併語義如下：

      - **鍵的唯一性 (Key Uniqueness)：** 在同一個節或文件的根層級中，相同的鍵不允許多次出現。合規的解析器必須在此情況下回報錯誤。
      - **合併時覆寫 (Override on Merge)：** 當使用淺層合併引用 `(section)` 合併一個節的內容時，在引用語句 _之後_ 定義的鍵，將會覆寫被合併節中的同名鍵。

## 附錄 A：文法摘要 (Grammar Summary)

此處使用擴充巴科斯範式 (EBNF) 提供 SDCL 的形式文法。

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

(* 註：INDENT 和 DEDENT 是概念性的 token，代表基於 TAB 字元的縮排層級增加/減少。*)
```
