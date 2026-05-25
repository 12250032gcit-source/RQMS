/* ── Auth guard + profile setup ───────────────────────────── */
window.onload = function () {
  const user = localStorage.getItem('rqms_user');
  if (!user) { window.location.href = 'index.html'; return; }

  const firstName = localStorage.getItem('rqms_first_name') || '';
  const lastName  = localStorage.getItem('rqms_last_name')  || '';
  const email     = localStorage.getItem('rqms_email')      || user;

  // Build nav profile chip
  buildNavProfile(firstName, lastName, email);

  // Pre-fill & lock booking form fields
  prefillAndLock('fname', firstName);
  prefillAndLock('lname', lastName);
  prefillAndLock('email', email);

  // Phone: pre-fill if stored, lock if available
  const storedPhone = localStorage.getItem('rqms_phone') || '';
  if (storedPhone) prefillAndLock('phone', storedPhone);

  loadAll();
};

function prefillAndLock(id, value) {
  const el = document.getElementById(id);
  if (!el || !value) return;
  el.value    = value;
  el.readOnly = true;
  el.classList.add('field-locked');
}

/* ── Nav profile chip ──────────────────────────────────────── */
function buildNavProfile(firstName, lastName, email) {
  const chip = document.getElementById('navProfileChip');
  if (!chip) return;

  // Build initials: first letter of first name + first letter of last name
  let initials = '?';
  if (firstName && lastName)      initials = (firstName[0] + lastName[0]).toUpperCase();
  else if (firstName)             initials = firstName.slice(0,2).toUpperCase();
  else if (lastName)              initials = lastName.slice(0,2).toUpperCase();
  else if (email)                 initials = email[0].toUpperCase();

  const fullName = [firstName, lastName].filter(Boolean).join(' ') || email;

  chip.innerHTML = `
    <div class="profile-chip" onclick="toggleProfileDropdown()" tabindex="0"
         onkeydown="if(event.key==='Enter'||event.key===' ')toggleProfileDropdown()">
      <div class="avatar">${initials}</div>
      <span class="chip-name">${firstName || email}</span>
      <svg class="chevron" viewBox="0 0 20 20" fill="currentColor" width="14" height="14">
        <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd"/>
      </svg>
    </div>
    <div class="profile-dropdown" id="profileDropdown">
      <div class="profile-header">
        <div class="avatar avatar-lg">${initials}</div>
        <div class="profile-info">
          <div class="profile-name">${fullName}</div>
          <div class="profile-email">${email}</div>
        </div>
      </div>
      <div class="profile-divider"></div>
      <button class="profile-menu-item" onclick="window.location.href='profile.html'">
        <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16">
          <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/>
        </svg>
        My Profile
      </button>
      <div class="profile-divider"></div>
      <button class="profile-menu-item" onclick="logout();window.location.href='index.html'">
        <svg viewBox="0 0 20 20" fill="currentColor" width="16" height="16">
          <path fill-rule="evenodd" d="M3 4.25A2.25 2.25 0 015.25 2h5.5A2.25 2.25 0 0113 4.25v2a.75.75 0 01-1.5 0v-2a.75.75 0 00-.75-.75h-5.5a.75.75 0 00-.75.75v11.5c0 .414.336.75.75.75h5.5a.75.75 0 00.75-.75v-2a.75.75 0 011.5 0v2A2.25 2.25 0 0110.75 18h-5.5A2.25 2.25 0 013 15.75V4.25z" clip-rule="evenodd"/>
          <path fill-rule="evenodd" d="M19 10a.75.75 0 00-.75-.75H8.704l1.048-1.08a.75.75 0 10-1.04-1.08l-2.5 2.57a.75.75 0 000 1.08l2.5 2.57a.75.75 0 101.04-1.08l-1.047-1.08H18.25A.75.75 0 0019 10z" clip-rule="evenodd"/>
        </svg>
        Sign Out
      </button>
    </div>
  `;
}

function toggleProfileDropdown() {
  const dd = document.getElementById('profileDropdown');
  if (!dd) return;
  dd.classList.toggle('open');
}

// Close dropdown when clicking outside
document.addEventListener('click', function(e) {
  const chip = document.querySelector('.profile-chip');
  const dd   = document.getElementById('profileDropdown');
  if (dd && chip && !chip.contains(e.target) && !dd.contains(e.target)) {
    dd.classList.remove('open');
  }
});

function logout() {
  ['rqms_user','rqms_role','rqms_first_name','rqms_last_name','rqms_email','rqms_phone']
    .forEach(k => localStorage.removeItem(k));
}

/* ── Queue actions ─────────────────────────────────────────── */
function loadAll() { loadUsers(); loadStats(); }

async function submitUser() {
  const fname = document.getElementById('fname').value.trim();
  const lname = document.getElementById('lname').value.trim();
  const email = document.getElementById('email').value.trim();
  const phone = document.getElementById('phone').value.trim();
  const note  = document.getElementById('note').value.trim();
  const time  = document.getElementById('time').value;
  const msg   = document.getElementById('msg');

  if (!fname || !lname || !email || !phone) {
    return showMsg(msg, 'Please fill in all required fields (name, email, phone).', 'error');
  }

  // Save phone for future sessions
  localStorage.setItem('rqms_phone', phone);

  try {
    const res  = await fetch('/user', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ first_name: fname, last_name: lname, email, phone, note, time, status: 'waiting' })
    });
    const data = await res.json();
    if (res.ok) {
      showMsg(msg, '✅ You have joined the queue! Staff will notify you.', 'success');
      // Only clear editable fields (note + time), NOT locked fields
      ['note','time'].forEach(id => document.getElementById(id).value = '');
      loadAll();
    } else {
      showMsg(msg, data.error || 'Failed to join queue.', 'error');
    }
  } catch (e) {
    showMsg(msg, 'Server error. Please try again.', 'error');
  }
}

async function loadStats() {
  try {
    const s   = await fetch('/queue/stats').then(r => r.json());
    const box = document.getElementById('statsBox');
    if (!box) return;
    box.innerHTML = `
      <div class="stat-card stat-waiting">⏳ Waiting<br><strong>${s.waiting}</strong></div>
      <div class="stat-card stat-seated">🪑 Seated<br><strong>${s.seated}</strong></div>
      <div class="stat-card stat-done">✅ Done<br><strong>${s.done}</strong></div>
      <div class="stat-card stat-cancelled">❌ Cancelled<br><strong>${s.cancelled}</strong></div>
    `;
  } catch (e) {}
}

async function loadUsers() {
  const container = document.getElementById('tableContainer');
  try {
    const users = await fetch('/users').then(r => r.json()) || [];
    if (!users.length) {
      container.innerHTML = '<div class="empty-state">No bookings yet. Be the first to join!</div>';
      return;
    }

    let html = `<table><thead><tr>
      <th>#</th><th>Name</th><th>Email</th><th>Phone</th>
      <th>Time</th><th>Note</th><th>Status</th><th>Table</th>
    </tr></thead><tbody>`;

    users.forEach(u => {
      const sc = {waiting:'status-waiting',seated:'status-seated',done:'status-done',cancelled:'status-cancelled'}[u.status]||'status-waiting';
      html += `<tr>
        <td>${u.id}</td>
        <td>${u.first_name} ${u.last_name}</td>
        <td>${u.email}</td>
        <td>${u.phone}</td>
        <td>${u.time||'—'}</td>
        <td>${u.note||'—'}</td>
        <td><span class="status-badge ${sc}">${u.status}</span></td>
        <td>${u.table_no||'—'}</td>
      </tr>`;
    });
    container.innerHTML = html + '</tbody></table>';
  } catch (e) {
    container.innerHTML = '<div class="empty-state">Could not load queue.</div>';
  }
}

function showMsg(el, text, type) {
  el.style.display = 'block';
  el.textContent   = text;
  el.className     = 'msg-box ' + (type === 'error' ? 'msg-error' : 'msg-success');
}
