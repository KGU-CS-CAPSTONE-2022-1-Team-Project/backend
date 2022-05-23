# API REFERENCE

- 처리결과는 http상태코드로 반환
- 세부내용은 message로 제공
- url: https://capston-blockapp.greenflamingo.dev:10321

## GET /owner/google

- 로그인 페이지로 리다이렉팅

## GET /owner/google/callback

- oauth callback
- redirecting으로 값을 전달

### Response

```json
{
  "access_token": "...."
}
```

## PUT /owner/youtuber

- 유튜브 채널 검증 및 주소등록
- 등록이 통과될 시 블록체인에 contract생성
  - contract처리 결과는 제공하지 않는다
- bearer token 필요

### Request

```json
{
  "address": "..."
}
```

### Response

```json
{
  "message": "success"
}
```

## PUT /owner/youtuber/:id

- contract URI
- contract URI에 존재하는 id를 통해 값을 조회

### Response

```json
{
  "title":    "...",
  "description": "...",
  "image":       "...",
  "external_url":         "...",
  "seller_fee_basis_points": 1000,
  "fee_recipient":           "..."
}
```

## POST /partner/nft

### Request

- form data형식으로 다음의 정보를 입력
  - image: 3mb미만의 jpg,jpeg,png파일 1개
  - name: nft이름
  - description: nft설명

### Response

```json
{
  "message": "success",
  "id":      "..."
}
```

## POST /partner/nft

### Request

- form data형식으로 다음의 정보를 입력
  - image: 3mb미만의 jpg,jpeg,png파일 1개
  - name: nft이름
  - description: nft설명

### Response

```json
{
  "message": "success",
  "id":      "..."
}
```

## GET /partner/nft/:id

- 등록할 때 반환된 id를 이용해 nft리소스를 조회
- NFT URI기능

### Response

```json
{
  "name":        "...",
  "description": "...",
  "image":       "url"
}
```

## GET /common/nickname/:address

- 주소를 요청하면 닉네임을 반환

### Response

```json
{
  "message": "success",
  "nickname": "..."
}
```

## POST /common/nickname/

- 닉네임 등록

### Request

```json
{
  "address": "...",
  "nickname": "..."
}
```

### Response

```json
{
  "message": "success"
}
```