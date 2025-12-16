let currentTab = 'search';
let searchParams = {};

const tabs = document.querySelectorAll('.tab-btn');
const tabContents = document.querySelectorAll('.tab-content');
const searchForm = document.getElementById('search-form');
const roomsGrid = document.getElementById('rooms-grid');
const bookingModal = document.getElementById('booking-modal');
const bookingForm = document.getElementById('booking-form');
const toast = document.getElementById('toast');

document.addEventListener('DOMContentLoaded', () => {
    setupTabs();
    setupForms();
    setupModals();
    setDefaultDates();
    loadAdminData();
});

function setupTabs() {
    tabs.forEach(btn => {
        btn.addEventListener('click', () => {
            const tabId = btn.dataset.tab;
            switchTab(tabId);
        });
    });
}

function switchTab(tabId) {
    currentTab = tabId;

    tabs.forEach(btn => {
        btn.classList.toggle('active', btn.dataset.tab === tabId);
    });

    tabContents.forEach(content => {
        content.classList.toggle('active', content.id === tabId);
    });

    if (tabId === 'admin') {
        loadAdminData();
    }
}

function setDefaultDates() {
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);

    const checkIn = document.querySelector('input[name="check_in"]');
    const checkOut = document.querySelector('input[name="check_out"]');

    if (checkIn && checkOut) {
        checkIn.valueAsDate = today;
        checkOut.valueAsDate = tomorrow;
    }
}

function setupForms() {
    searchForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(searchForm);
        searchParams = Object.fromEntries(formData.entries());
        await searchRooms(searchParams);
    });

    bookingForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await createBooking();
    });

    const addRoomForm = document.getElementById('add-room-form');
    if (addRoomForm) {
        addRoomForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            await addRoom(new FormData(addRoomForm));
        });
    }

    const notificationForm = document.getElementById('notification-form');
    if (notificationForm) {
        notificationForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            await sendNotification(new FormData(notificationForm));
        });
    }

    const addSpecialDateForm = document.getElementById('add-special-date-form');
    if (addSpecialDateForm) {
        addSpecialDateForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            await addSpecialDate(new FormData(addSpecialDateForm));
        });
    }
}

async function searchRooms(params) {
    try {
        const query = new URLSearchParams(params).toString();
        const res = await fetch(`/booking/rooms/search?${query}`);
        if (!res.ok) throw new Error('Не удалось загрузить номера');

        const rooms = await res.json();
        renderRooms(rooms);
    } catch (err) {
        showToast(err.message, 'error');
    }
}

function renderRooms(rooms) {
    if (!rooms || rooms.length === 0) {
        roomsGrid.innerHTML = '<p style="grid-column: 1/-1; text-align: center; padding: 20px;">Нет доступных номеров на выбранные даты</p>';
        return;
    }

    roomsGrid.innerHTML = rooms.map(item => `
        <div class="room-card">
            <div class="room-image">
                <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
                    <path d="M3 21h18M5 21V7l8-4 8 4v14M8 21v-9a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v9"/>
                </svg>
            </div>
            <div class="room-details">
                <div class="room-header">
                    <span class="room-type">${getRoomTypeName(item.room.room_type)}</span>
                    <span class="room-number">№${item.room.room_number}</span>
                </div>
                <p class="room-description">${item.room.description || 'Описание отсутствует'}</p>
                <div class="room-meta">
                    <div class="meta-item">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                            <circle cx="9" cy="7" r="4"></circle>
                            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                        </svg>
                        ${item.room.capacity} чел.
                    </div>
                </div>
                <div class="room-footer">
                    <div class="price">
                        ${formatPrice(item.total_price)} RUB
                        <span>за период</span>
                    </div>
                    <button onclick="openBookingModal(${JSON.stringify(item.room).replace(/"/g, '&quot;')})" class="btn-primary">
                        Забронировать
                    </button>
                </div>
            </div>
        </div>
    `).join('');
}

async function openBookingModal(room) {
    const modal = document.getElementById('booking-modal');
    modal.querySelector('input[name="room_id"]').value = room.id;
    document.getElementById('modal-room-number').textContent = room.room_number;
    document.getElementById('modal-room-type').textContent = getRoomTypeName(room.room_type);
    document.getElementById('modal-room-price').textContent = formatPrice(room.base_price);

    const checkIn = new Date(searchParams.check_in);
    const checkOut = new Date(searchParams.check_out);
    const nights = getNights(checkIn, checkOut);

    document.getElementById('modal-check-in').textContent = checkIn.toLocaleDateString();
    document.getElementById('modal-check-out').textContent = checkOut.toLocaleDateString();
    document.getElementById('modal-nights').textContent = nights;

    await calculatePrice(room.id, searchParams.check_in, searchParams.check_out);

    modal.classList.add('active');
}

async function calculatePrice(roomId, checkIn, checkOut) {
    try {
        const res = await fetch('/booking/price', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                room_id: roomId,
                check_in: checkIn,
                check_out: checkOut
            })
        });

        const data = await res.json();

        const breakdownHtml = data.daily_breakdown.map(day => `
            <div class="breakdown-item">
                <span>${new Date(day.date).toLocaleDateString()} (${day.reason})</span>
                <span>${formatPrice(day.day_price)} RUB</span>
            </div>
        `).join('');

        document.getElementById('price-breakdown').innerHTML = `
            ${breakdownHtml}
            <div class="breakdown-total">
                <span>Базовая цена:</span>
                <span>${formatPrice(data.base_price)} RUB/ночь</span>
            </div>
        `;

        document.getElementById('modal-total-price').textContent = formatPrice(data.total_price);
    } catch (err) {
        console.error(err);
    }
}

async function createBooking() {
    const formData = new FormData(bookingForm);
    const data = {
        room_id: parseInt(formData.get('room_id')),
        start_date: new Date(searchParams.check_in).toISOString(),
        end_date: new Date(searchParams.check_out).toISOString(),
        guest_info: {
            name: formData.get('name'),
            email: formData.get('email'),
            phone: formData.get('phone')
        }
    };

    try {
        const res = await fetch('/booking/', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!res.ok) throw new Error('Не удалось создать бронирование');

        showToast('Бронирование успешно создано!', 'success');
        closeModals();
        searchRooms(searchParams);
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function loadAdminData() {
    await Promise.all([
        loadAdminStats(),
        loadAdminRooms(),
        loadAdminBookings(),
        loadSpecialDates()
    ]);
}

async function loadAdminStats() {
    try {
        const res = await fetch('/admin/stats');
        const stats = await res.json();

        document.getElementById('hotel-stats').innerHTML = `
            <div class="stat-item">
                <div class="stat-label">Всего номеров</div>
                <div class="stat-value">${stats.total_rooms}</div>
            </div>
            <div class="stat-item">
                <div class="stat-label">Свободно</div>
                <div class="stat-value">${stats.available_rooms}</div>
            </div>
            <div class="stat-item">
                <div class="stat-label">Занято</div>
                <div class="stat-value">${stats.occupied_rooms}</div>
            </div>
            <div class="stat-item">
                <div class="stat-label">Выручка</div>
                <div class="stat-value">${formatPrice(stats.total_revenue)} ₽</div>
            </div>
        `;
    } catch (err) {
        console.error(err);
    }
}

async function loadSpecialDates() {
    try {
        const res = await fetch('/admin/dates');
        const dates = await res.json();

        const tbody = document.querySelector('#special-dates-table tbody');
        if (!dates) {
            tbody.innerHTML = '<tr><td colspan="4">Нет специальных дат</td></tr>';
            return;
        }

        tbody.innerHTML = dates.map(d => `
            <tr>
                <td>${new Date(d.date).toLocaleDateString()}</td>
                <td>${d.name}</td>
                <td>x${d.coefficient}</td>
                <td>
                    <button onclick="deleteSpecialDate(${d.id})" class="btn-icon" style="color: var(--danger)">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <polyline points="3 6 5 6 21 6"></polyline>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                    </button>
                </td>
            </tr>
        `).join('');
    } catch (err) {
        console.error(err);
    }
}

async function addSpecialDate(formData) {
    try {
        const data = Object.fromEntries(formData.entries());
        data.coefficient = parseFloat(data.coefficient);

        const res = await fetch('/admin/dates', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!res.ok) throw new Error('Не удалось добавить дату');

        showToast('Дата добавлена', 'success');
        closeModals();
        loadSpecialDates();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function deleteSpecialDate(id) {
    if (!confirm('Удалить эту дату?')) return;

    try {
        const res = await fetch(`/admin/dates/${id}`, { method: 'DELETE' });
        if (!res.ok) throw new Error('Не удалось удалить');
        loadSpecialDates();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function loadAdminRooms() {
    try {
        const res = await fetch('/admin/rooms');
        const rooms = await res.json();

        const tbody = document.querySelector('#admin-rooms-table tbody');
        if (!rooms || rooms.length === 0) {
            tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; padding: 20px;">Номеров не найдено</td></tr>';
            return;
        }

        tbody.innerHTML = rooms.map(room => `
            <tr>
                <td>${room.room_number}</td>
                <td>${getRoomTypeName(room.room_type)}</td>
                <td>${formatPrice(room.base_price)}</td>
                <td><span class="status-badge status-${room.status}">${getStatusName(room.status)}</span></td>
                <td>
                    <button onclick="deleteRoom(${room.id})" class="btn-icon" style="color: var(--danger)">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <polyline points="3 6 5 6 21 6"></polyline>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                    </button>
                </td>
            </tr>
        `).join('');
    } catch (err) {
        console.error(err);
    }
}

async function deleteRoom(id) {
    if (!confirm('Удалить этот номер? Все связанные бронирования будут удалены.')) return;

    try {
        const res = await fetch(`/admin/rooms/${id}`, { method: 'DELETE' });
        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.message || 'Не удалось удалить номер');
        }
        showToast('Номер удален', 'success');
        loadAdminRooms();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function loadAdminBookings() {
    try {
        const status = document.getElementById('booking-status-filter').value;
        const url = status ? `/admin/bookings?status=${status}` : '/admin/bookings';
        const res = await fetch(url);
        const bookings = await res.json();

        const tbody = document.querySelector('#admin-bookings-table tbody');
        if (!bookings || bookings.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" style="text-align: center; padding: 20px;">Бронирований не найдено</td></tr>';
            return;
        }

        tbody.innerHTML = bookings.map(b => `
            <tr>
                <td>#${b.id}</td>
                <td>
                    <div>${b.guest_info.name}</div>
                    <small style="color: var(--text-light)">${b.guest_info.email}</small>
                </td>
                <td>${b.room.room_number}</td>
                <td>${new Date(b.start_date).toLocaleDateString()} - ${new Date(b.end_date).toLocaleDateString()}</td>
                <td><span class="status-badge status-${b.status}">${getStatusName(b.status)}</span></td>
                <td>
                    ${b.status === 'pending' ? `
                        <button onclick="confirmBooking(${b.id})" class="btn-icon" style="color: var(--success)" title="Подтвердить">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="20 6 9 17 4 12"></polyline>
                            </svg>
                        </button>
                        <button onclick="cancelBooking(${b.id})" class="btn-icon" style="color: var(--danger)" title="Отменить">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"></line>
                                <line x1="6" y1="6" x2="18" y2="18"></line>
                            </svg>
                        </button>
                    ` : ''}
                </td>
            </tr>
        `).join('');
    } catch (err) {
        console.error(err);
    }
}

async function confirmBooking(id) {
    if (!confirm('Подтвердить бронирование?')) return;
    try {
        const res = await fetch(`/booking/${id}/confirm`, { method: 'PUT' });
        if (!res.ok) throw new Error('Не удалось подтвердить');
        showToast('Бронирование подтверждено', 'success');
        loadAdminBookings();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function cancelBooking(id) {
    if (!confirm('Отменить бронирование?')) return;
    try {
        const res = await fetch(`/booking/${id}/cancel`, { method: 'PUT' });
        if (!res.ok) throw new Error('Не удалось отменить');
        showToast('Бронирование отменено', 'success');
        loadAdminBookings();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function addRoom(formData) {
    try {
        const data = Object.fromEntries(formData.entries());
        data.base_price = parseFloat(data.base_price);
        data.capacity = parseInt(data.capacity);

        const res = await fetch('/admin/rooms', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!res.ok) throw new Error('Не удалось создать номер');

        showToast('Номер создан', 'success');
        closeModals();
        loadAdminRooms();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

async function sendNotification(formData) {
    try {
        const data = Object.fromEntries(formData.entries());
        const res = await fetch('/notification/send', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!res.ok) throw new Error('Не удалось отправить');

        showToast('Уведомление отправлено', 'success');
        document.getElementById('notification-form').reset();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

function getRoomTypeName(type) {
    const types = {
        'standard': 'Стандарт',
        'deluxe': 'Делюкс',
        'suite': 'Люкс',
        'family': 'Семейный'
    };
    return types[type] || type;
}

function getStatusName(status) {
    const statuses = {
        'available': 'Свободен',
        'occupied': 'Занят',
        'maintenance': 'Ремонт',
        'pending': 'Ожидает',
        'confirmed': 'Подтверждено',
        'cancelled': 'Отменено'
    };
    return statuses[status] || status;
}

function formatPrice(price) {
    return new Intl.NumberFormat('ru-RU').format(price);
}

function getNights(d1, d2) {
    return Math.ceil((d2 - d1) / (1000 * 60 * 60 * 24));
}

function showToast(message, type = 'success') {
    toast.textContent = message;
    toast.className = `toast show ${type}`;
    setTimeout(() => {
        toast.className = 'toast';
    }, 3000);
}

function setupModals() {
    document.querySelectorAll('.close-modal').forEach(btn => {
        btn.addEventListener('click', closeModals);
    });

    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            closeModals();
        }
    });
}

function closeModals() {
    document.querySelectorAll('.modal').forEach(m => m.classList.remove('active'));
}

function openAddRoomModal() {
    document.getElementById('add-room-modal').classList.add('active');
}

function openSpecialDateModal() {
    document.getElementById('add-special-date-modal').classList.add('active');
}
