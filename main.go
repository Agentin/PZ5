package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
)

func main() {
    _ = godotenv.Load()
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://postgres:YOUR_PASSWORD@localhost:5432/todo?sslmode=disable"
    }

    db, err := openDB(dsn)
    if err != nil {
        log.Fatalf("openDB error: %v", err)
    }
    defer db.Close()

    repo := NewRepo(db)

    // Демонстрация CreateMany
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    batchTitles := []string{"Задача из batch 1", "Задача из batch 2", "Задача из batch 3"}
    if err := repo.CreateMany(ctx, batchTitles); err != nil {
        log.Printf("CreateMany error: %v", err)
    } else {
        log.Println("Batch insert completed successfully")
    }

    // Демонстрация ListDone (невыполненные задачи)
    ctxListDone, cancelListDone := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancelListDone()

    pendingTasks, err := repo.ListDone(ctxListDone, false)
    if err != nil {
        log.Printf("ListDone (pending) error: %v", err)
    } else {
        fmt.Println("\n=== Pending Tasks ===")
        for _, t := range pendingTasks {
            fmt.Printf("#%d | %-24s | done=%-5v | %s\n", 
                t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
        }
    }

    // Демонстрация FindByID
    ctxFind, cancelFind := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancelFind()

    if task, err := repo.FindByID(ctxFind, 1); err != nil {
        log.Printf("FindByID error: %v", err)
    } else {
        fmt.Printf("\n=== Task #1 Details ===\n")
        fmt.Printf("ID: %d\nTitle: %s\nDone: %v\nCreated: %s\n", 
            task.ID, task.Title, task.Done, task.CreatedAt.Format(time.RFC3339))
    }

    // Демонстрация обычного списка всех задач
    ctxAll, cancelAll := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancelAll()

    allTasks, err := repo.ListTasks(ctxAll)
    if err != nil {
        log.Printf("ListTasks error: %v", err)
    } else {
        fmt.Println("\n=== All Tasks ===")
        for _, t := range allTasks {
            fmt.Printf("#%d | %-24s | done=%-5v | %s\n", 
                t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
        }
    }
}
