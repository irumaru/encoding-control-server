# APIインターフェス仕様

## [POST] /api/v1/job/create
### 関数
CreateJob  
### 概要
ジョブの作成  
### 入力
```
{
    "kind": "処理の種類",
    "priority": "優先順位",
    "client": "依頼クライアント"
}
```
### 出力(Job標準出力)
※excution_nodeは空白

## [POST] /api/v1/job/execution
### 関数
ExcutionJob
### 概要
実行可能なJobを探索, jobの状態をassignedへ変更, jobの情報をエンコーダーへ送信
### 入力
```
{
    "kind": "処理の種類",
    "execution_node": "エンコーダー名"
}
```
### 出力(Job標準出力)

## [GET] /api/v1/job/get
### 関数
GetJobs
### 概要
jobの一覧を取得
### 入力
なし
### 出力
```
{
    Job標準出力1,
    Job標準出力2
}
```

## [POST] /api/v1/job/{id}/update
### 関数
UpdateJob
### 概要
job状態の処理内容を更新
### 入力
```
{
    "status": "新たな状態",
    "detail": "省略可能: コメント"
}
```
### 出力(JobDetail標準出力)

## [GET] /api/v1/job/{id}/get
### 関数
GetJob
### 概要
対象のjobの状態を取得
### 入力
なし
### 出力(Job標準出力)

## [GET] /api/v1/job/{id}/get_details
### 関数
GetJobDetails
### 概要
対象のjobの詳細を取得
### 入力
なし
### 出力
```
{
    JobDetail標準出力1,
    JobDetail標準出力2
}
```

## [DELETE] /api/v1/job/{id}/delete
### 関数
DeleteJob
### 概要
対象のJobを削除
### 入力
なし
### 出力
なし


## 標準的な出力
### Job標準出力
```
{
    "id": "ID",
    "created_at": "作成時刻",
    "updated_at": "更新時刻",
    "kind": "種類",
    "priority": "優先順位",
    "status": "状態",
    "client": "依頼者",
    "excution_node": "実行ノード",
}
```
### JobDetail標準出力
```
{
    "id": "ID",
    "job_id": "JobのID",
    "created_at": "作成時刻",
    "updated_at": "更新時刻",
    "status": "状態",
    "detail": "コメント"
}
```
