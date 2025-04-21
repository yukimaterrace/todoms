# todoms API仕様書

このドキュメントは、todoms REST APIの詳細な仕様を提供します。

## 目次
- [認証エンドポイント](#認証エンドポイント)
  - [ユーザー登録](#ユーザー登録)
  - [ログイン](#ログイン)
  - [トークン更新](#トークン更新)
  - [現在のユーザー情報取得](#現在のユーザー情報取得)
- [TODOエンドポイント](#todoエンドポイント)
  - [全TODOアイテム取得](#全todoアイテム取得)
  - [特定のTODOアイテム取得](#特定のtodoアイテム取得)
  - [新規TODOアイテム作成](#新規todoアイテム作成)
  - [TODOアイテム更新](#todoアイテム更新)
  - [TODOアイテム削除](#todoアイテム削除)
- [エラーレスポンス一覧](#エラーレスポンス一覧)

## 認証エンドポイント

### ユーザー登録

**エンドポイント:** `POST /api/auth/signup`

**説明:** 新規ユーザーを登録します。

**リクエスト:**
```json
{
  "email": "user@example.com", 
  "password": "password123"
}
```

**リクエストパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| email | string | ✓ | ユーザーのメールアドレス (有効なメールアドレス形式) |
| password | string | ✓ | ユーザーのパスワード (6文字以上) |

**レスポンス:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com"
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| id | string | ユーザーの一意識別子 (UUID) |
| email | string | ユーザーのメールアドレス |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 201 | ユーザーが正常に作成された |
| 400 | リクエストボディが無効またはバリデーションエラー |
| 409 | メールアドレスがすでに使用されている |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "409-1",
  "message": "Email already exists"
}
```

### ログイン

**エンドポイント:** `POST /api/auth/login`

**説明:** 既存ユーザーの認証を行い、アクセストークンとリフレッシュトークンを発行します。

**リクエスト:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**リクエストパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| email | string | ✓ | ユーザーのメールアドレス |
| password | string | ✓ | ユーザーのパスワード |

**レスポンス:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| access_token | string | JWTアクセストークン (15分間有効) |
| refresh_token | string | JWTリフレッシュトークン (7日間有効) |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 200 | 認証に成功し、トークンが発行された |
| 400 | リクエストボディが無効またはバリデーションエラー |
| 401 | 無効なメールアドレスまたはパスワード |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "401-1",
  "message": "Invalid email or password"
}
```

### トークン更新

**エンドポイント:** `POST /api/auth/refresh`

**説明:** リフレッシュトークンを使用して新しいアクセストークンとリフレッシュトークンのペアを取得します。

**リクエスト:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**リクエストパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| refresh_token | string | ✓ | 有効なリフレッシュトークン |

**レスポンス:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| access_token | string | 新しいJWTアクセストークン |
| refresh_token | string | 新しいJWTリフレッシュトークン |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 200 | トークンの更新に成功 |
| 400 | リクエストボディが無効またはバリデーションエラー |
| 401 | 無効なトークン、期限切れのトークン、または無効なトークンタイプ |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "401-4",
  "message": "Token expired"
}
```

### 現在のユーザー情報取得

**エンドポイント:** `GET /api/auth/me`

**説明:** 現在認証されているユーザーの情報を取得します。

**認証:** 必要（Authorization: Bearer {access_token}）

**リクエスト:** リクエストボディなし

**レスポンス:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "user@example.com"
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| id | string | ユーザーの一意識別子 (UUID) |
| email | string | ユーザーのメールアドレス |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 200 | ユーザー情報の取得に成功 |
| 401 | 認証トークンがない、無効、または期限切れ |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "401-2",
  "message": "Missing authorization header"
}
```

## TODOエンドポイント

### 全TODOアイテム取得

**エンドポイント:** `GET /api/todos`

**説明:** 認証されたユーザーの全てのTODOアイテムを取得します。

**認証:** 必要（Authorization: Bearer {access_token}）

**リクエスト:** リクエストボディなし

**レスポンス:**
```json
{
  "todos": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "買い物に行く",
      "description": "牛乳とパンを購入する",
      "dueDate": "2025-05-01T15:00:00Z",
      "isCompleted": false,
      "createdAt": "2025-04-20T10:30:00Z",
      "updatedAt": "2025-04-20T10:30:00Z"
    },
    {
      "id": "223e4567-e89b-12d3-a456-426614174001",
      "title": "レポート作成",
      "description": null,
      "dueDate": null,
      "isCompleted": true,
      "createdAt": "2025-04-19T14:20:00Z",
      "updatedAt": "2025-04-20T09:15:00Z"
    }
  ]
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| todos | array | TODOアイテムの配列 |
| todos[].id | string | TODOアイテムの一意識別子 (UUID) |
| todos[].title | string | TODOアイテムのタイトル |
| todos[].description | string \| null | TODOアイテムの説明 (オプション) |
| todos[].dueDate | string \| null | 期限日時 (ISO8601形式、オプション) |
| todos[].isCompleted | boolean | 完了状態 |
| todos[].createdAt | string | 作成日時 (ISO8601形式) |
| todos[].updatedAt | string | 最終更新日時 (ISO8601形式) |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 200 | TODOアイテムの取得に成功 |
| 401 | 認証トークンがない、無効、または期限切れ |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "401-5",
  "message": "Invalid token"
}
```

### 特定のTODOアイテム取得

**エンドポイント:** `GET /api/todos/:id`

**説明:** 特定のTODOアイテムを取得します。IDで指定されたアイテムが認証されたユーザーのものである必要があります。

**認証:** 必要（Authorization: Bearer {access_token}）

**パスパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| id | string | ✓ | TODOアイテムの一意識別子 (UUID) |

**リクエスト:** リクエストボディなし

**レスポンス:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "買い物に行く",
  "description": "牛乳とパンを購入する",
  "dueDate": "2025-05-01T15:00:00Z",
  "isCompleted": false,
  "createdAt": "2025-04-20T10:30:00Z",
  "updatedAt": "2025-04-20T10:30:00Z"
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| id | string | TODOアイテムの一意識別子 (UUID) |
| title | string | TODOアイテムのタイトル |
| description | string \| null | TODOアイテムの説明 (オプション) |
| dueDate | string \| null | 期限日時 (ISO8601形式、オプション) |
| isCompleted | boolean | 完了状態 |
| createdAt | string | 作成日時 (ISO8601形式) |
| updatedAt | string | 最終更新日時 (ISO8601形式) |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 200 | TODOアイテムの取得に成功 |
| 400 | 無効なTODO ID形式 |
| 401 | 認証トークンがない、無効、または期限切れ |
| 403 | TODOアイテムにアクセスする権限がない |
| 404 | 指定されたIDのTODOアイテムが見つからない |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "404-1",
  "message": "Todo not found"
}
```

### 新規TODOアイテム作成

**エンドポイント:** `POST /api/todos`

**説明:** 認証されたユーザーの新しいTODOアイテムを作成します。

**認証:** 必要（Authorization: Bearer {access_token}）

**リクエスト:**
```json
{
  "title": "買い物に行く",
  "description": "牛乳とパンを購入する",
  "dueDate": "2025-05-01T15:00:00Z"
}
```

**リクエストパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| title | string | ✓ | TODOアイテムのタイトル |
| description | string | ✗ | TODOアイテムの説明 (オプション) |
| dueDate | string | ✗ | 期限日時 (ISO8601形式、オプション) |

**レスポンス:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "買い物に行く",
  "description": "牛乳とパンを購入する",
  "dueDate": "2025-05-01T15:00:00Z",
  "isCompleted": false,
  "createdAt": "2025-04-20T10:30:00Z",
  "updatedAt": "2025-04-20T10:30:00Z"
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| id | string | TODOアイテムの一意識別子 (UUID) |
| title | string | TODOアイテムのタイトル |
| description | string \| null | TODOアイテムの説明 (オプション) |
| dueDate | string \| null | 期限日時 (ISO8601形式、オプション) |
| isCompleted | boolean | 完了状態 (新規作成時は常にfalse) |
| createdAt | string | 作成日時 (ISO8601形式) |
| updatedAt | string | 最終更新日時 (ISO8601形式) |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 201 | TODOアイテムの作成に成功 |
| 400 | リクエストボディが無効またはバリデーションエラー |
| 401 | 認証トークンがない、無効、または期限切れ |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "400-2",
  "message": "Validation failed: Key: 'CreateTodoRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

### TODOアイテム更新

**エンドポイント:** `PUT /api/todos/:id`

**説明:** 特定のTODOアイテムを更新します。IDで指定されたアイテムが認証されたユーザーのものである必要があります。

**認証:** 必要（Authorization: Bearer {access_token}）

**パスパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| id | string | ✓ | 更新するTODOアイテムの一意識別子 (UUID) |

**リクエスト:**
```json
{
  "title": "買い物に行く（更新）",
  "description": "牛乳、パン、卵を購入する",
  "dueDate": "2025-05-02T15:00:00Z",
  "isCompleted": true
}
```

**リクエストパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| title | string | ✓ | TODOアイテムの新しいタイトル |
| description | string | ✗ | TODOアイテムの新しい説明 (オプション) |
| dueDate | string | ✗ | 新しい期限日時 (ISO8601形式、オプション) |
| isCompleted | boolean | ✓ | 新しい完了状態 |

**レスポンス:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "買い物に行く（更新）",
  "description": "牛乳、パン、卵を購入する",
  "dueDate": "2025-05-02T15:00:00Z",
  "isCompleted": true,
  "createdAt": "2025-04-20T10:30:00Z",
  "updatedAt": "2025-04-20T11:45:00Z"
}
```

**レスポンスフィールド:**
| フィールド | 型 | 説明 |
|----------|------|------------|
| id | string | TODOアイテムの一意識別子 (UUID) |
| title | string | 更新後のTODOアイテムのタイトル |
| description | string \| null | 更新後のTODOアイテムの説明 |
| dueDate | string \| null | 更新後の期限日時 (ISO8601形式) |
| isCompleted | boolean | 更新後の完了状態 |
| createdAt | string | 作成日時 (ISO8601形式) |
| updatedAt | string | 最終更新日時 (ISO8601形式) |

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 200 | TODOアイテムの更新に成功 |
| 400 | 無効なTODO ID形式またはリクエストボディが無効 |
| 401 | 認証トークンがない、無効、または期限切れ |
| 403 | TODOアイテムにアクセスする権限がない |
| 404 | 指定されたIDのTODOアイテムが見つからない |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "403-1",
  "message": "You don't have permission to access this todo"
}
```

### TODOアイテム削除

**エンドポイント:** `DELETE /api/todos/:id`

**説明:** 特定のTODOアイテムを削除します。IDで指定されたアイテムが認証されたユーザーのものである必要があります。

**認証:** 必要（Authorization: Bearer {access_token}）

**パスパラメータ:**
| パラメータ | 型 | 必須 | 説明 |
|----------|------|---------|------------|
| id | string | ✓ | 削除するTODOアイテムの一意識別子 (UUID) |

**リクエスト:** リクエストボディなし

**レスポンス:** レスポンスボディなし

**ステータスコード:**
| コード | 説明 |
|--------|------------|
| 204 | TODOアイテムの削除に成功 |
| 400 | 無効なTODO ID形式 |
| 401 | 認証トークンがない、無効、または期限切れ |
| 403 | TODOアイテムにアクセスする権限がない |
| 404 | 指定されたIDのTODOアイテムが見つからない |
| 500 | サーバーエラー |

**エラーレスポンスの例:**
```json
{
  "code": "400-10",
  "message": "Invalid todo ID format"
}
```

## エラーレスポンス一覧

すべてのエラーレスポンスは以下の形式で返されます:

```json
{
  "code": "[HTTPステータスコード]-[エラー番号]",
  "message": "エラーメッセージ"
}
```

### 400 Bad Request
| コード | メッセージ | 説明 |
|--------|-----------|------|
| 400-1 | Invalid request body | リクエストボディが無効 |
| 400-2 | Validation failed | バリデーションエラー （メッセージは具体的なエラー内容により変わる） |
| 400-10 | Invalid todo ID format | 無効なTODO ID形式 |

### 401 Unauthorized
| コード | メッセージ | 説明 |
|--------|-----------|------|
| 401-1 | Invalid email or password | 無効なメールアドレスまたはパスワード |
| 401-2 | Missing authorization header | 認証ヘッダーがない |
| 401-3 | Invalid authorization header format | 無効な認証ヘッダー形式 |
| 401-4 | Token expired | トークンが期限切れ |
| 401-5 | Invalid token | 無効なトークン |
| 401-6 | Invalid token type | 無効なトークンタイプ |

### 403 Forbidden
| コード | メッセージ | 説明 |
|--------|-----------|------|
| 403-1 | You don't have permission to access this todo | このTODOアイテムにアクセスする権限がない |

### 404 Not Found
| コード | メッセージ | 説明 |
|--------|-----------|------|
| 404-1 | Todo not found | 指定されたIDのTODOアイテムが見つからない |

### 409 Conflict
| コード | メッセージ | 説明 |
|--------|-----------|------|
| 409-1 | Email already exists | メールアドレスがすでに使用されている |

### 500 Internal Server Error
| コード | メッセージ | 説明 |
|--------|-----------|------|
| 500-1 | Failed to create user | ユーザーの作成に失敗 |
| 500-2 | Authentication failed | 認証に失敗 |
| 500-3 | Failed to get user claims | ユーザークレームの取得に失敗 |
| 500-4 | Invalid user ID format | 無効なユーザーID形式 |
| 500-10 | Failed to operate | 操作の実行に失敗 |