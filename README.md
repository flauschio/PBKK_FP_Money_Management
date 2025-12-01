<div align="center">
  <h1>Framework Based Programming - EF234501 (2025)</h1>
</div>

<p align="center">
  <b>Institut Teknologi Sepuluh Nopember</b><br>
  Sepuluh Nopember Institute of Technology
</p>

<p align="center">
  <img src="assets/Badge_ITS.png" width="50%">
</p>

<p align="justify">Source code to <a href="https://www.its.ac.id/informatika/wp-content/uploads/sites/44/2023/11/Module-Handbook-Bachelor-of-Informatics-Program-ITS.pdf">Framework Based Programming (EF234501)</a>'s final project. All solutions were created by <a href="https://github.com/aleahfaa">Iffa Amalia Sabrina</a> and <a href="https://github.com/flauschio">Danendra Ramadhan</a>.</p>

<div align="center">
  <table>
    <thead>
      <tr>
        <th align="center">NRP</th>
        <th align="center">Name</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td align="justify">5025221077</td>
        <td align="justify">Iffa Amalia Sabrina</td>
      </tr>
      <tr>
        <td align="justify">5025231165</td>
        <td align="justify">Danendra Ramadhan</td>
      </tr>
    </tbody>
  </table>
</div>

On behalf of:

**Agus Budi Raharjo, S.Kom., M.Kom., Ph.D.**

---

## Project Overview
Money management

## Task Distribution
1. Iffa Amalia Sabrina
    - Database Model using PostgreSQL
    - Category CRUD
    - Scheduled Transaction CRUD
    - Budget CRUD
2. Danendra Ramadhan
    - Database Implementation using Gorm
    - Authentication (Login, Register, Logout) using JWT
    - Transaction CRUD
    - Account (Create, View, Delete)

## Run Instruction (in Windows)
1. Download PostgreSQL
2. RUn installer
3. Set the password in `.env` file
4. Default port is `5432`
5. Connect to PostgreSQL `psql -U postgres`
6. Inside psql, create database: `CREATE DATABASE finance_manager;`
7. Exit psql `\q`
8. Run migrations `psql -U postgres -d finance_manager -f migrations\001_create_tables.sql`
9. Install Go dependencies `go mod download` and `go mod tidy`
10. Run the application `go run cmd/api/main.go` and it will start on `http://localhost:8080`
11. To view the database `psql -h localhost -p 5432 -U postgres -d finance_manager`
```sql
\l -- to see the list of all database
\dt -- to list all table in finance_manager
\d table_name -- see table structure
SELECT * FROM table_name;
\q -exit
```