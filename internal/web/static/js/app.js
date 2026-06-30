// ============================================================
// 全局变量
// ============================================================
let currentDatabaseId = null;
let selectedTables = [];
let allTables = [];
let ignoredColumns = [];
let columnOverrides = [];

// Tab2 自定义片段相关
let snippetTableColumns = [];   // 当前表的所有列信息
let snippetList = [];           // 已添加的片段列表
let snippetMergeEnabled = false; // 是否启用"并入生成"
let editingSnippetIndex = null;  // 当前正在编辑的片段索引

// QueryBuilder WHERE 状态
let whereRules = [];            // [{id, fieldIdx, operator}]
let whereLogic = 'AND';         // AND | OR
let whereRuleCounter = 0;

// Chip 选择状态
let selectedChips = {
    selectFields: new Set(),
    insertFields: new Set(),
    setFields: new Set(),
};
let orderBySelections = new Map(); // colIdx -> direction

// ============================================================
// 工具函数
// ============================================================
function showMessage(message, type = 'info') {
    const el = document.getElementById('message');
    el.textContent = message;
    el.className = `message ${type}`;
    el.style.display = 'block';
    const duration = type === 'success' ? 5000 : type === 'error' ? 4000 : 3000;
    setTimeout(() => { el.style.display = 'none'; }, duration);
}

function toPascalCase(str) {
    return str.split('_').map(w => w.charAt(0).toUpperCase() + w.slice(1).toLowerCase()).join('');
}

function capitalize(s) {
    return s ? s.charAt(0).toUpperCase() + s.slice(1) : '';
}

function copyCode(elementId) {
    const el = document.getElementById(elementId);
    navigator.clipboard.writeText(el.textContent).then(() => {
        showMessage('已复制到剪贴板', 'success');
    }).catch(() => {
        const ta = document.createElement('textarea');
        ta.value = el.textContent;
        document.body.appendChild(ta);
        ta.select();
        document.execCommand('copy');
        document.body.removeChild(ta);
        showMessage('已复制到剪贴板', 'success');
    });
}

// 客户端计算方法名（与后端逻辑保持一致，用于片段列表展示）
function computeMethodName(cfg) {
    if (cfg.methodName) return cfg.methodName;
    const whereFields = (cfg.whereFields || []).filter(
        f => f.operator !== 'IS NULL' && f.operator !== 'IS NOT NULL'
    );
    switch (cfg.operation) {
        case 'select': {
            if (whereFields.length === 0) return cfg.isBatch ? 'selectAll' : 'selectByFields';
            const parts = whereFields.map(f => capitalize(f.fieldName));
            return cfg.isBatch ? 'selectBy' + parts.join('And') + 'In' : 'selectBy' + parts.join('And');
        }
        case 'insert':
            return cfg.isBatch ? 'insertBatchByFields' : 'insertByFields';
        case 'delete': {
            if (whereFields.length === 0) return 'deleteByFields';
            const parts = whereFields.map(f => capitalize(f.fieldName));
            return cfg.isBatch ? 'deleteBy' + parts.join('And') + 'In' : 'deleteBy' + parts.join('And');
        }
        case 'update': {
            const setParts = (cfg.setFields || []).map(f => capitalize(f.fieldName));
            const whereParts = whereFields.map(f => capitalize(f.fieldName));
            let name = 'update';
            if (setParts.length > 0) name += setParts.join('And');
            if (whereParts.length > 0) name += 'By' + whereParts.join('And');
            if (cfg.isBatch) name += 'Batch';
            return name;
        }
    }
    return 'customMethod';
}

// ============================================================
// Tab 切换
// ============================================================
function switchTab(tabName) {
    document.querySelectorAll('.tab-panel').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
    document.getElementById('tabConfig').classList.toggle('active', tabName === 'config');
    document.getElementById('tabSnippet').classList.toggle('active', tabName === 'snippet');
    document.getElementById('tabBtnConfig').classList.toggle('active', tabName === 'config');
    document.getElementById('tabBtnSnippet').classList.toggle('active', tabName === 'snippet');
    if (tabName === 'snippet') refreshSnippetPanelState();
}

// ============================================================
// 数据库连接管理
// ============================================================
async function loadConnections() {
    try {
        const response = await fetch('/api/connections');
        const connections = await response.json();
        const list = document.getElementById('connectionList');
        list.innerHTML = '';
        connections.forEach(conn => {
            const item = document.createElement('div');
            item.className = 'connection-item';
            item.innerHTML = `
                <div class="connection-info" onclick="selectConnection(${conn.id}, '${conn.name}')">
                    <strong>${conn.name}</strong>
                    <small>${conn.dbType} - ${conn.host}:${conn.port}</small>
                </div>
                <div class="connection-actions">
                    <button class="btn-icon" onclick="event.stopPropagation(); editConnection(${conn.id})" title="编辑">
                        <svg width="16" height="16" fill="currentColor"><path d="M12.146 1.146a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1 0 .708l-10 10A.5.5 0 0 1 4.5 14H2a.5.5 0 0 1-.5-.5v-2.5a.5.5 0 0 1 .146-.354l10-10z"/></svg>
                    </button>
                    <button class="btn-icon btn-danger" onclick="event.stopPropagation(); deleteConnection(${conn.id})" title="删除">
                        <svg width="16" height="16" fill="currentColor"><path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5zm3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0V6z"/><path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1v1zM4.118 4L4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4H4.118zM2.5 3V2h11v1h-11z"/></svg>
                    </button>
                </div>`;
            list.appendChild(item);
        });
    } catch (error) {
        showMessage('加载连接列表失败: ' + error.message, 'error');
    }
}

async function selectConnection(id) {
    currentDatabaseId = id;
    selectedTables = [];
    allTables = [];
    ignoredColumns = [];
    columnOverrides = [];
    document.querySelectorAll('.connection-item').forEach(item => item.classList.remove('active'));
    event.currentTarget.closest('.connection-item').classList.add('active');
    await loadTables();
}

// ============================================================
// 表列表
// ============================================================
async function loadTables(filter = '') {
    if (!currentDatabaseId) return;
    const list = document.getElementById('tableList');
    list.innerHTML = '<div class="loading-placeholder"><span class="loading-spinner"></span>加载中...</div>';
    try {
        const response = await fetch('/api/tables', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ databaseId: currentDatabaseId, filter })
        });
        const tables = await response.json();
        allTables = tables;
        list.innerHTML = '';
        tables.forEach(tableName => {
            const item = document.createElement('div');
            item.className = 'table-item' + (selectedTables.includes(tableName) ? ' selected' : '');
            item.innerHTML = `
                <input type="checkbox"
                       ${selectedTables.includes(tableName) ? 'checked' : ''}
                       onchange="toggleTableSelection('${tableName}', this.checked)">
                <span class="table-item-name" onclick="toggleTableCheckbox('${tableName}')">${tableName}</span>`;
            list.appendChild(item);
        });
        updateSelectionCount();
    } catch (error) {
        showMessage('加载表列表失败: ' + error.message, 'error');
    }
}

function toggleTableCheckbox(tableName) {
    toggleTableSelection(tableName, !selectedTables.includes(tableName));
}

function toggleTableSelection(tableName, checked) {
    if (checked) {
        if (!selectedTables.includes(tableName)) selectedTables.push(tableName);
    } else {
        selectedTables = selectedTables.filter(t => t !== tableName);
    }
    updateDefaultEntityFields();
    loadTables(document.getElementById('tableFilter').value);
    if (document.getElementById('tabSnippet').classList.contains('active')) {
        refreshSnippetPanelState();
    }
}

function selectAllTables() {
    selectedTables = [...allTables];
    updateDefaultEntityFields();
    loadTables(document.getElementById('tableFilter').value);
}

function deselectAllTables() {
    selectedTables = [];
    updateDefaultEntityFields();
    loadTables(document.getElementById('tableFilter').value);
}

function updateSelectionCount() {
    const countEl = document.getElementById('selectionCount');
    countEl.textContent = selectedTables.length > 0 ? `(已选 ${selectedTables.length} 张)` : '';
}

// 需求1：默认值明文填写
function updateDefaultEntityFields() {
    const tableNameEl = document.getElementById('tableName');
    const domainNameEl = document.getElementById('domainObjectName');
    const mapperNameEl = document.getElementById('mapperName');
    if (selectedTables.length === 1) {
        const t = selectedTables[0];
        tableNameEl.value = t;
        if (!domainNameEl.dataset.userEdited) domainNameEl.value = toPascalCase(t);
        if (!mapperNameEl.dataset.userEdited) mapperNameEl.value = toPascalCase(t) + 'Mapper';
    } else if (selectedTables.length > 1) {
        tableNameEl.value = `(已选 ${selectedTables.length} 张表)`;
        if (!domainNameEl.dataset.userEdited) domainNameEl.value = '';
        if (!mapperNameEl.dataset.userEdited) mapperNameEl.value = '';
    } else {
        tableNameEl.value = '';
        if (!domainNameEl.dataset.userEdited) domainNameEl.value = '';
        if (!mapperNameEl.dataset.userEdited) mapperNameEl.value = '';
    }
}

// ============================================================
// 连接管理弹窗
// ============================================================
function showConnectionModal(connection = null) {
    const modal = document.getElementById('connectionModal');
    const title = document.getElementById('connectionModalTitle');
    if (connection) {
        title.textContent = '编辑数据库连接';
        document.getElementById('connectionId').value = connection.id;
        document.getElementById('connectionName').value = connection.name;
        document.getElementById('dbType').value = connection.dbType;
        document.getElementById('host').value = connection.host;
        document.getElementById('port').value = connection.port;
        document.getElementById('schema').value = connection.schema;
        document.getElementById('username').value = connection.username;
        document.getElementById('password').value = connection.password;
    } else {
        title.textContent = '新建数据库连接';
        document.getElementById('connectionForm').reset();
        document.getElementById('connectionId').value = '';
        document.getElementById('port').value = '3306';
        document.getElementById('host').value = 'localhost';
    }
    modal.style.display = 'block';
}

function hideConnectionModal() {
    document.getElementById('connectionModal').style.display = 'none';
}

async function testConnection() {
    const config = {
        dbType: document.getElementById('dbType').value,
        host: document.getElementById('host').value,
        port: document.getElementById('port').value,
        schema: document.getElementById('schema').value,
        username: document.getElementById('username').value,
        password: document.getElementById('password').value
    };
    try {
        const response = await fetch('/api/connections/test', {
            method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(config)
        });
        const result = await response.json();
        showMessage(result.success ? '连接成功!' : '连接失败: ' + result.message,
            result.success ? 'success' : 'error');
    } catch (error) {
        showMessage('测试连接失败: ' + error.message, 'error');
    }
}

async function saveConnection() {
    const id = document.getElementById('connectionId').value;
    const config = {
        name: document.getElementById('connectionName').value,
        dbType: document.getElementById('dbType').value,
        host: document.getElementById('host').value,
        port: document.getElementById('port').value,
        schema: document.getElementById('schema').value,
        username: document.getElementById('username').value,
        password: document.getElementById('password').value,
        encoding: 'utf8mb4'
    };
    if (!config.name || !config.host || !config.port || !config.schema || !config.username) {
        showMessage('请填写所有必填项', 'error'); return;
    }
    try {
        const url = id ? `/api/connections/${id}` : '/api/connections';
        const method = id ? 'PUT' : 'POST';
        const response = await fetch(url, {
            method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(config)
        });
        const result = await response.json();
        if (response.ok) {
            showMessage('保存成功!', 'success');
            hideConnectionModal();
            loadConnections();
        } else {
            showMessage('保存失败: ' + result.error, 'error');
        }
    } catch (error) {
        showMessage('保存失败: ' + error.message, 'error');
    }
}

async function editConnection(id) {
    try {
        const response = await fetch('/api/connections');
        const connections = await response.json();
        const connection = connections.find(c => c.id === id);
        if (!connection) { showMessage('连接不存在', 'error'); return; }
        showConnectionModal(connection);
    } catch (error) {
        showMessage('加载连接信息失败: ' + error.message, 'error');
    }
}

async function deleteConnection(id) {
    if (!confirm('确定要删除这个连接吗?')) return;
    try {
        const response = await fetch(`/api/connections/${id}`, { method: 'DELETE' });
        if (response.ok) {
            showMessage('删除成功', 'success');
            loadConnections();
            if (currentDatabaseId === id) {
                currentDatabaseId = null;
                document.getElementById('tableList').innerHTML = '<div class="empty-placeholder">请先选择数据库连接</div>';
            }
        } else {
            showMessage('删除失败: ' + (await response.json()).error, 'error');
        }
    } catch (error) {
        showMessage('删除失败: ' + error.message, 'error');
    }
}

// ============================================================
// 代码生成
// ============================================================
async function generateCode() {
    if (!currentDatabaseId) { showMessage('请先选择数据库连接', 'error'); return; }
    if (selectedTables.length === 0) { showMessage('请先选择表', 'error'); return; }
    if (snippetMergeEnabled && selectedTables.length > 1) {
        showMessage('使用自定义片段时仅支持单张表，请取消多余的表勾选', 'error'); return;
    }
    const config = {
        modelPackage: document.getElementById('modelPackage').value,
        modelPackageTargetFolder: document.getElementById('modelTargetFolder').value,
        daoPackage: document.getElementById('daoPackage').value,
        daoTargetFolder: document.getElementById('daoTargetFolder').value,
        mappingXMLPackage: document.getElementById('mapperPackage').value,
        mappingXMLTargetFolder: document.getElementById('mapperTargetFolder').value,
        encoding: document.getElementById('encoding').value,
        offsetLimit: document.getElementById('offsetLimit').checked,
        comment: document.getElementById('comment').checked,
        overrideXML: document.getElementById('overrideXML').checked,
        useLombokPlugin: document.getElementById('useLombokPlugin').checked,
        jsr310Support: document.getElementById('jsr310Support').checked,
        needToStringHashcodeEquals: document.getElementById('needToStringHashcodeEquals').checked,
        needConstructors: document.getElementById('needConstructors').checked,
        useJsonProperty: document.getElementById('useJsonProperty').checked,
        jsonPropertyUpperCase: document.getElementById('jsonPropertyUpperCase').checked,
        useBatchInsert: document.getElementById('useBatchInsert').checked,
        useBatchUpdate: document.getElementById('useBatchUpdate').checked,
        ignorePKOnInsert: document.getElementById('ignorePKOnInsert').checked,
        needForUpdate: document.getElementById('needForUpdate').checked,
        useTableNameAlias: document.getElementById('useTableNameAlias').checked,
        useActualColumnNames: document.getElementById('useActualColumnNames').checked,
        ignoredColumns, columnOverrides
    };
    const requestBody = { databaseId: currentDatabaseId, tableNames: selectedTables, config };
    if (snippetMergeEnabled && snippetList.length > 0) {
        requestBody.snippetConfigs = snippetList;
    }
    try {
        const hint = snippetMergeEnabled && snippetList.length > 0
            ? `正在生成代码并追加 ${snippetList.length} 个自定义片段...`
            : `正在生成 ${selectedTables.length} 张表的代码...`;
        showMessage(hint, 'info');
        const response = await fetch('/api/generate', {
            method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(requestBody)
        });
        const result = await response.json();
        if (response.ok && result.success) {
            showMessage('代码生成成功! 正在准备下载...', 'success');
            const a = document.createElement('a');
            a.href = `/api/download/${result.downloadId}`;
            a.download = `generated_${selectedTables.length}_tables.zip`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            setTimeout(() => showMessage(`已生成 ${result.tableCount} 张表, 共 ${result.files.length} 个文件`, 'info'), 1000);
        } else {
            showMessage('代码生成失败: ' + result.error, 'error');
        }
    } catch (error) {
        showMessage('代码生成失败: ' + error.message, 'error');
    }
}

async function saveConfig() {
    const name = prompt('请输入配置名称:');
    if (!name) return;
    const config = {
        name,
        projectFolder: '',
        modelPackage: document.getElementById('modelPackage').value,
        modelPackageTargetFolder: document.getElementById('modelTargetFolder').value,
        daoPackage: document.getElementById('daoPackage').value,
        daoTargetFolder: document.getElementById('daoTargetFolder').value,
        mappingXMLPackage: document.getElementById('mapperPackage').value,
        mappingXMLTargetFolder: document.getElementById('mapperTargetFolder').value,
        encoding: document.getElementById('encoding').value,
        offsetLimit: document.getElementById('offsetLimit').checked,
        comment: document.getElementById('comment').checked,
        overrideXML: document.getElementById('overrideXML').checked,
        useLombokPlugin: document.getElementById('useLombokPlugin').checked,
        jsr310Support: document.getElementById('jsr310Support').checked,
        needToStringHashcodeEquals: document.getElementById('needToStringHashcodeEquals').checked,
        needConstructors: document.getElementById('needConstructors').checked,
        useJsonProperty: document.getElementById('useJsonProperty').checked,
        jsonPropertyUpperCase: document.getElementById('jsonPropertyUpperCase').checked,
        useBatchInsert: document.getElementById('useBatchInsert').checked,
        useBatchUpdate: document.getElementById('useBatchUpdate').checked,
        ignorePKOnInsert: document.getElementById('ignorePKOnInsert').checked,
        needForUpdate: document.getElementById('needForUpdate').checked,
        useTableNameAlias: document.getElementById('useTableNameAlias').checked,
        useActualColumnNames: document.getElementById('useActualColumnNames').checked
    };
    try {
        const response = await fetch('/api/generator-configs', {
            method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(config)
        });
        showMessage(response.ok ? '配置保存成功!' : '配置保存失败: ' + (await response.json()).error,
            response.ok ? 'success' : 'error');
    } catch (error) {
        showMessage('配置保存失败: ' + error.message, 'error');
    }
}

// ============================================================
// 列定制
// ============================================================
async function showColumnModal() {
    if (!currentDatabaseId) { showMessage('请先选择数据库连接', 'error'); return; }
    if (selectedTables.length === 0) { showMessage('请先选择表', 'error'); return; }
    const selector = document.getElementById('columnTableSelector');
    selector.innerHTML = '';
    selectedTables.forEach(tableName => {
        const option = document.createElement('option');
        option.value = tableName;
        option.textContent = tableName;
        selector.appendChild(option);
    });
    await loadColumnsForTable(selectedTables[0]);
    document.getElementById('columnModal').style.display = 'block';
}

async function loadColumnsForTable(tableName) {
    if (!tableName) return;
    try {
        const response = await fetch('/api/columns', {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ databaseId: currentDatabaseId, tableName })
        });
        const columns = await response.json();
        const tbody = document.getElementById('columnTableBody');
        tbody.innerHTML = '';
        columns.forEach(col => {
            const isIgnored = ignoredColumns.includes(col.columnName);
            const override = columnOverrides.find(o => o.columnName === col.columnName) || {};
            const row = document.createElement('tr');
            row.innerHTML = `
                <td class="text-center"><input type="checkbox" class="col-ignore" data-column="${col.columnName}" ${isIgnored ? 'checked' : ''}></td>
                <td>${col.columnName}</td>
                <td>${col.dataType}</td>
                <td><input type="text" class="form-input col-property" data-column="${col.columnName}" value="${override.propertyName || ''}" placeholder="默认自动转换"></td>
                <td><input type="text" class="form-input col-javatype" data-column="${col.columnName}" value="${override.javaType || ''}" placeholder="默认自动推断" list="javaTypeList"></td>`;
            tbody.appendChild(row);
        });
        if (!document.getElementById('javaTypeList')) {
            const datalist = document.createElement('datalist');
            datalist.id = 'javaTypeList';
            datalist.innerHTML = `
                <option value="String"><option value="Integer"><option value="Long">
                <option value="Double"><option value="Float"><option value="BigDecimal">
                <option value="Boolean"><option value="Date"><option value="LocalDate">
                <option value="LocalDateTime"><option value="LocalTime">
                <option value="byte[]"><option value="Byte"><option value="Short">`;
            document.body.appendChild(datalist);
        }
    } catch (error) {
        showMessage('加载列信息失败: ' + error.message, 'error');
    }
}

function hideColumnModal() {
    document.getElementById('columnModal').style.display = 'none';
}

function applyColumnSettings() {
    ignoredColumns = [];
    columnOverrides = [];
    document.querySelectorAll('.col-ignore:checked').forEach(cb => ignoredColumns.push(cb.dataset.column));
    document.querySelectorAll('.col-property').forEach(input => {
        const propertyName = input.value.trim();
        const javaTypeInput = document.querySelector(`.col-javatype[data-column="${input.dataset.column}"]`);
        const javaType = javaTypeInput ? javaTypeInput.value.trim() : '';
        if (propertyName || javaType) columnOverrides.push({ columnName: input.dataset.column, propertyName, javaType });
    });
    hideColumnModal();
    let parts = [];
    if (ignoredColumns.length > 0) parts.push(`${ignoredColumns.length} 个字段将被忽略`);
    if (columnOverrides.length > 0) parts.push(`${columnOverrides.length} 个字段已自定义`);
    showMessage(parts.length > 0 ? `列设置已保存：${parts.join('，')}` : '列设置已保存，未做任何修改', 'success');
}

// ============================================================
// Tab2：刷新面板状态
// ============================================================
function refreshSnippetPanelState() {
    const warningMulti = document.getElementById('snippetMultiTableWarning');
    const warningNone = document.getElementById('snippetNoTableWarning');
    const panel = document.getElementById('snippetPanel');
    const countEl = document.getElementById('snippetSelectedCount');

    if (selectedTables.length === 0) {
        warningMulti.style.display = 'none';
        warningNone.style.display = 'block';
        panel.style.display = 'none';
    } else if (selectedTables.length > 1) {
        countEl.textContent = selectedTables.length;
        warningMulti.style.display = 'block';
        warningNone.style.display = 'none';
        panel.style.display = 'none';
    } else {
        warningMulti.style.display = 'none';
        warningNone.style.display = 'none';
        panel.style.display = 'block';
        const tableName = selectedTables[0];
        document.getElementById('snippetCurrentTable').textContent = tableName;
        document.getElementById('snippetCurrentModel').textContent = toPascalCase(tableName);
        loadSnippetTableColumns(tableName);
    }
}

async function loadSnippetTableColumns(tableName) {
    if (!currentDatabaseId || !tableName) return;
    try {
        const response = await fetch('/api/columns', {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ databaseId: currentDatabaseId, tableName })
        });
        snippetTableColumns = await response.json();
        // Reset chip / where state when table changes
        resetSnippetFieldState();
        renderSnippetFieldPanel();
    } catch (error) {
        showMessage('加载列信息失败: ' + error.message, 'error');
    }
}

function resetSnippetFieldState() {
    whereRules = [];
    whereRuleCounter = 0;
    whereLogic = 'AND';
    selectedChips = { selectFields: new Set(), insertFields: new Set(), setFields: new Set() };
    orderBySelections = new Map();
}

// ============================================================
// Tab2：字段面板渲染（QueryBuilder 风格）
// ============================================================
function renderSnippetFieldPanel() {
    const operation = document.getElementById('snippetOperation').value;
    const isBatch = document.getElementById('snippetIsBatch').checked;
    const container = document.getElementById('snippetFieldPanels');
    const batchHintEl = document.getElementById('snippetBatchHint');

    const batchHints = {
        select: 'IN 查询：WHERE 条件取第一个字段作为 IN 条件，返回列表',
        insert: '批量插入：使用 foreach 循环插入实体列表',
        delete: 'IN 删除：WHERE 条件取第一个字段作为 IN 条件',
        update: '批量更新：使用 foreach 循环更新实体列表'
    };
    if (isBatch) {
        batchHintEl.textContent = '💡 ' + (batchHints[operation] || '');
        batchHintEl.style.display = 'block';
    } else {
        batchHintEl.style.display = 'none';
    }

    let html = '';
    if (operation === 'select') {
        html += buildChipPanel('selectFields', '📤 SELECT 返回字段', '点击字段卡片切换选中状态，已选字段将出现在 SELECT 列表中');
        html += buildQueryBuilderPanel();
        html += buildOrderByChipPanel();
    } else if (operation === 'insert') {
        html += buildChipPanel('insertFields', '📥 INSERT 字段', '点击字段卡片切换选中状态，已选字段将出现在 INSERT 语句中');
    } else if (operation === 'delete') {
        html += buildQueryBuilderPanel();
    } else if (operation === 'update') {
        html += buildChipPanel('setFields', '✏️ SET 更新字段', '点击字段卡片切换选中状态，已选字段将被更新');
        html += buildQueryBuilderPanel();
    }

    container.innerHTML = html;
    restoreChipStates();
    renderWhereRules();
    syncCombinatorButtons();
    updateMethodNamePlaceholder();
}

function updateMethodNamePlaceholder() {
    const cfg = buildCurrentSnippetConfig();
    const autoName = computeMethodName(cfg);
    const input = document.getElementById('snippetMethodName');
    if (input) {
        input.placeholder = `留空自动生成`;
    }
    // 带动预览 badge
    let preview = document.getElementById('methodNamePreview');
    if (!preview) {
        const methodGroup = input && input.closest('.form-group');
        if (methodGroup) {
            preview = document.createElement('div');
            preview.id = 'methodNamePreview';
            preview.className = 'snippet-method-preview';
            methodGroup.appendChild(preview);
        }
    }
    if (preview) {
        const userVal = input && input.value.trim();
        if (userVal) {
            preview.innerHTML = `<span class="snippet-method-preview-label">方法名：</span><span class="snippet-method-preview-name">${userVal}</span>`;
        } else {
            preview.innerHTML = `<span class="snippet-method-preview-label">自动生成：</span><span class="snippet-method-preview-name">${autoName}</span>`;
        }
    }
}

// -----------------------------------------------------------------------
// Chip 选择器
// -----------------------------------------------------------------------
function buildChipPanel(panelId, title, hint) {
    const chips = snippetTableColumns.map((col, idx) => `
        <div class="field-chip" data-panel="${panelId}" data-col-idx="${idx}" onclick="toggleChip('${panelId}', ${idx})">
            <span class="chip-col-name">${col.columnName}</span>
            <span class="chip-col-type">${col.dataType}</span>
        </div>`).join('');
    return `
        <div class="snippet-field-panel">
            <div class="snippet-field-panel-title">${title}</div>
            <div class="snippet-field-panel-hint">${hint}</div>
            <div class="field-chips" id="${panelId}Chips">
                ${chips || '<div class="snippet-empty">暂无列信息</div>'}
            </div>
        </div>`;
}

function buildOrderByChipPanel() {
    const chips = snippetTableColumns.map((col, idx) => `
        <div class="field-chip orderby-chip" data-col-idx="${idx}" onclick="toggleOrderByChip(${idx})">
            <span class="chip-col-name">${col.columnName}</span>
            <span class="chip-col-type">${col.dataType}</span>
            <select class="orderby-dir-select" id="orderby_dir_${idx}"
                    onclick="event.stopPropagation()" onchange="setOrderByDirection(${idx}, this.value)"
                    style="display:none;">
                <option value="ASC">ASC ↑</option>
                <option value="DESC">DESC ↓</option>
            </select>
        </div>`).join('');
    return `
        <div class="snippet-field-panel">
            <div class="snippet-field-panel-title">🔤 ORDER BY 排序字段</div>
            <div class="snippet-field-panel-hint">点击字段卡片添加排序，展开后可选择方向</div>
            <div class="field-chips" id="orderByFieldsChips">
                ${chips || '<div class="snippet-empty">暂无列信息</div>'}
            </div>
        </div>`;
}

function toggleChip(panelId, colIdx) {
    if (selectedChips[panelId].has(colIdx)) {
        selectedChips[panelId].delete(colIdx);
    } else {
        selectedChips[panelId].add(colIdx);
    }
    const chip = document.querySelector(`#${panelId}Chips .field-chip[data-col-idx="${colIdx}"]`);
    if (chip) chip.classList.toggle('selected', selectedChips[panelId].has(colIdx));
    updateMethodNamePlaceholder();
}

function toggleOrderByChip(colIdx) {
    const chip = document.querySelector(`#orderByFieldsChips .field-chip[data-col-idx="${colIdx}"]`);
    const dirSelect = document.getElementById('orderby_dir_' + colIdx);
    if (orderBySelections.has(colIdx)) {
        orderBySelections.delete(colIdx);
        if (chip) chip.classList.remove('selected');
        if (dirSelect) dirSelect.style.display = 'none';
    } else {
        orderBySelections.set(colIdx, 'ASC');
        if (chip) chip.classList.add('selected');
        if (dirSelect) { dirSelect.value = 'ASC'; dirSelect.style.display = ''; }
    }
    updateMethodNamePlaceholder();
}

function setOrderByDirection(colIdx, direction) {
    if (orderBySelections.has(colIdx)) orderBySelections.set(colIdx, direction);
}

function restoreChipStates() {
    ['selectFields', 'insertFields', 'setFields'].forEach(panelId => {
        const container = document.getElementById(panelId + 'Chips');
        if (!container) return;
        selectedChips[panelId].forEach(colIdx => {
            const chip = container.querySelector(`.field-chip[data-col-idx="${colIdx}"]`);
            if (chip) chip.classList.add('selected');
        });
    });
    const obContainer = document.getElementById('orderByFieldsChips');
    if (obContainer) {
        orderBySelections.forEach((direction, colIdx) => {
            const chip = obContainer.querySelector(`.field-chip[data-col-idx="${colIdx}"]`);
            if (chip) chip.classList.add('selected');
            const dirSelect = document.getElementById('orderby_dir_' + colIdx);
            if (dirSelect) { dirSelect.value = direction; dirSelect.style.display = ''; }
        });
    }
}

function collectChipFields(panelId) {
    const container = document.getElementById(panelId + 'Chips');
    if (!container) return [];
    const result = [];
    container.querySelectorAll('.field-chip.selected').forEach(chip => {
        const colIdx = parseInt(chip.dataset.colIdx);
        const col = snippetTableColumns[colIdx];
        if (col) {
            const override = columnOverrides.find(o => o.columnName === col.columnName) || {};
            result.push({
                columnName: col.columnName,
                fieldName: override.propertyName || col.fieldName || snakeToCamel(col.columnName),
                jdbcType: override.jdbcType || col.jdbcType || col.dataType.toUpperCase(),
                javaType: override.javaType || col.javaType || 'Object'
            });
        }
    });
    return result;
}

function collectOrderByFields() {
    const container = document.getElementById('orderByFieldsChips');
    if (!container) return [];
    const result = [];
    container.querySelectorAll('.field-chip.selected').forEach(chip => {
        const colIdx = parseInt(chip.dataset.colIdx);
        const col = snippetTableColumns[colIdx];
        const direction = orderBySelections.get(colIdx) || 'ASC';
        if (col) {
            const override = columnOverrides.find(o => o.columnName === col.columnName) || {};
            result.push({
                columnName: col.columnName,
                fieldName: override.propertyName || col.fieldName || snakeToCamel(col.columnName),
                jdbcType: override.jdbcType || col.jdbcType || col.dataType.toUpperCase(),
                direction
            });
        }
    });
    return result;
}

function snakeToCamel(s) {
    return s.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
}

// -----------------------------------------------------------------------
// QueryBuilder WHERE 构建器
// -----------------------------------------------------------------------
const QB_OPERATORS = ['=', '!=', '>', '<', '>=', '<=', 'LIKE', 'IN', 'NOT IN', 'IS NULL', 'IS NOT NULL'];
const QB_OP_LABELS = { '=': '= 等于', '!=': '≠ 不等于', '>': '> 大于', '<': '< 小于',
    '>=': '≥ 大于等于', '<=': '≤ 小于等于', 'LIKE': '≈ 模糊匹配', 'IN': '∈ 包含', 'NOT IN': '∉ 不包含',
    'IS NULL': '∅ 为空', 'IS NOT NULL': '∈ 非空' };

function buildQueryBuilderPanel() {
    const isOr = whereLogic === 'OR';
    return `
        <div class="snippet-field-panel">
            <div class="snippet-field-panel-title">🔍 WHERE 条件构建器</div>
            <div class="snippet-field-panel-hint">参考 react-querybuilder 风格，左侧彩色块标识 AND/OR 分组，支持多种运算符</div>
            <div class="qb-container">
                <div class="qb-group ${isOr ? 'qb-group-or' : ''}">
                    <div class="qb-header">
                        <div class="qb-combinator-group">
                            <button class="qb-comb-btn ${!isOr ? 'active' : ''}" id="whereLogicAnd" onclick="setWhereLogic('AND')">AND</button>
                            <button class="qb-comb-btn ${isOr ? 'active' : ''}" id="whereLogicOr" onclick="setWhereLogic('OR')">OR</button>
                        </div>
                        <button class="qb-add-rule-btn" onclick="addWhereRule()">＋ 添加条件</button>
                    </div>
                    <div class="qb-rules-list" id="whereRulesList">
                        <!-- 由 renderWhereRules() 填充 -->
                    </div>
                </div>
            </div>
        </div>`;
}

function addWhereRule() {
    const id = whereRuleCounter++;
    whereRules.push({ id, fieldIdx: 0, operator: '=' });
    renderWhereRules();
    updateMethodNamePlaceholder();
}

function removeWhereRule(id) {
    whereRules = whereRules.filter(r => r.id !== id);
    renderWhereRules();
    updateMethodNamePlaceholder();
}

function updateWhereRule(id, key, rawValue) {
    const rule = whereRules.find(r => r.id === id);
    if (rule) rule[key] = key === 'fieldIdx' ? parseInt(rawValue) : rawValue;
    updateMethodNamePlaceholder();
}

function setWhereLogic(logic) {
    whereLogic = logic;
    syncCombinatorButtons();
    // Re-render connector labels
    renderWhereRules();
}

function syncCombinatorButtons() {
    const andBtn = document.getElementById('whereLogicAnd');
    const orBtn = document.getElementById('whereLogicOr');
    if (andBtn) andBtn.classList.toggle('active', whereLogic === 'AND');
    if (orBtn) orBtn.classList.toggle('active', whereLogic === 'OR');
    // 同时更新 qb-group 的左侧彩色边框
    const group = document.querySelector('.qb-group');
    if (group) group.classList.toggle('qb-group-or', whereLogic === 'OR');
}

function renderWhereRules() {
    const container = document.getElementById('whereRulesList');
    if (!container) return;
    if (whereRules.length === 0) {
        container.innerHTML = '<div class="qb-empty">暂无条件，点击上方"添加条件"按钮</div>';
        return;
    }
    const isOr = whereLogic === 'OR';
    const badgeClass = isOr ? 'or-badge' : '';
    const badgeText = isOr ? 'OR' : 'AND';
    const parts = whereRules.map((rule, idx) => {
        const fieldOptions = snippetTableColumns.map((col, i) =>
            `<option value="${i}" ${rule.fieldIdx === i ? 'selected' : ''}>${col.columnName}  (${col.dataType})</option>`
        ).join('');
        const opOptions = QB_OPERATORS.map(op =>
            `<option value="${op}" ${rule.operator === op ? 'selected' : ''}>${QB_OP_LABELS[op] || op}</option>`
        ).join('');
        const connector = idx > 0 ? `
            <div class="qb-rule-connector">
                <div class="qb-connector-line"></div>
                <span class="qb-connector-badge ${badgeClass}">${badgeText}</span>
                <div class="qb-connector-line"></div>
            </div>` : '';
        return `${connector}
            <div class="qb-rule" data-rule-id="${rule.id}">
                <span class="qb-rule-number">${idx + 1}.</span>
                <select class="qb-field-select" onchange="updateWhereRule(${rule.id}, 'fieldIdx', this.value)">
                    ${fieldOptions}
                </select>
                <select class="qb-op-select" onchange="updateWhereRule(${rule.id}, 'operator', this.value)">
                    ${opOptions}
                </select>
                <button class="qb-remove-btn" onclick="removeWhereRule(${rule.id})" title="删除此条件">✕</button>
            </div>`;
    });
    container.innerHTML = parts.join('');
}

function collectWhereConditions() {
    return whereRules.map(rule => {
        const col = snippetTableColumns[rule.fieldIdx];
        if (!col) return null;
        const override = columnOverrides.find(o => o.columnName === col.columnName) || {};
        return {
            columnName: col.columnName,
            fieldName: override.propertyName || col.fieldName || snakeToCamel(col.columnName),
            jdbcType: override.jdbcType || col.jdbcType || col.dataType.toUpperCase(),
            javaType: override.javaType || col.javaType || 'Object',
            operator: rule.operator || '='
        };
    }).filter(Boolean);
}

// ============================================================
// 片段操作
// ============================================================
function buildCurrentSnippetConfig() {
    const operation = document.getElementById('snippetOperation').value;
    const isBatch = document.getElementById('snippetIsBatch').checked;
    const methodName = document.getElementById('snippetMethodName').value.trim();
    const cfg = {
        operation, isBatch, methodName,
        whereLogic,
        selectFields: [], whereFields: [], orderByFields: [], insertFields: [], setFields: []
    };
    const whereConditions = collectWhereConditions();
    if (operation === 'select') {
        cfg.selectFields = collectChipFields('selectFields');
        cfg.whereFields = whereConditions;
        cfg.orderByFields = collectOrderByFields();
    } else if (operation === 'insert') {
        cfg.insertFields = collectChipFields('insertFields');
    } else if (operation === 'delete') {
        cfg.whereFields = whereConditions;
    } else if (operation === 'update') {
        cfg.setFields = collectChipFields('setFields');
        cfg.whereFields = whereConditions;
    }
    return cfg;
}

function addSnippet() {
    const cfg = buildCurrentSnippetConfig();
    // 校验
    const hasFields =
        cfg.selectFields.length > 0 || cfg.whereFields.length > 0 ||
        cfg.insertFields.length > 0 || cfg.setFields.length > 0 || cfg.orderByFields.length > 0;
    if (!hasFields) { showMessage('请至少配置一个字段或条件', 'error'); return; }
    
    if (!cfg.methodName) {
        cfg.methodName = computeMethodName(cfg);
    }
    
    // Check for duplicate method names
    const duplicateIdx = snippetList.findIndex((s, idx) => s.methodName === cfg.methodName && idx !== editingSnippetIndex);
    if (duplicateIdx !== -1) {
        showMessage(`方法名 '${cfg.methodName}' 已存在，请手动修改以避免冲突！`, 'error');
        return;
    }
    
    if (editingSnippetIndex !== null) {
        snippetList[editingSnippetIndex] = cfg;
        showMessage('片段已保存修改', 'success');
        cancelEditSnippet();
    } else {
        snippetList.push(cfg);
        showMessage(`片段已添加（共 ${snippetList.length} 个）`, 'success');
    }
    renderSnippetList();
}

function updateSnippetMethodName(idx, val) {
    if (snippetList[idx]) {
        const newName = val.trim();
        const duplicateIdx = snippetList.findIndex((s, i) => s.methodName === newName && i !== idx);
        if (duplicateIdx !== -1 && newName !== "") {
            showMessage(`方法名 '${newName}' 已存在，请使用其他名称！`, 'error');
            return;
        }
        snippetList[idx].methodName = newName;
        showMessage('方法名已更新', 'success');
    }
}

// 编辑片段：加载回表单，修改按钮状态
function editSnippet(idx) {
    const cfg = snippetList[idx];
    editingSnippetIndex = idx;
    loadSnippetIntoForm(cfg);
    
    const btnAdd = document.getElementById('btnAddSnippet');
    btnAdd.innerHTML = '💾 保存修改';
    btnAdd.classList.add('btn-edit-mode');
    
    let btnCancel = document.getElementById('btnCancelEdit');
    if (!btnCancel) {
        btnCancel = document.createElement('button');
        btnCancel.id = 'btnCancelEdit';
        btnCancel.className = 'btn btn-secondary';
        btnCancel.style.marginLeft = '10px';
        btnCancel.style.fontSize = '18px';
        btnCancel.style.padding = '15px 30px';
        btnCancel.style.borderRadius = '50px';
        btnCancel.innerHTML = '✕ 取消修改';
        btnCancel.onclick = cancelEditSnippet;
        btnAdd.parentNode.appendChild(btnCancel);
    }
    btnCancel.style.display = 'inline-block';
    
    showMessage('片段已加载到编辑区，修改后点击保存', 'info');
    document.getElementById('snippetPanel').scrollIntoView({ behavior: 'smooth' });
}

function cancelEditSnippet() {
    editingSnippetIndex = null;
    document.getElementById('btnAddSnippet').innerHTML = '＋ 添加当前片段';
    document.getElementById('btnAddSnippet').classList.remove('btn-edit-mode');
    const btnCancel = document.getElementById('btnCancelEdit');
    if (btnCancel) {
        btnCancel.style.display = 'none';
    }
    // reset form by re-rendering
    document.getElementById('snippetMethodName').value = '';
    resetSnippetFieldState();
    renderSnippetFieldPanel();
}

function loadSnippetIntoForm(cfg) {
    document.getElementById('snippetOperation').value = cfg.operation;
    document.getElementById('snippetIsBatch').checked = cfg.isBatch;
    document.getElementById('snippetMethodName').value = cfg.methodName || '';
    // 恢复 WHERE 状态
    whereRules = [];
    whereRuleCounter = 0;
    whereLogic = cfg.whereLogic || 'AND';
    // 恢复 chip 状态
    selectedChips = { selectFields: new Set(), insertFields: new Set(), setFields: new Set() };
    orderBySelections = new Map();
    // 恢复 WHERE rules
    (cfg.whereFields || []).forEach(f => {
        const colIdx = snippetTableColumns.findIndex(c => c.columnName === f.columnName);
        if (colIdx >= 0) {
            whereRules.push({ id: whereRuleCounter++, fieldIdx: colIdx, operator: f.operator || '=' });
        }
    });
    // 恢复 chip 面板
    const chipMap = { selectFields: 'selectFields', insertFields: 'insertFields', setFields: 'setFields' };
    Object.entries(chipMap).forEach(([cfgKey, panelId]) => {
        (cfg[cfgKey] || []).forEach(f => {
            const colIdx = snippetTableColumns.findIndex(c => c.columnName === f.columnName);
            if (colIdx >= 0) selectedChips[panelId].add(colIdx);
        });
    });
    // 恢复 ORDER BY
    (cfg.orderByFields || []).forEach(f => {
        const colIdx = snippetTableColumns.findIndex(c => c.columnName === f.columnName);
        if (colIdx >= 0) orderBySelections.set(colIdx, f.direction || 'ASC');
    });
    renderSnippetFieldPanel();
}

// 渲染已添加的片段列表
function renderSnippetList() {
    const container = document.getElementById('snippetItems');
    const countEl = document.getElementById('snippetCount');
    if (countEl) countEl.textContent = snippetList.length + ' 个';
    if (snippetList.length === 0) {
        container.innerHTML = '<div class="snippet-empty">暂未添加任何片段</div>';
        return;
    }
    const opLabels = { select: '查询', insert: '新增', delete: '删除', update: '更新' };
    const opColors = { select: '#3b82f6', insert: '#10b981', delete: '#ef4444', update: '#f97316' };
    container.innerHTML = snippetList.map((s, i) => {
        const label = `${opLabels[s.operation] || s.operation}${s.isBatch ? '(批量)' : ''}`;
        const autoName = computeMethodName(s);
        const isAutoName = !s.methodName || s.methodName === autoName;
        const displayName = s.methodName || autoName;
        const fieldCount =
            (s.selectFields || []).length + (s.insertFields || []).length +
            (s.setFields || []).length + (s.whereFields || []).length +
            (s.orderByFields || []).length;
        const badgeBg = opColors[s.operation] || '#667eea';
        return `
            <div class="snippet-item">
                <span class="snippet-item-badge" style="background:${badgeBg}">${label}</span>
                <div style="display:flex; flex-direction:column; gap:3px; flex:1; min-width:0;">
                    <div style="display:flex; align-items:center; gap:8px;">
                        <input type="text" class="snippet-item-method-input"
                            value="${displayName}"
                            onchange="updateSnippetMethodName(${i}, this.value)"
                            placeholder="方法名"
                            title="点击可修改方法名">
                        ${isAutoName ? '<span class="snippet-auto-label">🧠 自动生成</span>' : ''}
                    </div>
                    <span class="snippet-item-meta">${fieldCount} 个字段配置</span>
                </div>
                <div class="snippet-item-actions">
                    <button class="btn btn-sm btn-info" onclick="editSnippet(${i})" title="加载到编辑区修改">✏️ 编辑</button>
                    <button class="btn btn-sm btn-danger" onclick="removeSnippet(${i})">🗑️</button>
                </div>
            </div>`;
    }).join('');
}


function removeSnippet(idx) {
    snippetList.splice(idx, 1);
    renderSnippetList();
}

function clearSnippets() {
    if (snippetList.length === 0) return;
    if (!confirm(`确定清空全部 ${snippetList.length} 个自定义片段吗？`)) return;
    snippetList = [];
    renderSnippetList();
    if (snippetMergeEnabled) toggleSnippetMerge();
    showMessage('已清空所有片段', 'success');
}

function toggleSnippetMerge() {
    if (snippetList.length === 0 && !snippetMergeEnabled) {
        showMessage('请先添加至少一个自定义片段', 'error'); return;
    }
    snippetMergeEnabled = !snippetMergeEnabled;
    const btn = document.getElementById('btnSnippetMerge');
    const hint = document.getElementById('snippetMergeHint');
    if (snippetMergeEnabled) {
        btn.textContent = '🔗 并入生成（已启用）';
        btn.classList.remove('btn-success');
        btn.classList.add('btn-warning');
        hint.style.display = 'block';
    } else {
        btn.textContent = '🔗 并入生成（未启用）';
        btn.classList.remove('btn-warning');
        btn.classList.add('btn-success');
        hint.style.display = 'none';
    }
}

async function previewSnippet() {
    if (snippetList.length === 0) { showMessage('请先添加至少一个自定义片段', 'error'); return; }
    if (selectedTables.length !== 1) { showMessage('请先选择一张表', 'error'); return; }
    const tableName = selectedTables[0];
    const mapperName = toPascalCase(tableName) + 'Mapper';
    const modelPackage = document.getElementById('modelPackage').value || 'com.example.model';
    const modelType = modelPackage + '.' + toPascalCase(tableName);
    try {
        showMessage('正在生成预览...', 'info');
        const response = await fetch('/api/snippet/preview', {
            method: 'POST', headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ tableName, mapperName, modelType, snippetConfigs: snippetList })
        });
        const result = await response.json();
        if (response.ok && result.success) {
            document.getElementById('snippetJavaCode').textContent = result.javaCode;
            document.getElementById('snippetXmlCode').textContent = result.xmlCode;
            document.getElementById('snippetPreviewModal').style.display = 'block';
        } else {
            showMessage('预览失败: ' + result.error, 'error');
        }
    } catch (error) {
        showMessage('预览失败: ' + error.message, 'error');
    }
}

function hideSnippetPreviewModal() {
    document.getElementById('snippetPreviewModal').style.display = 'none';
}

// ============================================================
// 事件监听
// ============================================================
document.addEventListener('DOMContentLoaded', function () {
    loadConnections();
    document.getElementById('btnNewConnection').onclick = () => showConnectionModal();
    document.querySelectorAll('.close').forEach(el => {
        el.onclick = function (e) {
            e.stopPropagation();
            if (el.closest('#columnModal')) hideColumnModal();
            else if (el.closest('#snippetPreviewModal')) hideSnippetPreviewModal();
            else hideConnectionModal();
        };
    });
    document.querySelectorAll('.close-modal').forEach(el => {
        el.onclick = function () { hideConnectionModal(); };
    });
    document.getElementById('btnTestConnection').onclick = testConnection;
    document.getElementById('btnSaveConnection').onclick = saveConnection;
    document.getElementById('btnGenerate').onclick = generateCode;
    document.getElementById('btnSaveConfig').onclick = saveConfig;
    document.getElementById('tableFilter').oninput = e => loadTables(e.target.value);
    document.getElementById('dbType').onchange = e => {
        const ports = { 'MySQL': '3306', 'PostgreSQL': '5432', 'Oracle': '1521' };
        document.getElementById('port').value = ports[e.target.value] || '3306';
    };
    document.getElementById('useJsonProperty').onchange = e => {
        const lbl = document.getElementById('jsonPropertyOptionsLabel');
        lbl.style.display = e.target.checked ? 'flex' : 'none';
        if (!e.target.checked) document.getElementById('jsonPropertyUpperCase').checked = false;
    };
    document.addEventListener('keydown', e => {
        if (e.key === 'Escape') { hideConnectionModal(); hideColumnModal(); hideSnippetPreviewModal(); }
    });
    document.getElementById('btnCustomizeColumns').onclick = showColumnModal;
    document.getElementById('btnApplyColumns').onclick = applyColumnSettings;
    // 标记手动编辑
    document.getElementById('domainObjectName').addEventListener('input', function () {
        this.dataset.userEdited = this.value ? '1' : '';
    });
    document.getElementById('mapperName').addEventListener('input', function () {
        this.dataset.userEdited = this.value ? '1' : '';
    });
});
