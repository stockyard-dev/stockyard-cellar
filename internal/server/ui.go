package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Cellar</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--wine:#8b2252;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.wine{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.wine:hover{border-color:var(--leather)}
.wine-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.wine-name{font-size:.85rem;font-weight:700}
.wine-detail{font-size:.7rem;color:var(--cd);margin-top:.1rem}
.wine-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.wine-notes{font-size:.65rem;color:var(--cm);margin-top:.3rem;font-style:italic;padding:.3rem .5rem;border-left:2px solid var(--wine)}
.wine-actions{display:flex;gap:.3rem;flex-shrink:0}
.type-badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.type-badge.red{border-color:var(--wine);color:var(--wine)}.type-badge.white{border-color:var(--gold);color:var(--gold)}.type-badge.rose{border-color:#d4849a;color:#d4849a}.type-badge.sparkling{border-color:var(--cream);color:var(--cream)}
.stars{color:var(--gold);letter-spacing:1px;font-size:.65rem}
.qty{font-size:.6rem;padding:.1rem .3rem;background:var(--bg3)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.row3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> CELLAR</h1><button class="btn btn-p" onclick="openForm()">+ Add Wine</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search wines..." oninput="render()">
<select class="filter-sel" id="type-filter" onchange="render()"><option value="">All Types</option><option value="red">Red</option><option value="white">White</option><option value="rose">Rose</option><option value="sparkling">Sparkling</option></select>
</div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/wines').then(function(r){return r.json()});items=r.wines||[];renderStats();render();}
function renderStats(){var total=items.length;var bottles=items.reduce(function(s,w){return s+(w.quantity||0)},0);
var regions={};items.forEach(function(w){if(w.region)regions[w.region]=true});
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Labels</div></div><div class="st"><div class="st-v">'+bottles+'</div><div class="st-l">Bottles</div></div><div class="st"><div class="st-v">'+Object.keys(regions).length+'</div><div class="st-l">Regions</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var tf=document.getElementById('type-filter').value;var f=items;
if(tf)f=f.filter(function(w){return(w.type||'').toLowerCase()===tf});
if(q)f=f.filter(function(w){return(w.name||'').toLowerCase().includes(q)||(w.region||'').toLowerCase().includes(q)||(w.producer||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No wines in your cellar.</div>';return;}
var h='';f.forEach(function(w){
h+='<div class="wine"><div class="wine-top"><div style="flex:1">';
h+='<div class="wine-name">'+esc(w.name)+'</div>';
var detail=[];if(w.producer)detail.push(w.producer);if(w.vintage)detail.push(w.vintage);if(w.region)detail.push(w.region);
if(detail.length)h+='<div class="wine-detail">'+esc(detail.join(' &#183; '))+'</div>';
h+='</div><div class="wine-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+w.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+w.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div><div class="wine-meta">';
if(w.type)h+='<span class="type-badge '+w.type.toLowerCase()+'">'+esc(w.type)+'</span>';
if(w.quantity)h+='<span class="qty">'+w.quantity+' bottles</span>';
if(w.rating&&w.rating>0){var s="";for(var x=0;x<5;x++)s+=(x<w.rating?"&#9733;":"&#9734;");h+='<span class="stars">'+s+'</span>';}
if(w.price_paid)h+='<span>$'+(w.price_paid/100).toFixed(2)+'</span>';
h+='</div>';
if(w.notes)h+='<div class="wine-notes">'+esc(w.notes)+'</div>';
h+='</div>';});
document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Remove?'))return;await fetch(A+'/wines/'+id,{method:'DELETE'});load();}
function formHTML(wine){var i=wine||{name:'',type:'red',vintage:0,region:'',producer:'',quantity:1,rating:0,notes:'',price_paid:0};var isEdit=!!wine;
var h='<h2>'+(isEdit?'EDIT':'ADD')+' WINE</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'"></div>';
h+='<div class="row3"><div class="fr"><label>Type</label><select id="f-type"><option value="red"'+(i.type==='red'?' selected':'')+'>Red</option><option value="white"'+(i.type==='white'?' selected':'')+'>White</option><option value="rose"'+(i.type==='rose'?' selected':'')+'>Rose</option><option value="sparkling"'+(i.type==='sparkling'?' selected':'')+'>Sparkling</option></select></div>';
h+='<div class="fr"><label>Vintage</label><input id="f-vintage" type="number" value="'+(i.vintage||'')+'"></div>';
h+='<div class="fr"><label>Quantity</label><input id="f-qty" type="number" value="'+(i.quantity||1)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Region</label><input id="f-region" value="'+esc(i.region)+'"></div>';
h+='<div class="fr"><label>Producer</label><input id="f-producer" value="'+esc(i.producer)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Rating (1-5)</label><input id="f-rating" type="number" min="0" max="5" value="'+(i.rating||0)+'"></div>';
h+='<div class="fr"><label>Price ($)</label><input id="f-price" type="number" step="0.01" value="'+((i.price_paid||0)/100).toFixed(2)+'"></div></div>';
h+='<div class="fr"><label>Notes</label><input id="f-notes" value="'+esc(i.notes)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Add')+'</button></div>';
return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var w=null;for(var j=0;j<items.length;j++){if(items[j].id===id){w=items[j];break;}}if(!w)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(w);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var name=document.getElementById('f-name').value.trim();if(!name){alert('Name required');return;}
var body={name:name,type:document.getElementById('f-type').value,vintage:parseInt(document.getElementById('f-vintage').value)||0,region:document.getElementById('f-region').value.trim(),producer:document.getElementById('f-producer').value.trim(),quantity:parseInt(document.getElementById('f-qty').value)||1,rating:parseInt(document.getElementById('f-rating').value)||0,price_paid:Math.round(parseFloat(document.getElementById('f-price').value||0)*100),notes:document.getElementById('f-notes').value.trim()};
if(editId){await fetch(A+'/wines/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/wines',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
