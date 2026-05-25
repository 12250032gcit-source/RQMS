async function login() {
  const email = document.getElementById('email').value.trim();
  const pw    = document.getElementById('pw').value;
  const msg   = document.getElementById('msg');

  if (!email || !pw) return showMsg(msg, 'Please fill in all fields.', 'error');

  try {
    const res = await fetch('/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password: pw })
    });
    const data = await res.json();

    if (res.ok) {
      // Store full user profile for pre-filling booking form
      localStorage.setItem('rqms_user', email);
      localStorage.setItem('rqms_role', 'customer');
      localStorage.setItem('rqms_first_name', data.first_name || '');
      localStorage.setItem('rqms_last_name',  data.last_name  || '');
      localStorage.setItem('rqms_email',      data.email      || email);
      showMsg(msg, '✅ Login successful! Redirecting...', 'success');
      setTimeout(() => window.location.href = 'home.html', 1000);
    } else {
      showMsg(msg, data.error || 'Login failed.', 'error');
    }
  } catch (e) {
    showMsg(msg, 'Server error. Please try again.', 'error');
  }
}

function showMsg(el, text, type) {
  el.style.display = 'block';
  el.textContent = text;
  el.className = 'msg-box ' + (type === 'error' ? 'msg-error' : 'msg-success');
}
