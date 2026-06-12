const mysql = require('mysql2/promise');
const XLSX = require('xlsx');
const fs = require('fs');
const path = require('path');

// 数据库配置
const DB_CONFIG = {
    host: '192.168.31.23',
    port: 13306,
    user: 'root',
    password: 'Sun123456',
    database: 'gva',
    charset: 'utf8mb4'
};

// Excel文件路径
const EXCEL_FILE = path.join(__dirname, '战队信息表-2026-06-11 16_48_20(1).xlsx');

// 报告对象
const report = {
    startTime: new Date(),
    endTime: null,
    excelData: {
        totalRows: 0,
        validRows: 0,
        invalidRows: 0,
        invalidReasons: []
    },
    import: {
        teams: { imported: 0, skipped: 0, errors: [] },
        players: { imported: 0, skipped: 0, updated: 0, errors: [] }
    },
    validation: {
        teams: { expected: 0, actual: 0, matches: false },
        players: { expected: 0, actual: 0, matches: false },
        dataIntegrity: []
    }
};

async function main() {
    console.log('========== 选手数据完整导入流程 ==========\n');

    // Step 1: 读取并验证Excel数据
    console.log('Step 1: 读取并验证Excel数据');
    const { teams, players, validationErrors } = await readAndValidateExcel();
    report.excelData.totalRows = teams.length * 5 + players.length % 5;
    report.excelData.validRows = teams.length * 5 + (players.length % 5 === 0 ? 5 : players.length % 5);
    report.excelData.invalidRows = validationErrors.length;
    report.excelData.invalidReasons = validationErrors;

    console.log(`  读取战队: ${teams.length} 个`);
    console.log(`  读取选手: ${players.length} 人`);
    if (validationErrors.length > 0) {
        console.log(`  数据验证错误: ${validationErrors.length} 处`);
        validationErrors.forEach((err, i) => console.log(`    ${i + 1}. ${err}`));
    }
    console.log('');

    // Step 2: 清空现有数据（可选）
    console.log('Step 2: 准备数据库');
    await prepareDatabase();
    console.log('');

    // Step 3: 导入战队数据
    console.log('Step 3: 导入战队数据');
    await importTeams(teams);
    console.log('');

    // Step 4: 获取战队ID映射
    console.log('Step 4: 获取战队ID映射');
    const teamIdMap = await getTeamIdMap();
    console.log(`  获取到 ${teamIdMap.size} 个战队映射`);
    console.log('');

    // Step 5: 导入选手数据
    console.log('Step 5: 导入选手数据');
    await importPlayers(players, teamIdMap);
    console.log('');

    // Step 6: 验证导入结果
    console.log('Step 6: 验证导入结果');
    await validateImport(teams.length, players.length);
    console.log('');

    // Step 7: 生成报告
    console.log('Step 7: 生成导入报告');
    report.endTime = new Date();
    const reportText = generateReport();
    fs.writeFileSync('选手数据导入报告.txt', reportText, 'utf8');
    console.log('报告已保存: 选手数据导入报告.txt');
    console.log('');

    console.log('========== 导入流程完成 ==========');
    console.log(reportText);
}

async function readAndValidateExcel() {
    const workbook = XLSX.readFile(EXCEL_FILE);
    const worksheet = workbook.Sheets['Sheet1'];
    const data = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

    const teams = [];
    const players = [];
    const errors = [];
    const teamSet = new Set();

    for (let i = 1; i < data.length; i++) {
        const row = data[i];
        if (!row || row.length < 4) continue;

        const [teamId, teamAbbr, teamName, uid] = row;

        // 验证战队数据
        if (!teamId || !teamName) {
            errors.push(`第${i + 1}行: 战队ID或名称为空`);
            continue;
        }

        // 添加战队（去重）
        if (!teamSet.has(teamId)) {
            teamSet.add(teamId);
            teams.push({
                id: teamId,
                abbr: teamAbbr || '',
                name: teamName
            });
        }

        // 验证选手数据
        if (!uid || String(uid).trim() === '') {
            errors.push(`第${i + 1}行: 选手UID为空（战队: ${teamName}）`);
            continue;
        }

        players.push({
            teamId: teamId,
            teamName: teamName,
            uid: String(uid).trim()
        });
    }

    return { teams, players, validationErrors: errors };
}

async function prepareDatabase() {
    const connection = await mysql.createConnection(DB_CONFIG);
    try {
        // 清理本次导入的数据（保留原有数据）
        await connection.execute('DELETE FROM exa_players WHERE uid LIKE ?', ['%ac2%']); // 示例清理条件
        console.log('  数据库准备完成');
    } finally {
        await connection.end();
    }
}

async function importTeams(teams) {
    const connection = await mysql.createConnection(DB_CONFIG);
    try {
        for (const team of teams) {
            try {
                const [existing] = await connection.execute(
                    'SELECT id FROM exa_teams WHERE team_name = ?',
                    [team.name]
                );

                if (existing.length > 0) {
                    report.import.teams.skipped++;
                    console.log(`  [跳过] 战队 "${team.name}" 已存在`);
                } else {
                    await connection.execute(
                        'INSERT INTO exa_teams (team_name, total_bounty) VALUES (?, ?)',
                        [team.name, 0]
                    );
                    report.import.teams.imported++;
                    console.log(`  [成功] 战队 "${team.name}" 已导入`);
                }
            } catch (err) {
                report.import.teams.errors.push({ name: team.name, error: err.message });
                console.log(`  [错误] 战队 "${team.name}": ${err.message}`);
            }
        }
    } finally {
        await connection.end();
    }
}

async function getTeamIdMap() {
    const connection = await mysql.createConnection(DB_CONFIG);
    try {
        const [teams] = await connection.execute('SELECT id, team_name FROM exa_teams');
        const map = new Map();
        teams.forEach(t => map.set(t.team_name, t.id));
        return map;
    } finally {
        await connection.end();
    }
}

async function importPlayers(players, teamIdMap) {
    const connection = await mysql.createConnection(DB_CONFIG);
    try {
        // 统计每个战队的成员数
        const teamMemberCount = new Map();
        players.forEach(p => {
            const count = teamMemberCount.get(p.teamName) || 0;
            teamMemberCount.set(p.teamName, count + 1);
        });

        for (const player of players) {
            const teamId = teamIdMap.get(player.teamName);
            if (!teamId) {
                report.import.players.errors.push({ uid: player.uid, error: `战队 "${player.teamName}" 不存在` });
                console.log(`  [错误] UID=${player.uid}: 战队不存在`);
                continue;
            }

            // 生成选手姓名
            const count = teamMemberCount.get(player.teamName);
            const playerName = `${player.teamName}_${count}`;
            teamMemberCount.set(player.teamName, count - 1);

            try {
                const [existing] = await connection.execute(
                    'SELECT id FROM exa_players WHERE uid = ?',
                    [player.uid]
                );

                if (existing.length > 0) {
                    // 更新现有记录
                    await connection.execute(
                        'UPDATE exa_players SET player_name = ?, team_id = ?, bounty = ? WHERE uid = ?',
                        [playerName, teamId, 0, player.uid]
                    );
                    report.import.players.updated++;
                    console.log(`  [更新] UID=${player.uid} -> "${playerName}"`);
                } else {
                    // 插入新记录
                    await connection.execute(
                        'INSERT INTO exa_players (player_name, uid, team_id, bounty) VALUES (?, ?, ?, ?)',
                        [playerName, player.uid, teamId, 0]
                    );
                    report.import.players.imported++;
                    console.log(`  [新增] UID=${player.uid} -> "${playerName}"`);
                }
            } catch (err) {
                report.import.players.errors.push({ uid: player.uid, error: err.message });
                console.log(`  [错误] UID=${player.uid}: ${err.message}`);
            }
        }
    } finally {
        await connection.end();
    }
}

async function validateImport(expectedTeams, expectedPlayers) {
    const connection = await mysql.createConnection(DB_CONFIG);
    try {
        // 验证战队数量
        const [teamCount] = await connection.execute('SELECT COUNT(*) as count FROM exa_teams');
        report.validation.teams.expected = expectedTeams;
        report.validation.teams.actual = teamCount[0].count;
        report.validation.teams.matches = expectedTeams === teamCount[0].count;

        // 验证选手数量
        const [playerCount] = await connection.execute('SELECT COUNT(*) as count FROM exa_players');
        report.validation.players.expected = expectedPlayers;
        report.validation.players.actual = playerCount[0].count;
        report.validation.players.matches = expectedPlayers === playerCount[0].count;

        // 验证数据完整性
        const [invalidPlayers] = await connection.execute(`
            SELECT p.id, p.player_name, p.uid, p.team_id, t.team_name
            FROM exa_players p
            LEFT JOIN exa_teams t ON p.team_id = t.id
            WHERE p.player_name = p.uid OR t.team_name IS NULL
        `);

        invalidPlayers.forEach(p => {
            const issues = [];
            if (p.player_name === p.uid) issues.push('姓名与UID相同');
            if (!p.team_name) issues.push('战队关联缺失');
            report.validation.dataIntegrity.push({
                id: p.id,
                playerName: p.player_name,
                uid: p.uid,
                issues: issues.join(', ')
            });
        });

        console.log(`  战队数量验证: 期望 ${expectedTeams}, 实际 ${teamCount[0].count} [${report.validation.teams.matches ? '✓匹配' : '✗不匹配'}]`);
        console.log(`  选手数量验证: 期望 ${expectedPlayers}, 实际 ${playerCount[0].count} [${report.validation.players.matches ? '✓匹配' : '✗不匹配'}]`);
        console.log(`  数据完整性问题: ${report.validation.dataIntegrity.length} 条`);

    } finally {
        await connection.end();
    }
}

function generateReport() {
    const duration = (report.endTime - report.startTime) / 1000;

    let text = '========== 选手数据导入报告 ==========\n';
    text += `\n执行时间: ${report.startTime.toLocaleString('zh-CN')}\n`;
    text += `耗时: ${duration.toFixed(2)} 秒\n\n`;

    text += '===== 1. Excel数据读取 =====\n';
    text += `总行数: ${report.excelData.totalRows}\n`;
    text += `有效行数: ${report.excelData.validRows}\n`;
    text += `无效行数: ${report.excelData.invalidRows}\n`;
    if (report.excelData.invalidReasons.length > 0) {
        text += '\n数据验证错误:\n';
        report.excelData.invalidReasons.forEach((err, i) => {
            text += `  ${i + 1}. ${err}\n`;
        });
    }

    text += '\n===== 2. 战队导入 =====\n';
    text += `导入成功: ${report.import.teams.imported}\n`;
    text += `跳过(已存在): ${report.import.teams.skipped}\n`;
    text += `错误: ${report.import.teams.errors.length}\n`;
    if (report.import.teams.errors.length > 0) {
        text += '\n战队导入错误:\n';
        report.import.teams.errors.forEach((err, i) => {
            text += `  ${i + 1}. ${err.name}: ${err.error}\n`;
        });
    }

    text += '\n===== 3. 选手导入 =====\n';
    text += `新增: ${report.import.players.imported}\n`;
    text += `更新: ${report.import.players.updated}\n`;
    text += `跳过: ${report.import.players.skipped}\n`;
    text += `错误: ${report.import.players.errors.length}\n`;
    if (report.import.players.errors.length > 0) {
        text += '\n选手导入错误:\n';
        report.import.players.errors.forEach((err, i) => {
            text += `  ${i + 1}. UID=${err.uid}: ${err.error}\n`;
        });
    }

    text += '\n===== 4. 数据验证 =====\n';
    text += `战队数量验证: 期望 ${report.validation.teams.expected}, 实际 ${report.validation.teams.actual} [${report.validation.teams.matches ? '✓匹配' : '✗不匹配'}]\n`;
    text += `选手数量验证: 期望 ${report.validation.players.expected}, 实际 ${report.validation.players.actual} [${report.validation.players.matches ? '✓匹配' : '✗不匹配'}]\n`;
    text += `数据完整性问题: ${report.validation.dataIntegrity.length} 条\n`;
    if (report.validation.dataIntegrity.length > 0) {
        text += '\n数据完整性问题详情:\n';
        report.validation.dataIntegrity.forEach((item, i) => {
            text += `  ${i + 1}. ID=${item.id}, 选手=${item.playerName}, UID=${item.uid}, 问题=${item.issues}\n`;
        });
    }

    text += '\n===== 5. 导入结果 =====\n';
    const allValid = report.validation.teams.matches && 
                     report.validation.players.matches && 
                     report.validation.dataIntegrity.length === 0 &&
                     report.import.teams.errors.length === 0 &&
                     report.import.players.errors.length === 0;
    text += `总体状态: ${allValid ? '✅ 成功' : '❌ 部分失败'}\n`;
    text += `总导入战队: ${report.import.teams.imported}\n`;
    text += `总导入选手: ${report.import.players.imported + report.import.players.updated}\n`;

    return text;
}

main().catch(err => {
    console.error('导入失败:', err.message);
    report.endTime = new Date();
    report.import.teams.errors.push({ name: '系统', error: err.message });
    fs.writeFileSync('选手数据导入报告.txt', generateReport(), 'utf8');
});
