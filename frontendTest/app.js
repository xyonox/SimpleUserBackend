const apiBase = document.querySelector('#apiBase');
const message = document.querySelector('#message');
const result = document.querySelector('#result');

function url(path) { return `${apiBase.value.replace(/\/$/, '')}${path}`; }
function showMessage(text, error = false) { message.textContent = text; message.className = `message${error ? ' error' : ''}`; }
function showResult(value) { result.textContent = typeof value === 'string' ? value : JSON.stringify(value, null, 2); }
async function request(path, options = {}) {
  const response = await fetch(url(path), { ...options, credentials: 'include', headers: { 'Content-Type': 'application/json', ...(options.headers || {}) } });
  const text = await response.text();
  let body; try { body = text ? JSON.parse(text) : text; } catch { body = text; }
  if (!response.ok) throw new Error(`${response.status}: ${typeof body === 'string' ? body : JSON.stringify(body)}`);
  return body;
}

document.querySelectorAll('.tab').forEach(tab => tab.addEventListener('click', () => {
  document.querySelectorAll('.tab').forEach(item => item.classList.remove('active'));
  document.querySelectorAll('.form-panel').forEach(item => item.classList.remove('active'));
  tab.classList.add('active'); document.querySelector(`#${tab.dataset.panel}`).classList.add('active');
}));

document.querySelector('#registerPanel').addEventListener('submit', async event => {
  event.preventDefault(); const form = new FormData(event.currentTarget);
  try {
    const body = { name: form.get('name'), password_hash: form.get('password') };
    const data = await request('/user/create', { method: 'POST', body: JSON.stringify(body) });
    showMessage(data || 'Registrierung erfolgreich.'); event.currentTarget.reset();
  } catch (error) { showMessage(error.message, true); }
});

document.querySelector('#loginPanel').addEventListener('submit', async event => {
  event.preventDefault();
  try { await request('/user/login', { method: 'POST' }); showMessage('Test-Session gesetzt.'); document.querySelector('#authButton').click(); }
  catch (error) { showMessage(error.message, true); }
});

document.querySelector('#authButton').addEventListener('click', async () => {
  try { const data = await request('/user/auth'); showResult(data); showMessage('Authentifizierung erfolgreich.'); document.querySelector('.status-pill').classList.add('authenticated'); document.querySelector('#statusText').textContent = 'Authentifiziert'; }
  catch (error) { showResult(error.message); showMessage(error.message, true); }
});
document.querySelector('#usersButton').addEventListener('click', async () => { try { showResult(await request('/users')); } catch (error) { showResult(error.message); showMessage(error.message, true); } });
document.querySelector('#helloButton').addEventListener('click', async () => { try { showResult(await request('/helloworld')); showMessage('Backend erreichbar.'); } catch (error) { showMessage(error.message, true); } });
document.querySelector('#logoutButton').addEventListener('click', () => { document.cookie = 'session_token=; Max-Age=0; path=/'; showMessage('Cookie im Frontend gelöscht. Für echtes Logout zusätzlich Backend-Route einbauen.'); });
