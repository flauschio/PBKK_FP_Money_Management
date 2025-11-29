(function(){
  const form = document.getElementById('register-form')
  const msg = document.getElementById('msg')
  form.addEventListener('submit', async (e) => {
    e.preventDefault()
    msg.textContent = ''
    const payload = {
      name: document.getElementById('name').value,
      email: document.getElementById('email').value,
      password: document.getElementById('password').value,
    }
    try {
      const res = await fetch('/api/register', {
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
      msg.textContent = 'Registration successful. You can now login.'
      msg.style.color = 'green'
      setTimeout(() => { window.location = '/login' }, 1200)
    } catch (err) {
      msg.textContent = 'Network error'
      msg.style.color = 'red'
    }
  })
})();