async function signUp() {
  const fname = document.getElementById('fname').value.trim();
  const lname = document.getElementById('lname').value.trim();
  const email = document.getElementById('email').value.trim();
  const pw1   = document.getElementById('pw1').value;
  const pw2   = document.getElementById('pw2').value;
  const msg   = document.getElementById('msg');

  if (!fname || !lname || !email || !pw1) return showMsg(msg, 'All fields are required.', 'error');
  if (pw1 !== pw2) return showMsg(msg, 'Passwords do not match.', 'error');
  if (pw1.length < 6) return showMsg(msg, 'Password must be at least 6 characters.', 'error');

  try {
    const res = await fetch('/signin', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ first_name: fname, last_name: lname, email, password: pw1 })
    });
    const data = await res.json();

    if (res.ok) {
      showMsg(msg, '✅ Account created! Redirecting to login...', 'success');
      setTimeout(() => window.location.href = 'index.html', 1500);
    } else {
      showMsg(msg, data.error || 'Registration failed.', 'error');
    }
  } catch (e) {
    showMsg(msg, 'Server error.', 'error');
  }
}

function showMsg(el, text, type) {
  el.style.display = 'block';
  el.textContent = text;
  el.className = 'msg-box ' + (type === 'error' ? 'msg-error' : 'msg-success');
}
