## 01@List Pharmacy Specify Timestamp
#### GET `/pharmacy/v1/`

##### Request field (querystring)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
page | uint64 |    X     | - | 頁碼
row | uint64 |    X     | - | 筆數
specify_utc0_millisecond_timestamp | int64  |    X     | - | 指定時間戳（UTC+0 millisecond timestamp)

##### Pharmacy struct
field           |  type   | description
:--------------|:-------:|:----
uid | string  | pharmacy unique id
name | string  | pharmacy 名稱
cash_balance | float64 | pharmacy 現金餘額
created_time | string  | pharmacy 資料建立時間
day |  int64  | 開店日（1 = 星期一, 2 = 星期二, 3 = 星期三, 4 = 星期四, 5 = 星期五, 6 = 星期六, 0 = 星期日）
open_hour | float64 | 開店時段（24hour)
close_hour | float64 | 閉店時段（24hour)

##### Response field(JSON)
field           |    type    | description
:--------------|:----------:|:----
Count |   int64    | 總數
Row |   int64    | 筆數
Page |   int64    | 頁碼
pharmacies | []Pharmacy | pharmacy資料列, 參照 `Pharmacy struct`

## 02@Search For Pharmacies Or Masks Name
#### GET `/pharmacy/v1/mix`

##### Request field (querystring)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
page | uint64 |    X     | - | 頁碼
row | uint64 |    X     | - | 筆數
name | string |    X     | - | 查詢的字節

##### PharmacyProduct struct
field           |  type   | description
:--------------|:-------:|:----
uid | string  | pharmacy unique id
product_id | string  | product unique id
pharmacy_name | string  | pharmacy 名稱
product_name | string  | product 名稱
cash_balance | float64 | pharmacy 現金餘額
price | float64  | product 價格

##### Response field(JSON)
field           |       type        | description
:--------------|:-----------------:|:----
Count |       int64       | 總數
Row |       int64       | 筆數
Page |       int64       | 頁碼
pharmacy_products | []PharmacyProduct | pharmacyProduct, 參照 `PharmacyProduct struct`

## 03@List Product By Pharmacy
#### GET `/pharmacy/v1/{:pharmacy_uid}/product`

##### Request field (querystring)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
page | uint64 |    X     | - | 頁碼
row | uint64 |    X     | - | 筆數
sorted | string |    X     | - | 排序的欄位(support name / price)

##### Product struct
field           |  type   | description
:--------------|:-------:|:----
uid | string  | pharmacy unique id
product_id | string  | product unique id
name | string  | product 名稱
price | float64  | product 價格
created_time | string  | product 資料建立時間

##### Response field(JSON)
field           |   type    | description
:--------------|:---------:|:----
Count |   int64   | 總數
Row |   int64   | 筆數
Page |   int64   | 頁碼
products | []Product | product資料列, 參照 `Product struct`

## 04@List Pharmacies By Product Price Range
#### GET `/pharmacy/v1/product/price`

##### Request field (querystring)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
page | uint64 |    X     | - | 頁碼
row | uint64 |    X     | - | 筆數
min | int64  |    X     | - | 價格最小值
max | int64  |    X     | - | 價格最大值

##### Pharmacy struct
field           |  type   | description
:--------------|:-------:|:----
uid | string  | pharmacy unique id
name | string  | pharmacy 名稱
cash_balance | float64 | pharmacy 現金餘額
created_time | string  | pharmacy 資料建立時間

##### Response field(JSON)
field           |   type    | description
:--------------|:---------:|:----
Count |   int64   | 總數
Row |   int64   | 筆數
Page |   int64   | 頁碼
pharmacies | []Pharmacy | pharmacy, 參照 `Pharmacy struct`

## 05@List Top X Users Transaction Amount
#### GET `/transaction/v1/transaction/top`

##### Request field (querystring)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
top_number | int64  |    X     | - | 最排前幾名
utc0_millisecond_start_timestamp | int64  |    X     | - | 查詢範圍起始時間（UTC+0 millisecond timestamp）
utc0_millisecond_end_timestamp | int64  |    X     | - | 查詢範圍終止時間（UTC+0 millisecond timestamp）

##### TopTransactionAmountUser struct
field           |  type   | description
:--------------|:-------:|:----
uid | string  | user unique id
name | string  | user 名稱
transaction_amount | float64 | 總交易金額

##### Response field(JSON)
field           |            type            | description
:--------------|:--------------------------:|:----
top_transaction_amount_users | []TopTransactionAmountUser | topTransactionAmountUser, 參照 `TopTransactionAmountUser struct`

## 06@Get Transaction Total By Data Range
#### GET `/transaction/v1/transaction/product`

##### Request field (querystring)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
utc0_millisecond_start_timestamp | int64  |    X     | - | 查詢範圍起始時間（UTC+0 millisecond timestamp）
utc0_millisecond_end_timestamp | int64  |    X     | - | 查詢範圍終止時間（UTC+0 millisecond timestamp）

##### Response field(JSON)
field           |    type    | description
:--------------|:----------:|:----
total |   int64    | 交易總數
transaction_amount |  float64   | 交易總金額

## 07@Purchase
#### POST `/transaction/v1/purchase`

##### Request field (JSON)
field           |  type  | required | validate | description
:--------------|:------:|:--------:|:----:|:----
user_id | string |    O     | - | 購買人 user unique id
pharmacy_id | string |    O     | - | 購買店家 pharmacy unique id
product_id | string |    O     | - | 購買產品 product unique id
quantity |  int   |    O     | - | 購買產品數量

##### Response field(Text)
`ok`
