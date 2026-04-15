const startBtn = document.getElementById("start-btn");
const stopBtn = document.getElementById("stop-btn");
const statusText = document.getElementById("status-text");
const dot = document.getElementById("dot");
const logOutput = document.getElementById("log-output");
const queueBody = document.getElementById("queue-body");
const heroStatus = document.getElementById("hero-status");
const footerWpUrl = document.getElementById("footer-wp-url");
const footerAiEngine = document.getElementById("footer-ai-engine");

let ws;
let wsReconnectDelay = 3000;
const WS_MAX_DELAY = 30000;
let currentArticleId = null;
let queueInterval = null;

function initWS() {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const wsUrl = protocol + "://" + window.location.host + "/ws/logs";
    ws = new WebSocket(wsUrl);
    ws.onopen = () => { console.log("WS Connected"); wsReconnectDelay = 3000; };
    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (data.type === "log") {
                addLog(data.message, data.timestamp);
            } else if (data.type === "queue_update") {
                pollQueue();
            } else if (data.type === "status_update") {
                updateUI(data.data.running);
            }
        } catch (e) {
            // fallback for plain text logs
            addLog(event.data);
        }
    };
    ws.onclose = () => {
        setTimeout(initWS, wsReconnectDelay);
        wsReconnectDelay = Math.min(wsReconnectDelay * 2, WS_MAX_DELAY);
    };
    ws.onerror = (e) => { console.error("WS Error", e); ws.close(); };
}

function showToast(message, type = "info") {
    const container = document.getElementById("toast-container");
    const toast = document.createElement("div");
    toast.className = `toast toast-${type} show`;
    let icon = type === "success" ? "fa-check-circle" : (type === "error" ? "fa-exclamation-circle" : "fa-info-circle");
    toast.innerHTML = `<i class="fa-solid ${icon} mr-2"></i> <span class="text-sm font-semibold">${escapeHTML(message)}</span>`;
    container.appendChild(toast);
    setTimeout(() => {
        toast.classList.remove("show");
        setTimeout(() => toast.remove(), 400);
    }, 3000);
}

function switchView(viewName) {
    document.querySelectorAll('.view').forEach(el => {
        el.classList.add('hidden');
        el.style.display = 'none';
    });
    document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
    const targetView = document.getElementById('view-' + viewName);
    if (targetView) {
        targetView.classList.remove('hidden');
        targetView.style.display = (viewName === 'logs') ? 'flex' : 'block';
    }
    const navBtn = document.getElementById('nav-' + viewName);
    if (navBtn) navBtn.classList.add('active');

    if (queueInterval) clearInterval(queueInterval);
    
    if (viewName === 'queue') {
        pollQueue();
    }
    if (viewName === 'settings' || viewName === 'prompts' || viewName === 'dashboard') loadSettings();
    if (viewName === 'analytics') loadAnalytics();
}

async function loadSettings() {
    try {
        const r = await fetch("/api/settings");
        const settings = await r.json();
        if (!settings || Object.keys(settings).length === 0) return;

        if (document.getElementById('set-wp-url')) document.getElementById('set-wp-url').value = settings.wordpress?.url || "";
        if (document.getElementById('set-wp-user')) document.getElementById('set-wp-user').value = settings.wordpress?.username || "";
        if (document.getElementById('set-wp-pass')) document.getElementById('set-wp-pass').value = settings.wordpress?.app_password || "";
        if (document.getElementById('set-api-nvidia')) document.getElementById('set-api-nvidia').value = settings.ai?.nvidia_api_key || "";
        if (document.getElementById('set-api-modelslab')) document.getElementById('set-api-modelslab').value = settings.ai?.modelslab_api_key || "";
        if (document.getElementById('set-bot-mins')) document.getElementById('set-bot-mins').value = settings.bot?.run_interval_minutes || 10;
        if (document.getElementById('set-pub-mins')) document.getElementById('set-pub-mins').value = settings.bot?.publish_interval_minutes || 15;
        if (document.getElementById('set-adv-minlen')) document.getElementById('set-adv-minlen').value = settings.advanced?.min_article_len || 500;
        if (document.getElementById('set-bot-autopublish')) document.getElementById('set-bot-autopublish').checked = settings.bot?.auto_publish !== false;

        if (document.getElementById('set-prompt-system-nvidia')) {
            document.getElementById('set-prompt-system-nvidia').value = settings.prompts?.system_prompt_nvidia || "";
            document.getElementById('set-prompt-system-modelslab').value = settings.prompts?.system_prompt_modelslab || "";
        }

        if (footerWpUrl) footerWpUrl.innerText = settings.wordpress?.url || "Սահմանված չէ";
        if (footerAiEngine) {
            let textModel = "Gemma-3 (Nvidia NIM)";
            let imgModel = settings.ai?.modelslab_api_key ? 'ModelsLab SDXL' : 'Nvidia ImageGen';
            footerAiEngine.innerText = `${textModel} / ${imgModel}`;
        }

        const container = document.getElementById('topics-container');
        if (container) {
            container.innerHTML = '';
            if (settings.topics && settings.topics.length > 0) {
                settings.topics.forEach(t => addTopic(t));
            } else {
                addTopic();
            }
        }
    } catch(e) { console.error(e); }
}

function addTopic(data = {name:'', wp_category_id:1, rss_url:''}) {
    const container = document.getElementById('topics-container');
    if (!container) return;
    const div = document.createElement('div');
    div.className = 'grid grid-cols-1 md:grid-cols-12 gap-4 items-end bg-white/5 p-5 rounded-2xl border border-white/5 mb-4';
    div.innerHTML = `
        <div class="md:col-span-3"><label class="block text-[9px] font-bold uppercase text-slate-500 mb-2">Անուն</label><input type="text" class="topic-name w-full bg-black/40 border border-white/10 rounded-xl p-3 text-sm text-white" value="${escapeHTML(data.name)}"></div>
        <div class="md:col-span-2"><label class="block text-[9px] font-bold uppercase text-slate-500 mb-2">WP ID</label><input type="number" class="topic-cat w-full bg-black/40 border border-white/10 rounded-xl p-3 text-sm text-white" value="${escapeHTML(data.wp_category_id)}"></div>
        <div class="md:col-span-6"><label class="block text-[9px] font-bold uppercase text-slate-500 mb-2">RSS Հղում</label><input type="url" class="topic-rss w-full bg-black/40 border border-white/10 rounded-xl p-3 text-sm text-white" value="${escapeHTML(data.rss_url)}"></div>
        <div class="md:col-span-1 flex justify-center"><button type="button" class="text-slate-500 hover:text-rose-500 p-3" onclick="this.closest('.grid').remove()"><i class="fa-solid fa-trash"></i></button></div>
    `;
    container.appendChild(div);
}

async function saveSettings(e) {
    if (e) e.preventDefault();
    try {
        const current = await fetch("/api/settings").then(r => r.json());
        const topics = [];
        document.querySelectorAll('#topics-container .grid').forEach(row => {
            const name = row.querySelector('.topic-name').value.trim();
            const rss = row.querySelector('.topic-rss').value.trim();
            if (name || rss) {
                topics.push({
                    name: name,
                    wp_category_id: parseInt(row.querySelector('.topic-cat').value) || 1,
                    rss_url: rss
                });
            }
        });

        const cfg = {
            wordpress: {
                url: document.getElementById('set-wp-url').value,
                username: document.getElementById('set-wp-user').value,
                app_password: document.getElementById('set-wp-pass').value || current.wordpress?.app_password
            },
            ai: {
                tool: "nvidia",
                nvidia_api_key: document.getElementById('set-api-nvidia').value || current.ai?.nvidia_api_key,
                modelslab_api_key: document.getElementById('set-api-modelslab').value || current.ai?.modelslab_api_key
            },
            advanced: { min_article_len: parseInt(document.getElementById('set-adv-minlen').value) || 500 },
            bot: {
                run_interval_hours: 0,
                run_interval_minutes: parseInt(document.getElementById('set-bot-mins').value) || 10,
                publish_interval_minutes: parseInt(document.getElementById('set-pub-mins').value) || 15,
                auto_publish: document.getElementById('set-bot-autopublish').checked
            },
            topics: topics,
            prompts: {
                system_prompt_nvidia: document.getElementById('set-prompt-system-nvidia').value,
                system_prompt_modelslab: document.getElementById('set-prompt-system-modelslab').value
            }
        };

        const r = await fetch("/api/settings", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(cfg) });
        if (r.ok) { showToast("Կարգավորումները պահպանված են", "success"); loadSettings(); }
    } catch(e) { showToast("Սխալ պահպանելիս", "error"); }
}

async function pollQueue() {
    try {
        const r = await fetch("/api/queue");
        const queue = await r.json();
        if (!queueBody) return;
        queueBody.innerHTML = queue.length ? "" : '<tr><td colspan="4" class="py-12 text-center text-slate-500 italic">Հերթը դատարկ է</td></tr>';
        queue.forEach(item => {
            const tr = document.createElement('tr');
            tr.className = "border-b border-white/[0.02] hover:bg-white/[0.01] transition-colors cursor-pointer";
            tr.onclick = (e) => { if(!e.target.closest('button')) showPreview(item); };

            let status = item.status || 'pending';
            let isFailed = status === 'failed';
            let statusBadge = isFailed ? 'bg-rose-500/10 text-rose-400 border-rose-500/20' : 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';

            let safeTitle = escapeHTML(item.title);
            let safeImageUrl = escapeHTML(item.image_url);

            tr.innerHTML = `
                <td class="py-6 px-8"><img src="${safeImageUrl}" class="w-14 h-14 object-cover rounded-xl border border-white/10" onerror="this.src='/temp_image.jpg'"></td>
                <td class="py-6 px-8"><div class="text-sm font-bold text-white mb-1 line-clamp-1">${safeTitle}</div></td>
                <td class="py-6 px-8"><span class="px-3 py-1 rounded-full text-[9px] font-black uppercase border ${statusBadge}">${escapeHTML(status)}</span></td>
                <td class="py-6 px-8 text-right">
                    <button onclick="publishItem(${item.id})" class="text-emerald-400 hover:text-white mr-4"><i class="fa-solid fa-check"></i></button>
                    <button onclick="deleteItem(${item.id})" class="text-rose-400 hover:text-white"><i class="fa-solid fa-trash"></i></button>
                </td>
            `;
            queueBody.appendChild(tr);
        });
    } catch(e) {}
}

async function checkStatus() {
    try {
        const r = await fetch("/api/bot/status");
        const data = await r.json();
        updateUI(data.running);
    } catch (e) {}
}

function updateUI(running) {
    if (running) {
        startBtn.classList.add("hidden"); stopBtn.classList.remove("hidden");
        statusText.innerText = "Աշխատում է"; dot.classList.add("bg-emerald-500", "animate-pulse"); dot.classList.remove("bg-rose-500");
        if(heroStatus) heroStatus.innerText = "Active";
    } else {
        startBtn.classList.remove("hidden"); stopBtn.classList.add("hidden");
        statusText.innerText = "Կանգնեցված է"; dot.classList.remove("bg-emerald-500", "animate-pulse"); dot.classList.add("bg-rose-500");
        if(heroStatus) heroStatus.innerText = "Idle";
    }
}

async function toggleBot() {
    const r = await fetch("/api/bot/toggle", { method: "POST" });
    const data = await r.json();
    updateUI(data.status === "started");
}

async function deleteItem(id) {
    if(!confirm("Հեռացնե՞լ հերթից:")) return;
    const r = await fetch("/api/queue/delete", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({id: id}) });
    if(r.ok) pollQueue();
}

async function clearQueue() {
    if(!confirm("Մաքրե՞լ հերթի բոլոր հոդվածները (վերաշարադրված և սխալված):")) return;
    showToast("Մաքրվում է...", "info");
    const r = await fetch("/api/queue/clear", { method: "POST" });
    if(r.ok) {
        showToast("Հերթը մաքրված է", "success");
        pollQueue();
    } else {
        showToast("Սխալ մաքրելիս", "error");
    }
}

async function publishItem(id) {
    showToast("Հրապարակումը սկսվեց...", "info");
    const r = await fetch("/api/queue/publish", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({id: id}) });
    if(r.ok) { pollQueue(); }
}

function showPreview(item) {
    currentArticleId = item.id;
    document.getElementById('prev-title').innerText = item.title;
    // We trust this because it is sanitized on the backend via bluemonday.UGCPolicy
    document.getElementById('prev-content').innerHTML = item.rewritten_content?.String || item.content;
    document.getElementById('prev-image').src = item.image_url;
    document.getElementById('prev-image').style.display = item.image_url ? 'block' : 'none';
    document.getElementById('preview-modal').classList.remove('hidden');
    document.getElementById('prev-publish-btn').onclick = () => { publishItem(item.id); closePreview(); };
}

function openEditModal() {
    const title = document.getElementById('prev-title').innerText;
    const content = document.getElementById('prev-content').innerHTML;
    document.getElementById('edit-title').value = title;
    document.getElementById('edit-content').value = content;
    document.getElementById('edit-modal').classList.remove('hidden');
}

async function saveEdit() {
    if (!currentArticleId) return;
    const updated = {
        title: document.getElementById('edit-title').value,
        rewritten_content: document.getElementById('edit-content').value
    };
    const r = await fetch(`/api/queue/${currentArticleId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(updated)
    });
    if (r.ok) {
        showToast("Փոփոխությունները պահպանված են", "success");
        closeEditModal();
        pollQueue();
        document.getElementById('prev-title').innerText = updated.title;
        document.getElementById('prev-content').innerHTML = updated.rewritten_content;
    }
}

function closePreview() { document.getElementById('preview-modal').classList.add('hidden'); }
function closeEditModal() { document.getElementById('edit-modal').classList.add('hidden'); }
function togglePasswordVisibility(id) {
    const el = document.getElementById(id);
    const icon = document.getElementById('icon-' + id);
    if (el.type === 'password') { el.type = 'text'; icon.classList.replace('fa-eye', 'fa-eye-slash'); }
    else { el.type = 'password'; icon.classList.replace('fa-eye-slash', 'fa-eye'); }
}

async function loadAnalytics() {
    try {
        const r = await fetch("/api/analytics");
        const data = await r.json();
        document.getElementById('stat-published').innerText = data.total_published || 0;
        document.getElementById('stat-pending').innerText = data.pending_queue || 0;
        document.getElementById('stat-errors').innerText = data.errors || 0;
    } catch(e) {}
}

function escapeHTML(str) {
    if (str === null || str === undefined) return '';
    return String(str).replace(/[&<>'"]/g, 
        tag => ({
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            "'": '&#39;',
            '"': '&quot;'
        }[tag]));
}

function addLog(msg, ts) {
    if (!logOutput) return;
    const div = document.createElement("div");
    div.className = "mb-1 py-1 border-b border-white/[0.02]";
    let colorClass = msg.includes("SUCCESS") ? "text-emerald-400" : (msg.includes("ERROR") ? "text-rose-400" : "text-slate-400");
    let displayTs = ts || new Date().toLocaleTimeString();
    div.innerHTML = `<span class="text-slate-600 mr-3 text-[10px]">${escapeHTML(displayTs)}</span> <span class="${colorClass}">${escapeHTML(msg)}</span>`;
    logOutput.appendChild(div);
    logOutput.scrollTop = logOutput.scrollHeight;
}

// Init
initWS();
checkStatus();
loadSettings();
setInterval(checkStatus, 5000);
if (startBtn) startBtn.onclick = toggleBot;
if (stopBtn) stopBtn.onclick = toggleBot;
if (document.getElementById('settings-form')) document.getElementById('settings-form').onsubmit = saveSettings;
if (document.getElementById('prompts-form')) document.getElementById('prompts-form').onsubmit = saveSettings;
