package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-cellar/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/bottles",s.list)
s.mux.HandleFunc("POST /api/bottles",s.create)
s.mux.HandleFunc("GET /api/bottles/{id}",s.get)
s.mux.HandleFunc("PUT /api/bottles/{id}",s.update)
s.mux.HandleFunc("DELETE /api/bottles/{id}",s.del)
s.mux.HandleFunc("GET /api/stats",s.stats)
s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/cellar/"})})
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){
    q:=r.URL.Query().Get("q")
    filters:=map[string]string{}
    if v:=r.URL.Query().Get("type");v!=""{filters["type"]=v}
    if q!=""||len(filters)>0{wj(w,200,map[string]any{"bottles":oe(s.db.Search(q,filters))});return}
    wj(w,200,map[string]any{"bottles":oe(s.db.List())})
}
func(s *Server)create(w http.ResponseWriter,r *http.Request){if s.limits.MaxItems>0{if s.db.Count()>=s.limits.MaxItems{we(w,402,"Free tier limit reached. Upgrade at https://stockyard.dev/cellar/");return}};var e store.Bottle;json.NewDecoder(r.Body).Decode(&e);if e.Name==""{we(w,400,"name required");return};s.db.Create(&e);wj(w,201,s.db.Get(e.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){e:=s.db.Get(r.PathValue("id"));if e==nil{we(w,404,"not found");return};wj(w,200,e)}
func(s *Server)update(w http.ResponseWriter,r *http.Request){
    existing:=s.db.Get(r.PathValue("id"));if existing==nil{we(w,404,"not found");return}
    var patch store.Bottle;json.NewDecoder(r.Body).Decode(&patch);patch.ID=existing.ID;patch.CreatedAt=existing.CreatedAt
    if patch.Name==""{
patch.Name=existing.Name};if patch.Type==""{
patch.Type=existing.Type};if patch.Region==""{
patch.Region=existing.Region};if patch.Producer==""{
patch.Producer=existing.Producer};if patch.Notes==""{
patch.Notes=existing.Notes}
    if patch.Vintage==0{patch.Vintage=existing.Vintage};if patch.Quantity==0{patch.Quantity=existing.Quantity};if patch.Rating==0{patch.Rating=existing.Rating};if patch.PricePaid==0{patch.PricePaid=existing.PricePaid}
    s.db.Update(&patch);wj(w,200,s.db.Get(patch.ID))
}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"status":"ok","service":"cellar","bottles":s.db.Count()})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
