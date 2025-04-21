# todoms

## 概要

todosは、Go言語で開発されたTODO管理システムです。ユーザー認証機能を備えたREST APIを提供し、個人のタスク管理を簡単かつ効率的に行うことができます。

## 主な機能

- ユーザー登録・ログイン機能（JWT認証）
- TODOアイテムの作成・取得・更新・削除（CRUD操作）
- タスクの完了状態管理
- 期限日の設定と追跡

## 技術スタック

- **バックエンド**: Go言語
- **Webフレームワーク**: Echo
- **データベース**: PostgreSQL
- **認証**: JWT（JSON Web Token）
- **ロギング**: Zap
- **コンテナ化**: Docker + Docker Compose

## プロジェクト構造

```
.
├── config/            - アプリケーション設定
├── controller/        - HTTPリクエスト処理とルーティング
├── handler/           - 認証ハンドラーと共通処理
├── migration/         - データベースマイグレーションファイル
├── model/             - データモデルと構造体定義
├── repository/        - データアクセス層
├── service/           - ビジネスロジック
├── docker-compose.yml - Docker環境設定
├── main.go            - アプリケーションのエントリーポイント
└── Makefile           - ビルドと実行のためのコマンド
```

## セットアップ方法

### 前提条件

- Go 1.16以上
- Docker および Docker Compose
- Make（オプション）

### 環境構築

1. リポジトリをクローン:
   ```
   git clone [リポジトリURL]
   cd todoms
   ```

2. 依存パッケージのインストール:
   ```
   go mod download
   ```

3. Docker Composeでデータベース起動:
   ```
   docker-compose up -d
   ```

4. アプリケーション実行:
   ```
   go run main.go
   ```
   または、Makefileを使用:
   ```
   make run
   ```

### 環境変数

- `PORT`: APIサーバーのポート番号（デフォルト: 8080）
- `JWT_SECRET`: JWT署名用の秘密キー
- データベース接続情報（docker-compose.ymlで設定）

## API仕様

詳細なAPI仕様は[API_SPEC.md](API_SPEC.md)を参照してください。

### 認証エンドポイント

- `POST /api/auth/signup` - 新規ユーザー登録
- `POST /api/auth/login` - ログイン（アクセストークン発行）
- `POST /api/auth/refresh` - トークンの更新

### TODOエンドポイント（要認証）

- `GET /api/todos` - すべてのTODOアイテムを取得
- `GET /api/todos/:id` - 特定のTODOアイテムを取得
- `POST /api/todos` - 新しいTODOアイテムを作成
- `PUT /api/todos/:id` - 既存のTODOアイテムを更新
- `DELETE /api/todos/:id` - TODOアイテムを削除

## テスト

テストを実行するには:

```
go test -v ./...
```

または:

```
make test
```

## ライセンス

本プロジェクトはLICENSEファイルに記載されたライセンスの下で公開されています。
