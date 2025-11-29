const API_BASE = window.location.origin;
let categories = [];
let transactions = [];
let accounts = [];
document.addEventListener('DOMContentLoaded', () => {
    setupNavigation();
    setupLogout();
    updateUserUI();
    loadData();
    setupForms();
});

function updateUserUI() {
    const userNameEl = document.querySelector('.user-name');
    const userEmailEl = document.querySelector('.user-email');
    const logoutBtn = document.getElementById('logout-btn');
    let user = null;
    try {
        const raw = localStorage.getItem('user');
        if (raw) user = JSON.parse(raw);
    } catch (err) {
        user = null;
    }
    if (user && user.name) {
        if (userNameEl) userNameEl.textContent = user.name;
        if (userEmailEl) userEmailEl.textContent = user.email || '';
        if (logoutBtn) logoutBtn.style.display = '';
    } else {
        if (userNameEl) userNameEl.textContent = 'Anonymous User';
        if (userEmailEl) userEmailEl.textContent = 'No authentication';
        if (logoutBtn) logoutBtn.style.display = 'none';
    }
}

function setupLogout() {
    const btn = document.getElementById('logout-btn');
    if (!btn) return;
    btn.addEventListener('click', (e) => {
        e.preventDefault();
        try {
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
        } catch (err) {
            // ignore
        }
        // redirect to login page
        window.location = '/login';
    });
}
function setupNavigation() {
    document.querySelectorAll('.nav-item').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const section = e.currentTarget.dataset.section;
            showSection(section);
        });
    });
}
function showSection(section) {
    document.querySelectorAll('.nav-item').forEach(link => {
        link.classList.remove('active');
        if (link.dataset.section === section) {
            link.classList.add('active');
        }
    });
    document.querySelectorAll('.section-content').forEach(content => {
        content.classList.add('hidden');
    });
    document.getElementById(`${section}-section`).classList.remove('hidden');
    if (section === 'transactions') {
        renderTransactionsList();
    } else if (section === 'categories') {
        renderCategoriesList();
    }
}
async function loadData() {
    try {
        await Promise.all([
            loadCategories(),
            loadTransactions(),
            loadDashboardStats(),
            loadAccounts()
        ]);
    } catch (error) {
        console.error('Error loading data:', error);
        showNotification('Failed to load data', 'error');
    }
}
async function loadAccounts() {
    try {
        const response = await fetch(`${API_BASE}/api/accounts`);
        if (!response.ok) throw new Error('Failed to fetch accounts');
        accounts = await response.json() || [];
        renderAccountsList(accounts);
        updateAccountDropdown();
    } catch (error) {
        console.error('Error loading accounts:', error);
        showNotification('Failed to load accounts', 'error');
    }
}

function updateAccountDropdown() {
    const select = document.getElementById('transaction-account');
    if (!select) return;
    select.innerHTML = '<option value="">Select an account</option>';
    accounts.forEach(acc => {
        const option = document.createElement('option');
        option.value = acc.id;
        option.textContent = `${acc.bank_name} (${acc.amount < 0 ? '-' : ''}$${Math.abs(acc.amount).toFixed(2)})`;
        select.appendChild(option);
    });
}

function renderAccountsList(accounts) {
    const container = document.getElementById('accounts-list');
    if (!accounts || accounts.length === 0) {
        container.innerHTML = '<div style="padding: 40px; text-align: center; color: var(--text-secondary);">No accounts yet. Create one to get started!</div>';
        return;
    }
    container.innerHTML = accounts.map(acc => `
        <div class="account-card">
            <div class="account-header">
                <div class="account-icon">${acc.bank_name && acc.bank_name.toLowerCase().includes('credit') ? '💳' : '💰'}</div>
                <div>
                    <div class="account-name">${escapeHtml(acc.bank_name)}</nobr></div>
                    <div class="account-type">${acc.id ? 'Account' : ''}</div>
                </div>
            </div>
            <div style="display:flex; justify-content:space-between; align-items:center;">
                <div class="account-balance ${acc.amount < 0 ? 'text-red-600' : ''}">${acc.amount < 0 ? '-' : ''}$${Math.abs(acc.amount).toFixed(2)}</div>
                <div>
                    <button onclick="deleteAccount(${acc.id})" class="btn-icon" title="Delete">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

function openAccountModal() {
    document.getElementById('account-modal').classList.add('show');
}
function closeAccountModal() {
    document.getElementById('account-modal').classList.remove('show');
    document.getElementById('account-form').reset();
    document.getElementById('account-id').value = '';
    document.getElementById('account-modal-title').textContent = 'New Account';
}

async function saveAccount(e) {
    e.preventDefault();
    const id = document.getElementById('account-id').value;
    const name = document.getElementById('account-name').value.trim();
    const amount = parseFloat(document.getElementById('account-amount').value);
    if (!name) { showNotification('Please enter an account name', 'error'); return; }
    if (isNaN(amount)) { showNotification('Please enter a valid amount', 'error'); return; }
    try {
        const response = await fetch(`${API_BASE}/api/accounts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ bank_name: name, amount })
        });
        if (!response.ok) throw new Error('Failed to save account');
        await loadAccounts();
        closeAccountModal();
        showNotification('Account created successfully');
    } catch (error) {
        console.error('Error saving account:', error);
        showNotification('Failed to save account', 'error');
    }
}

async function deleteAccount(id) {
    if (!confirm('Are you sure you want to delete this account?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/accounts/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete account');
        await loadAccounts();
        showNotification('Account deleted successfully');
    } catch (error) {
        console.error('Error deleting account:', error);
        showNotification('Failed to delete account', 'error');
    }
}
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE}/api/categories`);
        if (!response.ok) throw new Error('Failed to fetch categories');
        categories = await response.json() || [];
        updateCategoryDropdown();
        renderCategoriesList();
    } catch (error) {
        console.error('Error loading categories:', error);
        showNotification('Failed to load categories', 'error');
    }
}
function updateCategoryDropdown() {
    const select = document.getElementById('transaction-category');
    select.innerHTML = '<option value="">Select a category</option>';
    categories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat.id;
        option.textContent = cat.name;
        select.appendChild(option);
    });
}
function renderCategoriesList() {
    const container = document.getElementById('categories-list');
    if (!categories || categories.length === 0) {
        container.innerHTML = '<div style="padding: 40px; text-align: center; color: var(--text-secondary);">No categories yet. Create one to get started!</div>';
        return;
    }
    container.innerHTML = categories.map(cat => `
        <div class="category-card">
            <div class="category-name">${escapeHtml(cat.name)}</div>
            <div class="category-stats">${getTransactionCount(cat.id)} transactions</div>
            <div class="category-actions">
                <button onclick="editCategory(${cat.id})" class="btn-icon" title="Edit">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                    </svg>
                </button>
                <button onclick="deleteCategory(${cat.id})" class="btn-icon delete" title="Delete">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                </button>
            </div>
        </div>
    `).join('');
}
function getTransactionCount(categoryId) {
    return transactions.filter(t => t.category_id === categoryId).length;
}
async function saveCategory(e) {
    e.preventDefault();
    const id = document.getElementById('category-id').value;
    const name = document.getElementById('category-name').value.trim();
    if (!name) {
        showNotification('Please enter a category name', 'error');
        return;
    }
    try {
        const url = id ? `${API_BASE}/api/categories/${id}` : `${API_BASE}/api/categories`;
        const method = id ? 'PUT' : 'POST';
        const response = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name })
        });
        if (!response.ok) throw new Error('Failed to save category');
        await loadCategories();
        closeCategoryModal();
        showNotification(`Category ${id ? 'updated' : 'created'} successfully`);
    } catch (error) {
        console.error('Error saving category:', error);
        showNotification('Failed to save category', 'error');
    }
}
async function deleteCategory(id) {
    if (!confirm('Are you sure you want to delete this category?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/categories/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Failed to delete category');
        await loadCategories();
        await loadTransactions();
        showNotification('Category deleted successfully');
    } catch (error) {
        console.error('Error deleting category:', error);
        showNotification('Failed to delete category', 'error');
    }
}
function editCategory(id) {
    const category = categories.find(c => c.id === id);
    if (!category) return;
    document.getElementById('category-id').value = category.id;
    document.getElementById('category-name').value = category.name;
    document.getElementById('category-modal-title').textContent = 'Edit Category';
    openCategoryModal();
}
async function loadTransactions() {
    try {
        const response = await fetch(`${API_BASE}/api/transactions`);
        if (!response.ok) throw new Error('Failed to fetch transactions');
        transactions = await response.json() || [];
        renderRecentTransactions();
        renderTransactionsList();
    } catch (error) {
        console.error('Error loading transactions:', error);
        showNotification('Failed to load transactions', 'error');
    }
}
function renderRecentTransactions() {
    const container = document.getElementById('recent-transactions');
    const recent = transactions.slice(0, 5);
    if (!recent.length) {
        container.innerHTML = '<tr><td colspan="4" style="text-align: center; padding: 40px; color: var(--text-secondary);">No transactions yet. Create one to get started!</td></tr>';
        return;
    }
    container.innerHTML = recent.map(t => createTransactionRow(t, false)).join('');
}
function renderTransactionsList() {
    const container = document.getElementById('transactions-list');
    if (!transactions.length) {
        container.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: 40px; color: var(--text-secondary);">No transactions yet. Create one to get started!</td></tr>';
        return;
    }
    container.innerHTML = transactions.map(t => createTransactionRow(t, true)).join('');
}
function createTransactionRow(transaction, showActions = false) {
    const isPositive = transaction.amount >= 0;
    const amountClass = isPositive ? 'text-green-600' : 'text-red-600';
    const amountPrefix = isPositive ? '+' : '';
    const date = new Date(transaction.created_at);
    const formattedDate = date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    return `
        <tr>
            <td><strong>${escapeHtml(transaction.name)}</strong></td>
            <td>${transaction.category_name || 'Uncategorized'}</td>
            <td>${formattedDate}</td>
            <td class="text-right ${amountClass}"><strong>${amountPrefix}$${Math.abs(transaction.amount).toFixed(2)}</strong></td>
            ${showActions ? `
                <td class="text-right">
                    <button onclick="editTransaction(${transaction.id})" class="btn-icon" title="Edit">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                        </svg>
                    </button>
                    <button onclick="deleteTransaction(${transaction.id})" class="btn-icon delete" title="Delete">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                    </button>
                </td>
            ` : ''}
        </tr>
    `;
}
async function saveTransaction(e) {
    e.preventDefault();
    const id = document.getElementById('transaction-id').value;
    const name = document.getElementById('transaction-name').value.trim();
    const amount = parseFloat(document.getElementById('transaction-amount').value);
    const categoryId = document.getElementById('transaction-category').value;
    const accountId = document.getElementById('transaction-account') ? document.getElementById('transaction-account').value : '';
    if (!name) {
        showNotification('Please enter a transaction name', 'error');
        return;
    }
    if (isNaN(amount)) {
        showNotification('Please enter a valid amount', 'error');
        return;
    }
    try {
        const url = id ? `${API_BASE}/api/transactions/${id}` : `${API_BASE}/api/transactions`;
        const method = id ? 'PUT' : 'POST';
        const body = {
            name,
            amount,
            category_id: categoryId ? parseInt(categoryId) : null,
            account_id: accountId ? parseInt(accountId) : null
        };
        const response = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });
        if (!response.ok) throw new Error('Failed to save transaction');
        await loadTransactions();
        await loadDashboardStats();
        closeTransactionModal();
        showNotification(`Transaction ${id ? 'updated' : 'created'} successfully`);
    } catch (error) {
        console.error('Error saving transaction:', error);
        showNotification('Failed to save transaction', 'error');
    }
}
async function deleteTransaction(id) {
    if (!confirm('Are you sure you want to delete this transaction?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/transactions/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Failed to delete transaction');
        await loadTransactions();
        await loadDashboardStats();
        showNotification('Transaction deleted successfully');
    } catch (error) {
        console.error('Error deleting transaction:', error);
        showNotification('Failed to delete transaction', 'error');
    }
}
function editTransaction(id) {
    const transaction = transactions.find(t => t.id === id);
    if (!transaction) return;
    document.getElementById('transaction-id').value = transaction.id;
    document.getElementById('transaction-name').value = transaction.name;
    document.getElementById('transaction-amount').value = transaction.amount;
    document.getElementById('transaction-category').value = transaction.category_id || '';
    const accSelect = document.getElementById('transaction-account');
    if (accSelect) accSelect.value = transaction.account_id || '';
    document.getElementById('transaction-modal-title').textContent = 'Edit Transaction';
    openTransactionModal();
}
async function loadDashboardStats() {
    try {
        const response = await fetch(`${API_BASE}/api/dashboard/stats`);
        if (!response.ok) throw new Error('Failed to fetch stats');
        const stats = await response.json();
        document.getElementById('total-income').textContent = `$${stats.total_income.toFixed(2)}`;
        document.getElementById('total-expenses').textContent = `$${stats.total_expenses.toFixed(2)}`;
        document.getElementById('net-balance').textContent = `$${stats.balance.toFixed(2)}`;
        document.getElementById('transaction-count').textContent = stats.transaction_count;
    } catch (error) {
        console.error('Error loading stats:', error);
    }
}
function openTransactionModal() {
    document.getElementById('transaction-modal').classList.add('show');
}
function closeTransactionModal() {
    document.getElementById('transaction-modal').classList.remove('show');
    document.getElementById('transaction-form').reset();
    document.getElementById('transaction-id').value = '';
    document.getElementById('transaction-modal-title').textContent = 'New Transaction';
}
function openCategoryModal() {
    document.getElementById('category-modal').classList.add('show');
}
function closeCategoryModal() {
    document.getElementById('category-modal').classList.remove('show');
    document.getElementById('category-form').reset();
    document.getElementById('category-id').value = '';
    document.getElementById('category-modal-title').textContent = 'New Category';
}
function setupForms() {
    document.getElementById('transaction-form').addEventListener('submit', saveTransaction);
    document.getElementById('category-form').addEventListener('submit', saveCategory);
    const accountForm = document.getElementById('account-form');
    if (accountForm) accountForm.addEventListener('submit', saveAccount);
}
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return String(text).replace(/[&<>"']/g, m => map[m]);
}
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 12px 20px;
        background: ${type === 'error' ? '#dc2626' : '#16a34a'};
        color: white;
        border-radius: 4px;
        font-size: 14px;
        font-weight: 500;
        z-index: 10000;
        animation: slideIn 0.2s ease;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    `;
    notification.textContent = message;
    document.body.appendChild(notification);
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.2s ease';
        setTimeout(() => notification.remove(), 200);
    }, 3000);
}