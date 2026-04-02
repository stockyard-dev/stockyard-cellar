package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Bottle struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Vintage int `json:"vintage"`
	Region string `json:"region"`
	Producer string `json:"producer"`
	Quantity int `json:"quantity"`
	Rating int `json:"rating"`
	Notes string `json:"notes"`
	PricePaid int `json:"price_paid"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"cellar.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS bottles(id TEXT PRIMARY KEY,name TEXT NOT NULL,type TEXT DEFAULT 'red',vintage INTEGER DEFAULT 0,region TEXT DEFAULT '',producer TEXT DEFAULT '',quantity INTEGER DEFAULT 1,rating INTEGER DEFAULT 0,notes TEXT DEFAULT '',price_paid INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Bottle)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO bottles(id,name,type,vintage,region,producer,quantity,rating,notes,price_paid,created_at)VALUES(?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Type,e.Vintage,e.Region,e.Producer,e.Quantity,e.Rating,e.Notes,e.PricePaid,e.CreatedAt);return err}
func(d *DB)Get(id string)*Bottle{var e Bottle;if d.db.QueryRow(`SELECT id,name,type,vintage,region,producer,quantity,rating,notes,price_paid,created_at FROM bottles WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Type,&e.Vintage,&e.Region,&e.Producer,&e.Quantity,&e.Rating,&e.Notes,&e.PricePaid,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Bottle{rows,_:=d.db.Query(`SELECT id,name,type,vintage,region,producer,quantity,rating,notes,price_paid,created_at FROM bottles ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Bottle;for rows.Next(){var e Bottle;rows.Scan(&e.ID,&e.Name,&e.Type,&e.Vintage,&e.Region,&e.Producer,&e.Quantity,&e.Rating,&e.Notes,&e.PricePaid,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Bottle)error{_,err:=d.db.Exec(`UPDATE bottles SET name=?,type=?,vintage=?,region=?,producer=?,quantity=?,rating=?,notes=?,price_paid=? WHERE id=?`,e.Name,e.Type,e.Vintage,e.Region,e.Producer,e.Quantity,e.Rating,e.Notes,e.PricePaid,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM bottles WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM bottles`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Bottle{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ?)"
        args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["type"];ok&&v!=""{where+=" AND type=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,type,vintage,region,producer,quantity,rating,notes,price_paid,created_at FROM bottles WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Bottle;for rows.Next(){var e Bottle;rows.Scan(&e.ID,&e.Name,&e.Type,&e.Vintage,&e.Region,&e.Producer,&e.Quantity,&e.Rating,&e.Notes,&e.PricePaid,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
