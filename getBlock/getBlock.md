# GetBlock.io请求编码对比

## 结论

不同的编码格式大小不一样，gzip压缩有助于节省带宽

|编码格式\传输方式|未压缩文件大小|gzip压缩文件大小|
|:---|---:|---:|
|jsonParsed|11MB|2.5MB|
|base64|6.1MB|1.9MB|
|base58|6.2MB|2.1MB|

如果将本地原始文件直接压缩，效果会更好，体积更小
|编码格式|原始文件大小|压缩后文件大小|
|---|---:|---:|
|jsonParsed|11MB|1.4MB|
|base64|6.1MB|1.2MB|
|base58|6.2MB|1.5MB|


## jsonParsed
```bash
curl -X POST "https://go.getblock.us/${GETBLOCK_ACCESS_TOKEN}/" \
  --output 355825332.jsonParsed.json \
  --header "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getBlock",
    "params": [
        355825332,
        {
            "encoding": "jsonParsed",
            "maxSupportedTransactionVersion": 0,
            "transactionDetails": "full",
            "rewards": true
        }
    ]
  }'
```
- jsonParsed
- 传输未压缩
- 文件大小：11MB

## jsonParsed + gzip
```bash
curl -X POST "https://go.getblock.us/${GETBLOCK_ACCESS_TOKEN}/" \
  --output 355825332.jsonParsed.gz \
  -H "Accept-Encoding: gzip" \
  --header "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getBlock",
    "params": [
        355825332,
        {
            "encoding": "jsonParsed",
            "maxSupportedTransactionVersion": 0,
            "transactionDetails": "full",
            "rewards": true
        }
    ]
  }'
```
- jsonParsed + gzip
- 传输gzip压缩
- 文件大小：2.5MB

## base64
```bash
curl -X POST "https://go.getblock.us/${GETBLOCK_ACCESS_TOKEN}/" \
  --output 355825332.base64.json \
  --header "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getBlock",
    "params": [
        355825332,
        {
            "encoding": "base64",
            "maxSupportedTransactionVersion": 0,
            "transactionDetails": "full",
            "rewards": true
        }
    ]
  }'
```
- base64
- 传输未压缩
- 文件大小：6.1MB

## base64 + gzip
```bash
curl -X POST "https://go.getblock.us/${GETBLOCK_ACCESS_TOKEN}/" \
  --output 355825332.base64.gz \
  -H "Accept-Encoding: gzip" \
  --header "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getBlock",
    "params": [
        355825332,
        {
            "encoding": "base64",
            "maxSupportedTransactionVersion": 0,
            "transactionDetails": "full",
            "rewards": true
        }
    ]
  }'
```
- base64 + gzip
- 传输gzip压缩
- 文件大小：1.9MB

## base58
```bash
curl -X POST "https://go.getblock.us/${GETBLOCK_ACCESS_TOKEN}/" \
  --output 355825332.base58.json \
  --header "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getBlock",
    "params": [
        355825332,
        {
            "encoding": "base58",
            "maxSupportedTransactionVersion": 0,
            "transactionDetails": "full",
            "rewards": true
        }
    ]
  }'
```
- base58
- 传输未压缩
- 文件大小：6.2MB

## base58 + gzip
```bash
curl -X POST "https://go.getblock.us/${GETBLOCK_ACCESS_TOKEN}/" \
  --output 355825332.base58.gz \
  -H "Accept-Encoding: gzip" \
  --header "Content-Type: application/json" \
  --data '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "getBlock",
    "params": [
        355825332,
        {
            "encoding": "base58",
            "maxSupportedTransactionVersion": 0,
            "transactionDetails": "full",
            "rewards": true
        }
    ]
  }'
```
- base58 + gzip
- 传输gzip压缩
- 文件大小：2.1MB
