(function(){
  const form = document.getElementById('login-form')
  const msg = document.getElementById('msg')
  form.addEventListener('submit', async (e) => {
    e.preventDefault()
    msg.textContent = ''
    const payload = {
      email: document.getElementById('email').value,
      password: document.getElementById('password').value,
    }
    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      const data = await res.json()
      if (!res.ok) {
        msg.textContent = data.error || JSON.stringify(data)
        msg.style.color = 'red'
        return
      }
      // Save tokens and user info then redirect to home
      if (data.access_token) {
        localStorage.setItem('access_token', data.access_token)
      }
      if (data.refresh_token) {
        localStorage.setItem('refresh_token', data.refresh_token)
      }
      if (data.user) {
        try {
          localStorage.setItem('user', JSON.stringify(data.user))
        } catch (err) {
          // ignore storage errors
        }
      }
      msg.textContent = 'Login successful'
      msg.style.color = 'green'
      setTimeout(() => { window.location = '/' }, 700)
    } catch (err) {
      msg.textContent = 'Network error'
      msg.style.color = 'red'
    }
  })
})();