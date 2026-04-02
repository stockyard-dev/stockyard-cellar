package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-cellar/internal/server";"github.com/stockyard-dev/stockyard-cellar/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="10220"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./cellar-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("cellar: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Cellar — wine and spirits inventory\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("cellar: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
