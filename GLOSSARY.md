# HomeBox 领域上下文

自托管的家庭物品管理系统：一个 Group（群组/家庭）拥有物品，并设置全局货币用于展示物品价值。金额在数据库中以纯数字存储，货币含义由所属 Group 的货币代码决定。

## Language

**货币代码（Currency code）**:
货币的规范标识符，采用 ISO 4217 标准（CNY、USD、EUR）。是 Group 设置中存储、并在 API 中传输的值；约定数据库存小写、API 返回大写。人民币的规范代码是 CNY。
_Avoid_: RMB（人民币通用缩写，不是合法 ISO 4217 代码）、currency

**货币符号（Currency symbol）**:
金额展示时由货币代码经 Intl.NumberFormat 推导出的符号（CNY → ¥、USD → $）。不存储、不传输、不参与领域逻辑。
_Avoid_: currency code、手写符号映射

**货币名称（Currency name）**:
货币选择下拉框中展示的人类可读名称（如 "Chinese Yuan"）。仅展示用，不得与存储代码混用。
_Avoid_: 货币代码、currency

**Group**:
物品与设置（含货币）的归属单元。货币是 Group 级属性；Item 不持有货币字段。
_Avoid_: 家庭、租户、账户

**Item**:
库存物品条目，含购买价（purchasePrice）、原价（originalPrice）等纯数字金额字段，其货币含义由所属 Group 的货币代码决定。
_Avoid_: 资产、货物
