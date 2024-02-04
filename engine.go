package main

import ("fmt"
        "database/sql"
        "context"
        "io"
        "log"
        "net"
        "time"
)

func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string) (err error) {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    defer func() {
        if err == nil {
            err = tx.Commit()
        }
        if err != nil {
            tx.Rollback()
        }
    }()
    _, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
    if err != nil {
        return err
    }
    // use tx to do more database work
    return nil
}

func handleConn(c net.Conn) {
    defer c.Close()
    for {
        _, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
        if err != nil {
            return // e.g., Client disconnected
        }
        time.Sleep(1 * time.Second)
    }
}

func main() {

    fmt.Println("ValidatorEngine v0.1, started.")

    listener, err := net.Listen("tcp", "localhost:8000")
    if err != nil {
        log.Fatal(err)
    }
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Print(err) // e.g. connection aborted
            continue
        }
        handleConn(conn) // Handle one connection at a time
    }
}
