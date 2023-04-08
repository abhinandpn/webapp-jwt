
## WEB-APP JWT authentication

Simple web application written in golang using postgres database

## API Reference

#### Get all items

```http
  GET /api/items
```

| Package     | Link                      | Description                |
| :--------   | :-------                  | :------------------------- |
| `GORM `     | `https://gorm.io/`        | **Required**. Your API key |
| `GIN `      | `https://gin-gonic.com/`  | **Required**. Your API key |
| `CRYPTO `   | `https://pkg.go.dev/golang.org/x/crypto@v0.8.0`                  | **Required**. Your API key |
| `JWT `      | `string`                  | **Required**. Your API key |
| `GODOTENV ` | `string`                  | **Required**. Your API key |

#### Get item

```http
  GET /api/items/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### add(num1, num2)

Takes two numbers and returns the sum.



