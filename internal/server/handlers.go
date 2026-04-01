package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-cellar/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){q:=r.URL.Query().Get("q");typ:=r.URL.Query().Get("type");list,_:=s.db.List(q,typ);if list==nil{list=[]store.Bottle{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var b store.Bottle;json.NewDecoder(r.Body).Decode(&b);if b.Name==""{writeError(w,400,"name required");return};if b.Quantity==0{b.Quantity=1};s.db.Create(&b);writeJSON(w,201,b)}
func(s *Server)handleDrink(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Drink(id);writeJSON(w,200,map[string]string{"status":"enjoyed"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
