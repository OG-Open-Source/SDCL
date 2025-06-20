# SDCL (OGATA 的标准数据字符存储语言) 规范

## 1. 范围

本文件定义了 SDCL（OGATA 的标准数据字符存储语言）的语法与核心特性。SDCL 是一种为配置文件的简洁性与清晰度而设计的人类可读数据格式。

## 2. 文件结构

- **字符编码：** SDCL 文件**必须**使用 UTF-8 编码。
- **换行符：** 换行符（`\n`，LF）是标准的行分隔符。解析器**必须**忽略回车符（`\r`，CR）。
- **缩进：** 缩进**必须**使用制表符（`\t`）。**不允许**使用空格进行缩进。

## 3. 核心组件

### 3.1. 键值对

- SDCL 的基本结构是键值对。
- **语法：** `KEY VALUE`
- 键与其值由一个或多个空格分隔。
- **键（Key）：**
  - **不得**包含空格。
  - 点号（`.`）是键的有效部分，不代表层级。
  - 键区分大小写。
- **值（Value）：**
  - 值从键之后的第一个非空格字符开始，按字面解析直到行尾。
  - 一个值被视为 `string` 且**必须**用双引号 (`""`) 包裹，除非它是 `number`、`boolean`、`null` 或引用。
  - 例如，`"My Awesome App"` 是一个单一的字符串值。
  - 为了避免与其他数据类型混淆，字符串值强制使用双引号。

### 3.2. 数据类型

SDCL 为清晰与结构化定义了以下几种数据类型：

- **`string` 字符串：** 一串由双引号 (`""`) 包裹的字符序列。所有文字数据都必须使用，以区别于 `number`、`boolean`、`null` 或引用等其他类型。
- **`number` 数字：** 表示整数与浮点数。
- **`boolean` 布尔值：** `true` 或 `false`。
- **`null` 空值：** 表示值的缺失。
- **`object` 对象：** 一组键值对的集合，由花括号（`{}`）包围。
- **`array` 数组：** 一个值的有序列表，由方括号（`[]`）包围。
- **`insertion` 插入：** 一种包含或合并数据的机制，由括号（`()`）表示。

### 3.3. 注释

- 注释以井号（`#`）开头，并延伸至行尾。
- 注释**必须**位于其自己的专用行上。

## 4. 高级功能

### 4.1. 对象（字典/映射）

- 对象用于将相关的键值对分组。
- **语法：**
  ```sdcl
  key: {
  	inner_key1 "value1"
  	inner_key2 "value2"
  }
  ```
- 起始花括号 `{` 必须与键和冒号在同一行，结束花括号 `}` 必须在新的一行。

### 4.2. 数组（列表）

- 数组是值的有序集合。
- **语法：**
  ```sdcl
  key: [
  	"value1"
  	"一个字符串值"
  	123
  	true
  ]
  ```
- 起始方括号 `[` 必须与键和冒号在同一行，结束方括号 `]` 必须在新的一行。

### 4.3. 引用与包含

SDCL 支持引用以重用数据。

- **内部引用（插入）：**

  - **语法：** `(path.to.value)`
  - **目的：** 插入另一个键的值。若在对象的顶层使用，则会合并被引用对象的内容。

- **外部引用：**
  - **语法：** `.[source].(key)`
  - **环境变量：** `.[env].(VAR_NAME)`
  - **文件包含：** `.[path/to/file.sdcl].(key_in_file)`

## 5. Front Matter

- SDCL 支持一个可选的“front matter”块，用于在文件开头嵌入元数据，由 `---` 包围。
- `---` 内的内容会被当作 SDCL 解析，之后的内容则被忽略。

## 6. 语法示例

```sdcl
# 基本键值对
app.name "My Awesome App"
app.version "1.0.0"
is_enabled true

# 对象定义
database: {
	host "localhost"
	port 5432
	user "admin"
	password .[env].(DB_PASS) # 外部引用
}

# 数组定义
features: [
	"User Authentication"
	"Data Processing"
	"reporting"
]

# 内容包含
base_settings: {
	timeout 30
	retries 3
}

production_settings: {
	(base_settings) # 在此合并 base_settings
	timeout 60      # 覆盖 timeout 的值
}

# 引用一个值
admin_user (database.user)
```
