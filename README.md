## 何これ？
FANZA同人で大幅割引中のやつ一覧ツール  
https://twitter.com/tetsujp84/status/1326580965771149312?s=20

割引中一覧  
https://docs.google.com/spreadsheets/d/1sWbWfYzzNfUEQFx6u79JteCPJ5IVO_gKdaetLGArSYY/edit#gid=0

## 仕組み
[FANZA（DMM）のAPI](https://affiliate.dmm.com/api/) を利用  
Go言語  
AWS Lambdaで定期実行してます  

## 動かすためには
- FANZA API利用のためDMMアフィリエイトに登録する
- SpreadSheetへの書き込みのためOAuthの設定をする
    - 参考 https://medium.com/veltra-engineering/how-to-use-google-sheets-api-with-golang-9e50ee9e0abc
- SpreadSheetを作成しシートIDを取得しておく
- 定期実行のためAWS Lambdaのセットアップをする
    - 不要ならmain.goからAWSに関する制御部分を削除する
