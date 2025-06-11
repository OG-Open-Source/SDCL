## **SDCL 解釋器 Parser 模組詳細開發計劃 (Golang)**

本計劃將專注於 SDCL 解釋器的 **Parser (語法分析器)** 模組，詳細闡述其開發步驟、核心挑戰及 Golang 實現細節。Parser 模組的職責是將 Lexer 產生的 Token 序列轉換為抽象語法樹 (AST)。

### **1. Parser 模組概覽**

- **職責**：
  - 接收來自 Lexer 的 token.Token 序列。
  - 根據 SDCL 語法規則，驗證 Token 序列的合法性。
  - 構建一個 ast.Program（抽象語法樹）來表示 SDCL 文件的結構。
  - 處理語法層面的錯誤，並生成精確的錯誤報告。
- **輸入**：*lexer.Lexer 實例。
- **輸出**：*ast.Program 實例（如果解析成功），或 error（如果解析失敗）。
- **關鍵挑戰**：SDCL 規範的嚴格性、多種語法結構（物件、陣列、引用、標籤、前置區塊）、以及對重複鍵和表單混用的強制檢查。

### **2. Parser 核心結構 (parser/parser.go)**

package parser

import (  
 "errors"  
 "fmt"  
 "sdcl/ast" // 引入 AST 模組  
 "sdcl/lexer" // 引入 Lexer 模組  
 "sdcl/token" // 引入 Token 模組  
)

// FileLoader 接口定義了載入外部 SDCL 檔案的方法  
// 這在 Parser 階段主要用於外部 SDCL 檔案引用的語法檢查，實際載入和解決在 Resolver 階段  
type FileLoader interface {  
 LoadSdclFile(filePath string) (string, error)  
}

// Parser 結構體，包含了解析器所需的所有狀態  
type Parser struct {  
 l *lexer.Lexer // 詞法分析器實例  
 currentToken token.Token // 當前正在處理的 Token  
 peekToken token.Token // 下一個 Token (前瞻 Token)  
 errors []error // 儲存解析過程中遇到的所有錯誤  
 isCompactForm *bool // 指示文件是否為緊湊型表單 (nil: 未確定, true: 緊湊, false: 擴展)  
 fileLoader FileLoader // 外部檔案載入器，用於處理 .XXX.sdcl.KEY 的語法  
}

// New 構造函數，創建一個新的 Parser 實例  
func New(l *lexer.Lexer, loader FileLoader) *Parser {  
 p := &Parser{  
 l: l,  
 errors: make([]error, 0),  
 fileLoader: loader, // 傳遞檔案載入器  
 }  
 // 初始化 currentToken 和 peekToken  
 p.nextToken()  
 p.nextToken()  
 return p  
}

// nextToken 將 currentToken 推進到 peekToken，並從 Lexer 獲取新的 peekToken  
func (p *Parser) nextToken() {  
 p.currentToken = p.peekToken  
 p.peekToken = p.l.NextToken()  
}

// currentTokenIs 檢查當前 Token 是否為指定類型  
func (p *Parser) currentTokenIs(t token.TokenType) bool {  
 return p.currentToken.Type == t  
}

// peekTokenIs 檢查下一個 Token 是否為指定類型  
func (p *Parser) peekTokenIs(t token.TokenType) bool {  
 return p.peekToken.Type == t  
}

// expectPeek 檢查下一個 Token 是否為預期類型。如果是，則推進 Token 並返回 true；否則，記錄錯誤並返回 false。  
func (p *Parser) expectPeek(t token.TokenType) bool {  
 if p.peekTokenIs(t) {  
 p.nextToken()  
 return true  
 } else {  
 p.peekError(t) // 記錄錯誤  
 return false  
 }  
}

// peekError 記錄一個預期 Token 類型不匹配的錯誤  
func (p *Parser) peekError(t token.TokenType) {  
 msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d, column %d",  
 t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)  
 p.errors = append(p.errors, errors.New(msg))  
}

// addError 記錄一個自定義的解析錯誤  
func (p *Parser) addError(err error) {  
 p.errors = append(p.errors, err)  
}

// Errors 返回所有記錄的解析錯誤  
func (p *Parser) Errors() []error {  
 return p.errors  
}

### **3. Parser 的主要解析流程與挑戰**

#### **3.1 ParseProgram()：文件級解析**

這是 Parser 的入口點，負責解析整個 SDCL 文件。

- **功能**：
  1. **前置區塊 (Front Matter) 處理**：
     - 檢查文件是否以 --- 開頭。
     - 如果存在，則進入 Front Matter 解析模式：解析其內部內容作為一個獨立的 SDCL 結構（通常是一個物件），直到遇到第二個 ---。
     - Parser 應將 Front Matter 內容儲存到 ast.Program 的一個專用欄位中。
     - **挑戰**：在 Front Matter 內部，需要遞迴地調用 Parser 的其他解析方法。一旦 Front Matter 結束，--- 後的任何內容應被 Parser 忽略（交由其他工具處理，如 Markdown 解析器）。
  2. **文件表單 (Form) 判斷**：
     - 在解析主體內容之前，調用 determineForm() 方法來判斷整個 SDCL 文件是 Expanded Form 還是 Compact Form。
     - **挑戰**：判斷邏輯通常基於第一個物件或陣列的結構（例如，{ 後是立即換行還是空格）。一旦判斷，這個狀態 (p.isCompactForm) 將在整個解析過程中被強制執行。
  3. **主體內容解析**：
     - 進入一個循環，不斷調用 parseStatement() 直到遇到 EOF。
     - 將所有解析出的 Statement 收集到 ast.Program 的 Statements 列表中。
  4. **錯誤收集**：返回 *ast.Program 或所有收集到的錯誤。

#### **3.2 determineForm()：文件表單判斷**

這個輔助方法負責在文件開頭判斷文件的整體表單。

- **功能**：
  1. 跳過任何頂層註解或空白行。
  2. 找到第一個非註解、非空白的 Token。
  3. 如果這個 Token 是 { (物件開始) 或 [ (陣列開始)：
     - 查看 peekToken：如果它是換行符 n 或 Eof，則極有可能是 Expanded Form。
     - 如果 peekToken 是其他任何非空白字符，則可能是 Compact Form。
  4. 設置 p.isCompactForm 的值。
  5. **挑戰**：這個判斷必須足夠魯棒，能夠在複雜的開頭情況下正確識別表單。一旦設置，parseObject() 和 parseArray() 必須根據此狀態嚴格執行其解析邏輯。

#### **3.3 parseStatement()：語句級解析**

解析 SDCL 文件中的一個頂層語句。

- **功能**：
  - 根據 currentToken 的類型，分派給不同的子解析函數。
  - 主要處理兩種頂層語句：
    1. **鍵值對聲明**：以 KeyName 開頭（例如 app.name: "My App"）。
    2. **內容包含聲明**：以 LParen (即 () 開頭（例如 (common_settings) 或 ((user_data))）。
  - **挑戰**：需要明確區分這兩種語句的開頭，並調用正確的解析器。

#### **3.4 parseKeyValueStatement()：鍵值對解析**

解析形如 key: value 的鍵值對。

- **功能**：
  1. **解析鍵路徑**：調用 parseKeyPath() 來處理點分隔的鍵名（例如 app.settings.debug）。
  2. 期望下一個 Token 是 token.Colon (:), 如果不是則報錯。
  3. **解析值**：調用 parseExpression() 來解析 : 後面的值（可以是字面量、物件、陣列、引用等）。
  4. 返回一個 *ast.KeyValueStatement 節點。
  5. **挑戰**：**重複鍵檢測**：在解析物件內部時，Parser 必須在將鍵值對添加到 ObjectLiteral 的 Elements 列表之前，檢查其鍵是否已在當前物件範圍內存在。如果重複，則記錄錯誤並終止解析。這通常需要一個臨時的 map 來跟踪已見過的鍵。

#### **3.5 parseKeyPath()：鍵路徑解析**

解析點分隔的鍵名。

- **功能**：
  1. 期望當前 Token 是 token.KeyName。
  2. 讀取第一個鍵段。
  3. 進入循環：如果下一個 Token 是 .，則消費 .，然後期望下一個 Token 還是 KeyName，並將其作為新的鍵段加入。
  4. 返回一個 *ast.KeyPath 節點，其中包含所有鍵段（例如 []string{"app", "settings", "debug"}）。
  5. **挑戰**：確保正確處理單個鍵名和多個點分隔的鍵名。

#### **3.6 parseObject()：物件字面量解析**

解析形如 { ... } 的 SDCL 物件。

- **功能**：
  1. 期望當前 Token 是 token.LBrace ({)。
  2. 根據 p.isCompactForm 的狀態，決定解析其內部元素的方式：
     - **Compact Form**：元素之間預期由**空格**分隔。Parser 將在同一行上連續解析鍵值對或內容包含，直到遇到 }。
     - **Expanded Form**：預期每個鍵值對或內容包含出現在**新行**上。Parser 將跳過換行符並解析每個元素。
  3. 在迴圈中，不斷調用 parseKeyValueStatement() 或 parseInclusionStatement() 來解析物件的內部元素。
  4. 期望最後一個 Token 是 token.RBrace (})。
  5. 返回一個 *ast.ObjectLiteral 節點。
  6. **挑戰**：這是 Parser 中最複雜的函數之一。它需要精確地處理兩種表單的元素分隔，並在解析每個元素時進行上述的重複鍵檢測。

#### **3.7 parseArray()：陣列字面量解析**

解析形如 [ ... ] 的 SDCL 陣列。

- **功能**：
  1. 期望當前 Token 是 token.LBracket ([).
  2. 同樣根據 p.isCompactForm 的狀態，決定解析其內部元素的方式（空格分隔或換行分隔）。
  3. 在迴圈中，不斷調用 parseExpression() 來解析每個陣列元素。
  4. 期望最後一個 Token 是 token.RBracket (]).
  5. 返回一個 *ast.ArrayLiteral 節點。
  6. **挑戰**：與 parseObject() 類似，需要精確地處理兩種表單的元素分隔。

#### **3.8 parseExpression()：表達式解析**

這是遞迴下降的核心，根據當前 Token 類型分派到不同的解析函數。

- **功能**：
  - 根據 currentToken 的類型，調用對應的 parseLiteral() 族函數（例如 parseStringLiteral(), parseNumberLiteral() 等）。
  - 如果 Token 是 {，調用 parseObject()。
  - 如果 Token 是 [，調用 parseArray()。
  - 如果 Token 是 (，調用 parseReferenceOrInclusion()。
  - 如果 Token 是 <，調用 parseTaggedLiteral()。
  - 如果 Token 是 .env. 或 .XXX.sdcl.，則解析為對應的外部引用表達式。
- **挑戰**：確保所有可能的 SDCL 值類型都能被正確識別和解析為對應的 AST 節點。

#### **3.9 parseReferenceOrInclusion()：引用與包含解析**

解析形如 (path) 或 ((path)) 的結構。

- **功能**：
  1. 期望當前 Token 是 token.LParen (()。
  2. 檢查下一個 Token：
     - 如果下一個 Token 也是 (，則表示 ((path))，設置 withKey = true，並推進 Token。
     - 否則，表示 (path)，設置 withKey = false。
  3. 解析內部路徑（調用 parsePathExpression()）。
  4. 期望下一個 Token 是 token.RParen ())。
  5. 返回一個 *ast.PathExpression (如果解析的是值引用) 或 *ast.InclusionStatement (如果解析的是內容包含)。
- **挑戰**：精確地區分這三種語法，並將其表示為 AST 中正確的節點類型。

#### **3.10 parseTaggedLiteral()：顯式型別標籤解析**

解析形如 <type>value</type> 的結構。

- **功能**：
  1. 期望當前 Token 是 <type> 開頭的標籤（例如 token.LtIntGt）。
  2. 提取 tagType（例如 "int"）。
  3. 調用 parseExpression() 來解析 <type> 和 </type> 之間的值。
  4. 期望下一個 Token 是對應的 </type> 結束標籤。
  5. **挑戰**：
     - **類型驗證**：在解析完內部 value 後，Parser 必須驗證 value 的實際類型是否與 tagType 所聲明的類型一致。例如，如果 tagType 是 "int"，而 value 是 ast.StringLiteral，則應記錄錯誤。
     - **結束標籤匹配**：確保結束標籤與開始標籤類型一致。

#### **3.11 錯誤處理策略**

- Parser 應該有自己的 errors []error 列表，用於收集解析過程中遇到的所有錯誤。
- 每個解析函數都應該返回 (result, error)，在發生錯誤時，將錯誤資訊（包含 token.Line 和 token.Column）添加到 p.errors 列表中，並返回 nil, err 或一個部分 AST 節點和錯誤。
- ParseProgram() 在結束時，應檢查 p.errors 列表，如果非空，則返回一個包含所有錯誤的 error。

### **4. 測試策略 (針對 Parser 模組)**

- **語法正確性測試**：
  - 為所有 SDCL 規範中的有效語法結構編寫測試用例，包括物件、陣列、各類字面量、不同形式的引用和包含。
  - 測試巢狀結構的深度，確保遞迴解析正確。
  - 測試 Expanded Form 和 Compact Form 的單獨使用，確保 Parser 能正確識別並解析。
- **錯誤語法測試**：
  - **重複鍵**：測試同一物件範圍內重複鍵的輸入，確保 Parser 能準確捕捉並報告錯誤。
  - **逗號分隔符**：測試在陣列或物件中使用逗號的情況，確保 Parser 報告錯誤。
  - **行尾註解**：測試 key: "value" # comment 這種情況，確保 Parser 報告錯誤。
  - **表單混用**：測試在同一文件中混用 Expanded Form 和 Compact Form 的情況，確保 Parser 報告錯誤。
  - **無效鍵名**：測試鍵名包含空格或引號的情況。
  - **不匹配的括號/引號/標籤**：測試語法不匹配的情況。
  - **無效的顯式類型標籤**：測試 <int>"hello"</int> 這樣值與標籤類型不匹配的情況。
- **錯誤報告測試**：
  - 對於所有錯誤情況，驗證 Parser 報告的錯誤訊息是否清晰、準確，是否包含正確的行號和列號。

Parser 是 SDCL 解釋器的大腦，其設計和實現的健壯性直接決定了整個解釋器的質量。投入足夠的時間在設計、實現細節和全面的測試上，是確保成功的關鍵。
