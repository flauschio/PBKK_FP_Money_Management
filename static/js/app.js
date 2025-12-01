const API_BASE = window.location.origin;
let categories = [];
let transactions = [];
let accounts = [];
let budgets = [];
let scheduledTransactions = [];

function authHeaders() {
    const token = localStorage.getItem('access_token');
    return {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
    };
}

document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('access_token');
    if (!token && window.location.pathname === '/') {
        window.location = '/login';
        return;
    }
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
        const refresh = localStorage.getItem('refresh_token');
        // call server to revoke refresh token
        if (refresh) {
            fetch('/api/logout', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refresh })
            }).catch(() => {});
        }
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
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
    } else if (section === 'accounts') {
        loadAccounts();
    } else if (section === 'budgets') {
        loadBudgets();
    } else if (section === 'scheduled') {
        loadScheduledTransactions();
    }
}
async function loadData() {
    try {
        await Promise.all([
            loadCategories(),
            loadTransactions(),
            loadDashboardStats(),
            loadAccounts(),
            loadBudgets(),
            loadScheduledTransactions()
        ]);
    } catch (error) {
        console.error('Error loading data:', error);
    }
}
async function loadAccounts() {
    try {
        const response = await fetch(`${API_BASE}/api/accounts`, {
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to fetch accounts');
        accounts = await response.json() || [];
        renderAccountsList(accounts);
        updateAccountDropdown();
        updateScheduledAccountDropdown();
    } catch (error) {
        console.error('Error loading accounts:', error);
    }
}
function renderAccountsList(accountsData) {
    const container = document.getElementById('accounts-list');
    if (!container) return;
    if (!accountsData || accountsData.length === 0) {
        container.innerHTML = '<p style="padding:20px; text-align:center; color:var(--text-secondary);">No accounts yet. Create one to get started!</p>';
        return;
    }
    container.innerHTML = accountsData.map(account => `
        <div class="account-card">
            <div class="account-header">
                <div class="account-icon">🏦</div>
                <div>
                    <div class="account-name">${account.bank_name}</div>
                </div>
            </div>
            <div class="account-balance">$${account.amount.toFixed(2)}</div>
            <div style="margin-top:12px; display:flex; gap:8px;">
                <button onclick="deleteAccount(${account.id})" class="btn-icon delete">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>
        </div>
    `).join('');
}
async function deleteAccount(id) {
    if (!confirm('Are you sure you want to delete this account?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/accounts/${id}`, {
            method: 'DELETE',
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to delete account');
        await loadAccounts();
        await loadDashboardStats();
    } catch (error) {
        console.error('Error deleting account:', error);
        alert('Failed to delete account');
    }
}
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE}/api/categories`, {
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to fetch categories');
        categories = await response.json() || [];
        updateCategoryDropdown();
        updateBudgetCategoryDropdown();
        updateScheduledCategoryDropdown();
        renderCategoriesList();
    } catch (error) {
        console.error('Error loading categories:', error);
    }
}
function renderCategoriesList() {
    const container = document.getElementById('categories-list');
    if (!container) return;
    if (categories.length === 0) {
        container.innerHTML = '<p style="padding:20px; text-align:center; color:var(--text-secondary);">No categories yet. Create one to organize your transactions!</p>';
        return;
    }
    container.innerHTML = categories.map(category => `
        <div class="category-card">
            <div class="category-name">${category.name}</div>
            <div class="category-actions">
                <button onclick="editCategory(${category.id})" class="btn-icon">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button onclick="deleteCategory(${category.id})" class="btn-icon delete">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </div>
        </div>
    `).join('');
}
function updateCategoryDropdown() {
    const select = document.getElementById('transaction-category');
    if (!select) return;
    select.innerHTML = '<option value="">None</option>' + 
        categories.map(cat => `<option value="${cat.id}">${cat.name}</option>`).join('');
}
function updateBudgetCategoryDropdown() {
    const select = document.getElementById('budget-category');
    if (!select) return;
    select.innerHTML = '<option value="">Select a category</option>' + 
        categories.map(cat => `<option value="${cat.id}">${cat.name}</option>`).join('');
}
function updateScheduledCategoryDropdown() {
    const select = document.getElementById('scheduled-category');
    if (!select) return;
    select.innerHTML = '<option value="">None</option>' + 
        categories.map(cat => `<option value="${cat.id}">${cat.name}</option>`).join('');
}
function updateAccountDropdown() {
    const select = document.getElementById('transaction-account');
    if (!select) return;
    select.innerHTML = '<option value="">None</option>' + 
        accounts.map(acc => `<option value="${acc.id}">${acc.bank_name}</option>`).join('');
}
function updateScheduledAccountDropdown() {
    const select = document.getElementById('scheduled-account');
    if (!select) return;
    select.innerHTML = '<option value="">None</option>' + 
        accounts.map(acc => `<option value="${acc.id}">${acc.bank_name}</option>`).join('');
}
function editCategory(id) {
    const category = categories.find(c => c.id === id);
    if (!category) return;
    document.getElementById('category-id').value = category.id;
    document.getElementById('category-name').value = category.name;
    document.getElementById('category-modal-title').textContent = 'Edit Category';
    document.getElementById('category-modal').classList.add('show');
}
async function deleteCategory(id) {
    if (!confirm('Are you sure you want to delete this category?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/categories/${id}`, {
            method: 'DELETE',
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to delete category');
        await loadCategories();
    } catch (error) {
        console.error('Error deleting category:', error);
        alert('Failed to delete category');
    }
}
async function loadTransactions() {
    try {
        const response = await fetch(`${API_BASE}/api/transactions`, {
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to fetch transactions');
        transactions = await response.json() || [];
        renderTransactionsList();
        renderRecentTransactions();
    } catch (error) {
        console.error('Error loading transactions:', error);
    }
}
function renderRecentTransactions() {
    const tbody = document.getElementById('recent-transactions');
    if (!tbody) return;
    const recent = transactions.slice(0, 5);
    if (recent.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" style="text-align:center; padding:20px; color:var(--text-secondary);">No transactions yet</td></tr>';
        return;
    }
    tbody.innerHTML = recent.map(t => `
        <tr>
            <td>${t.name}</td>
            <td>${t.category?.name || 'Uncategorized'}</td>
            <td>${new Date(t.created_at).toLocaleDateString()}</td>
            <td class="text-right ${t.amount >= 0 ? 'text-green-600' : 'text-red-600'}">
                ${t.amount >= 0 ? '+' : ''}$${Math.abs(t.amount).toFixed(2)}
            </td>
        </tr>
    `).join('');
}
function renderTransactionsList() {
    const tbody = document.getElementById('transactions-list');
    if (!tbody) return;
    if (transactions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; padding:20px; color:var(--text-secondary);">No transactions yet. Create one to get started!</td></tr>';
        return;
    }
    tbody.innerHTML = transactions.map(t => `
        <tr>
            <td>${t.name}</td>
            <td>${t.category?.name || 'Uncategorized'}</td>
            <td>${new Date(t.created_at).toLocaleDateString()}</td>
            <td class="text-right ${t.amount >= 0 ? 'text-green-600' : 'text-red-600'}">
                ${t.amount >= 0 ? '+' : ''}$${Math.abs(t.amount).toFixed(2)}
            </td>
            <td class="text-right">
                <button onclick="editTransaction(${t.id})" class="btn-icon">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button onclick="deleteTransaction(${t.id})" class="btn-icon delete">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </td>
        </tr>
    `).join('');
}
function editTransaction(id) {
    const transaction = transactions.find(t => t.id === id);
    if (!transaction) return;
    document.getElementById('transaction-id').value = transaction.id;
    document.getElementById('transaction-name').value = transaction.name;
    document.getElementById('transaction-amount').value = transaction.amount;
    document.getElementById('transaction-category').value = transaction.category_id || '';
    document.getElementById('transaction-account').value = transaction.account_id || '';
    document.getElementById('transaction-modal-title').textContent = 'Edit Transaction';
    document.getElementById('transaction-modal').classList.add('show');
}
async function deleteTransaction(id) {
    if (!confirm('Are you sure you want to delete this transaction?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/transactions/${id}`, {
            method: 'DELETE',
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to delete transaction');
        await loadData();
    } catch (error) {
        console.error('Error deleting transaction:', error);
        alert('Failed to delete transaction');
    }
}
async function loadBudgets() {
    try {
        const response = await fetch(`${API_BASE}/api/budgets`, {
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to fetch budgets');
        budgets = await response.json() || [];
        renderBudgetsList();
    } catch (error) {
        console.error('Error loading budgets:', error);
    }
}
function renderBudgetsList() {
    const tbody = document.getElementById('budgets-list');
    if (!tbody) return;
    if (budgets.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" style="text-align:center; padding:20px; color:var(--text-secondary);">No budgets set. Create one to track your spending!</td></tr>';
        return;
    }
    tbody.innerHTML = budgets.map(b => {
        const categoryName = categories.find(c => c.id === b.category_id)?.name || 'Unknown';
        const progressColor = b.percentage > 0 ? 'green' : 'red';
        const progressWidth = Math.min(100, Math.max(0, 100 - b.percentage));
        return `
        <tr>
            <td>${categoryName}</td>
            <td>$${b.amount.toFixed(2)}</td>
            <td class="text-red-600">$${b.spent.toFixed(2)}</td>
            <td class="${b.remaining >= 0 ? 'text-green-600' : 'text-red-600'}">$${b.remaining.toFixed(2)}</td>
            <td>
                <div style="display: flex; align-items: center; gap: 8px;">
                    <div style="flex: 1; height: 8px; background: var(--bg-tertiary); border-radius: 4px; overflow: hidden;">
                        <div style="height: 100%; width: ${progressWidth}%; background: ${progressColor}; transition: width 0.3s;"></div>
                    </div>
                    <span style="font-size: 12px; color: var(--text-secondary); min-width: 45px;">${b.percentage.toFixed(0)}%</span>
                </div>
            </td>
            <td style="text-transform: capitalize;">${b.criteria}</td>
            <td class="text-right">
                <button onclick="editBudget(${b.id})" class="btn-icon">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button onclick="deleteBudget(${b.id})" class="btn-icon delete">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </td>
        </tr>
    `}).join('');
}
function editBudget(id) {
    const budget = budgets.find(b => b.id === id);
    if (!budget) return;
    document.getElementById('budget-id').value = budget.id;
    document.getElementById('budget-category').value = budget.category_id;
    document.getElementById('budget-amount').value = budget.amount;
    document.getElementById('budget-criteria').value = budget.criteria;
    document.getElementById('budget-modal-title').textContent = 'Edit Budget';
    document.getElementById('budget-modal').classList.add('show');
}
async function deleteBudget(id) {
    if (!confirm('Are you sure you want to delete this budget?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/budgets/${id}`, {
            method: 'DELETE',
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to delete budget');
        await loadBudgets();
    } catch (error) {
        console.error('Error deleting budget:', error);
        alert('Failed to delete budget');
    }
}
async function loadScheduledTransactions() {
    try {
        const response = await fetch(`${API_BASE}/api/scheduled`, {
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to fetch scheduled transactions');
        scheduledTransactions = await response.json() || [];
        renderScheduledTransactionsList();
    } catch (error) {
        console.error('Error loading scheduled transactions:', error);
    }
}
function renderScheduledTransactionsList() {
    const tbody = document.getElementById('scheduled-list');
    if (!tbody) return;
    if (scheduledTransactions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" style="text-align:center; padding:20px; color:var(--text-secondary);">No scheduled transactions. Create one to automate recurring payments!</td></tr>';
        return;
    }
    tbody.innerHTML = scheduledTransactions.map(st => {
        const categoryName = st.category_id ? (categories.find(c => c.id === st.category_id)?.name || 'Unknown') : 'None';
        const accountName = st.account_id ? (accounts.find(a => a.id === st.account_id)?.bank_name || 'Unknown') : 'None';
        const nextDate = new Date(st.repeat_at).toLocaleDateString();
        return `
        <tr>
            <td>${st.name}</td>
            <td class="${st.amount >= 0 ? 'text-green-600' : 'text-red-600'}">
                ${st.amount >= 0 ? '+' : ''}$${Math.abs(st.amount).toFixed(2)}
            </td>
            <td style="text-transform: capitalize;">${st.repetition}</td>
            <td>${nextDate}</td>
            <td>${categoryName}</td>
            <td>${accountName}</td>
            <td class="text-right">
                <button onclick="editScheduledTransaction(${st.id})" class="btn-icon">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
                <button onclick="deleteScheduledTransaction(${st.id})" class="btn-icon delete">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                </button>
            </td>
        </tr>
    `}).join('');
}
function editScheduledTransaction(id) {
    const st = scheduledTransactions.find(s => s.id === id);
    if (!st) return;
    document.getElementById('scheduled-id').value = st.id;
    document.getElementById('scheduled-name').value = st.name;
    document.getElementById('scheduled-amount').value = st.amount;
    document.getElementById('scheduled-repetition').value = st.repetition;
    document.getElementById('scheduled-repeat-at').value = new Date(st.repeat_at).toISOString().split('T')[0];
    document.getElementById('scheduled-category').value = st.category_id || '';
    document.getElementById('scheduled-account').value = st.account_id || '';
    document.getElementById('scheduled-modal-title').textContent = 'Edit Scheduled Transaction';
    document.getElementById('scheduled-modal').classList.add('show');
}
async function deleteScheduledTransaction(id) {
    if (!confirm('Are you sure you want to delete this scheduled transaction?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/scheduled/${id}`, {
            method: 'DELETE',
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to delete scheduled transaction');
        await loadScheduledTransactions();
    } catch (error) {
        console.error('Error deleting scheduled transaction:', error);
        alert('Failed to delete scheduled transaction');
    }
}
async function processScheduledTransactions() {
    if (!confirm('Process all due scheduled transactions now?')) return;
    try {
        const response = await fetch(`${API_BASE}/api/scheduled/process`, {
            method: 'POST',
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to process scheduled transactions');
        const result = await response.json();
        alert(`Processed ${result.processed} scheduled transaction(s)`);
        await loadData();
    } catch (error) {
        console.error('Error processing scheduled transactions:', error);
        alert('Failed to process scheduled transactions');
    }
}
async function loadDashboardStats() {
    try {
        const response = await fetch(`${API_BASE}/api/dashboard/stats`, {
            headers: authHeaders()
        });
        if (!response.ok) throw new Error('Failed to fetch stats');
        const stats = await response.json();
        document.getElementById('total-income').textContent = `$${stats.total_income.toFixed(2)}`;
        document.getElementById('total-expenses').textContent = `$${stats.total_expenses.toFixed(2)}`;
        document.getElementById('net-balance').textContent = `$${stats.balance.toFixed(2)}`;
        document.getElementById('transaction-count').textContent = stats.transaction_count;
    } catch (error) {
        console.error('Error loading dashboard stats:', error);
    }
}
function setupForms() {
    const transactionForm = document.getElementById('transaction-form');
    if (transactionForm) {
        transactionForm.addEventListener('submit', handleTransactionSubmit);
    }
    const categoryForm = document.getElementById('category-form');
    if (categoryForm) {
        categoryForm.addEventListener('submit', handleCategorySubmit);
    }
    const accountForm = document.getElementById('account-form');
    if (accountForm) {
        accountForm.addEventListener('submit', handleAccountSubmit);
    }
    const budgetForm = document.getElementById('budget-form');
    if (budgetForm) {
        budgetForm.addEventListener('submit', handleBudgetSubmit);
    }
    const scheduledForm = document.getElementById('scheduled-form');
    if (scheduledForm) {
        scheduledForm.addEventListener('submit', handleScheduledSubmit);
    }
}
async function handleTransactionSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('transaction-id').value;
    const name = document.getElementById('transaction-name').value;
    const amount = parseFloat(document.getElementById('transaction-amount').value);
    const categoryId = document.getElementById('transaction-category').value;
    const accountId = document.getElementById('transaction-account').value;
    const data = {
        name,
        amount,
        category_id: categoryId ? parseInt(categoryId) : null,
        account_id: accountId ? parseInt(accountId) : null
    };
    if (amount < 0 && categoryId) {
        try {
            const budgetCheck = await fetch(`${API_BASE}/api/budgets/check`, {
                method: 'POST',
                headers: authHeaders(),
                body: JSON.stringify(data)
            });
            const budgetResult = await budgetCheck.json();
            if (budgetResult.exceeded) {
                const confirmMsg = `⚠️ WARNING: Budget Exceeded!\n\nBudget: $${budgetResult.budget.toFixed(2)}\nCurrent Spent: $${budgetResult.spent.toFixed(2)}\nNew Total: $${budgetResult.new_total.toFixed(2)}\n\nDo you want to continue?`;
                if (!confirm(confirmMsg)) {
                    return;
                }
            }
        } catch (err) {
            console.error('Error checking budget:', err);
        }
    }
    try {
        const url = id ? `${API_BASE}/api/transactions/${id}` : `${API_BASE}/api/transactions`;
        const method = id ? 'PUT' : 'POST';
        const response = await fetch(url, {
            method,
            headers: authHeaders(),
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Failed to save transaction');
        closeTransactionModal();
        await loadData();
    } catch (error) {
        console.error('Error saving transaction:', error);
        alert('Failed to save transaction');
    }
}
async function handleCategorySubmit(e) {
    e.preventDefault();
    const id = document.getElementById('category-id').value;
    const name = document.getElementById('category-name').value;
    try {
        const url = id ? `${API_BASE}/api/categories/${id}` : `${API_BASE}/api/categories`;
        const method = id ? 'PUT' : 'POST';
        const response = await fetch(url, {
            method,
            headers: authHeaders(),
            body: JSON.stringify({ name })
        });
        if (!response.ok) throw new Error('Failed to save category');
        closeCategoryModal();
        await loadCategories();
    } catch (error) {
        console.error('Error saving category:', error);
        alert('Failed to save category');
    }
}
async function handleAccountSubmit(e) {
    e.preventDefault();
    const name = document.getElementById('account-name').value;
    const amount = parseFloat(document.getElementById('account-amount').value);
    try {
        const response = await fetch(`${API_BASE}/api/accounts`, {
            method: 'POST',
            headers: authHeaders(),
            body: JSON.stringify({ bank_name: name, amount })
        });
        if (!response.ok) throw new Error('Failed to create account');
        closeAccountModal();
        await loadAccounts();
    } catch (error) {
        console.error('Error creating account:', error);
        alert('Failed to create account');
    }
}
async function handleBudgetSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('budget-id').value;
    const categoryId = parseInt(document.getElementById('budget-category').value);
    const amount = parseFloat(document.getElementById('budget-amount').value);
    const criteria = document.getElementById('budget-criteria').value;
    const data = {
        category_id: categoryId,
        amount,
        criteria
    };
    try {
        const url = id ? `${API_BASE}/api/budgets/${id}` : `${API_BASE}/api/budgets`;
        const method = id ? 'PUT' : 'POST';
        const response = await fetch(url, {
            method,
            headers: authHeaders(),
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Failed to save budget');
        closeBudgetModal();
        await loadBudgets();
    } catch (error) {
        console.error('Error saving budget:', error);
        alert('Failed to save budget');
    }
}
async function handleScheduledSubmit(e) {
    e.preventDefault();
    const id = document.getElementById('scheduled-id').value;
    const name = document.getElementById('scheduled-name').value;
    const amount = parseFloat(document.getElementById('scheduled-amount').value);
    const repetition = document.getElementById('scheduled-repetition').value;
    const repeatAt = document.getElementById('scheduled-repeat-at').value + 'T00:00:00Z';
    const categoryId = document.getElementById('scheduled-category').value;
    const accountId = document.getElementById('scheduled-account').value;
    const data = {
        name,
        amount,
        repetition,
        repeat_at: repeatAt,
        category_id: categoryId ? parseInt(categoryId) : null,
        account_id: accountId ? parseInt(accountId) : null
    };
    try {
        const url = id ? `${API_BASE}/api/scheduled/${id}` : `${API_BASE}/api/scheduled`;
        const method = id ? 'PUT' : 'POST';
        const response = await fetch(url, {
            method,
            headers: authHeaders(),
            body: JSON.stringify(data)
        });
        if (!response.ok) throw new Error('Failed to save scheduled transaction');
        closeScheduledModal();
        await loadScheduledTransactions();
    } catch (error) {
        console.error('Error saving scheduled transaction:', error);
        alert('Failed to save scheduled transaction');
    }
}
function openTransactionModal() {
    document.getElementById('transaction-id').value = '';
    document.getElementById('transaction-name').value = '';
    document.getElementById('transaction-amount').value = '';
    document.getElementById('transaction-category').value = '';
    document.getElementById('transaction-account').value = '';
    document.getElementById('transaction-modal-title').textContent = 'New Transaction';
    document.getElementById('transaction-modal').classList.add('show');
}
function closeTransactionModal() {
    document.getElementById('transaction-modal').classList.remove('show');
}
function openCategoryModal() {
    document.getElementById('category-id').value = '';
    document.getElementById('category-name').value = '';
    document.getElementById('category-modal-title').textContent = 'New Category';
    document.getElementById('category-modal').classList.add('show');
}
function closeCategoryModal() {
    document.getElementById('category-modal').classList.remove('show');
}
function openAccountModal() {
    document.getElementById('account-name').value = '';
    document.getElementById('account-amount').value = '0';
    document.getElementById('account-modal').classList.add('show');
}
function closeAccountModal() {
    document.getElementById('account-modal').classList.remove('show');
}
function openBudgetModal() {
    document.getElementById('budget-id').value = '';
    document.getElementById('budget-category').value = '';
    document.getElementById('budget-amount').value = '';
    document.getElementById('budget-criteria').value = 'monthly';
    document.getElementById('budget-modal-title').textContent = 'New Budget';
    document.getElementById('budget-modal').classList.add('show');
}
function closeBudgetModal() {
    document.getElementById('budget-modal').classList.remove('show');
}
function openScheduledModal() {
    document.getElementById('scheduled-id').value = '';
    document.getElementById('scheduled-name').value = '';
    document.getElementById('scheduled-amount').value = '';
    document.getElementById('scheduled-repetition').value = 'monthly';
    document.getElementById('scheduled-repeat-at').value = '';
    document.getElementById('scheduled-category').value = '';
    document.getElementById('scheduled-account').value = '';
    document.getElementById('scheduled-modal-title').textContent = 'New Scheduled Transaction';
    document.getElementById('scheduled-modal').classList.add('show');
}
function closeScheduledModal() {
    document.getElementById('scheduled-modal').classList.remove('show');
}