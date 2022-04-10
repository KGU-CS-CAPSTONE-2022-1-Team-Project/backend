# 백엔드

## 주요기능

- 회원가입
  - 구글로그인을 통한 회원가입
    - 정보
      - 지갑
      - 팬아트 제작자
    - 스트리머 여부를 확인
      - 가입자가 직접 선택
        - 만약 스트리머면 인증절차 수행
- 로그인
    - 1초에 최대 3회 이상 인증시 3분간 ip밴 수행

- 스트리머 인증 
  - 유튜브 api를 이용해 해당 채널이 존재하는지 확인
  - 스트리머 인증 조건(유튜브 수익창출 허락 기준 따름)
  - [스트리머 인증 조건](https://developers.google.com/youtube/v3/docs/channels?hl=ko)
    1. 가입 후 3개월 이상
    2. 구독자 최소 1000명 이상
    3. 만18세 이상
  - 스트리머 인증 완료시 해당 스트리머에 대한 스마트 컨트랙 제작

- 2차 창작자의 NFT민팅
    - 원하는 스트리머를 선택
    - 팬아트 민팅 허용 요청
      - 그림
      - 가격
    - 그림에 대한 유사도 검사 수행
    - 유사도 검사를 통과하면 스트리머 컨트랙에 민팅 요청
    - 민팅완료된 주소를 2차창작자에게 반환
