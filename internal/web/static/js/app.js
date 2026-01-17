// 全局变量
let currentDatabaseId = null;
let currentTableName = null;

// 显示消息提示
function showMessage(message, type = 'info') {
    const messageEl = document.getElementById('message');
    messageEl.textContent = message;
    messageEl.className = `message ${type}`;
    messageEl.style.display = 'block';

    // 根据类型设置不同的显示时长
    const duration = type === 'success' ? 5000 : type === 'error' ? 4000 : 3000;

    setTimeout(() => {
        messageEl.style.display = 'none';
    }, duration);
}

function toPascalCase(str) {
    return str.split('_').map(word =>
        word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
    ).join('');
}

// 加载数据库连接列表
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
                </div>
            `;
            list.appendChild(item);
        });
    } catch (error) {
        showMessage('加载连接列表失败: ' + error.message, 'error');
    }
}

// 选择数据库连接
async function selectConnection(id, name) {
    currentDatabaseId = id;

    // 高亮选中的连接
    document.querySelectorAll('.connection-item').forEach(item => {
        item.classList.remove('active');
    });
    event.currentTarget.closest('.connection-item').classList.add('active');

    // 加载表列表
    await loadTables();
}

// 加载表列表
async function loadTables(filter = '') {
    if (!currentDatabaseId) return;

    try {
        const response = await fetch('/api/tables', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                databaseId: currentDatabaseId,
                filter: filter
            })
        });

        const tables = await response.json();

        const list = document.getElementById('tableList');
        list.innerHTML = '';

        tables.forEach(tableName => {
            const item = document.createElement('div');
            item.className = 'table-item';
            item.textContent = tableName;
            item.onclick = () => selectTable(tableName);
            list.appendChild(item);
        });
    } catch (error) {
        showMessage('加载表列表失败: ' + error.message, 'error');
    }
}

// 选择表
function selectTable(tableName) {
    currentTableName = tableName;

    // 填充表单
    document.getElementById('tableName').value = tableName;
    document.getElementById('domainObjectName').value = toPascalCase(tableName);
    document.getElementById('mapperName').value = toPascalCase(tableName) + 'Mapper';
}

// 新建/编辑连接
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

// 测试连接
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
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });

        const result = await response.json();

        if (result.success) {
            showMessage('连接成功!', 'success');
        } else {
            showMessage('连接失败: ' + result.message, 'error');
        }
    } catch (error) {
        showMessage('测试连接失败: ' + error.message, 'error');
    }
}

// 保存连接
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

    // 验证必填项
    if (!config.name || !config.dbType || !config.host ||
        !config.port || !config.schema || !config.username) {
        showMessage('请填写所有必填项', 'error');
        return; // 验证失败，不关闭对话框
    }

    try {
        const url = id ? `/api/connections/${id}` : '/api/connections';
        const method = id ? 'PUT' : 'POST';

        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });

        const result = await response.json();

        if (response.ok) {
            showMessage('保存成功!', 'success');
            hideConnectionModal(); // 只在成功时关闭
            loadConnections();
        } else {
            showMessage('保存失败: ' + result.error, 'error');
            // 失败时不关闭对话框，让用户修改
        }
    } catch (error) {
        showMessage('保存失败: ' + error.message, 'error');
        // 失败时不关闭对话框
    }
}

// 编辑连接
async function editConnection(id) {
    try {
        const response = await fetch('/api/connections');
        const connections = await response.json();
        const connection = connections.find(c => c.id === id);

        if (!connection) {
            showMessage('连接不存在', 'error');
            return;
        }

        // 填充表单字段
        document.getElementById('connectionId').value = connection.id; // Keep this for identifying the connection
        document.getElementById('connectionName').value = connection.name || '';
        document.getElementById('dbType').value = connection.dbType || 'MySQL'; // Default to MySQL if not set
        document.getElementById('host').value = connection.host || '';
        document.getElementById('port').value = connection.port || '';
        document.getElementById('schema').value = connection.schema || '';
        document.getElementById('username').value = connection.username || '';
        document.getElementById('password').value = connection.password || '';
        document.getElementById('encoding').value = connection.encoding || 'utf8mb4';

        // Update modal title
        document.getElementById('connectionModalTitle').textContent = '编辑数据库连接';

        // Show modal
        showConnectionModal(connection); // Pass connection object to showConnectionModal
    } catch (error) {
        showMessage('加载连接信息失败: ' + error.message, 'error');
    }
}

// 删除连接
async function deleteConnection(id) {
    if (!confirm('确定要删除这个连接吗?')) return;

    try {
        const response = await fetch(`/api/connections/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            showMessage('删除成功', 'success');
            loadConnections();

            // 如果删除的是当前选中的连接，清空表列表
            if (currentDatabaseId === id) {
                currentDatabaseId = null;
                document.getElementById('tableList').innerHTML = '';
            }
        } else {
            const result = await response.json();
            showMessage('删除失败: ' + result.error, 'error');
        }
    } catch (error) {
        showMessage('删除失败: ' + error.message, 'error');
    }
}

// 生成代码
async function generateCode() {
    if (!currentDatabaseId) {
        showMessage('请先选择数据库连接', 'error');
        return;
    }

    if (!currentTableName) {
        showMessage('请先选择表', 'error');
        return;
    }

    const config = {
        modelPackage: document.getElementById('modelPackage').value,
        modelPackageTargetFolder: document.getElementById('modelTargetFolder').value,
        daoPackage: document.getElementById('daoPackage').value,
        daoTargetFolder: document.getElementById('daoTargetFolder').value,
        mappingXMLPackage: document.getElementById('mapperPackage').value,
        mappingXMLTargetFolder: document.getElementById('mapperTargetFolder').value,
        tableName: document.getElementById('tableName').value,
        domainObjectName: document.getElementById('domainObjectName').value,
        mapperName: document.getElementById('mapperName').value,
        generateKeys: document.getElementById('generateKeys').value,
        encoding: document.getElementById('encoding').value,
        offsetLimit: document.getElementById('offsetLimit').checked,
        comment: document.getElementById('comment').checked,
        overrideXML: document.getElementById('overrideXML').checked,
        useLombokPlugin: document.getElementById('useLombokPlugin').checked,
        jsr310Support: document.getElementById('jsr310Support').checked
    };

    try {
        const response = await fetch('/api/generate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                databaseId: currentDatabaseId,
                config: config
            })
        });

        const result = await response.json();

        if (response.ok && result.success) {
            showMessage('代码生成成功! 正在准备下载...', 'success');

            // 自动触发下载
            const downloadUrl = `/api/download/${result.downloadId}`;
            const a = document.createElement('a');
            a.href = downloadUrl;
            a.download = `${currentTableName}_generated.zip`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);

            setTimeout(() => {
                showMessage(`已生成文件: ${result.files.join(', ')}`, 'info');
            }, 1000);
        } else {
            showMessage('代码生成失败: ' + result.error, 'error');
        }
    } catch (error) {
        showMessage('代码生成失败: ' + error.message, 'error');
    }
}

// 保存配置
async function saveConfig() {
    const name = prompt('请输入配置名称:');
    if (!name) return;

    const config = {
        name: name,
        projectFolder: document.getElementById('projectFolder').value,
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
        jsr310Support: document.getElementById('jsr310Support').checked
    };

    try {
        const response = await fetch('/api/generator-configs', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config)
        });

        if (response.ok) {
            showMessage('配置保存成功!', 'success');
        } else {
            const result = await response.json();
            showMessage('配置保存失败: ' + result.error, 'error');
        }
    } catch (error) {
        showMessage('配置保存失败: ' + error.message, 'error');
    }
}

// 事件监听
document.addEventListener('DOMContentLoaded', function () {
    // 加载连接列表
    loadConnections();

    // 新建连接按钮
    document.getElementById('btnNewConnection').onclick = () => showConnectionModal();

    // 模态框关闭按钮
    document.querySelectorAll('.close').forEach(el => {
        el.onclick = function (e) {
            e.stopPropagation();
            hideConnectionModal();
        };
    });

    document.querySelectorAll('.close-modal').forEach(el => {
        el.onclick = function (e) {
            e.stopPropagation();
            hideConnectionModal();
        };
    });

    // 不再监听模态框背景点击事件，避免失焦关闭
    // 用户只能通过关闭按钮或ESC键关闭对话框

    // 测试连接按钮
    document.getElementById('btnTestConnection').onclick = testConnection;

    // 保存连接按钮
    document.getElementById('btnSaveConnection').onclick = saveConnection;

    // 生成代码按钮
    document.getElementById('btnGenerate').onclick = generateCode;

    // 保存配置按钮
    document.getElementById('btnSaveConfig').onclick = saveConfig;

    // 表过滤
    document.getElementById('tableFilter').oninput = (e) => {
        loadTables(e.target.value);
    };

    // 数据库类型变化时更新默认端口
    document.getElementById('dbType').onchange = (e) => {
        const port = e.target.value === 'MySQL' ? '3306' : '5432';
        document.getElementById('port').value = port;
    };

    // ESC键关闭模态框
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape') {
            hideConnectionModal();
        }
    });
});
