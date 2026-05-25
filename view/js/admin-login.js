async function staffLogin() {
  const email = document.getElementById('email').value.trim();
  const pw    = document.getElementById('pw').value;
  const msg   = document.getElementById('msg');

  if (!email || !pw) return showMsg(msg, 'Please fill in all fields.', 'error');

  try {
    const res = await fetch('/staff/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password: pw })
    });
    const data = await res.json();

    if (res.ok) {
      localStorage.setItem('rqms_staff', email);
      localStorage.setItem('rqms_role', data.role);
      showMsg(msg, '✅ Login successful! Redirecting...', 'success');
      setTimeout(() => window.location.href = 'admin-dashboard.html', 1000);
    } else {
      showMsg(msg, data.error || 'Login failed. Check credentials.', 'error');
    }
  } catch(e) {
    showMsg(msg, 'Server error.', 'error');
  }
}

function showMsg(el, text, type) {
  el.style.display = 'block';
  el.textContent = text;
  el.className = 'msg-box ' + (type === 'error' ? 'msg-error' : 'msg-success');
}
