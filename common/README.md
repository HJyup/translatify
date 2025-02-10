# Common Service

### **1. Environment Handling (`utils/env.go`)**


- `EnvString(key string) string` → Retrieves an environment variable or logs an error if missing.

  
  ```go
  value := utils.EnvString("DATABASE_URL")
  ```

### **2. JSON Helpers (`utils/json.go`)**

- `WriteJSON(w http.ResponseWriter, status int, data interface{})` → Sends a JSON response.

  
  ```go
  utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Success"})
  ```
- `ReadJSON(r *http.Request, data interface{}) error` → Parses JSON request body into a struct.

  
  ```go
  var requestBody struct { Name string }
  err := utils.ReadJSON(r, &requestBody)
  ```
- `WriteError(w http.ResponseWriter, status int, message string)` → Sends a JSON error response.

  
  ```go
  utils.WriteError(w, http.StatusBadRequest, "Invalid request")
  ```

### **3. Database Helpers (`utils/scan.go`)**


- `RowScanner` interface → Standardizes database row scanning.

  
  ```go
  type RowScanner interface {
      Scan(dest ...interface{}) error
  }
  ```

