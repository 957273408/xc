const XLSX = require('xlsx');
const path = require('path');

const filePath = path.join(__dirname, '战队信息表-2026-06-11 16_48_20(1).xlsx');
console.log('读取文件:', filePath);

const workbook = XLSX.readFile(filePath);
const worksheet = workbook.Sheets['Sheet1'];
const data = XLSX.utils.sheet_to_json(worksheet, { header: 1 });

// 解析数据
const headers = data[0]; // ["战队编号","战队缩写","战队名","UID"]
console.log('\n表头:', headers);

// 收集所有战队信息
const teamsMap = new Map();
const players = [];

for (let i = 1; i < data.length; i++) {
    const row = data[i];
    if (row && row.length >= 4) {
        const [teamId, teamAbbr, teamName, uid] = row;
        
        if (teamId && teamName) {
            // 收集战队信息
            if (!teamsMap.has(teamId)) {
                teamsMap.set(teamId, {
                    id: teamId,
                    abbr: teamAbbr || '',
                    name: teamName,
                    memberCount: 0
                });
            }
            teamsMap.get(teamId).memberCount++;
            
            // 收集选手信息
            if (uid) {
                players.push({
                    teamId: teamId,
                    teamName: teamName,
                    uid: uid
                });
            }
        }
    }
}

// 输出战队信息
console.log('\n========== 战队信息 (共' + teamsMap.size + '个) ==========');
const teams = Array.from(teamsMap.values()).sort((a, b) => a.id - b.id);
teams.forEach(team => {
    console.log(`战队${team.id}: ${team.name} (${team.abbr}) - ${team.memberCount}名成员`);
});

// 输出选手信息
console.log('\n========== 选手信息 (共' + players.length + '名) ==========');
players.forEach((player, index) => {
    console.log(`选手${index + 1}: UID=${player.uid}, 战队=${player.teamName}`);
});

// 统计缺失数据
const playersWithUid = players.filter(p => p.uid && p.uid.trim() !== '');
console.log('\n========== 数据完整性检查 ==========');
console.log(`总选手数: ${players.length}`);
console.log(`有效UID选手数: ${playersWithUid.length}`);
console.log(`缺失UID选手数: ${players.length - playersWithUid.length}`);

// 检查每个战队的选手数量
console.log('\n========== 战队成员统计 ==========');
teams.forEach(team => {
    console.log(`${team.name}: ${team.memberCount}名选手`);
});

// 生成可用于导入的数据格式（JSON）
const exportData = {
    teams: teams.map(t => ({
        teamId: t.id,
        teamName: t.name,
        teamAbbr: t.abbr,
        memberCount: t.memberCount
    })),
    players: playersWithUid.map((p, idx) => ({
        playerIndex: idx + 1,
        uid: p.uid,
        teamId: p.teamId,
        teamName: p.teamName
    }))
};

console.log('\n========== JSON导出数据 ==========');
console.log(JSON.stringify(exportData, null, 2));
