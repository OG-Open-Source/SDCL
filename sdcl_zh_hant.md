# SDCL (OGATA 的標準資料字元儲存語言) 規範

## 1. 範圍

本文件定義了 SDCL（OGATA 的標準資料字元儲存語言）的語法與核心特性。SDCL 是一種為設定檔的簡潔性與清晰度而設計的人類可讀資料格式。

## 2. 文件結構

- **字元編碼：** SDCL 檔案**必須**使用 UTF-8 編碼。
- **換行符：** 換行符（`\n`，LF）是標準的行分隔符。解析器**必須**忽略回車符（`\r`，CR）。
- **縮排：** 縮排**必須**使用定位字元（`\t`）。**不允許**使用空格進行縮排。

## 3. 核心元件

### 3.1. 鍵值對

- SDCL 的基本結構是鍵值對。
- **語法：** `KEY VALUE`
- 鍵與其值由一個或多個空格分隔。
- **鍵（Key）：**
  - **不得**包含空格。
  - 點號（`.`）是鍵的有效部分，不代表階層。
  - 鍵區分大小寫。
- **值（Value）：**
  - 值從鍵之後的第一個非空格字元開始，按字面解析直到行尾。
  - 一個值被視為 `string` 且**必須**用雙引號 (`""`) 包裹，除非它是 `number`、`boolean`、`null` 或引用。
  - 例如，`"My Awesome App"` 是一個單一的字串值。
  - 為了避免與其他資料類型混淆，字串值強制使用雙引號。

### 3.2. 資料類型

SDCL 為清晰與結構化定義了幾種資料類型：

- **`string` 字串：** 一串由雙引號 (`""`) 包裹的字元序列。所有文字資料都必須使用，以區別於 `number`、`boolean`、`null` 或引用等其他類型。
- **`number` 數字：** 表示整數與浮點數。
- **`boolean` 布林值：** `true` 或 `false`。
- **`null` 空值：** 表示值的缺失。
- **`object` 物件：** 一組鍵值對的集合，由花括號（`{}`）包圍。
- **`array` 陣列：** 一個值的有序列表，由方括號（`[]`）包圍。
- **`insertion` 插入：** 一種包含或合併資料的機制，由括號（`()`）表示。

### 3.3. 註釋

- 註釋以井號（`#`）開頭，並延伸至行尾。
- 註釋**必須**位於其自己的專用行上。

## 4. 進階功能

### 4.1. 物件（字典/映射）

- 物件用於將相關的鍵值對分組。
- **語法：**
  ```sdcl
  key: {
  	inner_key1 "value1"
  	inner_key2 "value2"
  }
  ```
- 起始花括號 `{` 必須與鍵和冒號在同一行，結束花括號 `}` 必須在新的一行。

### 4.2. 陣列（列表）

- 陣列是值的有序集合。
- **語法：**
  ```sdcl
  key: [
  	"value1"
  	"一個字串值"
  	123
  	true
  ]
  ```
- 起始方括號 `[` 必須與鍵和冒號在同一行，結束方括號 `]` 必須在新的一行。

### 4.3. 引用與包含

SDCL 支援引用以重用資料。

- **內部引用（插入）：**

  - **語法：** `(path.to.value)`
  - **目的：** 插入另一個鍵的值。若在物件的頂層使用，則會合併被引用物件的內容。

- **外部引用：**
  - **語法：** `.[source].(key)`
  - **環境變數：** `.[env].(VAR_NAME)`
  - **檔案包含：** `.[path/to/file.sdcl].(key_in_file)`

## 5. Front Matter

- SDCL 支援一個可選的「front matter」區塊，用於在檔案開頭嵌入元資料，由 `---` 包圍。
- `---` 內的內容會被當作 SDCL 解析，之後的內容則被忽略。

## 6. 語法範例

```sdcl
# 基本鍵值對
app.name "My Awesome App"
app.version "1.0.0"
is_enabled true

# 物件定義
database: {
	host "localhost"
	port 5432
	user "admin"
	password .[env].(DB_PASS) # 外部引用
}

# 陣列定義
features: [
	"User Authentication"
	"Data Processing"
	"reporting"
]

# 內容包含
base_settings: {
	timeout 30
	retries 3
}

production_settings: {
	(base_settings) # 在此合併 base_settings
	timeout 60      # 覆蓋 timeout 的值
}

# 引用一個值
admin_user (database.user)
```
