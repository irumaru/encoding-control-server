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
実行時にStatusを更新する。

### Jobの完了

Jobの実行が完了した場合、完了マークとする。
完了時にStatusを更新する。

## Route

### Route/Runner

- /api/v1/runner
  - GET　ランナー一覧を取得
  - POST　ランナーを追加
- /api/v1/runner/:id
  - POST　ランナーを更新
  - DELETE　ランナーを削除

### Route/Job

- /api/v1/job
  - GET　Job一覧を取得
  - POST　Jobを追加
- /api/v1/job/:id
  - POST　Jobを更新
  - DELETE　Jobを削除

### Route/Runner/Job

Runnerごとに、スケジュールされたJobを取得

- /api/v1/runner/:id/job
  - GET　スケジュールされたJobを取得

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
  - Drop
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
  - EpgStation
- Option
- RunnerID
- Name
