# job-queueコントローラー

## 概要

Queueに追加したJobを複数のJob実行ノードで実行する。
その際、Jobと実行ノードの優先順位を考慮する。

## 処理の流れ

### Job実行ノードの登録

Webから新しいノードを登録。

### Jobの追加

Jobを作成する際はREST API経由でQueueに追加する。

### Jobのスケジュール

実行可能なノードとJobを探し、実行可能なものがある場合、優先順位を考慮してスケジュールする。

### Jobの実行

実行可能なJobがある場合、Jobを実行する。

### Jobの完了

Jobの実行が完了した場合、完了マークとする。

## Route

### Route/Runner

- /api/v1/runner
  - GET　ランナー一覧を取得
  - POST　ランナーを追加
- /api/v1/runner/:id
  - POST　ランナーを更新

### Route/Runner/Job

- /api/v1/runner/:id/job
  - GET　実行可能なJobを取得
- /api/v1/runner/:id/job/:id
  - POST　JobのStatusを更新
    - Running
    - Completion
    - Error

### Route/Job

- /api/v1/job
  - GET　Job一覧を取得
  - POST　Jobを追加

## 状態

### Runner

Job実行ノード

- ID
- CreatedAt
- UpdatedAt
- Priority
  - -127～127 優先順位
- Status
  - Ready
  - Maintenance
  - Error
- Name

### Job

Queueに追加されたJob

- ID
- CreatedAt
- UpdatedAt
- Priority
  - -127～127 優先順位
- Status
  - Waiting
  - Scheduled
  - Running
  - Completion
  - Error
- Kind
- Option
- Runner
- Name
