# **Golang Server with MongoDB and JSON CRUD Operations**

## **Project Overview**

This project involves the development of a web-based application using **Go (Golang)** and **MongoDB**. The application provides CRUD functionality with a focus on **Filtering**, **Sorting**, and **Pagination** for managing data efficiently. It also includes **structured logging**, **error handling**, **rate limiting** for production-level robustness.

The server will handle **POST** and **GET** requests, interacting with a **MongoDB** database, and integrate a simple **HTML interface** to interact with the server. In addition to basic CRUD operations, the application will allow users to filter, sort, and paginate data. Administrators can send email notifications related to product updates, and the application will handle structured logs and errors effectively.

---

## **Team Members**

- **Dautov Alisher Team Leader**
- **Dautov Alisher Backend**
- **Dautov Alisher Database Manager**

---

## **Project Description**

The goal of this project is to develop a web server in Go that communicates with a MongoDB database. The server will be able to:

- Handle **POST** and **GET** HTTP requests, processing and validating **JSON** data.
- Perform **CRUD operations** (Create, Read, Update, Delete) on data stored in **MongoDB**.
- Integrate **Filtering**, **Sorting**, and **Pagination** mechanisms to display data in an optimized and user-friendly manner.
- Implement **structured logging** using a third-party library like **logrus**, logging all user actions with timestamps.
- Add **error handling** mechanisms for catching and responding to errors throughout the system.
- Implement **rate limiting** to prevent overloading the server with excessive requests.
- Provide an **administrative interface** for sending email notifications to users regarding important updates (e.g., order status promotions, etc.).

---

## **Features**

1. **POST and GET requests** with JSON data.
2. **MongoDB integration** for storing and managing data.
3. **CRUD operations** (Create, Read, Update, Delete).
4. **Filtering, Sorting, and Pagination** capabilities.
5. **Structured logging, error handling, and rate limiting**.

---

## **Technologies Used**

- **Go** - Programming language used for the server.
- **MongoDB** - NoSQL database used for data storage.
- **HTML** - Frontend interface for interacting with the server.
- **Postman** - Tool for sending test requests to the server.

---

## **Setup and Installation**

1. Clone the repository:
    ```bash
    git clone https://github.com/freboy1/Online_Store.git
   
2. Open Project:
    ```bash
    cd Online_Store
3. Run Main.go:
    ```bash
    go run main.go
