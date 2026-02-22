# Getting Started

# User Management API (7hunt-be-rest-api)
โปรเจกต์ RESTful API สำหรับบริหารจัดการข้อมูลผู้ใช้งาน (User Management) พัฒนาด้วย Go (Golang), Gin Framework และฐานข้อมูล MongoDB
โปรเจกต์นี้ถูกพัฒนาขึ้นเพื่อรองรับระบบ Login/Register, จัดการผู้ใช้งาน รวมถึงมีการสาธิตแนวทางการเขียน Unit Test ที่ครอบคลุม (Clean Architecture & Testing) และการดีพลอยด้วย Docker

(A simple RESTful API for managing user information using Go (1.24+), Gin, and MongoDB.
This project demonstrates clean architecture, unit testing with mock frameworks (testify), and Docker-based deployment.)

### ความต้องการของระบบ (System Requirement)
* Go Version 1.24 or higher
* Docker Engine or Docker Desktop
* Docker compose

### วิธีการติดตั้งระบบและฐานข้อมูล
1.) ทำการ Clone Repository ด้วยคำสั่ง

    git clone https://github.com/[your-repo]/7hunt-be-rest-api

2.) ตั้งค่า environment variables โดยคัดลอกไฟล์

    cp .env.example .env 
    (หรือสร้างไฟล์ .env และกำหนดค่า MONGODB_URI, JWT_SECRET_KEY, ฯลฯ)

3.) รัน Docker หรือเปิด Docker Desktop

4.) รันคำสั่ง 

    docker-compose up --build -d 
    
5.) จากนั้นรอจนกว่า docker จะ build และ รันระบบเสร็จสิ้นสามารถใช้งานได้ที่ `http://localhost:8080`

### หากต้องการรันระบบผ่าน Go โดยตรง (ไม่ต้องใช้ Docker)
1.) ต้องมั่นใจว่ามี MongoDB รันอยู่ (จะจำลองทดสอบผ่านระบบ Local หรือ MongoDB Atlas ก็ได้) และอัพเดตค่า `MONGODB_URI` ใน `.env`
2.) ใช้คำสั่งรันระบบ:

    go run cmd/main.go

### วิธีการรัน Unit Tests และตรวจสอบ Coverage

โปรเจกต์นี้ใช้ `github.com/stretchr/testify` สำหรับการเขียน Test และ Mocking

1.) ก่อนรัน Test ต้องแน่ใจว่าติดตั้ง module เรียบร้อย
    
    go mod tidy

2.) รัน Test ทั้งหมดในโปรเจกต์

    go test ./... -v

3.) ตรวจสอบ Code Coverage รวม 

    go test ./... -cover

4.) ดูรายงาน Code Coverage แบบกราฟิกผ่าน Browser (HTML)

    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out

---

### Example API requests and expected responses.

| Method | Endpoint              | Description                      |
|:-------|:----------------------|:---------------------------------|
| POST   | /auth/register        | Create a new user (Registration) |
| POST   | /auth/login           | Authenticate user and get Token  |
| GET    | /api/users            | Get all users                    |
| GET    | /api/users/:userId    | Get a specific user by ID        |
| PUT    | /api/users/:userId    | Update a user by ID              |
| DELETE | /api/users/:userId    | Delete a user by ID              |


#### API Endpoint
API นี้ใช้สำหรับจัดการบัญชีผู้ใช้งานทั่วไป (Create, Read, Update, Delete) รวมถึงระบบ Authentication สำหรับการยืนยันตัวตน

1.) Register a New User

Request

    POST http://localhost:8080/auth/register
    Content-Type: application/json

Request Body

    {
        "name": "John Doe",
        "email": "johndoe@example.com",
        "password": "password123"
    }

Expected Response (Success)

    {
        "StatusCode": 200,
        "Message": "Success"
    }

2.) User Login

Request

    POST http://localhost:8080/auth/login
    Content-Type: application/json

Request Body

    {
        "email": "johndoe@example.com",
        "password": "password123"
    }

Expected Response (Success)

    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }

3.) Get All Users (Protected)

Request

    GET http://localhost:8080/api/users
    Authorization: Bearer <your_jwt_token_here>

Expected Response

    {
        "users": [
            {
                "id": "60a7c4f1c9e77c001f8d9b1a",
                "name": "John Doe",
                "email": "johndoe@example.com",
                "created_at": "2026-02-22T08:00:00Z"
            }
        ],
        "user_count": 1
    }

4.) Get User By ID (Protected)

Request

    GET http://localhost:8080/api/users/60a7c4f1c9e77c001f8d9b1a
    Authorization: Bearer <your_jwt_token_here>

Expected Response

    {
        "id": "60a7c4f1c9e77c001f8d9b1a",
        "name": "John Doe",
        "email": "johndoe@example.com",
        "created_at": "2026-02-22T08:00:00Z"
    }

5.) Update User (Protected)

Request

    PUT http://localhost:8080/api/users/60a7c4f1c9e77c001f8d9b1a
    Content-Type: application/json
    Authorization: Bearer <your_jwt_token_here>

Request Body

    {
        "name": "John Doe Updated",
        "email": "johndoe_updated@example.com"
    }

6.) Delete User (Protected)

Request

    DELETE http://localhost:8080/api/users/60a7c4f1c9e77c001f8d9b1a
    Authorization: Bearer <your_jwt_token_here>

---

### Error Example

ระบบมีการใช้ Standard Error Response เพื่อให้ Client สามารถรับทราบถึงข้อผิดพลาดที่เกิดขึ้น

Expected Error Response Structure:

    {
        "StatusCode": 400,
        "Error": {
            "ErrorDesc": "Bad Request",
            "ErrorValidate": "email is required | password is required" // มีเฉพาะในกรณี Validate Error
        }
    }

If `email` already exists on Registration

    POST /auth/register
    {
        "name": "Duplicate User",
        "email": "existing@example.com",
        "password": "password123"
    }

Expected Response

    {
        "StatusCode": 409,
        "Error": {
            "ErrorDesc": "Email already exists"
        }
    }

If Validation Error on Login

    POST /auth/login
    {
        "email": "johndoe@example.com"
        // Missing Password
    }

Expected Response

    {
        "StatusCode": 400,
        "Error": {
            "ErrorDesc": "Validate Failed",
            "ErrorValidate": "Password is a required field"
        }
    }
