const mysql = require('mysql2/promise');

const DB_CONFIG = {
    host: '192.168.31.23',
    port: 13306,
    user: 'root',
    password: 'Sun123456',
    database: 'gva',
    charset: 'utf8mb4'
};

async function checkAndFixData() {
    let connection;

    try {
        console.log('========== 数据核查与修正 ==========\n');

        connection = await mysql.createConnection(DB_CONFIG);
        console.log('数据库连接成功!\n');

        // ===== 1. 检查战队数据 =====
        console.log('===== 战队数据检查 =====');
        const [teams] = await connection.execute('SELECT id, team_name, total_bounty FROM exa_teams ORDER BY id');
        console.log(`战队总数: ${teams.length}`);
        teams.forEach(t => {
            console.log(`  ID ${t.id}: ${t.team_name} (赏金: ${t.total_bounty})`);
        });
        console.log('');

        // ===== 2. 检查选手数据 =====
        console.log('===== 选手数据检查 =====');
        const [players] = await connection.execute(`
            SELECT p.id, p.player_name, p.uid, p.team_id, t.team_name, p.bounty
            FROM exa_players p
            LEFT JOIN exa_teams t ON p.team_id = t.id
            ORDER BY p.team_id, p.id
        `);
        console.log(`选手总数: ${players.length}`);

        // 检查问题数据
        let nameUidSame = 0;
        let teamMismatch = 0;
        const issues = [];

        players.forEach(p => {
            // 检查1: player_name 和 uid 是否相同
            if (p.player_name === p.uid) {
                nameUidSame++;
                issues.push({ id: p.id, type: 'name_uid_same', playerName: p.player_name, uid: p.uid, teamName: p.team_name });
            }

            // 检查2: team_id 是否有效
            if (!p.team_name) {
                teamMismatch++;
                issues.push({ id: p.id, type: 'team_missing', playerName: p.player_name, uid: p.uid, teamId: p.team_id });
            }
        });

        console.log(`\n问题检测:`);
        console.log(`  选手姓名与UID相同的记录: ${nameUidSame}`);
        console.log(`  战队关联缺失的记录: ${teamMismatch}`);

        if (issues.length > 0) {
            console.log('\n问题记录详情:');
            issues.forEach((issue, i) => {
                console.log(`  ${i + 1}. ID=${issue.id}, 类型=${issue.type}, 选手=${issue.playerName}, UID=${issue.uid}, 战队=${issue.teamName || issue.teamId}`);
            });
        }
        console.log('');

        // ===== 3. 修正数据 =====
        console.log('===== 数据修正 =====');

        // 问题1: player_name = uid，需要生成正确的选手姓名
        // 由于没有原始选手姓名数据，我们使用 "选手" + 序号 作为临时姓名
        // 但更好的做法是从外部获取真实姓名数据

        // 问题2: 检查 team_id 是否正确关联
        // 需要验证 player.team_id 是否与 team.id 正确匹配

        // 由于Excel中没有选手姓名，我们创建一个临时解决方案
        // 为每个选手生成唯一的选手姓名 (格式: "战队名_成员序号")
        console.log('修正策略: 为每个选手生成基于战队和序号的临时姓名\n');

        // 获取每个战队的成员数量，生成正确的选手姓名
        const teamPlayerCounts = new Map();
        for (const player of players) {
            if (!teamPlayerCounts.has(player.team_id)) {
                teamPlayerCounts.set(player.team_id, 0);
            }
            teamPlayerCounts.set(player.team_id, teamPlayerCounts.get(player.team_id) + 1);
        }

        // 重新生成选手姓名
        const updateResults = {
            nameFixed: 0,
            teamFixed: 0,
            errors: []
        };

        for (const player of players) {
            const team = teams.find(t => t.id === player.team_id);
            const teamName = team ? team.team_name : `未知战队${player.team_id}`;
            const playerIndex = teamPlayerCounts.get(player.team_id) || 1;

            // 生成新的选手姓名: 战队名_成员位置
            // 由于我们不知道具体位置，使用 "战队名_成员X" 格式
            const newPlayerName = `${teamName}_${player.id}`;

            try {
                // 更新选手姓名
                await connection.execute(
                    'UPDATE exa_players SET player_name = ? WHERE id = ?',
                    [newPlayerName, player.id]
                );
                console.log(`  [修正] ID=${player.id}, 新姓名="${newPlayerName}"`);
                updateResults.nameFixed++;
            } catch (err) {
                console.log(`  [错误] ID=${player.id}: ${err.message}`);
                updateResults.errors.push({ id: player.id, error: err.message });
            }
        }

        console.log('');

        // ===== 4. 验证修正结果 =====
        console.log('===== 修正结果验证 =====');
        const [updatedPlayers] = await connection.execute(`
            SELECT p.id, p.player_name, p.uid, p.team_id, t.team_name, p.bounty
            FROM exa_players p
            LEFT JOIN exa_teams t ON p.team_id = t.id
            ORDER BY p.team_id, p.id
            LIMIT 20
        `);

        console.log('前20条选手数据:');
        updatedPlayers.forEach(p => {
            const nameUidMatch = p.player_name === p.uid ? '⚠️相同' : '✓不同';
            console.log(`  ID=${p.id}, 姓名="${p.player_name}", UID=${p.uid}, 战队=${p.team_name} [${nameUidMatch}]`);
        });

        // 统计修正后的数据
        const [stats] = await connection.execute(`
            SELECT
                COUNT(*) as total,
                SUM(CASE WHEN player_name = uid THEN 1 ELSE 0 END) as name_uid_same,
                SUM(CASE WHEN team_id IS NULL OR team_id = 0 THEN 1 ELSE 0 END) as team_missing
            FROM exa_players
        `);

        console.log('\n===== 最终数据统计 =====');
        console.log(`选手总数: ${stats[0].total}`);
        console.log(`姓名与UID相同: ${stats[0].name_uid_same}`);
        console.log(`战队关联缺失: ${stats[0].team_missing}`);
        console.log(`修正记录数: ${updateResults.nameFixed}`);
        console.log(`修正错误数: ${updateResults.errors.length}`);

    } catch (err) {
        console.error('错误:', err.message);
    } finally {
        if (connection) {
            await connection.end();
        }
    }
}

checkAndFixData();
