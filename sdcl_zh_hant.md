# SDCL (OGATA 的標準資料字元儲存語言) 規範

## 1. 範圍

本文件定義了 SDCL（OGATA 的標準資料字元儲存語言）的語法、語義和核心特性。SDCL 旨在作為一種人類可讀且機器可解析的資料儲存格式，專為配置管理和資料引用而設計。本規範涵蓋了 SDCL 文件的結構、資料類型、引用機制、註釋規則以及展開態和壓縮態的表示。

本規範旨在供 SDCL 解析器開發者、資料架構師以及需要使用 SDCL 進行資料交換或配置管理的應用程式開發者使用。SDCL 的設計汲取了 JSON 的簡潔性、YAML 的可讀性以及 TOML 對配置文件的友好性，並引入了強大的內部和外部引用機制，以促進配置的模組化和重用，同時保持嚴格的語法以簡化解析並減少歧義。

本規範不涉及特定程式語言的 SDCL 解析器實現細節，也不涉及網路傳輸協議。

## 2. 術語和定義

本文件使用了以下術語和定義，並按字母順序排列：

- **物件（Object）：** 由一組鍵值對組成的無序集合，用花括號 `{}` 表示。
- **陣列（Array）：** 由一系列值組成的有序集合，用方括號 `[]` 表示。
- **標量值（Scalar Value）：** 一個不可再分的單一資料單元，如字串、數字、布林值、空值、日期、時間、日期時間、國家代碼或 Base64 編碼資料。
- **內容包含（Content Inclusion）：** 一種引用機制，用於將指定路徑的物件或陣列的內容（不含其鍵名）或整個結構（包含其鍵名）嵌入到目前範圍中。
- **環境變數（Environment Variable）：** 在 SDCL 文件外部，由作業系統或執行環境提供的具名變數。
- **解析器（Parser）：** 一個軟體組件，負責讀取 SDCL 文件，驗證其語法，並將其轉換為應用程式可用的內部資料結構。
- **鍵（Key）：** 用於識別資料項的名稱。
- **路徑（Path）：** 一個由點（`.`）分隔的鍵序列，用於在 SDCL 文件中唯一標識一個巢狀的值或結構。
- **引用（Reference）：** 一種機制，允許一個值直接連結到同一檔案中或外部來源（如環境變數或其他 SDCL 檔案）的另一個值或內容。
- **結構化類型（Structured Type）：** 指物件或陣列，它們可以包含其他值（包括標量值或其他結構化類型）。
- **值（Value）：** 與鍵相關聯的資料。
- **字面值（Literal）：** 源代碼中值的固定表示；例如，`123` 是一個數字字面值，`"hello"` 是一個字串字面值。
- **展開態（Expanded Form）：** 使用換行符和縮排來清晰地定義物件和陣列層次結構的 SDCL 格式。
- **壓縮態（Compact Form）：** 在單行上定義物件和陣列，使用空格作為元素之間分隔符的 SDCL 格式。

## 3. 符號和縮寫

- **ISO 3166-1 alpha-2：** 國際標準化組織定義的國家代碼標準。
- **ISO 8601：** 國際標準化組織發布的日期和時間表示法。
- **RFC 4648：** Base64 編碼的相關標準。
- **SDCL：** OGATA 的標準資料字元儲存語言（Standard Data Character Storage Language）。
- **UTF-8：** 一種普遍使用的 Unicode 字元編碼方案。

## 4. 規範

### 4.1. 資料類型和引用規則

SDCL 定義了一組精確的基本資料類型，每種類型都有明確的引用要求，以確保明確的解析和資料完整性。

- **字串（Strings）：** 一系列 Unicode 字元。所有字串值**必須**用雙引號（`""`）括起來。不支援多行字串；所有字串內容必須位於一個邏輯行上。

  - 範例: `"hello world"`、`"user@example.com"`、`"123 Main St."`

- **數字（Numbers）：** 表示整數和浮點數值。數字**不得**用引號括起來。

  - **整數（Integers）：** 整數（例如，`123`、`-45`）。
  - **浮點數（Floating-Point Numbers）：** 帶有小數部分的數字（例如，`3.14`、`-0.001`）。
  - 範例: `123`、`3.14`、`-100`、`0.5`

- **布林值（Booleans）：** 表示邏輯真值。只識別 `true` 或 `false`，並且它們**不得**用引號括起來。

  - 範例: `true`、`false`

- **空值（Null）：** 表示值的缺失。只識別 `null`，並且它**不得**用引號括起來。

  - 範例: `null`

- **日期（Dates）：** 表示日曆日期。值**必須**遵循 `YYYY-MM-DD` 格式，並且**不得**用引號括起來。

  - 範例: `2025-05-27`、`1999-12-31`

- **時間（Times）：** 表示一天中的時間。值**必須**遵循 `HH:MM:SS` 格式（24 小時制），並且**不得**用引號括起來。

  - 範例: `14:30:00`、`09:15:05`

- **日期時間（Datetimes）：** 表示一個特定的時間點，結合了日期和時間。值**必須**遵循 ISO 8601 格式（例如，UTC 為 `YYYY-MM-DDTHH:MM:SSZ` 或帶偏移量為 `YYYY-MM-DDTHH:MM:SS+HH:MM`），並且**不得**用引號括起來。

  - 範例: `2025-05-27T14:30:00Z`、`2023-10-27T10:00:00+08:00`

- **國家（Countries）：** 表示 ISO 3166-1 alpha-2 國家代碼。值**必須**由恰好兩個大寫 ASCII 字母組成，並且**不得**用引號括起來。

  - 範例: `TW`、`US`、`JP`

- **Base64：** 表示 Base64 編碼的二進位資料。值**必須**遵循 Base64 編碼規則，並且**不得**用引號括起來。常用於嵌入小型二進位資料，如圖片或密鑰。
  - 範例: `SGVsbG8gV29ybGQ=`（Base64 編碼的 "Hello World"）

**引用規則摘要：**

- **不帶引號的類型：** `int`、`float`、`bool`、`null`、`country`、`date`、`time`、`datetime` 和 `base64` 值**不得**用雙引號（`""`）括起來。
- **帶引號的類型：** 只有 `string` 值**必須**用雙引號（`""`）括起來。

### 4.1.10 字元編碼

SDCL 文件預設應使用 UTF-8 編碼。所有字串資料均被視為 Unicode 字串。

### 4.2. 嚴格性與語法

SDCL 強制執行一套嚴格的語法規則，以保證明確的解析和一致的資料表示。遵守這些規則對於有效的 SDCL 文件至關重要。

- **鍵命名約定：**

  - 鍵名，特別是點分隔路徑中的最後一個部分（例如，`XXX.YYY` 中的 `YYY`），**不得**包含空格或點（`.`）符號。這確保了清晰的層次結構。鍵名本身**不得**用引號括起來。
  - 完整路徑（例如，`XXX.YYY.ZZZ`）用於表示巢狀結構，其中每個點表示更深一層巢狀。
  - 在同一物件範圍內，**不得**出現相同的鍵名。如果解析器在同一物件範圍內遇到重複的鍵名，必須拋出解析錯誤並終止處理。

- **字串值分隔：**

  - 所有字面字串值**必須**用雙引號（`""`）括起來。
  - 不支援單引號（`''`）作為字串分隔符。
  - 不支援多行字串格式（例如，Python 的三引號 `"""` 或 YAML 的塊標量 `|`/`>`）。所有字串內容必須定義在一個邏輯行上。

- **結構元素：**

  - **物件（表）：** 由用花括號（`{}`）括起來的鍵值對表示。物件中的每個鍵值對在展開態中應位於新行，或在壓縮態中用至少一個空格分隔。
  - **陣列：** 由用方括號（`[]`）括起來的值序列表示。陣列中的元素在展開態中應位於新行（若為複雜元素如物件）或用至少一個空格分隔（若為簡單標量值），在壓縮態中則用至少一個空格分隔。

- **鍵值指派語法：**

  - SDCL 對於指派標量值與結構化類型（物件和陣列）給鍵時，使用不同的運算子。此區別增強了清晰度並防止了歧義。
  - **標量指派 (`=`):** 指派標量值（例如，字串、數字、布林值）時，**必須**使用等號 (`=`)。
    - 範例: `key = "value"`, `version = 1.0`
  - **結構指派 (`:`):** 指派物件 (`{...}`) 或陣列 (`[...]`) 給鍵時，**必須**使用冒號 (`:`)。
    - 範例: `settings: { ... }`, `features: [ ... ]`

- **物件範圍定義：**
  - 邏輯上屬於一個物件的鍵值對**必須**明確地用花括號（`{}`）括起來。這清晰地定義了資料的範圍和層次結構。

### 4.3. 分隔符與註釋

- **分隔符：** SDCL 採用特定的分隔符來劃分資料元素和結構：

  - **換行符：** 主要用於在多行上下文（例如，在物件或陣列中）中分隔不同的鍵值對或元素。每個新行表示結構化塊中的一個新條目或新元素。
  - **空格：** 用於在單行上下文（特別是在陣列中）中分隔元素。多個空格被視為一個分隔符。
  - **逗號（`,`）：** 明確**禁止**作為陣列元素或鍵值對的分隔符。禁止使用逗號是為了簡化語法，減少與 JSON 等格式的混淆，並推動更清晰的換行/空格分隔風格。此嚴格規則增強了清晰度並避免了與其他資料格式的歧義。

- **註釋：** SDCL 支援包含註釋以用於文件和清晰度：
  - **單行註釋：** 以井號（`#`）開頭。從 `#` 到行尾的所有文本都被視為註釋並被解析器忽略。
  - **位置：** 註釋**必須**出現在其自己的專用行上。
  - **禁止行尾註釋：** 不允許註釋跟隨鍵值對或任何其他資料元素在同一行上（例如，`key = "value" # comment` 是無效的）。此規則確保了一致的格式並簡化了解析。

### 4.4. SDCL 獨特功能

SDCL 透過幾個獨特的功能脫穎而出，這些功能旨在增強配置和資料檔案中的資料管理、模組化和可重用性。

- **顯式類型標籤（可選）：**

  - SDCL 提供可選的 `<type></type>` 標籤來顯式宣告值的預期格式或類型。雖然解析器通常可以根據字面表示和引用規則推斷類型，顯式類型標籤在以下情況特別有用：
    1. 增強可讀性，明確標示值的預期類型
    2. 在某些情況下，避免潛在的歧義（例如，一個未引用的值可能被誤認為特定類型）
    3. 作為一種驗證手段，確保值符合預期格式。
  - **支援的標籤：** `<int></int>`、`<float></float>`、`<bool></bool>`、`<str></str>`、`<date></date>`、`<time></time>`、`<datetime></datetime>`、`<country></country>`、`<base64></base64>`。
  - **用法：** 對於原始類型（例如，`int`、`float`、`bool`、`null`、`string`、`country`、`base64`），這些標籤是可選的。如果存在，解析器將優先使用顯式標籤，否則推斷類型。
  - **範例：**

  ```sdcl
  # 顯式標記的日期時間
  startDate = <datetime>2024-05-26T18:30:00Z</datetime>
  # 顯式標記的國家
  userCountry = <country>TW</country>
  # 顯式標記的 Base64
  imageData = <base64>T0dBVEEncyBTdGFuZGFyZCBEYXRhIENoYXJhY3RlciBTdG9yYWdlIExhbmd1YWdl</base64>
  # 這些在功能上等同於其未標記的、正確引用的對應項（例如，`startDate = 2024-05-26T18:30:00Z`）。
  ```

- **基於路徑的引用與內容包含：**

  - SDCL 利用點表示法（`.`）來表示文件中的層次路徑，其中 `XXX.YYY = value` 被解析為 `XXX: { YYY = value }`。此表示法是 SDCL 強大引用和包含機制的核心：

  1.  **值引用（在鍵值對中）：**

      - **語法：** `key = (path.to.value)`
      - **目的：** 此語法允許鍵的值動態引用位於*同一 SDCL 檔案*中指定路徑的另一個值。這建立了一個即時連結：如果引用的值發生變更，引用鍵的值將自動更新。
      - **範例：**

      ```sdcl
      base.config.defaultLogLevel = "DEBUG"
      # service.logging.level 將會是 "DEBUG"
      service.logging.level = (base.config.defaultLogLevel)
      ```

  2.  **內容包含（物件/陣列，不帶鍵名）：**

      - **語法：** `(path.to.object_or_array)` （獨立於一行，不帶指派運算子和值）
      - **目的：** 此機制用於將已定義物件或陣列的*內容*（巢狀元素）直接嵌入到目前範圍中。它的功能類似於 YAML 的合併鍵（`<<: *alias`），但用於直接內容注入而不包含來源鍵名本身。如果引用的路徑不存在，或解析到的值不是預期的類型（例如，對於物件包含，路徑指向一個標量值），解析器必須拋出錯誤。
      - **範例：**

      ```sdcl
      common_resource_limits: {
        cpu = "500m"
        memory = "512Mi"
      }
      service.resources.limits: {
        # 插入 cpu = "500m" 和 memory = "512Mi"
        (common_resource_limits)
        # 覆寫 common_resource_limits 的 cpu 值
        cpu = "250m"
      }
      # service.resources.limits 的結果為 { cpu = "250m", memory = "512Mi" }
      # 此為 JSON 語法，SDCL 不使用逗號作為分隔符
      ```

  3.  **內容包含（物件/陣列，帶鍵名）：**
      - **語法：** `((path.to.object_or_array))`
      - **目的：** 此語法用於包含*整個*物件或陣列，包括其鍵名和所有內容，從指定路徑到目前範圍。這適用於嵌入完整的、命名的資料結構。此機制僅支援物件和陣列，不支援簡單的鍵值對。
      - **範例：**
      ```sdcl
      users: [
        { id = 1 name = "Alice" }
        { id = 2 name = "Bob" }
      ]
      user_data_container: {
        # 插入完整的 users 陣列
        ((users))
      }
      # user_data_container 的結果為 { users: [ { id = 1, name = "Alice" }, { id = 2, name = "Bob" } ] }
      # 此為 JSON 語法，SDCL 不使用逗號作為分隔符
      ```

- **外部引用：**

  - **語法：** `.env.KEY`
  - **目的：** 以點（`.`）開頭且未用雙引號（`""`）括起來的值被特別解釋為對外部環境變數的引用。這允許 SDCL 配置動態地從執行環境中拉取值。解析器應提供機制來讀取執行環境中的環境變數。應注意，將敏感資訊儲存在環境變數中可能存在安全風險，需謹慎評估。
  - **支援的格式：** 目前，只支援 `.env.KEY` 和 `.XXX.sdcl.KEY` 格式，其中 `KEY` 是要引用的環境變數的名稱，`XXX` 是要引用的 SDCL 檔案名稱。
  - **範例：**

  ```sdcl
  # 引用名為 API_KEY 的環境變量
  api.key = .env.API_KEY

  # 從 config.sdcl 中，引用 database.port 對應值
  my.external.port = .config.sdcl.database.port
  ```

  - **`.XXX.sdcl.KEY` 引用：** 此語法用於引用另一個 SDCL 檔案（`XXX.sdcl`）中的值。解析器需要定義一個檔案解析策略。通常，這可能涉及：
    - 相對於目前 SDCL 檔案的路徑。
    - 預先定義的搜尋路徑列表，例如：解析器可以支援配置一個或多個基礎目錄作為引用外部 SDCL 檔案的搜尋路徑。
      解析器必須防止循環引用（例如，A 引用 B，B 引用 A），如果檢測到循環引用，則應拋出錯誤。如果引用的外部檔案不存在或無法讀取，或者檔案內指定的路徑不存在，解析器必須拋出錯誤。

- **內容覆蓋/合併：**

  - SDCL 中的“最後定義獲勝”原則僅適用於透過**內容包含**或**值引用**機制引入的內容。這表示當一個鍵的值或內容透過引用被引入到當前作用域時，如果該鍵在當前作用域中已經存在，則新引入的值將覆蓋舊值。
  - **直接定義的重複鍵名：** 在同一作用域內，直接定義的重複鍵名是**不允許**的。如果解析器在同一作用域內遇到直接定義的重複鍵名，必須拋出解析錯誤並終止處理。
  - 範例：

  ```sdcl
  # 透過內容包含進行覆蓋
  default_settings: {
    timeout = 1000
    retries = 3
  }
  service_config: {
    (default_settings)
    retries = 5 # 覆蓋 default_settings 中的 retries
  }
  # 結果: service_config: { timeout = 1000, retries = 5 }

  # 透過值引用進行覆蓋
  base.url = "http://example.com"
  api.endpoint = (base.url)
  base.url = "http://new-example.com" # 覆蓋 base.url，api.endpoint 也會更新
  # 結果: api.endpoint = "http://new-example.com"

  # 無效：直接定義的重複鍵名
  # invalid_config: {
  #   key1 = "value1"
  #   key1 = "value2" # 錯誤：重複的鍵名
  # }
  ```

### 4.5. 展開態與壓縮態

SDCL 支援兩種主要形式來表示資料結構：展開態（Expanded Form）和壓縮態（Compact Form）。這兩種形式提供了靈活性，以適應不同的可讀性和空間效率需求。

- **展開態（Expanded Form）：**

  - 展開態使用換行符和縮排來清晰地定義物件和陣列的層次結構。每個鍵值對或陣列元素通常位於其自己的行上，增強了可讀性，尤其適用於複雜或巢狀的資料結構。
  - **縮排：** 在展開態中，縮排主要用於增強人類可讀性，並非語法強制要求。解析器應主要依賴花括號 `{}`、方括號 `[]` 以及換行符（或物件/陣列內元素間的空格）來確定結構和分隔。然而，推薦使用一致的縮排風格（例如，2 個空格或 4 個空格）以提高可維護性。
  - **範例：**

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

- **壓縮態（Compact Form）：**
  - 壓縮態允許在單行上定義物件和陣列，使用空格作為元素之間的分隔符。這種形式在空間效率方面更為緊湊，適用於簡單的資料結構或當需要最小化檔案大小時。
  - **範例：**
  ```sdcl
  application:{name="My App" version=1.0 settings:{debug=true logLevel="INFO"}}
  features:["userManagement" "reporting" "notifications"]
  ```
- **不可混用：**
  - 在單一 SDCL 文件中，展開態和壓縮態**不得**混用。文件必須完全遵循其中一種形式，以確保解析的一致性和可預測性。
  - **選擇：** 展開態推薦用於複雜、巢狀的資料結構或需要高度可讀性的情況。壓縮態適用於簡單的資料結構、單行配置或需要最小化檔案大小的場景。

### 4.6. Front Matter（元資料區塊）

SDCL 支援一個可選的「front matter」區塊，這是一種源於靜態網站生成器的熟悉語法，允許在檔案開頭嵌入一個用於元資料的 SDCL 文件。此功能可讓 SDCL 與其他文件格式（如 Markdown）無縫整合。

- **語法與分隔符:**
- SDCL front matter 區塊**必須**從檔案的第一行開始，且**必須**由包含三個連字號 (`---`) 的行作為開頭與結尾的分隔。
- 兩個 `---` 分隔符之間的內容會被當作標準的 SDCL 文件進行解析。

- **解析規則:**
- 當解析器在第一行遇到 `---` 時，**必須**進入 front-matter 解析模式。
- 解析器將讀取並解析其內容，直到遇到結尾的 `---` 行。
- 結尾 `---` 之後的所有內容都被視為檔案的主要內容，且**必須**被 SDCL 解析器忽略。這些內容可以由其他工具（例如 Markdown 解析器）處理。
- 如果未找到結尾的 `---`，解析器**應**拋出錯誤。

- **使用情境:**
- 此功能主要用於將元資料（如標題、作者、日期）與主要文件（如以 Markdown 撰寫的部落格文章或報告）關聯起來。

- **範例:**

```markdown
---
# 此區塊會被當作 SDCL 解析
title = "一個包含 Front Matter 的範例"
author = "OG-Open-Source"
tags: [ "SDCL" "metadata" "example" ]
---

# 主要文件內容

這部分的內容在 SDCL front matter 之外，會被分開處理。
```

## 5. SDCL 語法範例

本節提供 SDCL 語法的綜合範例，涵蓋了各種資料類型、結構和引用機制。

```sdcl
# 這是 SDCL 中的單行註釋。
# 所有註釋必須位於其自己的專用行上。

# SDCL 應用程式範例配置

# 基本鍵值對 (key = value) 與結構指派 (key: {...} 或 key: [...])
application: {
  name = "My SDCL App"
  version = 1.0
  enabled = true
  debugMode = false
}

# 數字和空值
server: {
  port = 8080
  timeout = 30.5
  maxConnections = 100
  logLevel = "INFO"
  adminEmail = null
}

# 日期、時間、日期時間、國家和 Base64 類型
event: {
  startDate = <datetime>2024-05-26T18:30:00Z</datetime>
  # 不帶顯式標籤的日期時間
  endDate = 2024-05-27T09:00:00Z
  eventDate = <date>2024-05-26</date>
  # 不帶顯式標籤的時間
  eventTime = 18:30:00
  # 不帶顯式標籤的國家
  originCountry = TW
  # 帶顯式標籤的國家
  destinationCountry = <country>US</country>
  # 帶顯式標籤的 Base64
  profileImage = <base64>T0dBVEEncyBTdGFuZGFyZCBEYXRhIENoYXJhY3RlciBTdG9yYWdlIExhbmd1YWdl</base64>
  # 不帶顯式標籤的 Base64
  documentContent = VGhpcyBpcyBhIHRlc3QgZG9jdW1lbnQu
}

# 使用點表示法的巢狀配置
database.type = "PostgreSQL" # 等同 database: { type = "PostgreSQL" }
database.host = "localhost"
database.port = 5432
database.user = "admin"
database.password = "secure_password_123"

# 直接使用 key: { ... } 定義巢狀物件
database.connectionPool: {
  maxSize = 20
  idleTimeout = 60000
}

# 字串陣列
features.enabledFeatures: [
  "userManagement"
  "reporting"
  "notifications"
]

# 物件陣列
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

# 值引用範例（內部引用）
# 首先定義一個基本值
base.config.defaultLogLevel = "DEBUG"

# 現在，引用它
# 此值將為 "DEBUG"，如果 base.config.defaultLogLevel 變更，它將自動更新。
service.logging.level = (base.config.defaultLogLevel)

# 外部引用範例
# 引用名為 API_KEY 的環境變數。
api.key = .env.API_KEY

# 引用外部 SDCL 檔案中的鍵
my.external.port = .config.sdcl.database.port

# 特定服務的配置
service.api: {
  baseUrl = "https://api.example.com/v1"
  apiKey = "your_api_key_here"
  rateLimit: {
    requestsPerMinute = 100
    burst = 10
  }
}

# 內容包含範例（物件，不帶鍵名）
# 此處包含 database.connectionPool 的內容。
db_settings: {
  (database.connectionPool)
  # 覆蓋 database.connectionPool 中的 'maxSize' 值。
  maxSize = 25
}

# 內容包含範例（陣列，不帶鍵名）
additional_features: [ "adminPanel" "analyticsDashboard" ]
all_features: [
  "userManagement"
  # 包含來自 'additional_features' 陣列的元素。
  (additional_features)
]

# 內容包含範例（物件，帶鍵名）
user_data_container: {
  # 包含整個 'users' 陣列，包括其鍵 'users'。
  ((users))
}

# 內容包含範例（物件，帶鍵名）用於巢狀物件
service_limits: {
  # 包含整個 'rateLimit' 物件，包括其鍵 'rateLimit'。
  ((service.api.rateLimit))
}

# 值引用範例（獨立）
# 如果 'database.port' 是 5432，那麼 'my.referenced.port' 將為 5432。
my.referenced.port = (database.port)

# 物件內容包含（獨立）
direct_db_settings: {
  (database.connectionPool)
  # 覆蓋 'idleTimeout' 的值。
  idleTimeout = 70000
}

# 陣列內容包含（獨立）
direct_user_list: [
  (users)
]

# 展開態和壓縮態範例
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

# 無效語法範例
# 無效：行尾註釋
key = "value" # 這是無效的行尾註釋

# 無效：使用冒號指派標量值
# key: "value"

# 無效：使用等號指派結構
# my_object = { key = "value" }

# 無效：陣列中使用逗號
# invalid_array: [ "a", "b" ]

# 無效：字串使用單引號
# invalid_string: 'test'

# 無效：在物件範圍內重複的鍵
# duplicate_key_object: {
#   myKey = "value1"
#   myKey = "value2"
# }

# --- Front Matter 範例 ---
# '---' 分隔線之間的內容是有效的 SDCL。
---
title = "我的文件"
date = 2025-06-10
tags: [ "tech" "specs" "sdcl" ]
---
# 文件的這一部分不會被 SDCL 解析器解析。
```

## 6. 結論

本文件作為 SDCL 的綜合設計規範，包含了所有商定的功能、語法規則和範例。它反映了為定義一種健壯且直觀的資料儲存語言所做的協作努力。SDCL 透過其獨特的特性組合，旨在成為一個強大、靈活且易於使用的資料儲存語言，特別適用於現代應用程式的配置管理和資料表示需求。
