let allTables = [];

window.onload = function () {
  const staff = localStorage.getItem('rqms_staff');
  if (!staff) { window.location.href = 'admin-login.html'; return; }

  const el = document.getElementById('navStaff');
  if (el) {
    const role = localStorage.getItem('rqms_role') || 'staff';
    el.textContent = '👤 ' + staff + ' (' + role + ')';
  }

  loadStats();
  loadQueue();
  loadTables();
};

function adminLogout() {
  localStorage.removeItem('rqms_staff');
  localStorage.removeItem('rqms_role');
}

function showTab(name, btn) {
  document.querySelectorAll('.tab-content').forEach(el => el.style.display = 'none');
  document.querySelectorAll('.tab-btn').forEach(el => el.classList.remove('active'));
  document.getElementById('tab-' + name).style.display = 'block';
  if (btn) btn.classList.add('active');
  if (name === 'tables') loadTables();
  if (name === 'staff')  loadStaff();
  if (name === 'queue')  loadQueue();
}

async function loadStats() {
  try {
    const s = await fetch('/queue/stats').then(r => r.json());
    document.getElementById('statsGrid').innerHTML = `
      <div class="stat-card stat-waiting">⏳ Waiting<br><strong>${s.waiting}</strong></div>
      <div class="stat-card stat-seated">🪑 Seated<br><strong>${s.seated}</strong></div>
      <div class="stat-card stat-done">✅ Done<br><strong>${s.done}</strong></div>
      <div class="stat-card stat-cancelled">❌ Cancelled<br><strong>${s.cancelled}</strong></div>
      <div class="stat-card stat-total">📊 Total<br><strong>${s.total}</strong></div>
    `;
  } catch (e) {}
}

async function loadQueue() {
  const filter    = document.getElementById('filterStatus').value;
  const container = document.getElementById('queueTable');
  try {
    let users = await fetch('/users').then(r => r.json()) || [];
    if (filter) users = users.filter(u => u.status === filter);

    if (users.length === 0) {
      container.innerHTML = '<div class="empty-state">No entries found.</div>';
      return;
    }

    let html = `<table><thead><tr>
      <th>#</th><th>Name</th><th>Phone</th><th>Email</th>
      <th>Time</th><th>Note</th><th>Status</th><th>Table</th><th>Actions</th>
    </tr></thead><tbody>`;

    users.forEach(u => {
      const sc = {waiting:'status-waiting',seated:'status-seated',done:'status-done',cancelled:'status-cancelled'}[u.status]||'status-waiting';
      // Escape table_no for safe inline JS string usage
      const tn = (u.table_no || '').replace(/'/g, "\\'");
      html += `<tr>
        <td>${u.id}</td>
        <td>${escHtml(u.first_name)} ${escHtml(u.last_name)}</td>
        <td>${escHtml(u.phone)}</td>
        <td>${escHtml(u.email)}</td>
        <td>${u.time||'—'}</td>
        <td>${escHtml(u.note||'—')}</td>
        <td><span class="status-badge ${sc}">${u.status}</span></td>
        <td>${escHtml(u.table_no||'—')}</td>
        <td><div class="action-btns">
          ${u.status==='waiting'?`<button class="btn btn-success btn-sm" onclick="openSeatModal(${u.id})">🪑 Seat</button>`:''}
          ${u.status==='seated'?`<button class="btn btn-primary btn-sm" onclick="updateStatus(${u.id},'done','${tn}')">✅ Done</button>`:''}
          ${u.status!=='cancelled'&&u.status!=='done'?`<button class="btn btn-warning btn-sm" onclick="updateStatus(${u.id},'cancelled','${tn}')">❌ Cancel</button>`:''}
          <button class="btn btn-danger btn-sm" onclick="deleteEntry(${u.id})">🗑</button>
        </div></td>
      </tr>`;
    });
    container.innerHTML = html + '</tbody></table>';
  } catch (e) {
    container.innerHTML = '<div class="empty-state">Could not load queue.</div>';
  }
}

// Bug fix: pass the current table_no so the backend can free it on cancel/done
async function updateStatus(id, status, tableNo) {
  await fetch('/queue/status', {
    method:'PUT', headers:{'Content-Type':'application/json'},
    body: JSON.stringify({id, status, table_no: tableNo})
  });
  loadQueue(); loadStats(); loadTables();
}

async function deleteEntry(id) {
  if (!confirm('Delete this entry?')) return;
  await fetch('/user?id='+id, {method:'DELETE'});
  loadQueue(); loadStats();
}

function openSeatModal(userId) {
  const available = allTables.filter(t => t.status === 'available');
  if (!available.length) { alert('No available tables right now.'); return; }
  document.getElementById('seatUserId').value = userId;
  const sel = document.getElementById('seatTableNo');
  sel.innerHTML = '';
  available.forEach(t => {
    const opt = document.createElement('option');
    opt.value = t.table_no;
    opt.textContent = `${t.table_no} (cap: ${t.capacity})`;
    sel.appendChild(opt);
  });
  document.getElementById('seatModal').style.display = 'flex';
}

function closeModal() { document.getElementById('seatModal').style.display = 'none'; }

async function confirmSeat() {
  const id      = parseInt(document.getElementById('seatUserId').value);
  const tableNo = document.getElementById('seatTableNo').value;
  await updateStatus(id, 'seated', tableNo);
  closeModal();
}

async function loadTables() {
  try {
    allTables = await fetch('/tables').then(r => r.json()) || [];
    const grid = document.getElementById('tablesGrid');
    if (!grid) return;
    grid.innerHTML = allTables.map(t => {
      const cls = {available:'table-card-available',occupied:'table-card-occupied',reserved:'table-card-reserved'}[t.status]||'table-card-available';
      const tn = t.table_no.replace(/'/g, "\\'");
      return `<div class="table-card ${cls}">
        <div class="table-number">${escHtml(t.table_no)}</div>
        <div class="table-cap">👥 ${t.capacity}</div>
        <div class="table-status">${t.status}</div>
        <div class="table-card-actions">
          ${t.status!=='available'?`<button class="btn btn-success btn-sm" onclick="setTableStatus('${tn}','available')">Free</button>`:''}
          ${t.status==='available'?`<button class="btn btn-warning btn-sm" onclick="setTableStatus('${tn}','reserved')">Reserve</button>`:''}
        </div>
      </div>`;
    }).join('');
  } catch(e) {}
}

async function setTableStatus(tableNo, status) {
  await fetch('/tables/status', {
    method:'PUT', headers:{'Content-Type':'application/json'},
    body: JSON.stringify({table_no: tableNo, status})
  });
  loadTables();
}

async function loadStaff() {
  const container = document.getElementById('staffTable');
  try {
    const list = await fetch('/staff/list').then(r => r.json()) || [];
    if (!list.length) { container.innerHTML = '<div class="empty-state">No staff yet.</div>'; return; }
    let html = `<table><thead><tr><th>#</th><th>Name</th><th>Email</th><th>Role</th></tr></thead><tbody>`;
    list.forEach(s => {
      html += `<tr><td>${s.id}</td><td>${escHtml(s.first_name)} ${escHtml(s.last_name)}</td><td>${escHtml(s.email)}</td>
        <td><span class="status-badge" style="background:#e9d8fd;color:#553c9a;">${escHtml(s.role)}</span></td></tr>`;
    });
    container.innerHTML = html + '</tbody></table>';
  } catch(e) {
    container.innerHTML = '<div class="empty-state">Could not load staff.</div>';
  }
}

// Prevent XSS from DB values rendered into HTML
function escHtml(s) {
  if (!s) return '';
  return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}
