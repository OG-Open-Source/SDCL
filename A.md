## **SDCL 解釋器開發規劃 (Golang)**

本規劃將詳細闡述使用 Golang 從零開始構建 SDCL (OGATA's Standard Data Character Storage Language) 解釋器的步驟，並遵循 Go 語言的慣例及小駝峰命名法。

### **1. SDCL 規範深度解析**

在開始任何程式碼編寫之前，必須**徹底理解 SDCL 規範**。這不僅僅是閱讀，更是要內化其所有規則和隱含的行為。這包括：

- **資料型別字面量 (Literals)**：SDCL 定義了多種基本資料型別，包括 string, number, boolean, null, date, time, datetime, country, base64。
  - **字串**：必須使用雙引號 "" 包裹，且不支援多行字串。轉義字符（如 ", n, uXXXX）的處理是關鍵。
  - **數字**：不帶引號，可以是整數或浮點數。
  - **布林值 (Boolean)**：true 或 false，不帶引號。
  - **空值 (Null)**：null，不帶引號。
  - **日期、時間、日期時間**：嚴格遵循 ISO 8601 格式，例如 YYYY-MM-DD、HH:MM:SS、YYYY-MM-DDTHH:MM:SSZ 等，且不帶引號。這需要強大的正規表達式或自定義的日期時間解析邏輯進行驗證。
  - **國家碼 (Country)**：ISO 3166-1 alpha-2 格式（兩個大寫英文字母），不帶引號。
  - **Base64**：Base64 編碼的字串，不帶引號。
- **鍵名 (Key Naming)**：
  - 鍵名本身不允許包含空格或點 (.) 符號，也不允許用引號包裹。
  - 同一物件範圍內**不允許出現重複的鍵名**，這是解釋器必須嚴格執行的錯誤檢查點。
  - 點分隔符 (.) 用於表示巢狀結構，例如 app.settings.debug: true 會被解析為 app: { settings: { debug: true } }。解釋器需要將這種扁平的點路徑轉換為深度巢狀的 Go map 結構。
- **結構元素 (Structural Elements)**：
  - **物件 {}**：由鍵值對組成，鍵值對之間可以通過換行（Expanded Form）或空格（Compact Form）分隔。
  - **陣列 []**：由值序列組成，值之間同樣可以通過換行或空格分隔。
- **分隔符 (Delimiters)**：**嚴禁使用逗號 (,)** 作為陣列元素或鍵值對的分隔符。這是 SDCL 與 JSON 的主要區別之一，解釋器必須檢查並報告此類語法錯誤。
- **註解 (Comments)**：
  - 單行註解以井字號 (#) 開頭。
  - **嚴格限制**：註解必須出現在其獨立的行上，不允許跟隨在鍵值對或其他數據元素之後（即不允許行尾註解，如 key: "value" # comment）。
- **顯式型別標籤 (Explicit Type Tags)**：
  - <type>value</type> 形式，例如 <int>123</int>。
  - 這些標籤是可選的，主要用於增強可讀性或解決字面量歧義。
  - 解釋器需要解析這些標籤，並在解析其內部 value 時進行類型驗證，確保 value 的實際類型與 <type> 標籤所聲明的類型一致。
- **路徑式引用與內容包含 (Referencing & Content Inclusion)**：SDCL 最強大的特性之一。
  - **值引用 (Value Reference)**：key: (path.to.value)。它允許一個鍵的值動態地引用同一個 SDCL 文件中其他路徑上的值。這需要在解析 AST 後的單獨階段進行處理。
  - **內容包含 (無鍵名)**：(path.to.object_or_array)。用於將另一個物件或陣列的*內容*（不包含其自身的鍵名）直接嵌入到當前物件或陣列中。
  - **內容包含 (有鍵名)**：((path.to.object_or_array))。用於將*整個*帶有鍵名的物件或陣列嵌入到當前物件中。
  - 這些引用的作用域、解析順序、以及對循環引用的嚴格檢測和錯誤報告，是 Resolver 模組的複雜核心。
- **外部引用 (External References)**：
  - **環境變數**：.env.KEY。直接從執行環境中獲取指定名稱的環境變數的值。
  - **外部 SDCL 檔案**：.XXX.sdcl.KEY。引用另一個 SDCL 文件中的值。這要求解釋器能夠載入外部檔案，並遞迴地解析它們。
  - **挑戰**：檔案載入策略（相對路徑解析、搜索路徑）、對循環引用的防範（例如 A 引用 B，B 引用 A），以及對檔案不存在或路徑不存在的錯誤處理。
- **內容覆蓋與合併 (Overriding & Merging)**：
  - SDCL 遵循「後定義覆蓋 (last definition wins)」原則，但僅限於通過**內容包含**或**值引用**機制引入的內容。
  - **直接定義的重複鍵是絕對不允許的**，必須拋出解析錯誤。
- **表單 (Forms)**：SDCL 嚴格區分兩種文件表單。
  - Expanded Form (擴展型)：使用換行和縮排清晰定義層次結構。
  - Compact Form (緊湊型)：物件和陣列定義在單行，元素間用空格分隔。
  - **關鍵限制**：一個 SDCL 文件中**不允許混用**這兩種表單。解釋器必須在解析初期判斷文件所屬表單，並在整個解析過程中強制執行該表單的規則。
- **前置區塊 (Front Matter)**：
  - 一個可選的元數據區塊，用三條連字符 --- 作為分隔符。
  - 如果存在，必須位於文件最開始。解釋器應將其內容解析為一個獨立的 SDCL 結構，並忽略其後的任何文本內容。這常見於將 SDCL 作為 Markdown 文件元數據的場景。

### **2. 核心架構**

SDCL 解釋器將採用經典的編譯器兩階段設計，並額外增加一個解決引用階段：

1. **詞法分析器 (Lexer)**：
   - **職責**：將原始 SDCL 字串分解為一系列有意義的詞法單元 (Tokens)。
   - **輸出**：一個 []token.Token 序列。
2. **語法分析器 (Parser)**：
   - **職責**：接收 Lexer 產生的 Token 序列，根據 SDCL 的語法規則構建一個抽象語法樹 (AST)。AST 是原始程式碼的結構化表示，但尚未解決所有引用。
   - **輸出**：一個 *ast.Program 結構，代表整個文件的 AST。
3. **解析器 (Resolver)**：
   - **職責**：遍歷 Parser 構建的 AST，解決所有內部和外部引用、執行內容包含合併、處理類型轉換，並產生最終的、可以直接使用的 Golang 資料結構。
   - **輸出**：一個 object.SdclObject (即 map[string]interface{}) 或 object.SdclArray (即 []interface{}) 的巢狀結構。

### **3. 模組詳述與 Golang 實現考量**

所有模組、類型、函數和變數都將遵循 Golang 的小駝峰命名規範。

#### **3.1 token 模組**

這個模組的核心是定義 SDCL 語言的最小構成單位。

- **檔案**：token/token.go
- **設計考量**：
  - 使用 type TokenType string 而非 int 類型，可以讓 Token 類型在調試和錯誤訊息中更具可讀性。
  - Token 結構體應包含 Type (TokenType)、Literal (原始字串值)、以及 Line 和 Column (用於精確錯誤報告)。
  - 所有 SDCL 特定的關鍵字、分隔符、字面量類型和特殊前綴都應該有其對應的 TokenType 常數。

#### **3.2 lexer 模組**

Lexer 是解釋器的第一道關卡，負責將原始輸入字串轉換為 Parser 可理解的 Token 序列。

- **檔案**：lexer/lexer.go, lexer/lexer_test.go
- **核心組件**：
  - Lexer 結構體：包含 input (原始 SDCL 字串)、position (當前讀取位置)、readPosition (下一個字符的位置)、ch (當前字符)、line (行號)、column (列號)。
  - New(input string) *Lexer：構造函數，初始化 Lexer 並讀取第一個字符。
  - readChar()：推進 Lexer 的讀取位置，更新 ch、position、line、column。
  - peekChar()：查看下一個字符而不推進位置。
  - NextToken()：核心方法，負責識別並返回下一個 Token。
- **實現細節與挑戰**：
  - **跳過空白字符**：skipWhitespace() 應處理空格、Tab、換行符。
  - **單字符 Token**：直接匹配 {, }, [, ], :, (, )。
  - **--- 處理**：需要檢查連續的三個 -，並區分它與負數的開始。
  - **註解處理**：readComment() 應讀取 # 後到行尾的所有字符，並返回 token.Comment。Lexer 不負責驗證註解的位置是否合法，這將是 Parser 的職責。
  - **字串處理**：readString() 必須能正確識別雙引號包裹的字串，並處理所有 SDCL 規範中定義的轉義字符。這是一個常見的錯誤源，需要仔細實現。
  - **數字處理**：readNumber() 應識別整數和浮點數，包括負號和小數點、指數部分。
  - **關鍵字/特殊字面量識別**：
    - readIdentifier()：先讀取一個連續的非空白、非特殊符號的字符序列。
    - lookupIdentifier(literal string) token.Token：這是最複雜的部分。它會判斷讀取到的 literal 是 true, false, null 這些關鍵字，還是 YYYY-MM-DD 格式的日期，HH:MM:SS 格式的時間，ISO 國家碼，或者 Base64 編碼的數據。
    - **挑戰**：由於這些都是不帶引號的，lookupIdentifier 必須有嚴格的優先順序和驗證邏輯。例如，它應該先嘗試匹配日期、時間、國家碼和 Base64（通常通過正規表達式或專門的格式檢查函數），如果都不匹配，再判斷是否為 true/false/null 關鍵字。如果這些都不是，那麼它才是一個 KeyName。這種歧義處理是 Lexer 的主要難點。
  - **顯式型別標籤**：readTypeTag() 應識別 <type> 和 </type> 標籤，例如 <int>、</str> 等。
  - **外部引用前綴**：readExternalReferencePrefix() 應識別 .env. 和 .XXX.sdcl. 形式。.XXX.sdcl. 中的 XXX 是一個動態的檔案名部分，Lexer 應能正確地將其作為一個整體識別出來。
  - **錯誤處理**：當遇到無法識別的字符時，應返回 token.Illegal 類型，並在 Token 中包含錯誤發生的位置信息。

#### **3.3 ast 模組**

AST (Abstract Syntax Tree) 是 SDCL 程式碼的抽象表示，它移除了語法上的細節（如括號、逗號等），只保留了程式碼的邏輯結構。

- **檔案**：ast/ast.go
- **設計考量**：
  - Node 接口：所有 AST 節點的基礎，定義了 TokenLiteral() 和 String() 方法，用於調試和表示。
  - Statement 接口：代表 SDCL 語句，如鍵值對聲明或內容包含聲明。
  - Expression 接口：代表 SDCL 值或可計算的表達式，如字面量、路徑引用。
  - Program 結構體：AST 的根節點，包含一系列頂層 Statement。
  - **KeyValueStatement**：表示 key: value 對。其 Key 應是一個 *KeyPath 結構，以便處理點分隔的鍵名。
  - **KeyPath**：表示 a.b.c 這種鍵名路徑，內部儲存為 []string{"a", "b", "c"}。
  - **InclusionStatement**：表示 (path) 或 ((path)) 內容包含。需要包含 Path (一個 *PathExpression) 和 WithKey 布林值來區分兩種包含類型。
  - **字面量表達式 (Literal Expressions)**：為每種 SDCL 資料型別 (StringLiteral, NumberLiteral, BooleanLiteral, NullLiteral, DateLiteral, TimeLiteral, DatetimeLiteral, CountryLiteral, Base64Literal) 定義獨立的結構體，它們都實現 Expression 接口。
  - **PathExpression**：表示 (path.to.value) 中的內部路徑，儲存為 []string 片段。
  - **EnvReferenceExpression 和 SdclReferenceExpression**：分別表示環境變數引用和外部 SDCL 檔案引用。外部檔案引用需要包含 FileName 和 KeyPath。
  - **TaggedLiteralExpression**：表示 <type>value</type> 結構，包含標籤類型 (TagType) 和內部值表達式 (Value)。
  - **ObjectLiteral 和 ArrayLiteral**：表示物件和陣列的 AST 節點，包含其內部元素，並標記是否為 IsCompact (緊湊型) 表單。

#### **3.4 parser 模組**

Parser 的職責是將 Lexer 產生的 Token 序列轉換為結構化的 AST。

- **檔案**：parser/parser.go, parser/parser_test.go
- **核心組件**：
  - Parser 結構體：包含 Lexer 實例、currentToken (當前處理的 Token)、peekToken (下一個 Token)、errors (錯誤列表)、isCompactForm (文件表單狀態)。
  - New(l *lexer.Lexer, loader FileLoader) *Parser：構造函數，初始化 Parser。
  - nextToken()：推進 currentToken 和 peekToken。
  - expectPeek(t token.TokenType) bool：檢查 peekToken 是否為預期類型，如果是則推進，否則記錄錯誤。
  - currentTokenIs(t token.TokenType) bool / peekTokenIs(t token.TokenType) bool：檢查當前/下一個 Token 的類型。
  - addError(err error)：統一的錯誤記錄機制。
  - FileLoader 接口：抽象了檔案載入邏輯，使得 Parser 可以獨立於具體的檔案系統實現，便於測試和擴展。
  - ParseProgram() (*ast.Program, error)：Parser 的入口點，解析整個 SDCL 文件。
- **實現細節與挑戰**：
  - **前置區塊 (Front Matter) 處理**：
    - ParseProgram 應首先檢查文件開頭是否為 ---。
    - 如果存在，進入 Front Matter 解析模式，解析其內部內容（這本身就是一個 SDCL 結構），直到遇到第二個 ---。
    - Front Matter 內容應被儲存在 ast.Program 的一個專門欄位中。
    - --- 之後的任何內容都應被 Parser 忽略。
  - **文件表單 (Form) 判斷**：
    - determineForm() 方法應在解析初期被呼叫。它需要檢查第一個物件或陣列的開頭符號（{ 或 [）之後，是立即跟隨換行還是空格。
    - 一旦確定了 isCompactForm，Parser 必須在整個解析過程中強制執行該表單的規則，如果發現混用，則拋出解析錯誤。
  - **parseStatement()**：解析頂層語句。SDCL 的頂層語句可以是鍵值對 (key: value) 或內容包含 ((path), ((path)))。
  - **parseKeyValueStatement()**：解析 key: value。這需要調用 parseKeyPath() 來解析鍵路徑，並調用 parseExpression() 來解析值。
  - **parseObject()**：
    - 期望當前 Token 是 {。
    - 根據 isCompactForm 狀態，決定是按空格（緊湊型）還是按換行（擴展型）分隔元素。
    - 在迴圈中，解析每個鍵值對 (KeyValueStatement) 或內容包含 (InclusionStatement)。
    - **重複鍵檢測**：在構建 ObjectLiteral 時，必須在每個 KeyValueStatement 被解析後，立即檢查其鍵名是否與同一物件範圍內的已存在鍵重複。如果重複，立即記錄錯誤並終止解析。
  - **parseArray()**：
    - 期望當前 Token 是 [。
    - 同樣根據 isCompactForm 狀態處理元素分隔。
    - 在迴圈中，解析每個陣列元素（這些元素本身也是 Expression）。
  - **parseExpression()**：核心的遞迴下降解析方法。它會根據當前 Token 的類型，分派給不同的解析子函數。
  - **parseLiteral() 族函數**：為每種 SDCL 字面量類型（字串、數字、布林值、日期、時間等）提供特定的解析函數，將其轉換為對應的 AST 節點。
  - **parsePathExpression()**：解析 a.b.c 這樣點分隔的路徑，將其拆分為 []string{"a", "b", "c"}。
  - **parseReferenceOrInclusion()**：處理 ( 開頭的 Token。它需要區分 (path) 和 ((path))，並解析其內部的路徑。
  - **parseTaggedLiteral()**：解析 <type>value</type> 結構。在解析 value 後，需要驗證其類型是否與 tagType 一致，並處理結束標籤 </type>。
  - **錯誤報告**：所有解析函數都應返回 (result, error)，並且在發生錯誤時，將錯誤信息（包含精確的行號和列號）記錄到 Parser 的 errors 列表中。

#### **3.5 resolver 模組**

Resolver 是解釋器的最終階段，負責將 Parser 生成的 AST 轉換為可直接使用的 Golang 數據結構，同時解決所有引用和執行數據合併。

- **檔案**：resolver/resolver.go
- **核心組件**：
  - Resolver 結構體：
    - rootData object.SdclObject：儲存最終解析並解決後的 SDCL 數據。
    - loadingFiles map[string]bool：用於檢測外部 SDCL 檔案的循環引用。
    - fileLoader parser.FileLoader：共享的檔案載入器實例。
    - errors []error：記錄解決過程中的錯誤。
  - New(loader parser.FileLoader) *Resolver：構造函數。
  - ResolveProgram(program *ast.Program, currentFilePath string) (object.SdclObject, error)：Resolver 的入口點。
- **實現細節與挑戰**：
  - **兩階段解決 (或多階段)**：為了解決前向引用和相互引用，Resolver 可能需要至少兩次遍歷：
    1. **第一遍 (構建骨架)**：遍歷 AST，將所有直接定義的鍵值對（包括空的物件和陣列）以及引用/包含的**佔位符**填充到 rootData 這個 map[string]interface{} 結構中。這確保了所有路徑在第二遍解決時都能被查找。
    2. **第二遍 (解決引用與合併)**：遞迴遍歷 rootData。當遇到一個 AST 引用節點（例如 *ast.PathExpression、*ast.EnvReferenceExpression 等）時，執行其解決邏輯：
       - **內部引用 (*ast.PathExpression)**：使用 lookupPath(data interface{}, segments []string) 函數在 rootData 中查找對應路徑的值。
       - **環境變數引用 (*ast.EnvReferenceExpression)**：調用 os.Getenv(key) 獲取環境變數。如果變數不存在或為空，根據規範拋出錯誤。
       - **外部 SDCL 檔案引用 (*ast.SdclReferenceExpression)**：
         - 計算外部檔案的絕對路徑。
         - **循環引用檢測**：在 loadingFiles 中檢查該絕對路徑是否已存在。如果存在，立即拋出循環引用錯誤。如果不存在，則將其加入 loadingFiles，並在函數結束時移除 (使用 defer)。
         - 使用 fileLoader 載入檔案內容。
         - **遞迴解析**：為外部檔案內容創建新的 lexer 和 parser 實例，並遞迴調用 parser.ParseProgram() 獲取其 AST。
         - **遞迴解決**：為外部檔案的 AST 創建新的 resolver 實例，並遞迴調用 extResolver.ResolveProgram()。將當前檔案路徑傳遞給遞迴調用，以便循環引用檢測。
         - 從解析後的外部數據中查找 KeyPath 指定的值。
       - **內容包含 (*ast.InclusionStatement)**：
         - 解析其 Path，獲取被包含的來源物件或陣列。
         - 如果 WithKey 為 true (((path)))，將來源物件/陣列連同其「根鍵名」（即路徑的最後一個片段）嵌入到當前物件中。
         - 如果 WithKey 為 false ((path))：
           - 如果來源是被包含到**物件**中，將來源物件的所有鍵值對合併到當前物件。此處必須應用「後定義覆蓋」原則。
           - 如果來源是被包含到**陣列**中，將來源陣列的所有元素追加到當前陣列。
       - **類型轉換與驗證**：在解決引用或最終將值設置到 rootData 時，執行必要的類型轉換（例如，將字串數字轉換為 float64 或 int，日期字串轉換為 time.Time 對象）。
  - **lookupPath(data interface{}, segments []string) (interface{}, error)**：一個輔助函數，用於在巢狀 map 或 slice 中根據路徑片段查找值。
  - **setNestedValue(data object.SdclObject, segments []string, value interface{})**：一個輔助函數，用於根據路徑片段在巢狀 map 中設置值，並處理不存在的巢狀結構的創建。
  - **錯誤傳遞與處理**：所有解析和解決函數都應返回 (result, error)，確保錯誤能層層上報，並最終由 main 函數統一處理和列印。

#### **3.6 object 模組**

這個模組定義了 SDCL 最終解析並解決後的 Golang 數據結構。

- **檔案**：object/object.go
- **設計考量**：
  - SdclObject：簡單地定義為 map[string]interface{}，這是 Golang 中最自然的物件表示方式。
  - SdclArray：簡單地定義為 []interface{}，這是 Golang 中最自然的陣列表示方式。
  - 可以選擇為 date, time, datetime, country, base64 定義更具體、帶有驗證邏輯的自定義 Go 類型（例如 type SdclDate time.Time），以提供更強的類型安全和數據處理能力。這取決於對數據模型精確度的要求。

#### **3.7 main 模組**

這個模組是解釋器的入口點，負責整合所有其他模組，提供一個可執行的 CLI 介面。

- **檔案**：main.go
- **核心功能**：
  - 處理命令行參數，獲取 SDCL 文件路徑。
  - 讀取 SDCL 文件內容。
  - 創建 lexer.Lexer 實例。
  - 創建 parser.Parser 實例，並傳入一個 parser.FileLoader 接口的實現（例如 BasicFileLoader，它會從檔案系統讀取檔案）。
  - 調用 parser.ParseProgram() 獲取 AST。
  - 創建 resolver.Resolver 實例，並傳入相同的 parser.FileLoader。
  - 調用 resolver.ResolveProgram() 解決引用並獲取最終數據。
  - 處理並列印所有解析和解決過程中發生的錯誤。
  - 將最終的 Go 數據結構列印到標準輸出，可以使用 fmt.Printf("%+vn", resolvedData) 進行簡單輸出，或者使用 encoding/json 包將其轉換為 JSON 格式以便於查看。

### **7. 開發與測試策略**

- **逐步開發**：從最簡單的 Lexer 開始，然後是能夠解析基本字面量的 Parser，接著是物件和陣列，最後才是 SDCL 特有的複雜特性（引用、包含、表單判斷、Front Matter）。
- **單元測試**：為每個模組（token, lexer, ast, parser, resolver, object）編寫詳盡的單元測試。
  - **Lexer 測試**：確保所有 Token 類型都能正確識別，包括邊界情況和錯誤輸入。
  - **Parser 測試**：測試所有 SDCL 語法結構的解析，包括嵌套、合法和非法的重複鍵、錯誤的註解位置、不允許的逗號、以及正確的表單判斷。
  - **Resolver 測試**：這是最關鍵的測試。需要測試所有引用類型（內部、環境變數、外部檔案），內容包含的合併邏輯（特別是「後定義覆蓋」），以及對循環引用的檢測。
- **集成測試**：編寫端到端的集成測試，從讀取整個 SDCL 檔案到最終生成數據，驗證整個流程的正確性。
- **錯誤處理測試**：確保在各種錯誤情況下（語法錯誤、引用錯誤、檔案不存在等）都能產生清晰、準確的錯誤訊息，包含行號和列號。
- **規範遵循**：在開發過程中，頻繁參照 SDCL 規範，確保每一個細節都得到正確實現。

這個詳細的規劃將作為你使用 Golang 構建 SDCL 解釋器的路線圖，涵蓋了從底層詞法分析到高層語意解析的所有關鍵點和挑戰。
