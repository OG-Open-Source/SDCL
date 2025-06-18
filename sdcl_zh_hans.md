# SDCL (OGATA 的标准数据字符存储语言) 规范

## 1. 范围

本文件定义了 SDCL（OGATA 的标准数据字符存储语言）的语法、语义和核心特性。SDCL 旨在作为一种人类可读且机器可解析的数据存储格式，专为配置管理和数据引用而设计。本规范涵盖了 SDCL 文件的结构、数据类型、引用机制、注释规则以及展开态和压缩态的表示。

本规范旨在供 SDCL 解析器开发者、数据架构师以及需要使用 SDCL 进行数据交换或配置管理的应用程序开发者使用。SDCL 的设计汲取了 JSON 的简洁性、YAML 的可读性以及 TOML 对配置文件的友好性，并引入了强大的内部和外部引用机制，以促进配置的模块化和重用，同时保持严格的语法以简化解析并减少歧义。

本规范不涉及特定编程语言的 SDCL 解析器实现细节，也不涉及网络传输协议。

## 2. 术语和定义

本文件使用了以下术语和定义，并按字母顺序排列：

- **对象（Object）：** 由一组键值对组成的无序集合，用花括号 `{}` 表示。
- **数组（Array）：** 由一系列值组成的有序集合，用方括号 `[]` 表示。
- **标量值（Scalar Value）：** 一个不可再分的单一数据单元，如字符串、数字、布尔值、空值、日期、时间、日期时间、国家代码或 Base64 编码数据。
- **内容包含（Content Inclusion）：** 一种引用机制，用于将指定路径的对象或数组的内容（不含其键名）或整个结构（包含其键名）嵌入到目前范围中。
- **环境变量（Environment Variable）：** 在 SDCL 文件外部，由操作系统或执行环境提供的具名变量。
- **解析器（Parser）：** 一个软件组件，负责读取 SDCL 文件，验证其语法，并将其转换为应用程序可用的内部数据结构。
- **键（Key）：** 用于识别数据项的名称。
- **路径（Path）：** 一个由点（`.`）分隔的键序列，用于在 SDCL 文件中唯一标识一个嵌套的值或结构。
- **引用（Reference）：** 一种机制，允许一个值直接链接到同一文件或外部来源（如环境变量或其他 SDCL 文件）的另一个值或内容。
- **结构化类型（Structured Type）：** 指对象或数组，它们可以包含其他值（包括标量值或其他结构化类型）。
- **值（Value）：** 与键相关联的数据。
- **字面值（Literal）：** 源代码中值的固定表示；例如，`123` 是一个数字字面值，`"hello"` 是一个字符串字面值。
- **展开态（Expanded Form）：** 使用换行符和缩进来清晰地定义对象和数组层次结构的 SDCL 格式。
- **压缩态（Compact Form）：** 在单行上定义对象和数组，使用空格作为元素之间分隔符的 SDCL 格式。

## 3. 符号和缩写

- **ISO 3166-1 alpha-2：** 国际标准化组织定义的国家代码标准。
- **ISO 8601：** 国际标准化组织发布的日期和时间表示法。
- **RFC 4648：** Base64 编码的相关标准。
- **SDCL：** OGATA 的标准数据字符存储语言。
- **UTF-8：** 一种普遍使用的 Unicode 字符编码方案。

## 4. 规范

### 4.1. 数据类型和引用规则

SDCL 定义了一组精确的基本数据类型，每种类型都有明确的引用要求，以确保明确的解析和数据完整性。

- **字符串（Strings）：** 一系列 Unicode 字符。所有字符串值**必须**用双引号（`""`）括起来。不支持多行字符串；所有字符串内容必须位于一个逻辑行上。

  - 示例: `"hello world"`、`"user@example.com"`、`"123 Main St."`

- **数字（Numbers）：** 表示整数和浮点数值。数字**不得**用引号括起来。

  - **整数（Integers）：** 整数（例如，`123`、`-45`）。
  - **浮点数（Floating-Point Numbers）：** 带有小数部分的数字（例如，`3.14`、`-0.001`）。
  - 示例: `123`、`3.14`、`-100`、`0.5`

- **布尔值（Booleans）：** 表示逻辑真值。只识别 `true` 或 `false`，并且它们**不得**用引号括起来。

  - 示例: `true`、`false`

- **空值（Null）：** 表示值的缺失。只识别 `null`，并且它**不得**用引号括起来。

  - 示例: `null`

- **日期（Dates）：** 表示日历日期。值**必须**遵循 `YYYY-MM-DD` 格式，并且**不得**用引号括起来。

  - 示例: `2025-05-27`、`1999-12-31`

- **时间（Times）：** 表示一天中的时间。值**必须**遵循 `HH:MM:SS` 格式（24 小时制），并且**不得**用引号括起来。

  - 示例: `14:30:00`、`09:15:05`

- **日期时间（Datetimes）：** 表示一个特定的时间点，结合了日期和时间。值**必须**遵循 ISO 8601 格式（例如，UTC 为 `YYYY-MM-DDTHH:MM:SSZ` 或带偏移量为 `YYYY-MM-DDTHH:MM:SS+HH:MM`），并且**不得**用引号括起来。

  - 示例: `2025-05-27T14:30:00Z`、`2023-10-27T10:00:00+08:00`

- **国家（Countries）：** 表示 ISO 3166-1 alpha-2 国家代码。值**必须**由恰好两个大写 ASCII 字母组成，并且**不得**用引号括起来。

  - 示例: `TW`、`US`、`JP`

- **Base64：** 表示 Base64 编码的二进制数据。值**必须**遵循 Base64 编码规则，并且**不得**用引号括起来。常用于嵌入小型二进制数据，如图片或密钥。
  - 示例: `SGVsbG8gV29ybGQ=`（Base64 编码的 "Hello World"）

**引用规则摘要：**

- **不带引号的类型：** `int`、`float`、`bool`、`null`、`country`、`date`、`time`、`datetime` 和 `base64` 值**不得**用双引号（`""`）括起来。
- **带引号的类型：** 只有 `string` 值**必须**用双引号（`""`）括起来。

### 4.1.10 字符编码

SDCL 文件默认应使用 UTF-8 编码。所有字符串数据均被视为 Unicode 字符串。

### 4.2. 严格性与语法

SDCL 强制执行一套严格的语法规则，以保证明确的解析和一致的数据表示。遵守这些规则对于有效的 SDCL 文件至关重要。

- **键命名约定：**

  - 键名，特别是点分隔路径中的最后一个部分（例如，`XXX.YYY` 中的 `YYY`），**不得**包含空格或点（`.`）符号。这确保了清晰的层次结构。键名本身**不得**用引号括起来。
  - 完整路径（例如，`XXX.YYY.ZZZ`）用于表示嵌套结构，其中每个点表示更深一层嵌套。
  - 在同一对象范围内，**不得**出现相同的键名。如果解析器在同一对象范围内遇到重复的键名，必须抛出解析错误并终止处理。

- **字符串值分隔：**

  - 所有字面字符串值**必须**用双引号（`""`）括起来。
  - 不支持单引号（`''`）作为字符串分隔符。
  - 不支持多行字符串格式（例如，Python 的三引号 `"""` 或 YAML 的块标量 `|`/`>`）。所有字符串内容必须定义在一个逻辑行上。

- **结构元素：**

  - **对象（表）：** 由用花括号（`{}`）括起来的键值对表示。对象中的每个键值对在展开态中应位于新行，或在压缩态中用至少一个空格分隔。
  - **数组：** 由用方括号（`[]`）括起来的值序列表示。数组中的元素在展开态中应位于新行（若为复杂元素如对象）或用至少一个空格分隔（若为简单标量值），在压缩态中则用至少一个空格分隔。

- **键值指派语法：**

  - SDCL 对于指派标量值与结构化类型（对象和数组）给键时，使用不同的运算符。此区别增强了清晰度并防止了歧义。
  - **标量指派 (`=`):** 指派标量值（例如，字符串、数字、布尔值）时，**必须**使用等号 (`=`)。
    - 示例: `key = "value"`, `version = 1.0`
  - **结构指派 (`:`):** 指派对象 (`{...}`) 或数组 (`[...]`) 给键时，**必须**使用冒号 (`:`)。
    - 示例: `settings: { ... }`, `features: [ ... ]`

- **对象范围定义：**
  - 逻辑上属于一个对象的键值对**必须**明确地用花括号（`{}`）括起来。这清晰地定义了数据的范围和层次结构。

### 4.3. 分隔符与注释

- **分隔符：** SDCL 采用特定的分隔符来划分数据元素和结构：

  - **换行符：** 主要用于在多行上下文（例如，在对象或数组中）中分隔不同的键值对或元素。每个新行表示结构化块中的一个新条目或新元素。
  - **空格：** 用于在单行上下文（特别是在数组中）中分隔元素。多个空格被视为一个分隔符。
  - **逗号（`,`）：** 明确**禁止**作为数组元素或键值对的分隔符。禁止使用逗号是为了简化语法，减少与 JSON 等格式的混淆，并推动更清晰的换行/空格分隔风格。此严格规则增强了清晰度并避免了与其他数据格式的歧义。

- **注释：** SDCL 支持包含注释以用于文件和清晰度：
  - **单行注释：** 以井号（`#`）开头。从 `#` 到行尾的所有文本都被视为注释并被解析器忽略。
  - **位置：** 注释**必须**出现在其自己的专用行上。
  - **禁止行尾注释：** 不允许注释跟随键值对或任何其他数据元素在同一行上（例如，`key = "value" # comment` 是无效的）。此规则确保了一致的格式并简化了解析。

### 4.4. SDCL 独特功能

SDCL 通过几个独特的功能脱颖而出，这些功能旨在增强配置和数据文件中的数据管理、模块化和可重用性。

- **显式类型标签（可选）：**

  - SDCL 提供可选的 `<type></type>` 标签来显式声明值的预期格式或类型。虽然解析器通常可以根据字面表示和引用规则推断类型，显式类型标签在以下情况特别有用：
    1. 增强可读性，明确标示值的预期类型
    2. 在某些情况下，避免潜在的歧义（例如，一个未引用的值可能被误认为特定类型）
    3. 作为一种验证手段，确保值符合预期格式。
  - **支持的标签：** `<int></int>`、`<float></float>`、`<bool></bool>`、`<str></str>`、`<date></date>`、`<time></time>`、`<datetime></datetime>`、`<country></country>`、`<base64></base64>`。
  - **用法：** 对于原始类型（例如，`int`、`float`、`bool`、`null`、`string`、`country`、`base64`），这些标签是可选的。如果存在，解析器将优先使用显式标签，否则推断类型。
  - **示例：**

  ```sdcl
  # 显式标记的日期时间
  startDate = <datetime>2024-05-26T18:30:00Z</datetime>
  # 显式标记的国家
  userCountry = <country>TW</country>
  # 显式标记的 Base64
  imageData = <base64>T0dBVEEncyBTdGFuZGFyZCBEYXRhIENoYXJhY3RlciBTdG9yYWdlIExhbmd1YWdl</base64>
  # 这些在功能上等同于其未标记的、正确引用的对应项（例如，`startDate = 2024-05-26T18:30:00Z`）。
  ```

- **基于路径的引用与内容包含：**

  - SDCL 利用点表示法（`.`）来表示文件中的层次路径，其中 `XXX.YYY = value` 被解析为 `XXX: { YYY = value }`。此表示法是 SDCL 强大引用和包含机制的核心：

  1.  **值引用（在键值对中）：**

      - **语法：** `key = (path.to.value)`
      - **目的：** 此语法允许键的值动态引用位于*同一 SDCL 文件*中指定路径的另一个值。这建立了一个即时链接：如果引用的值发生变更，引用键的值将自动更新。
      - **示例：**

      ```sdcl
      base.config.defaultLogLevel = "DEBUG"
      # service.logging.level 将会是 "DEBUG"
      service.logging.level = (base.config.defaultLogLevel)
      ```

  2.  **内容包含（对象/数组，不带键名）：**

      - **语法：** `(path.to.object_or_array)` （独立于一行，不带指派运算符和值）
      - **目的：** 此机制用于将已定义对象或数组的*内容*（嵌套元素）直接嵌入到目前范围中。它的功能类似于 YAML 的合并键（`<<: *alias`），但用于直接内容注入而不包含来源键名本身。如果引用的路径不存在，或解析到的值不是预期的类型（例如，对于对象包含，路径指向一个标量值），解析器**必须**抛出错误。
      - **示例：**

      ```sdcl
      common_resource_limits: {
        cpu = "500m"
        memory = "512Mi"
      }
      service.resources.limits: {
        # 插入 cpu = "500m" 和 memory = "512Mi"
        (common_resource_limits)
        # 覆盖 common_resource_limits 的 cpu 值
        cpu = "250m"
      }
      # service.resources.limits 的结果为 { cpu = "250m", memory = "512Mi" }
      # 此为 JSON 语法，SDCL 不使用逗号作为分隔符
      ```

  3.  **内容包含（对象/数组，带键名）：**
      - **语法：** `((path.to.object_or_array))`
      - **目的：** 此语法用于包含*整个*对象或数组，包括其键名和所有内容，从指定路径到目前范围。这适用于嵌入完整的、命名的结构化数据。此机制仅支持对象和数组，不支持简单的键值对。
      - **示例：**
      ```sdcl
      users: [
        { id = 1 name = "Alice" }
        { id = 2 name = "Bob" }
      ]
      user_data_container: {
        # 插入完整的 users 数组
        ((users))
      }
      # user_data_container 的结果为 { users: [ { id = 1, name = "Alice" }, { id = 2, name = "Bob" } ] }
      # 此为 JSON 语法，SDCL 不使用逗号作为分隔符
      ```

- **外部引用：**

  - **语法：** `.env.KEY`
  - **目的：** 以点（`.`）开头且未用双引号（`""`）括起来的值被特别解释为对外部环境变量的引用。这允许 SDCL 配置动态地从执行环境中拉取值。解析器应提供机制来读取执行环境中的环境变量。应注意，将敏感信息存储在环境变量中可能存在安全风险，需谨慎评估。
  - **支持的格式：** 目前，只支持 `.env.KEY` 和 `.XXX.sdcl.KEY` 格式，其中 `KEY` 是要引用的环境变量的名称，`XXX` 是要引用的 SDCL 文件名称。
  - **示例：**

  ```sdcl
  # 引用名为 API_KEY 的环境变量
  api.key = .env.API_KEY

  # 从 config.sdcl 中，引用 database.port 对应值
  my.external.port = .config.sdcl.database.port
  ```

  - **`.XXX.sdcl.KEY` 引用：** 此语法用于引用另一个 SDCL 文件（`XXX.sdcl`）中的值。解析器需要定义一个文件解析策略。通常，这可能涉及：
    - 相对于目前 SDCL 文件的路径。
    - 预先定义的搜索路径列表，例如：解析器可以支持配置一个或多个基础目录作为引用外部 SDCL 文件的搜索路径。
      解析器**必须**防止循环引用（例如，A 引用 B，B 引用 A），如果检测到循环引用，则应抛出错误。如果引用的外部文件不存在或无法读取，或者文件内指定的路径不存在，解析器**必须**抛出错误。

- **内容覆盖/合并：**

  - SDCL 中的“最后定义获胜”原则仅适用于通过**内容包含**或**值引用**机制引入的内容。这意味着当一个键的值或内容通过引用被引入到当前作用域时，如果该键在当前作用域中已经存在，则新引入的值将覆盖旧值。
  - **直接定义的重复键名：** 在同一作用域内，直接定义的重复键名是**不允许**的。如果解析器在同一作用域内遇到直接定义的重复键名，必须抛出解析错误并终止处理。
  - 示例：

  ```sdcl
  # 通过内容包含进行覆盖
  default_settings: {
    timeout = 1000
    retries = 3
  }
  service_config: {
    (default_settings)
    retries = 5 # 覆盖 default_settings 中的 retries
  }
  # 结果: service_config: { timeout = 1000, retries = 5 }

  # 通过值引用进行覆盖
  base.url = "http://example.com"
  api.endpoint = (base.url)
  base.url = "http://new-example.com" # 覆盖 base.url，api.endpoint 也会更新
  # 结果: api.endpoint = "http://new-example.com"

  # 无效：直接定义的重复键名
  # invalid_config: {
  #   key1 = "value1"
  #   key1 = "value2" # 错误：重复的键名
  # }
  ```

### 4.5. 展开态与压缩态

SDCL 支持两种主要形式来表示数据结构：展开态（Expanded Form）和压缩态（Compact Form）。这两种形式提供了灵活性，以适应不同的可读性和空间效率需求。

- **展开态（Expanded Form）：**

  - 展开态使用换行符和缩进来清晰地定义对象和数组的层次结构。每个键值对或数组元素通常位于其自己的行上，增强了可读性，尤其适用于复杂或嵌套的数据结构。
  - **缩进：** 在展开态中，缩进主要用于增强人类可读性，并非语法强制要求。解析器应主要依赖花括号 `{}`、方括号 `[]` 以及换行符（或对象/数组内元素间的空格）来确定结构和分隔。然而，推荐使用一致的缩进风格（例如，2 个空格或 4 个空格）以提高可维护性。
  - **示例：**

  ```sdcl
  application: {
    name = "My App"
    version = 1.0
    settings: {
      debug = true
      logLevel = "INFO"
    }
  }
  ```

- **压缩态（Compact Form）：**
  - 压缩态允许在单行上定义对象和数组，使用空格作为元素之间分隔符。这种形式在空间效率方面更为紧凑，适用于简单的数据结构或当需要最小化文件大小时。
  - **示例：**
  ```sdcl
  application:{name="My App" version=1.0 settings:{debug=true logLevel="INFO"}}
  features:["userManagement" "reporting" "notifications"]
  ```
- **不可混用：**
  - 在单一 SDCL 文件中，展开态和压缩态**不得**混用。文件**必须**完全遵循其中一种形式，以确保解析的一致性和可预测性。
  - **选择：** 展开态推荐用于复杂、嵌套的数据结构或需要高度可读性的情况。压缩态适用于简单的数据结构、单行配置或需要最小化文件大小的场景。

### 4.6. Front Matter（元数据块）

SDCL 支持一个可选的“front matter”块，这是一种源于静态网站生成器的熟悉语法，允许在文件开头嵌入一个用于元数据的 SDCL 文档。此功能可让 SDCL 与其他文件格式（如 Markdown）无缝集成。

- **语法与分隔符:**
- SDCL front matter 块**必须**从文件的第一行开始，且**必须**由包含三个连字符 (`---`) 的行作为开头与结尾的分隔。
- 两个 `---` 分隔符之间的内容会被当作标准的 SDCL 文件进行解析。

- **解析规则:**
- 当解析器在第一行遇到 `---` 时，**必须**进入 front-matter 解析模式。
- 解析器将读取并解析其内容，直到遇到结尾的 `---` 行。
- 结尾 `---` 之后的所有内容都被视为文件的主要内容，且**必须**被 SDCL 解析器忽略。这些内容可以由其他工具（例如 Markdown 解析器）处理。
- 如果未找到结尾的 `---`，解析器**应**抛出错误。

- **使用场景:**
- 此功能主要用于将元数据（如标题、作者、日期）与主要文件（如以 Markdown 撰写的博客文章或报告）关联起来。

- **示例:**

  ```markdown
  ---
  # 此区块会被当作 SDCL 解析
  title = "一个包含 Front Matter 的示例"
  author = "OG-Open-Source"
  tags: [ "SDCL" "metadata" "example" ]
  ---

  # 主要文件内容

  这部分的内容在 SDCL front matter 之外，会被分开处理。
  ```

## 5. SDCL 语法示例

本节提供 SDCL 语法的综合示例，涵盖了各种数据类型、结构和引用机制。

```sdcl
# 这是 SDCL 中的单行注释。
# 所有注释必须位于其自己的专用行上。

# SDCL 应用程序示例配置

# 基本键值对 (key = value) 与结构指派 (key: {...} 或 key: [...])
application: {
  name = "My SDCL App"
  version = 1.0
  enabled = true
  debugMode = false
}

# 数字和空值
server: {
  port = 8080
  timeout = 30.5
  maxConnections = 100
  logLevel = "INFO"
  adminEmail = null
}

# 日期、时间、日期时间、国家和 Base64 类型
event: {
  startDate = <datetime>2024-05-26T18:30:00Z</datetime>
  # 不带显式标签的日期时间
  endDate = 2024-05-27T09:00:00Z
  eventDate = <date>2024-05-26</date>
  # 不带显式标签的时间
  eventTime = 18:30:00
  # 不带显式标签的国家
  originCountry = TW
  # 带显式标签的国家
  destinationCountry = <country>US</country>
  # 带显式标签的 Base64
  profileImage = <base64>T0dBVEEncyBTdGFuZGFyZCBEYXRhIENoYXJhY3RlciBTdG9yYWdlIExhbmd1YWdl</base64>
  # 不带显式标签的 Base64
  documentContent = VGhpcyBpcyBhIHRlc3QgZG9jdW1lbnQu
}

# 使用点表示法的嵌套配置
database.type = "PostgreSQL" # 等同 database: { type = "PostgreSQL" }
database.host = "localhost"
database.port = 5432
database.user = "admin"
database.password = "secure_password_123"

# 直接使用 key: { ... } 定义嵌套对象
database.connectionPool: {
  maxSize = 20
  idleTimeout = 60000
}

# 字符串数组
features.enabledFeatures: [
  "userManagement"
  "reporting"
  "notifications"
]

# 对象数组
users: [
  {
    id = 1
    name = "Alice"
    email = "alice@example.com"
  }
  {
    id = 2
    name = "Bob"
    email = "bob@example.com"
  }
]

# 值引用示例（内部引用）
# 首先定义一个基本值
base.config.defaultLogLevel = "DEBUG"

# 现在，引用它
# 此值将为 "DEBUG"；如果 base.config.defaultLogLevel 变更，它将自动更新。
service.logging.level = (base.config.defaultLogLevel)

# 外部引用示例
# 引用名为 API_KEY 的环境变量。
api.key = .env.API_KEY

# 引用外部 SDCL 文件中的键
my.external.port = .config.sdcl.database.port

# 特定服务的配置
service.api: {
  baseUrl = "https://api.example.com/v1"
  apiKey = "your_api_key_here"
  rateLimit: {
    requestsPerMinute = 100
    burst = 10
  }
}

# 内容包含示例（对象，不带键名）
# 此处包含 database.connectionPool 的内容。
db_settings: {
  (database.connectionPool)
  # 覆盖 database.connectionPool 中的 'maxSize' 值。
  maxSize = 25
}

# 内容包含示例（数组，不带键名）
additional_features: [ "adminPanel" "analyticsDashboard" ]
all_features: [
  "userManagement"
  # 包含来自 'additional_features' 数组的元素。
  (additional_features)
]

# 内容包含示例（对象，带键名）
user_data_container: {
  # 包含整个 'users' 数组，包括其键 'users'。
  ((users))
}

# 内容包含示例（对象，带键名）用于嵌套对象
service_limits: {
  # 包含整个 'rateLimit' 对象，包括其键 'rateLimit'。
  ((service.api.rateLimit))
}

# 值引用示例（独立）
# 如果 'database.port' 是 5432，那么 'my.referenced.port' 将为 5432。
my.referenced.port = (database.port)

# 对象内容包含（独立）
direct_db_settings: {
  (database.connectionPool)
  # 覆盖 'idleTimeout' 的值。
  idleTimeout = 70000
}

# 数组内容包含（独立）
direct_user_list: [
  (users)
]

# 展开态和压缩态示例
expanded_example: {
  key1 = "value1"
  nested_object: {
    nested_key = 123
    another_key = true
  }
  array_example: [
    "item1"
    "item2"
  ]
}

compact_example:{name="My App" version=1.0 settings:{debug=true logLevel="INFO"}}
features:["userManagement" "reporting" "notifications"]

# 无效语法示例
# 无效：行尾注释
key = "value" # 这是无效的行尾注释

# 无效：使用冒号指派标量值
# key: "value"

# 无效：使用等号指派结构
# my_object = { key = "value" }

# 无效：数组中使用逗号
# invalid_array: [ "a", "b" ]

# 无效：字符串使用单引号
# invalid_string: 'test'

# 无效：在对象范围内重复的键
# duplicate_key_object: {
#   myKey = "value1"
#   myKey = "value2"
# }

# --- Front Matter 示例 ---
# '---' 分隔线之间的内容是有效的 SDCL。
---
title = "我的文件"
date = 2025-06-10
tags: [ "tech" "specs" "sdcl" ]
---
# 文件的这一部分不会被 SDCL 解析器解析。
```

## 6. 结论

本文件作为 SDCL 的综合设计规范，包含了所有商定的功能、语法规则和示例。它反映了为定义一种健壮且直观的数据存储语言所做的协作努力。SDCL 通过其独特的功能组合，旨在成为一个强大、灵活且易于使用的数据存储语言，特别适用于现代应用程序的配置管理和数据表示需求。
