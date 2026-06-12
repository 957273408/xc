const XLSX = require('xlsx');
const fs = require('fs');
const path = require('path');

const filePath = path.join(__dirname, '战队信息表-2026-06-11 16_48_20(1).xlsx');
const outputPath = path.join(__dirname, '数据报告.txt');

let report = '========== 战队信息表数据分析报告 ==========\n\n';

const workbook = XLSX.readFile(filePath);
const worksheet = workbook.Sheets['Sheet1'];
const data = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

// 解析数据
const headers = data[0];
report += '表头: ' + JSON.stringify(headers) + '\n\n';

// 收集所有战队信息
const teamsMap = new Map();
const players = [];

for (let i = 1; i < data.length; i++) {
    const row = data[i];
    if (row && row.length >= 4) {
        const [teamId, teamAbbr, teamName, uid] = row;

        if (teamId && teamName) {
            if (!teamsMap.has(teamId)) {
                teamsMap.set(teamId, {
                    id: teamId,
                    abbr: teamAbbr || '',
                    name: teamName,
                    members: []
                });
            }

            if (uid && uid.toString().trim() !== '') {
                teamsMap.get(teamId).members.push(uid.toString());
                players.push({
                    teamId: teamId,
                    teamName: teamName,
                    uid: uid.toString()
                });
            }
        }
    }
}

const teams = Array.from(teamsMap.values()).sort((a, b) => a.id - b.id);

// 输出战队信息
report += '========== 战队列表 (共' + teams.length + '个) ==========\n\n';
teams.forEach(team => {
    report += `战队${team.id}: ${team.name} (${team.abbr})\n`;
    report += `  成员数量: ${team.members.length}\n`;
    report += `  成员UID: ${team.members.join(', ')}\n\n`;
});

// 统计信息
report += '\n========== 数据完整性检查 ==========\n';
report += `总战队数: ${teams.length}\n`;
report += `总选手数: ${players.length}\n`;

// 检查有问题的战队
let issueCount = 0;
teams.forEach(team => {
    if (team.members.length !== 5) {
        report += `警告: 战队"${team.name}"只有${team.members.length}名成员，少于预期的5名\n`;
        issueCount++;
    }
});

if (issueCount === 0) {
    report += '所有战队成员数量正常（每队5人）\n';
}

// 生成JSON格式的完整数据
report += '\n\n========== JSON数据导出 ==========\n';

const exportData = {
    teams: teams.map(t => ({
        teamId: t.id,
        teamName: t.name,
        teamAbbr: t.abbr,
        memberCount: t.members.length
    })),
    players: players.map((p, idx) => ({
        playerIndex: idx + 1,
        uid: p.uid,
        teamId: p.teamId,
        teamName: p.teamName
    }))
};

report += '\n--- 战队数据 ---\n';
report += JSON.stringify(exportData.teams, null, 2);
report += '\n\n--- 选手数据 ---\n';
report += JSON.stringify(exportData.players, null, 2);

// 生成SQL INSERT语句
report += '\n\n========== SQL INSERT语句 (供参考) ==========\n';
report += '\n--- 战队表 (exa_teams) ---\n';
exportData.teams.forEach(t => {
    report += `INSERT INTO exa_teams (team_name, total_bounty) VALUES ('${t.teamName}', 0);\n`;
});

report += '\n--- 选手表 (exa_players) ---\n';
report += '注意: 需要先插入战队获取ID，然后关联选手。以下为选手数据:\n';
exportData.players.forEach((p, idx) => {
    report += `选手${idx + 1}: UID=${p.uid}, 战队ID=${p.teamId} (${p.teamName})\n`;
});

// 写入文件
fs.writeFileSync(outputPath, report, 'utf8');
console.log('报告已生成:', outputPath);

// 同时输出到控制台
console.log(report);
